package api

import (
	"encoding/json"
	"errors"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/appctx"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/discovery/filter"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/match"
	"net/http"
)

type Discovery struct {
	cache   app.Cache[string, *app.User]
	filter  filter.Filter
	matcher match.Matcher
}

func NewDiscoveryHandler(cache app.Cache[string, *app.User], matcher match.Matcher, f filter.Filter) Discovery {
	return Discovery{cache: cache, matcher: matcher, filter: f}
}

func (d Discovery) GetEligibleUsersHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	users := d.cache.GetAll(r.Context())
	currentUser := appctx.GetUserFromCtx(r.Context())
	eligibleUsers := d.filter.Apply(r.Context(), currentUser.Pref, users)
	if len(eligibleUsers) == 0 {
		response(
			w,
			enc,
			&app.ErrorOutput{Message: "no users found in your area, make sure your location is updated using `x-location-header`"},
			http.StatusOK)
	}
	respBody := struct {
		Users []*app.User `json:"results"`
	}{eligibleUsers}
	w.WriteHeader(http.StatusOK)
	_ = enc.Encode(respBody)
	return
}

func (d Discovery) PostSwipeHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var req *match.SwipeInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	resp, err := d.matcher.Match(r.Context(), req)
	if err != nil {
		if errors.As(err, &app.ErrBadRequest{}) {
			response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
			return
		}
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	if resp == nil {
		response(w, enc, nil, http.StatusNoContent)
		return
	}
	response(w, enc, resp, http.StatusOK)
	return
}
