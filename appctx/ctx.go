package appctx

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
)

type userCtxKey struct{}
type locCtxKey struct{}

func GetUserFromCtx(ctx context.Context) *app.User {
	u, ok := ctx.Value(&userCtxKey{}).(*app.User)
	if !ok {
		return nil
	}
	return u
}

func AddUserToCtx(ctx context.Context, u *app.User) context.Context {
	return context.WithValue(ctx, &userCtxKey{}, u)
}

func GetLocFromCtx(ctx context.Context) *app.Location {
	l, ok := ctx.Value(&locCtxKey{}).(*app.Location)
	if !ok {
		return nil
	}
	return l
}

func AddLocToCtx(ctx context.Context, loc *app.Location) context.Context {
	return context.WithValue(ctx, &locCtxKey{}, loc)
}
