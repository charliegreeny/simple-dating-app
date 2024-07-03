package user

import (
	"context"
	"errors"
	"fmt"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type service struct {
	db      *gorm.DB
	log     *zap.Logger
	cache   app.Cache[string, *Output]
	uuidGen app.UUIDGenerator
}

const duplicateKeyErrCode = 1062
const foreignKeyErrCode = 1452

func NewService(db *gorm.DB, l *zap.Logger, cache app.Cache[string, *Output]) app.GetterCreator[*Input, *Output] {
	return &service{db: db, log: l, uuidGen: app.DefaultUUIDGenerator, cache: cache}
}

func (s service) Get(_ context.Context, id string) (*Output, error) {
	s.log.Debug("Getting user", zap.String("id", id))
	var e *Entity
	r := s.db.First(&e, "id = ?", id)
	if r.Error != nil {
		if errors.Is(r.Error, gorm.ErrRecordNotFound) {
			return nil, app.ErrNotFound{Message: fmt.Sprintf("no user found for id %s", id)}
		}
		return nil, fmt.Errorf("db error getting user: %w", r.Error)
	}
	return e.toOutput(), nil
}

func (s service) Create(ctx context.Context, input *Input) (*Output, error) {
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
		if ok && sqlErr.Number == duplicateKeyErrCode {
			return nil, app.ErrBadRequest{Message: "email already in use"}
		}
		if ok && sqlErr.Number == foreignKeyErrCode {
			return nil, app.ErrBadRequest{Message: "gender provided not supported"}
		}
		return nil, err
	}
	o := e.toOutput()
	go s.updateCache(ctx, o)
	return o, nil
}

func (s service) updateCache(ctx context.Context, o *Output) {
	func(o *Output) {
		err1 := s.cache.Add(ctx, o.ID, o)
		if err1 != nil {
			s.log.Error("could not add new user to cache", zap.Error(err1))
		}
	}(o)
}
