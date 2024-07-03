package token

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/config"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
)

var errInvalidToken = app.ErrUnauthorized{Message: "invalid token"}

type Login struct {
	jwtUserCache app.Cache[string, *user.Output] `name:"jwtCache"`
	u            app.GetterCreator[*user.Input, *user.Output]
}

func NewLogin(params *config.TokenParams) *Login {
	return &Login{jwtUserCache: params.Cache, u: params.UserGetterCreator}
}

func (l Login) LoginUser(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	u, err := l.u.Get(ctx, input.ID)
	if err != nil {
		return nil, err
	}
	if input.Password != u.Password {
		return nil, app.ErrUnauthorized{Message: "invalid password"}
	}
	token, err := createToken(u.ID)
	if err != nil {
		return nil, err
	}
	err = l.jwtUserCache.Add(ctx, token, u)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{Token: token}, nil
}
