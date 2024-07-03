package user

import (
	"context"
	"fmt"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"gorm.io/gorm"
)

type Cache struct {
	c  map[string]*Output
	db *gorm.DB
}

func NewCache(db *gorm.DB) app.Cache[string, *Output] {
	c := &Cache{db: db, c: map[string]*Output{}}
	c.fetchAll()
	return c
}

func (c *Cache) fetchAll() app.Cache[string, *Output] {
	//TODO add cron job to repeat this process every x mins
	var users []*Entity
	c.db.Find(&users)
	for _, u := range users {
		c.c[u.ID] = u.toOutput()
	}
	return c
}

func (c *Cache) GetAll(ctx context.Context) []*Output {
	var users []*Output
	for _, v := range c.c {
		if v.ID == ctx.Value(&app.UserCtxKey{}).(*Output).ID {
			continue
		}
		users = append(users, v)
	}
	return users
}

func (c *Cache) Get(_ context.Context, id string) (*Output, error) {
	//TODO replace service getter with this or combine to, if not in cache get from db
	u, ok := c.c[id]
	if !ok {
		return nil, app.ErrNotFound{Message: fmt.Sprintf("user id '%s' not found", id)}
	}
	return u, nil
}

func (c *Cache) Add(_ context.Context, id string, u *Output) error {
	c.c[id] = u
	return nil
}
