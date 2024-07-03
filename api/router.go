package api

import (
	appMiddleware "github.com/charliegreeny/simple-dating-app/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func Router(auth *appMiddleware.Auth, u *UserHandler, l *LoginHandler, d *DiscoveryHandler) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(appMiddleware.ResponseHeader)

	r.Post("/login", l.login)

	r.Group(func(r chi.Router) {
		r.Use(auth.Auth)
		r.Post("/swipe", d.PostSwipeHandler)
		r.Get("/discover", d.GetEligibleUsersHandler)
	})

	r.Mount("/user", userRoutes(r, u))

	return http.ListenAndServe(":8080", r)
}

func userRoutes(r *chi.Mux, u *UserHandler) http.Handler {
	r.Post("/create", u.createUserHandler)
	r.Get("/{id}", u.getUserHandler)
	return r
}
