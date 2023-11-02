package role

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/oklog/ulid/v2"
)

type roleRoute struct {
	mutate         MutationData
	read           ReadData
	rolePermission RolePermissionService
}

func NewRoute(
	mutate MutationData,
	read ReadData,
	rolePermission RolePermissionService,
) *roleRoute {
	return &roleRoute{
		mutate:         mutate,
		read:           read,
		rolePermission: rolePermission,
	}
}

func (p *roleRoute) Routes() *chi.Mux {
	r := chi.NewMux()
	r.Post("/", p.createRole)
	r.Get("/", p.getAllRole)
	r.Get("/{id}", p.getOneRole)
	r.Patch("/{id}", p.updateRole)
	r.Delete("/{id}", p.deleteRole)
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

func (p *roleRoute) deleteRole(
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

	if err := p.mutate.DeleteRole(ctx, id); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeMessage(w, http.StatusOK, "success delete Role")
}

func (p *roleRoute) updateRole(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var body createRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	data, err := p.mutate.EditRole(ctx, id, body.Name, body.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeData(w, http.StatusOK, data, nil)
}

type createRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

func (c createRoleRequest) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.Name, validation.Required),
	)
}

func (p *roleRoute) createRole(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body createRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if err := body.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	data, err := p.mutate.CreateRole(ctx, body.Name, body.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeData(w, http.StatusCreated, data, nil)
}

func (p *roleRoute) getOneRole(
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
		if errors.Is(err, ErrRoleNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeData(w, http.StatusOK, data, nil)
}

func (p *roleRoute) getAllRole(
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
	writeData(w, http.StatusOK, data.Roles, meta)
}

type permissionRoleRoute struct {
	mutate         MutationData
	read           ReadData
	rolePermission RolePermissionService
}

func NewPermissionRoute(
	mutate MutationData,
	read ReadData,
	rolePermission RolePermissionService,
) *permissionRoleRoute {
	return &permissionRoleRoute{
		mutate:         mutate,
		read:           read,
		rolePermission: rolePermission,
	}
}

func (p *permissionRoleRoute) Routes() *chi.Mux {
	r := chi.NewMux()
	r.Get("/{id}/role", p.getListRole)
	r.Post("/", p.assignPermission)
	r.Delete("/", p.deletePermission)
	r.Get("/{id}/permission", p.getPermission)
	return r
}

func (p *permissionRoleRoute) getPermission(
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

	data, err := p.rolePermission.GetPermission(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var meta struct {
		Total int `json:"total"`
	}
	meta.Total = data.Count
	writeData(w, http.StatusOK, data.Permissions, meta)
}

func (p *permissionRoleRoute) getListRole(
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

	data, err := p.rolePermission.GetRoleByPermission(ctx, id)
	if err != nil {
		if errors.Is(err, ErrPermissionNotFound) {
			writeError(w, http.StatusNotFound, err)
			return
		}
		writeError(w, http.StatusBadRequest, err)
		return
	}
	var meta struct {
		Total int `json:"total"`
	}
	meta.Total = data.Count
	writeData(w, http.StatusOK, data.Roles, meta)
}

type assignPermissionRequest struct {
	PermissionId ulid.ULID `json:"permission_id" validate:"required"`
	RoleId       ulid.ULID `json:"role_id" validate:"required"`
}

func (c assignPermissionRequest) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.PermissionId, validation.Required),
		validation.Field(&c.RoleId, validation.Required),
	)
}

func (p *permissionRoleRoute) assignPermission(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body assignPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	err := p.rolePermission.AssignPermisson(ctx, body.RoleId, body.PermissionId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeMessage(w, http.StatusCreated, "success assign a permission")
}

func (p *permissionRoleRoute) deletePermission(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body assignPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	err := p.rolePermission.DeletePermission(ctx, body.RoleId, body.PermissionId)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeMessage(w, http.StatusCreated, "success remove a permission")
}
