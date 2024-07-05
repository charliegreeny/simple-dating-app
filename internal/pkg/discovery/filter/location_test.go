package filter

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/appctx"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCache struct{}

func (m mockCache) Get(_ context.Context, _ string) (*app.User, error) {
	return &app.User{
		Loc: &app.Location{Lat: 53.43571503300695, Long: -2.260886698797288},
	}, nil
}

func (m mockCache) GetAll(ctx context.Context) []*app.User {
	return nil
}

func (m mockCache) Add(ctx context.Context, key string, v *app.User) error {
	return nil
}

func emptySlice() []*app.User {
	var s []*app.User
	return s
}

func Test_location_Apply(t *testing.T) {
	type args struct {
		current *app.Preference
		users   []*app.User
	}
	tests := []struct {
		ctx  context.Context
		name string
		args args
		want []*app.User
	}{
		{
			name: "user loc (from cache) within distance returns 1 user",
			ctx:  appctx.AddUserToCtx(context.Background(), &app.User{UserOutput: &app.UserOutput{ID: "id"}}),
			args: args{
				current: &app.Preference{MaxDistance: 10},
				users:   []*app.User{{Loc: &app.Location{Lat: 53.43635019762404, Long: -2.264665168705832}}},
			},
			want: []*app.User{{Loc: &app.Location{Lat: 53.43635019762404, Long: -2.264665168705832}, DistanceFrom: 1}},
		},
		{
			name: "user (from ctx) within distance returns 1 user",
			ctx:  appctx.AddLocToCtx(context.Background(), &app.Location{Lat: 53.43571503300695, Long: -2.260886698797288}),
			args: args{
				current: &app.Preference{MaxDistance: 10},
				users:   []*app.User{{Loc: &app.Location{Lat: 51.5032, Long: -0.1195}}},
			},
			want: emptySlice(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := location{
				userCache: &mockCache{},
			}
			assert.Equalf(t, tt.want, l.Apply(tt.ctx, tt.args.current, tt.args.users),
				"Apply(%v, %v, %v)", tt.ctx, tt.args.current, tt.args.users)
		})
	}
}
