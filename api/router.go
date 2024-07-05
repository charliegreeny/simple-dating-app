package api

import (
	handler "github.com/charliegreeny/simple-dating-app/api/handler"
	customMiddleware "github.com/charliegreeny/simple-dating-app/api/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func Router(auth *customMiddleware.Auth, location *customMiddleware.Location, u handler.User, l handler.Login, d handler.Discovery, p handler.Preference) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(customMiddleware.ResponseHeader)

	r.Post("/login", l.Login)
	r.Post("/user/create", u.CreateUserHandler)

	r.Group(func(r chi.Router) {
		r.Use(auth.Auth)
		r.Use(location.SetLocationCtx)
		r.Post("/swipe", d.PostSwipeHandler)
		r.Get("/discover", d.GetEligibleUsersHandler)
		r.Get("/user", u.GetUserHandler)
		r.Patch("/user/preference", p.UpdatePreferenceHandler)
	})

	return http.ListenAndServe(":8080", r)
}
