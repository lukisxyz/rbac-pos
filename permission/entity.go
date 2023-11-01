package permission

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Permission struct {
	Id          ulid.ULID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Url         string    `json:"url"`
	CreatedAt   time.Time `json:"created_at"`
}

func newPermission(name, desc, url string) Permission {
	id := ulid.Make()
	return Permission{
		Id:          id,
		Name:        name,
		Description: desc,
		Url:         url,
		CreatedAt:   time.Now(),
	}
}
