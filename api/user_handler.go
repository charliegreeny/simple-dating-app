package api

import (
	"encoding/json"
	"errors"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type UserHandler struct {
	validator *validator.Validate
	app.GetterCreator[*user.Input, *user.Output]
}

func NewUserHandler(validator *validator.Validate, creator app.GetterCreator[*user.Input, *user.Output]) *UserHandler {
	return &UserHandler{validator: validator, GetterCreator: creator}
}

func (u UserHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	id := chi.URLParam(r, "id")
	resp, err := u.Get(r.Context(), id)
	if err != nil {
		if errors.As(err, &app.ErrNotFound{}) {
			response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusNotFound)
			return
		}
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	response(w, enc, resp, http.StatusOK)
	return
}

func (u UserHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var req *user.Input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	err := u.validator.Struct(req)
	if err != nil {
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	resp, err := u.Create(r.Context(), req)
	if err != nil {
		if errors.As(err, &app.ErrBadRequest{}) {
			response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
			return
		}
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	response(w, enc, resp, http.StatusCreated)
}

func response(w http.ResponseWriter, enc *json.Encoder, resp any, statusCode int) {
	w.WriteHeader(statusCode)
	_ = enc.Encode(resp)
	return
}
