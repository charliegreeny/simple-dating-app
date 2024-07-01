package token

import (
	"github.com/charliegreeny/simple-dating-app/internal/model"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
)

var errInvalidToken = model.ErrUnauthorized{Message: "invalid token"}

type Cache struct {
	c map[string]*user.Output
	u model.GetterCreator[*user.Input, *user.Output]
}

func NewCache(userGetter model.GetterCreator[*user.Input, *user.Output]) *Cache {
	return &Cache{c: make(map[string]*user.Output), u: userGetter}
}

func (c Cache) LoginUser(input *LoginInput) (*LoginOutput, error) {
	u, err := c.u.Get(input.ID)
	if err != nil {
		return nil, err
	}
	if input.Password != u.Password {
		return nil, model.ErrUnauthorized{Message: "invalid password"}
	}
	token, err := createToken(u.ID)
	if err != nil {
		return nil, err
	}
	c.c[token] = u
	return &LoginOutput{Token: token}, nil
}

func (c Cache) GetLoginUser(jwt string) (*user.Output, error) {
	ok := verifyToken(jwt)
	if !ok {
		return nil, errInvalidToken
	}
	u, ok := c.c[jwt]
	if !ok {
		return nil, errInvalidToken
	}
	return u, nil
}
