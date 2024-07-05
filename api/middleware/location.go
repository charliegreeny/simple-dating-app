package middleware

import (
	"encoding/json"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/appctx"
	"go.uber.org/zap"
	"net/http"
	"strings"
)

type Location struct {
	service app.EntityService[*app.Location, *app.Location]
	log     *zap.Logger
}

func NewLocation(l app.EntityService[*app.Location, *app.Location], log *zap.Logger) *Location {
	return &Location{service: l, log: log}
}

func (l Location) SetLocationCtx(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		locInput := r.Header.Get("x-user-location")
		if locInput == "" {
			next.ServeHTTP(w, r)
			return
		}
		var loc *app.Location
		if err := json.NewDecoder(strings.NewReader(locInput)).Decode(&loc); err != nil {
			l.log.Warn("Failed to parse x-user-location header", zap.String("location", locInput))
			next.ServeHTTP(w, r)
			return
		}
		ctx := r.Context()
		u := appctx.GetUserFromCtx(ctx)
		if u == nil {
			l.log.Warn("No user in ctx", zap.String("location", locInput))
			next.ServeHTTP(w, r)
			return
		}
		loc.UserID = u.ID
		next.ServeHTTP(w, r.WithContext(appctx.AddLocToCtx(ctx, loc)))
		go func() {
			_, err := l.service.Update(ctx, loc)
			if err != nil {
				l.log.Warn("Failed to update location", zap.String("location", locInput))
			}
		}()
	}
	return http.HandlerFunc(fn)
}
