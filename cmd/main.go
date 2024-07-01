package main

import (
	"github.com/charliegreeny/simple-dating-app/config"
	"github.com/charliegreeny/simple-dating-app/internal/api"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/token"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			config.NewSqlDb,
			validator.New,
			api.NewUserHandler,
			api.NewLoginHandler,
			token.NewCache,
			user.NewService,
			zap.NewProduction,
		),
		fx.Invoke(api.Router)).
		Run()
}
