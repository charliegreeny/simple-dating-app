package preference

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
	db        *gorm.DB
	userCache app.Cache[string, *app.User]
	log       *zap.Logger
}

func NewService(db *gorm.DB,
	userCache app.Cache[string, *app.User],
	log *zap.Logger) app.EntityService[*app.Preference, *app.Preference] {
	return &service{db: db, userCache: userCache, log: log}
}

func (p service) Create(ctx context.Context, input *app.Preference) (*app.Preference, error) {
	u := appctx.GetUserFromCtx(ctx)
	if input.UserID == "" {
		if u == nil {
			return nil, app.ErrNotFound{Message: "could not find logged in user"}
		}
		input.UserID = u.ID
	}
	if u != nil {
		go func(u *app.User, input *app.Preference) {
			u.Pref = input
			err := p.userCache.Add(ctx, input.UserID, u)
			if err != nil {
				p.log.Warn("failed to update user location in cache", zap.Error(err))
			}
		}(u, input)
	}
	err := p.db.Create(input).Error
	if err != nil {
		var sqlErr *mysql.MySQLError
		ok := errors.As(err, &sqlErr)
		if ok && sqlErr.Number == app.ForeignKeyErrCode {
			return nil, app.ErrBadRequest{
				Message: fmt.Sprintf("could not add preferencs for id %s: %s", input.UserID, err.Error())}
		}
		return nil, err
	}
	return input, nil
}

func (p service) Get(ctx context.Context, id string) (*app.Preference, error) {
	u, err := p.userCache.Get(ctx, id)
	if err != nil || u.Pref == nil {
		p.log.Info("could not find preferences in cache, fetching from db",
			zap.String("id", id), zap.Error(err))
		var e *app.Preference
		r := p.db.First(&e, "user_id = ?", id)
		if r.Error != nil {
			if errors.Is(r.Error, gorm.ErrRecordNotFound) {
				return nil, app.ErrNotFound{Message: fmt.Sprintf("no perferences found for id %s", id)}
			}
			return nil, fmt.Errorf("db error getting preference: %w", r.Error)
		}
		return e, nil
	}
	return u.Pref, nil
}

func (p service) Update(ctx context.Context, input *app.Preference) (*app.Preference, error) {
	user := appctx.GetUserFromCtx(ctx)
	if user == nil {
		return nil, errors.New("invalid data: could not find user updating preferences")
	}
	if input == user.Pref {
		return user.Pref, nil
	}
	input.UserID = user.ID
	user.Pref = setAllFields(input, user.Pref)
	go func(u *app.User) {
		err := p.userCache.Add(ctx, input.UserID, u)
		if err != nil {
			p.log.Warn("failed to update user location in cache", zap.Error(err))
		}
	}(user)
	err := p.db.Save(&user.Pref).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, app.ErrNotFound{
				Message: fmt.Sprintf("no user to update service db for id %s: %s", input.UserID, err.Error())}
		}
		var sqlErr *mysql.MySQLError
		ok := errors.As(err, &sqlErr)
		if ok && sqlErr.Number == app.ForeignKeyErrCode {
			return nil, app.ErrBadRequest{
				Message: fmt.Sprintf("could not update service in for for i %s: %s", input.UserID, err.Error())}
		}
		return nil, err
	}
	return user.Pref, nil
}

func setAllFields(input *app.Preference, current *app.Preference) *app.Preference {
	if input.MaxAge == nil {
		input.MaxAge = current.MaxAge
	}
	if input.Gender == "" {
		input.Gender = current.Gender
	}
	if input.MinAge == 0 {
		input.MinAge = current.MinAge
	}
	if input.MaxDistance == 0 {
		input.MaxDistance = current.MaxDistance
	}
	return input
}

func DefaultPreferences(userID, userGender string) *app.Preference {
	g := ""
	if userGender == app.Male {
		g = app.Female
	}
	if userGender == app.Female {
		g = app.Male
	}
	return &app.Preference{
		UserID:      userID,
		Gender:      g,
		MinAge:      18,
		MaxDistance: 100,
	}
}
