package domain

import (
	"encoding/json"
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

func NewPermission(name, desc, url string) Permission {
	id := ulid.Make()
	return Permission{
		Id:          id,
		Name:        name,
		Description: desc,
		Url:         url,
		CreatedAt:   time.Now(),
	}
}

func (a *Permission) MarshalJSON() ([]byte, error) {
	var j struct {
		Id   ulid.ULID `json:"id"`
		Name string    `json:"name"`
		Desc string    `json:"description"`
		Url  string    `json:"url"`
	}

	j.Id = a.Id
	j.Name = a.Name
	j.Desc = a.Description
	j.Url = a.Url

	return json.Marshal(j)
}
