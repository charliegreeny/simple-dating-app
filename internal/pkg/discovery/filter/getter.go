package filter

import "github.com/charliegreeny/simple-dating-app/app"

func GetFilters(userCache app.Cache[string, *app.User]) []Filter {
	return []Filter{
		newGender(),
		newAge(),
		newLocations(userCache),
	}
}
