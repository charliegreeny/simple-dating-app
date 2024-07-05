package middleware

import (
	"encoding/json"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/appctx"
	"net/http"
)

type Auth struct {
	tokenCache app.Cache[string, string]
	userCache  app.Cache[string, *app.User]
}

func NewAuth(tokenCache app.Cache[string, string], userCache app.Cache[string, *app.User]) *Auth {
	return &Auth{tokenCache: tokenCache, userCache: userCache}
}

func (a Auth) Auth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		jwt := r.Header.Get("Authorization")
		if jwt == "" {
			enc := json.NewEncoder(w)
			w.WriteHeader(http.StatusUnauthorized)
			_ = enc.Encode(app.ErrorOutput{Message: "'Authorization' headers needs to be provided"})
			return
		}
		ctx := r.Context()
		uID, err := a.tokenCache.Get(ctx, jwt)
		if err != nil {
			enc := json.NewEncoder(w)
			w.WriteHeader(http.StatusUnauthorized)
			_ = enc.Encode(app.ErrorOutput{Message: "user needs to log in to use this service"})
			return
		}
		u, err := a.userCache.Get(ctx, uID)
		if err != nil {
			enc := json.NewEncoder(w)
			w.WriteHeader(http.StatusUnauthorized)
			_ = enc.Encode(app.ErrorOutput{Message: "could not find logged in user"})
			return
		}
		next.ServeHTTP(w, r.WithContext(appctx.AddUserToCtx(ctx, u)))
	}
	return http.HandlerFunc(fn)
}
