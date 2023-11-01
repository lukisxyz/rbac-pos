package permission

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/oklog/ulid/v2"
)

type permissionRoute struct {
	mutate MutationData
	read   ReadData
}

func NewRoute(
	mutate MutationData,
	read ReadData,
) *permissionRoute {
	return &permissionRoute{
		mutate: mutate,
		read:   read,
	}
}

func (p *permissionRoute) Routes() *chi.Mux {
	r := chi.NewMux()
	r.Post("/", p.createPermission)
	r.Get("/", p.getAllPermission)
	r.Get("/{id}", p.getOnePermission)
	r.Patch("/{id}", p.updatePermission)
	r.Delete("/{id}", p.deletePermission)
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

func (p *permissionRoute) deletePermission(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	if err := p.mutate.DeletePermission(ctx, id); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeMessage(w, http.StatusOK, "success delete permission")
}

func (p *permissionRoute) updatePermission(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var body createPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	data, err := p.mutate.EditPermission(ctx, id, body.Name, body.Description, body.Url)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeData(w, http.StatusOK, data, nil)
}

type createPermissionRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Url         string `json:"url" validate:"required"`
}

func (c createPermissionRequest) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.Name, validation.Required),
		validation.Field(&c.Url, validation.Required),
	)
}

func (p *permissionRoute) createPermission(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body createPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	data, err := p.mutate.CreatePermission(ctx, body.Name, body.Description, body.Url)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeData(w, http.StatusCreated, data, nil)
}

func (p *permissionRoute) getOnePermission(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	data, err := p.read.GetOneById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeData(w, http.StatusOK, data, nil)
}

func (p *permissionRoute) getAllPermission(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	data, err := p.read.GetAll(ctx)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var meta struct {
		Total int `json:"total"`
	}
	meta.Total = data.Count
	writeData(w, http.StatusOK, data.Permissions, meta)
}
