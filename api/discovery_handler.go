package api

import (
	"encoding/json"
	"errors"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/match"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
	"net/http"
)

type DiscoveryHandler struct {
	app.Cache[string, *user.Output]
	match.Matcher
}

func NewDiscoveryHandler(cache app.Cache[string, *user.Output], matcher match.Matcher) *DiscoveryHandler {
	return &DiscoveryHandler{Cache: cache, Matcher: matcher}
}

func (d DiscoveryHandler) GetEligibleUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := d.GetAll(r.Context())
	response(w, json.NewEncoder(w), users, http.StatusOK)
	return
}

func (d DiscoveryHandler) PostSwipeHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var req *match.SwipeInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	resp, err := d.Match(r.Context(), req)
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
