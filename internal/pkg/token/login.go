package token

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/config"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
)

type Login struct {
	jwtUserCache app.Cache[string, *user.Output]
	userGetter   app.IDGetter[*user.Output]
}

func NewLogin(params *config.TokenParams) *Login {
	return &Login{jwtUserCache: params.Cache, userGetter: params.UserGetterCreator}
}

func (l Login) LoginUser(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	u, err := l.userGetter.Get(ctx, input.ID)
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
