package api

import (
	"encoding/json"
	"github.com/charliegreeny/simple-dating-app/internal/model"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
	"net/http"
)

type userHandler struct {
	model.Creator[*user.Input, *user.Output]
}

func (u userHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var req *user.Input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	resp, err := u.Create(req)
	if err != nil {
		response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusInternalServerError)
	}
	response(w, enc, resp, http.StatusCreated)
}

func response(w http.ResponseWriter, enc *json.Encoder, resp any, statusCode int) {
	w.WriteHeader(statusCode)
	_ = enc.Encode(resp)
	return
}
