package filter

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
	"slices"
)

type age struct {
}

func newAge() Filter {
	return &age{}
}

func (a age) Apply(_ context.Context, preference *app.Preference, users []*app.User) []*app.User {
	return slices.DeleteFunc(users, func(u *app.User) bool {
		if preference.MaxAge == nil {
			return u.Age < preference.MinAge
		}
		return u.Age > *preference.MaxAge || u.Age < preference.MinAge
	})
}
