package api

import (
	"encoding/json"
	"errors"
	"github.com/charliegreeny/simple-dating-app/internal/model"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type UserHandler struct {
	validator *validator.Validate
	model.GetterCreator[*user.Input, *user.Output]
}

func NewUserHandler(validator *validator.Validate, creator model.GetterCreator[*user.Input, *user.Output]) *UserHandler {
	return &UserHandler{validator: validator, GetterCreator: creator}
}

func (u UserHandler) getUserHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	id := chi.URLParam(r, "id")
	resp, err := u.Get(id)
	if err != nil {
		if errors.As(err, &model.ErrNotFound{}) {
			response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusNotFound)
			return
		}
		response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	response(w, enc, resp, http.StatusOK)
	return
}

func (u UserHandler) createUserHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var req *user.Input
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	err := u.validator.Struct(req)
	if err != nil {
		response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	resp, err := u.Create(req)
	if err != nil {
		if errors.As(err, &model.ErrBadRequest{}) {
			response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
			return
		}
		response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	response(w, enc, resp, http.StatusCreated)
}

func response(w http.ResponseWriter, enc *json.Encoder, resp any, statusCode int) {
	w.WriteHeader(statusCode)
	_ = enc.Encode(resp)
	return
}
