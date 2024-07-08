package location

import (
	"context"
	"errors"
	"fmt"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/appctx"
	"github.com/go-sql-driver/mysql"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type service struct {
	userCache app.Cache[string, *app.User]
	db        *gorm.DB
	log       *zap.Logger
}

func NewService(userCache app.Cache[string, *app.User], db *gorm.DB, log *zap.Logger) app.EntityService[*app.Location, *app.Location] {
	return &service{userCache: userCache, db: db, log: log}
}

func (s service) Create(ctx context.Context, input *app.Location) (*app.Location, error) {
	u := appctx.GetUserFromCtx(ctx)
	input.UserID = u.ID
	err := s.db.Create(input).Error
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
	go func(u *app.User, location *app.Location) {
		u.Loc = input
		s.updateUserCache(u, ctx, input)
	}(u, input)
	return input, nil
}

func (s service) Get(ctx context.Context, id string) (*app.Location, error) {
	u, err := s.userCache.Get(ctx, id)
	if err != nil || u.Loc == nil {
		s.log.Debug("could not find location in cache, querying db", zap.String("id", id), zap.Error(err))
		var loc *app.Location
		err = s.db.Find(&loc, "user_id = ?", id).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, app.ErrNotFound{Message: fmt.Sprintf("no location record found for id %s", id)}
			}
			return nil, fmt.Errorf("db error getting location: %w", err)
		}
	}
	return u.Loc, nil
}

func (s service) Update(ctx context.Context, input *app.Location) (*app.Location, error) {
	u := appctx.GetUserFromCtx(ctx)
	if u == nil {
		return nil, errors.New("invalid data: could not find user updating location")
	}
	input.UserID = u.ID
	if u.Loc == nil {
		return s.Create(ctx, input)
	}
	if input == u.Loc {
		return u.Loc, nil
	}
 	u.Loc = input
	go func(u *app.User) {
		s.updateUserCache(u, ctx, input)
	}(u)
	err := s.db.Save(&input).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app.ErrNotFound{
				Message: fmt.Sprintf("no user  to update location db for id %s: %s", input.UserID, err.Error())}
		}
		var sqlErr *mysql.MySQLError
		ok := errors.As(err, &sqlErr)
		if ok && sqlErr.Number == app.ForeignKeyErrCode {
			return nil, app.ErrBadRequest{
				Message: fmt.Sprintf("could not update location for id in db %s: %s", input.UserID, err.Error())}
		}
		return nil, err
	}
	return input, nil
}

func (s service) updateUserCache(u *app.User, ctx context.Context, input *app.Location) {
	err := s.userCache.Add(ctx, input.UserID, u)
	if err != nil {
		s.log.Warn("failed to update user location in cache", zap.Error(err))
	}
}
