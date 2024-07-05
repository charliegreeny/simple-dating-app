package mock

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/stretchr/testify/mock"
)

type UserCache struct {
	mock.Mock
}

func (u *UserCache) Get(ctx context.Context, key string) (*app.User, error) {
	u.Called(ctx, key)
	args := u.Called(ctx, key)
	return args.Get(0).(*app.User), args.Error(1)
}

func (u *UserCache) GetAll(ctx context.Context) []*app.User {
	u.Called(ctx)
	args := u.Called(ctx)
	return args.Get(0).([]*app.User)
}

func (u *UserCache) Add(ctx context.Context, key string, v *app.User) error {
	u.Called(ctx, key, v)
	args := u.Called(ctx, key, v)
	return args.Error(0)
}
