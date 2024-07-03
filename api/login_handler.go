package api

import (
	"encoding/json"
	"errors"
	"github.com/charliegreeny/simple-dating-app/internal/app"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/token"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type LoginHandler struct {
	validator  *validator.Validate
	loginCache *token.Login
}

func NewLoginHandler(c *token.Login, validate *validator.Validate) *LoginHandler {
	return &LoginHandler{loginCache: c, validator: validate}
}

func (l LoginHandler) login(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var req *token.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	err := l.validator.Struct(req)
	if err != nil {
		response(w, enc, app.ErrorOutput{Message: "password required"}, http.StatusBadRequest)
		return
	}
	respBody, err := l.loginCache.LoginUser(r.Context(), req)
	if err != nil {
		if errors.As(err, &app.ErrUnauthorized{}) {
			response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusUnauthorized)
			return
		}
		if errors.As(err, &app.ErrNotFound{}) {
			response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusNotFound)
			return
		}
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	response(w, enc, respBody, http.StatusOK)
}
