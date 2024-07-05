package filter

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/stretchr/testify/assert"
	"testing"
)

func i(i int) *int {
	return &i
}

func Test_age_Apply(t *testing.T) {
	user25yo := &app.User{UserOutput: &app.UserOutput{Age: 25}}
	user50yo := &app.User{UserOutput: &app.UserOutput{Age: 50}}
	user18yo := &app.User{UserOutput: &app.UserOutput{Age: 18}}
	user30yo := &app.User{UserOutput: &app.UserOutput{Age: 30}}

	tests := []struct {
		name    string
		perf    *app.Preference
		want    []*app.User
		wantLen int
	}{
		{
			name:    "max age 29 min age 18 returns 2 users",
			perf:    &app.Preference{MaxAge: i(29), MinAge: 18},
			want:    []*app.User{user25yo, user18yo},
			wantLen: 2,
		},
		{
			name:    "max age 100 min age 18 returns 4 users",
			perf:    &app.Preference{MaxAge: i(100), MinAge: 18},
			want:    []*app.User{user25yo, user50yo, user18yo, user30yo},
			wantLen: 2,
		},
		{
			name:    "no max age, min age 50 returns 4 users",
			perf:    &app.Preference{MinAge: 50},
			want:    []*app.User{user50yo},
			wantLen: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stubUsers := []*app.User{user25yo, user50yo, user18yo, user30yo}
			a := age{}

			assert.Equalf(t, tt.want, a.Apply(context.Background(), tt.perf, stubUsers),
				"Apply(%v, %v, %v)", context.Background(), tt.perf, stubUsers)
		})
	}
}
