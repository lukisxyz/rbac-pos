package account

import (
	"encoding/json"
	"time"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	Id        ulid.ULID `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

type AccountRole struct {
	RoleId    ulid.ULID `json:"role_id"`
	AccountId ulid.ULID `json:"account_id"`
}

func newAccount(email, pwd string) (Account, error) {
	id := ulid.Make()
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	return Account{
		Id:        id,
		Email:     email,
		Password:  string(bytes),
		CreatedAt: time.Now(),
	}, err
}

func (a *Account) MarshalJSON() ([]byte, error) {
	var j struct {
		Id    ulid.ULID `json:"id"`
		Email string    `json:"email"`
	}

	j.Id = a.Id
	j.Email = a.Email

	return json.Marshal(j)
}
