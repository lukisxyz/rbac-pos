package account

import (
	"encoding/json"
	"errors"
	"net/http"
	"pos/utils/httpresponse"

	"github.com/go-chi/chi/v5"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/oklog/ulid/v2"
)

type accountRoute struct {
	mutate MutationData
	read   ReadData
}

func NewRoute(
	mutate MutationData,
	read ReadData,
) *accountRoute {
	return &accountRoute{
		mutate: mutate,
		read:   read,
	}
}

type publicAccountRoute struct {
	mutate MutationData
	read   ReadData
}

func NewPublicRoute(
	mutate MutationData,
	read ReadData,
) *publicAccountRoute {
	return &publicAccountRoute{
		mutate: mutate,
		read:   read,
	}
}

func (p *publicAccountRoute) Routes() *chi.Mux {
	r := chi.NewMux()
	r.Post("/", p.createAccount)
	return r
}

func (p *accountRoute) Routes() *chi.Mux {
	r := chi.NewMux()
	r.Get("/", p.getAllAccount)
	r.Get("/{id}", p.getOneAccount)
	r.Patch("/password/{id}", p.updatePassword)
	r.Delete("/{id}", p.deleteAccount)
	return r
}

func (p *accountRoute) deleteAccount(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	if err := p.mutate.DeleteAccount(ctx, id); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpresponse.WriteMessage(w, http.StatusOK, "success delete Account")
}

type updatePasswordRequest struct {
	Password string `json:"password" `
}

func (c updatePasswordRequest) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.Password, validation.Required, validation.Length(8, 32)),
	)
}

func (p *accountRoute) updatePassword(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	var body updatePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	data, err := p.mutate.EditAccount(ctx, id, body.Password)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpresponse.WriteData(w, http.StatusOK, data, nil)
}

type createAccountRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c createAccountRequest) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.Email, validation.Required, is.Email),
		validation.Field(&c.Password, validation.Required, validation.Length(8, 32)),
	)
}

func (p *publicAccountRoute) createAccount(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	data, err := p.mutate.CreateAccount(ctx, body.Email, body.Password)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpresponse.WriteData(w, http.StatusCreated, data, nil)
}

func (p *accountRoute) getOneAccount(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	data, err := p.read.GetOneById(ctx, id)
	if err != nil {
		if errors.Is(err, ErrAccountNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err)
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpresponse.WriteData(w, http.StatusOK, data, nil)
}

func (p *accountRoute) getAllAccount(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	data, err := p.read.GetAll(ctx)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	var meta struct {
		Total int `json:"total"`
	}
	meta.Total = data.Count
	httpresponse.WriteData(w, http.StatusOK, data.Accounts, meta)
}

type accountRoleRoute struct {
	mutate         MutationData
	read           ReadData
	accountRoleSvc RoleAccountService
}

func NewRoleRoute(
	mutate MutationData,
	read ReadData,
	accountRoleSvc RoleAccountService,
) *accountRoleRoute {
	return &accountRoleRoute{
		mutate:         mutate,
		read:           read,
		accountRoleSvc: accountRoleSvc,
	}
}

func (p *accountRoleRoute) Routes() *chi.Mux {
	r := chi.NewMux()
	r.Get("/{id}/role", p.getListRole)
	r.Post("/", p.assignRole)
	r.Delete("/", p.deleteRole)
	r.Get("/{id}/account", p.getRole)
	return r
}

func (p *accountRoleRoute) getRole(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	data, err := p.accountRoleSvc.GetAccount(ctx, id)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err)
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	var meta struct {
		Total int `json:"total"`
	}
	meta.Total = data.Count
	httpresponse.WriteData(w, http.StatusOK, data.Accounts, meta)
}

func (p *accountRoleRoute) getListRole(
	w http.ResponseWriter,
	r *http.Request,
) {
	idStr := chi.URLParam(r, "id")
	id, err := ulid.Parse(idStr)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	data, err := p.accountRoleSvc.GetRoleByAccount(ctx, id)
	if err != nil {
		if errors.Is(err, ErrRoleNotFound) {
			httpresponse.WriteError(w, http.StatusNotFound, err)
			return
		}
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	var meta struct {
		Total int `json:"total"`
	}
	meta.Total = data.Count
	httpresponse.WriteData(w, http.StatusOK, data.Roles, meta)
}

type assignRoleRequest struct {
	AccountId ulid.ULID `json:"account_id" validate:"required"`
	RoleId    ulid.ULID `json:"role_id" validate:"required"`
}

func (c assignRoleRequest) Validate() error {
	return validation.ValidateStruct(
		&c,
		validation.Field(&c.AccountId, validation.Required),
		validation.Field(&c.RoleId, validation.Required),
	)
}

func (p *accountRoleRoute) assignRole(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body assignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	err := p.accountRoleSvc.AssignRole(ctx, body.RoleId, body.AccountId)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpresponse.WriteMessage(w, http.StatusCreated, "success assign a role")
}
func (p *accountRoleRoute) deleteRole(
	w http.ResponseWriter,
	r *http.Request,
) {
	var body assignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := body.Validate(); err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	ctx := r.Context()

	err := p.accountRoleSvc.DeleteRole(ctx, body.RoleId, body.AccountId)
	if err != nil {
		httpresponse.WriteError(w, http.StatusBadRequest, err)
		return
	}
	httpresponse.WriteMessage(w, http.StatusOK, "success remove a role")
}
