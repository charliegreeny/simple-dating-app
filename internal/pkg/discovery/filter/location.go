package filter

import (
	"context"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/appctx"
	"github.com/jftuga/geodist"
)

type location struct {
	userCache app.Cache[string, *app.User]
}

func newLocations(cache app.Cache[string, *app.User]) Filter {
	return &location{userCache: cache}
}

func (l location) Apply(ctx context.Context, current *app.Preference, users []*app.User) []*app.User {
	loc := appctx.GetLocFromCtx(ctx)
	if loc == nil {
		u, err := l.userCache.Get(ctx, appctx.GetUserFromCtx(ctx).ID)
		if err != nil || u.Loc == nil {
			return users
		}
		loc = u.Loc
	}
	userCord := geodist.Coord{
		Lat: loc.Lat,
		Lon: loc.Long,
	}
	var filteredUsers []*app.User
	for _, user := range users {
		if user.Loc == nil {
			continue
		}
		_, distKm := geodist.HaversineDistance(userCord, geodist.Coord{Lat: user.Loc.Lat, Lon: user.Loc.Long})
		intDistKm := int(distKm)
		if intDistKm > current.MaxDistance {
			continue
		}
		if intDistKm == 0 {
			intDistKm = 1
		}
		user.DistanceFrom = intDistKm
		filteredUsers = append(filteredUsers, user)
	}
	return filteredUsers
}
