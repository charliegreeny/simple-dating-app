package token

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
)

type Login struct {
	jwtUserCache app.Cache[string, string]
	userCache    app.Cache[string, *app.User]
}

func NewLogin(jwtCache app.Cache[string, string], userCache app.Cache[string, *app.User]) *Login {
	return &Login{jwtUserCache: jwtCache, userCache: userCache}
}

func (l Login) LoginUser(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	u, err := l.userCache.Get(ctx, input.ID)
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
	err = l.jwtUserCache.Add(ctx, token, input.ID)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{Token: token}, nil
}
