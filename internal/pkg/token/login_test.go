package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/charliegreeny/simple-dating-app/app"
	appMock "github.com/charliegreeny/simple-dating-app/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockJWTCache struct {
	mock.Mock
}

func (m *mockJWTCache) Get(ctx context.Context, key string) (string, error) {
	m.Called(ctx, key)
	args := m.Called(ctx, key)
	return args.String(0), args.Error(1)
}

func (m *mockJWTCache) GetAll(ctx context.Context) []string {
	//TODO implement me
	panic("implement me")
}

func (m *mockJWTCache) Add(ctx context.Context, key string, v string) error {
	m.Called(ctx, key, v)
	args := m.Called(ctx, key, v)
	return args.Error(0)
}

func TestLogin_LoginUser(t *testing.T) {
	tests := []struct {
		name             string
		input            *LoginInput
		stubUser         *app.UserOutput
		stubUserCacheErr error
		stubJWTCacheErr  error
		cacheCalls       int
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "successfully login user and returns token and nil error",
			input: &LoginInput{
				ID:       "id",
				Password: "pwd",
			},
			stubUser: &app.UserOutput{
				ID:       "id",
				Password: "pwd",
			},
			stubUserCacheErr: nil,
			stubJWTCacheErr:  nil,
			cacheCalls:       1,
			wantErr:          assert.NoError,
		},
		{
			name: "wrong password and returns no token and ErrUnauthorized",
			input: &LoginInput{
				ID:       "id",
				Password: "pwd",
			},
			stubUser: &app.UserOutput{
				ID:       "id",
				Password: "strongerPassword",
			},
			stubUserCacheErr: nil,
			stubJWTCacheErr:  nil,
			cacheCalls:       0,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, app.ErrUnauthorized{}) && err.Error() == "invalid password"
			},
		},
		{
			name: "ErrNotFound from userGetter and returns no token and ErrNotFound",
			input: &LoginInput{
				ID:       "id",
				Password: "pwd",
			},
			stubUser:         nil,
			stubUserCacheErr: app.ErrNotFound{Message: "user not found"},
			stubJWTCacheErr:  nil,
			cacheCalls:       0,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return errors.Is(err, app.ErrNotFound{}) && err.Error() == "user not found"
			},
		},
		{
			name: "Error returned from Cache and returns no token and error",
			input: &LoginInput{
				ID:       "id",
				Password: "pwd",
			},
			stubUser: &app.UserOutput{
				ID:       "id",
				Password: "pwd",
			},
			stubUserCacheErr: nil,
			stubJWTCacheErr:  errors.New("cache error"),
			cacheCalls:       1,
			wantErr:          assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCacheMock := &appMock.UserCache{}
			userCacheMock.On("Get", context.Background(), tt.input.ID).
				Return(tt.stubUser, tt.stubUserCacheErr).Once()

			jwtUserCache := &mockJWTCache{}
			jwtUserCache.On("Add", mock.Anything, mock.Anything,
				mock.MatchedBy(func(user *app.UserOutput) bool {
					return assert.Equal(t, tt.stubUser, user)
				}),
			).Return(tt.stubJWTCacheErr).Times(tt.cacheCalls)
			l := Login{
				jwtUserCache: jwtUserCache,
				userCache:    userCacheMock,
			}
			got, err := l.LoginUser(context.Background(), tt.input)
			if !tt.wantErr(t, err) {
				return
			}
			if err == nil {
				assert.NotEqual(t, got, nil, fmt.Sprintf("response is nil"))
				assert.NotEqual(t, "", got.Token, fmt.Sprintf("token return empty"))
			}
		})
	}
}
