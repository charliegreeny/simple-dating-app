package api

import (
	"encoding/json"
	"errors"
	"github.com/charliegreeny/simple-dating-app/internal/model"
	"github.com/charliegreeny/simple-dating-app/internal/pkg/token"
	"net/http"
)

type LoginHandler struct {
	loginCache *token.Cache
}

func NewLoginHandler(c *token.Cache) *LoginHandler {
	return &LoginHandler{loginCache: c}
}

func (l LoginHandler) login(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var req *token.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	respBody, err := l.loginCache.LoginUser(req)
	if err != nil {
		if errors.As(err, &model.ErrUnauthorized{}) {
			response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusUnauthorized)
			return
		}
		if errors.As(err, &model.ErrNotFound{}) {
			response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusNotFound)
			return
		}
		response(w, enc, model.ErrorOutput{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	response(w, enc, respBody, http.StatusOK)
}
