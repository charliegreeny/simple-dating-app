package token

import (
	"context"
	"errors"
	"github.com/charliegreeny/simple-dating-app/app"
)

type cache struct {
	c map[string]string
}

func NewCache() app.Cache[string, string] {
	return &cache{c: map[string]string{}}
}

func (c cache) Get(_ context.Context, jwt string) (string, error) {
	uId, ok := c.c[jwt]
	if !ok || !verifyToken(jwt) {
		return "", errors.New("invalid token")
	}
	return uId, nil
}

func (c cache) GetAll(context.Context) []string {
	return nil
}

func (c cache) Add(_ context.Context, jwt string, id string) error {
	c.c[jwt] = id
	return nil
}
