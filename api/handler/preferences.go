package api

import (
	"encoding/json"
	"errors"
	"github.com/charliegreeny/simple-dating-app/app"
	"net/http"
)

type Preference struct {
	service app.EntityService[*app.Preference, *app.Preference]
}

func NewPreference(service app.EntityService[*app.Preference, *app.Preference]) Preference {
	return Preference{service: service}
}

func (p Preference) UpdatePreferenceHandler(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	var req *app.Preference
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusBadRequest)
		return
	}
	updateP, err := p.service.Update(r.Context(), req)
	if err != nil {
		if errors.As(err, &app.ErrNotFound{}) {
			response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusNotFound)
			return
		}
		response(w, enc, app.ErrorOutput{Message: err.Error()}, http.StatusInternalServerError)
		return
	}
	response(w, enc, updateP, http.StatusOK)
}
