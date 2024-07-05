package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/appctx"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/preference"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type service struct {
	db      *gorm.DB
	log     *zap.Logger
	cache   app.Cache[string, *app.User]
	uuidGen app.UUIDGenerator
	p       app.EntityService[*app.Preference, *app.Preference]
}

func NewService(
	db *gorm.DB,
	l *zap.Logger,
	cache app.Cache[string, *app.User],
	p app.EntityService[*app.Preference, *app.Preference]) app.GetterCreator[*Input, *app.UserOutput] {

	return &service{
		db:      db,
		log:     l,
		uuidGen: app.DefaultUUIDGenerator,
		p:       p,
		cache:   cache,
	}
}

func (s service) Get(ctx context.Context, id string) (*app.UserOutput, error) {
	s.log.Debug("Getting user", zap.String("id", id))
	user, err := s.cache.Get(ctx, id)
	if err != nil {
		s.log.Info("user not in cache, fetching from db", zap.String("id", id), zap.Error(err))
		var e *Entity
		r := s.db.First(&e, "id = ?", id)
		if r.Error != nil {
			if errors.Is(r.Error, gorm.ErrRecordNotFound) {
				return nil, app.ErrNotFound{Message: fmt.Sprintf("no user found for id %s", id)}
			}
			return nil, fmt.Errorf("db error getting user: %w", r.Error)
		}
		return e.ToOutput(), nil
	}
	return user.UserOutput, err
}

func (s service) Create(ctx context.Context, input *Input) (*app.UserOutput, error) {
	dob, err := time.Parse(time.DateOnly, input.Dob)
	if err != nil {
		return nil, app.ErrBadRequest{Message: "invalid date format provided - needs to be yyyy-mm-dd"}
	}
	if calcAge(dob) < 18 {
		return nil, app.ErrForbidden{Message: "users must be at least 18 years old"}
	}
	e := &Entity{
		ID:          s.uuidGen(),
		Name:        input.Name,
		Gender:      input.Gender,
		DateOfBirth: dob,
		Email:       input.Email,
		Password:    input.Password,
	}
	s.log.Debug("Creating user", zap.String("id", e.ID))
	err = s.db.Create(e).Error
	if err != nil {
		var sqlErr *mysql.MySQLError
		ok := errors.As(err, &sqlErr)
		if ok && sqlErr.Number == app.DuplicateKeyErrCode {
			return nil, app.ErrBadRequest{Message: "email already in use"}
		}
		if ok && sqlErr.Number == app.ForeignKeyErrCode {
			return nil, app.ErrBadRequest{Message: "gender provided not supported"}
		}
		return nil, err
	}
	o := e.ToOutput()
	go func(ctx context.Context, o *app.UserOutput) {
		s.updateCache(ctx, s.createUser(ctx, o))
	}(ctx, o)
	return o, nil
}

func (s service) createUser(ctx context.Context, o *app.UserOutput) *app.User {
	u := &app.User{
		UserOutput: o,
	}
	loc := appctx.GetLocFromCtx(ctx)
	if loc != nil {
		u.Loc = loc
	}
	p, err := s.p.Create(ctx, preference.DefaultPreferences(o.ID, o.Gender))
	if err != nil {
		s.log.Error("could not add default preferences for new user", zap.String("id", u.ID), zap.Error(err))
	}
	u.Pref = p
	return u
}

func (s service) updateCache(ctx context.Context, u *app.User) {
	err := s.cache.Add(ctx, u.ID, u)
	if err != nil {
		s.log.Error("could not add new user to cache", zap.Error(err))
	}
}
