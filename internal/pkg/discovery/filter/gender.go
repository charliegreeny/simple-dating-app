package filter

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
	"slices"
)

type gender struct{}

func newGender() Filter {
	return &gender{}
}

func (g gender) Apply(_ context.Context, preference *app.Preference, users []*app.User) []*app.User {
	return slices.DeleteFunc(users, func(u *app.User) bool {
		return u.Gender != preference.Gender
	})
}
