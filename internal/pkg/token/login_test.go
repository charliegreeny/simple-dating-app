package token

import (
	"context"
	"errors"
	"fmt"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type mockCache struct {
	mock.Mock
}

func (m *mockCache) Get(_ context.Context, key string) (*user.Output, error) {
	//TODO implement me
	panic("implement me")
}

func (m *mockCache) GetAll(_ context.Context) []*user.Output {
	//TODO implement me
	panic("implement me")
}

func (m *mockCache) Add(ctx context.Context, key string, u *user.Output) error {
	args := m.Called(ctx, key, u)
	return args.Error(0)
}

type mockGetter struct {
	mock.Mock
}

func (m *mockGetter) Get(_ context.Context, id string) (*user.Output, error) {
	args := m.Called(id)
	return args.Get(0).(*user.Output), args.Error(1)
}

func TestLogin_LoginUser(t *testing.T) {
	tests := []struct {
		name          string
		input         *LoginInput
		stubUser      *user.Output
		stubGetterErr error
		stubCacheErr  error
		cacheCalls    int
		wantErr       assert.ErrorAssertionFunc
	}{
		{
			name: "successfully login user and returns token and nil error",
			input: &LoginInput{
				ID:       "id",
				Password: "pwd",
			},
			stubUser: &user.Output{
				ID:       "id",
				Password: "pwd",
			},
			stubGetterErr: nil,
			stubCacheErr:  nil,
			cacheCalls:    1,
			wantErr:       assert.NoError,
		},
		{
			name: "wrong password and returns no token and ErrUnauthorized",
			input: &LoginInput{
				ID:       "id",
				Password: "pwd",
			},
			stubUser: &user.Output{
				ID:       "id",
				Password: "strongerPassword",
			},
			stubGetterErr: nil,
			stubCacheErr:  nil,
			cacheCalls:    0,
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
			stubUser:      nil,
			stubGetterErr: app.ErrNotFound{Message: "user not found"},
			stubCacheErr:  nil,
			cacheCalls:    0,
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
			stubUser: &user.Output{
				ID:       "id",
				Password: "pwd",
			},
			stubGetterErr: nil,
			stubCacheErr:  errors.New("cache error"),
			cacheCalls:    1,
			wantErr:       assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getterMock := &mockGetter{}
			getterMock.On("Get", mock.Anything).Return(tt.stubUser, tt.stubGetterErr).Once()

			jwtUserCache := &mockCache{}
			jwtUserCache.On("Add", mock.Anything, mock.Anything,
				mock.MatchedBy(func(user *user.Output) bool {
					return assert.Equal(t, tt.stubUser, user)
				}),
			).Return(tt.stubCacheErr).Times(tt.cacheCalls)
			l := Login{
				jwtUserCache: jwtUserCache,
				userGetter:   getterMock,
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
