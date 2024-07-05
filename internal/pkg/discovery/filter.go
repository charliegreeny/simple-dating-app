package discovery

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/discovery/filter"
)

type filterCombined struct {
	filters []filter.Filter
}

func NewFilterCombined(userCache app.Cache[string, *app.User]) filter.Filter {
	return &filterCombined{filters: filter.GetFilters(userCache)}
}

func (fc filterCombined) Apply(ctx context.Context, current *app.Preference, users []*app.User) []*app.User {
	u := users
	for _, f := range fc.filters {
		u = f.Apply(ctx, current, u)
	}
	return u
}
