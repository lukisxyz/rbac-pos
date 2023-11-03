package domain

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
)

type Role struct {
	Id          ulid.ULID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type RolePermission struct {
	PermissionId ulid.ULID `json:"permission_id"`
	RoleId       ulid.ULID `json:"role_id"`
}

func NewRole(name, desc string) Role {
	id := ulid.Make()
	return Role{
		Id:          id,
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now(),
	}
}

type ReadRoleResponse struct {
	Role
	TotalPermissions int      `json:"total_permission"`
	Permissions      []string `json:"permissions"`
}

func (a *ReadRoleResponse) MarshalJSON() ([]byte, error) {
	var j struct {
		Id              ulid.ULID `json:"id"`
		Name            string    `json:"name"`
		Desc            string    `json:"description"`
		TotalPermission int       `json:"total_permission"`
		Permissions     []string  `json:"permissions"`
	}

	j.Id = a.Id
	j.Name = a.Name
	j.Desc = a.Description
	j.TotalPermission = a.TotalPermissions
	j.Permissions = a.Permissions

	return json.Marshal(j)
}
