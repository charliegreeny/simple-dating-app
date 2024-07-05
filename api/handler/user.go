package api

import (
	"encoding/json"
	"errors"
	"github.com/charliegreeny/simple-dating-app/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/user/service"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type User struct {
	validator *validator.Validate
	app.GetterCreator[*service.Input, *app.UserOutput]
}

func NewUserHandler(validator *validator.Validate, creator app.GetterCreator[*service.Input, *app.UserOutput]) User {
	return User{validator: validator, GetterCreator: creator}
}

func (u User) GetUserHandler(w http.ResponseWriter, r *http.Request) {
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

func (u User) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var req *service.Input
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
