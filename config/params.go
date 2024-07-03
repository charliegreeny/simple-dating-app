package config

import (
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
)

type TokenParams struct {
	Cache             app.Cache[string, *user.Output]
	UserGetterCreator app.GetterCreator[*user.Input, *user.Output]
}
