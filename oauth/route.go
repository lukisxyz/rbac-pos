package oauth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type accountRoute struct {
	svc ServiceOAuth
}

func NewRoute(
	svc ServiceOAuth,
) *accountRoute {
	return &accountRoute{
		svc: svc,
	}
}

func (p *accountRoute) Routes() *chi.Mux {
	r := chi.NewMux()
	r.Post("/login", p.login)
	return r
}

func writeMessage(w http.ResponseWriter, status int, msg string) {
	var j struct {
		Msg string `json:"message"`
	}

	j.Msg = msg

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		http.Error(w, err.Error(), status)
		return
	}
}

func writeData(w http.ResponseWriter, status int, data, meta any) {
	var j struct {
		Data any `json:"data"`
		Meta any `json:"meta"`
	}

	j.Data = data
	j.Meta = meta

	w.Header().Add("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(j); err != nil {
		http.Error(w, err.Error(), status)
		return
	}
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeMessage(w, status, err.Error())
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c loginRequest) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(8, 32)),
	)
}

func (p *accountRoute) login(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body loginRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	data, err := p.svc.Login(ctx, body.Email, body.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeData(w, http.StatusOK, data, nil)
}
