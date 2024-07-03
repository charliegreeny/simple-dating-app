package api

import (
	"encoding/json"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
	"net/http"
)

type DiscoveryHandler struct {
	app.Cache[string, *user.Output]
}

func NewDiscoveryHandler(cache app.Cache[string, *user.Output]) *DiscoveryHandler {
	return &DiscoveryHandler{Cache: cache}
}

func (d DiscoveryHandler) GetEligibleUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := d.GetAll(r.Context())
	response(w, json.NewEncoder(w), users, http.StatusOK)
	return
}

//func (d DiscoveryHandler) PostSwipeHandler(writer http.ResponseWriter, request *http.Request) {
//
//}
