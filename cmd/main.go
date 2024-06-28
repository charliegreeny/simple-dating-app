package main

import (
	"github.com/charliegreeny/simple-dating-app/internal/api"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(),
		fx.Invoke(
			api.Router)).
		Run()
}
