package api

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
)

func Router(u *UserHandler, l *LoginHandler) error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(responseHeader)

	r.Post("/login", l.login)

	r.Mount("/user", userRoutes(r, u))

	return http.ListenAndServe(":8080", r)
}

func userRoutes(r *chi.Mux, u *UserHandler) http.Handler {
	r.Post("/create", u.createUserHandler)
	r.Get("/{id}", u.getUserHandler)
	return r
}

func responseHeader(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
