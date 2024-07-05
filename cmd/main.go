package main

import (
	"github.com/charliegreeny/simple-dating-app/api"
	handler "github.com/charliegreeny/simple-dating-app/api/handler"
	"github.com/charliegreeny/simple-dating-app/api/middleware"
	"github.com/charliegreeny/simple-dating-app/config"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/discovery"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/location"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/match"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/preference"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/token"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user/cache"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user/service"
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
			middleware.NewLocation,
			handler.NewUserHandler,
			handler.NewLoginHandler,
			handler.NewDiscoveryHandler,
			handler.NewPreference,
			token.NewLogin,
			token.NewCache,
			cache.NewCache,
			service.NewService,
			match.NewMatcher,
			preference.NewService,
			location.NewService,
			discovery.NewFilterCombined,
		),
		fx.Invoke(api.Router)).
		Run()
}
