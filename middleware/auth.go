package middleware

import (
	"context"
	"encoding/json"
	"github.com/charliegreeny/simple-dating-app/config"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
	"net/http"
)

type Auth struct {
	c app.Cache[string, *user.Output]
}

func NewAuth(params *config.TokenParams) *Auth {
	return &Auth{c: params.Cache}
}

func (a Auth) Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		jwt := r.Header.Get("Authorization")
		if jwt == "" {
			enc := json.NewEncoder(w)
			w.WriteHeader(http.StatusUnauthorized)
			_ = enc.Encode(app.ErrorOutput{Message: "authorised headers needs to be provided"})
			return
		}
		u, err := a.c.Get(r.Context(), jwt)
		if err != nil {
			enc := json.NewEncoder(w)
			w.WriteHeader(http.StatusUnauthorized)
			_ = enc.Encode(app.ErrorOutput{Message: "user needs to log in to use this service"})
			return
		}
		ctx := r.Context()
		next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, user.CtxKey{}, u)))
	}
	return http.HandlerFunc(fn)
}
