package token

import (
	"context"
	"errors"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
)

type cache struct {
	c map[string]*user.Output `name:"jwtCache"`
}

func NewCache() app.Cache[string, *user.Output] {
	return &cache{c: map[string]*user.Output{}}
}

func (c cache) Get(_ context.Context, jwt string) (*user.Output, error) {
	u, ok := c.c[jwt]
	if !ok || !verifyToken(jwt) {
		return nil, errors.New("jwt not valid")
	}
	return u, nil
}

func (c cache) GetAll(context.Context) []*user.Output {
	return nil
}

func (c cache) Add(_ context.Context, jwt string, u *user.Output) error {
	c.c[jwt] = u
	return nil
}
