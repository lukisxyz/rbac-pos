package role

import (
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
