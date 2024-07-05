package filter

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
)

type Filter interface {
	Apply(ctx context.Context, current *app.Preference, users []*app.User) []*app.User
}
