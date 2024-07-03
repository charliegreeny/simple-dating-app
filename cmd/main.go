package main

import (
	"github.com/charliegreeny/simple-dating-app/api"
	"github.com/charliegreeny/simple-dating-app/config"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/token"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
	"github.com/charliegreeny/simple-dating-app/middleware"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			zap.NewProduction,
			validator.New,
			config.NewSqlDb,
			middleware.NewAuth,
			api.NewUserHandler,
			api.NewLoginHandler,
			api.NewDiscoveryHandler,
			token.NewLogin,
			user.NewCache,
			newTokenParam,
			user.NewService,
		),
		fx.Invoke(api.Router)).
		Run()
}

func newTokenParam(gc app.GetterCreator[*user.Input, *user.Output]) *config.TokenParams {
	return &config.TokenParams{
		Cache:             token.NewCache(),
		UserGetterCreator: gc,
	}
}
