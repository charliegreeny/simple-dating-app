package cache

import (
	"context"
	"fmt"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/appctx"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user/service"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type cache struct {
	c   map[string]*app.User
	db  *gorm.DB
	log *zap.Logger
}

func NewCache(db *gorm.DB, l *zap.Logger) (app.Cache[string, *app.User], error) {
	c := &cache{
		db:  db,
		c:   map[string]*app.User{},
		log: l,
	}
	err := c.fetchAll()
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (c *cache) fetchAll() error {
	//TODO add cron job to repeat this process every x mins
	type joinRow struct {
		ID          string    `gorm:"column:id; primaryKey"`
		Name        string    `gorm:"column:name"`
		Gender      string    `gorm:"column:gender"`
		DateOfBirth time.Time `gorm:"column:date_of_birth"`
		Email       string    `gorm:"column:email"`
		Password    string    `gorm:"column:password"`
		PrefGender  string    `gorm:"column:preference_gender"`
		MinAge      int       `gorm:"column:min_age"`
		MaxAge      *int      `gorm:"column:max_age"`
		MaxDistance int       `gorm:"column:max_distance"`
		Lat         float64   `gorm:"column:lat"`
		Long        float64   `gorm:"column:long"`
	}
	var joinRows []joinRow
	r := c.db.Table("users").Select("*").Joins("INNER JOIN preferences ON preferences.user_id = users.id").
		Joins("INNER JOIN locations ON locations.user_id = users.id").Scan(&joinRows)
	if r.Error != nil {
		return r.Error
	}
	for _, row := range joinRows {
		c.c[row.ID] = &app.User{
			UserOutput: service.Entity{
				ID:          row.ID,
				Name:        row.Name,
				Gender:      row.Gender,
				DateOfBirth: row.DateOfBirth,
				Email:       row.Email,
				Password:    row.Password,
			}.ToOutput(),
			Pref: &app.Preference{
				UserID:      row.ID,
				Gender:      row.PrefGender,
				MinAge:      row.MinAge,
				MaxAge:      row.MaxAge,
				MaxDistance: row.MaxDistance,
			},
			Loc: &app.Location{
				UserID: row.ID,
				Lat:    row.Lat,
				Long:   row.Long,
			},
		}
	}
	return nil
}

func (c *cache) GetAll(ctx context.Context) []*app.User {
	var users []*app.User
	for _, v := range c.c {
		if v.ID == appctx.GetUserFromCtx(ctx).ID {
			continue
		}
		users = append(users, v)
	}
	return users
}

func (c *cache) Get(_ context.Context, id string) (*app.User, error) {
	//TODO replace service getter with this or combine to, if not in cache get from db
	u, ok := c.c[id]
	if !ok {
		return nil, app.ErrNotFound{Message: fmt.Sprintf("user id '%s' not found", id)}
	}
	return u, nil
}

func (c *cache) Add(_ context.Context, id string, u *app.User) error {
	c.c[id] = u
	return nil
}
