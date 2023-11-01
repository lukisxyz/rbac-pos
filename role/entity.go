package role

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

func newRole(name, desc string) Role {
	id := ulid.Make()
	return Role{
		Id:          id,
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now(),
	}
}

func (a *Role) MarshalJSON() ([]byte, error) {
	var j struct {
		Id   ulid.ULID `json:"id"`
		Name string    `json:"name"`
		Desc string    `json:"description"`
	}

	j.Id = a.Id
	j.Name = a.Name
	j.Desc = a.Description

	return json.Marshal(j)
}
