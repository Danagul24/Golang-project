package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	Admin  Role = "admin"
	Client Role = "client"
)

type User struct {
	ID                int    `json:"id" db:"id"`
	Name              string `json:"name" db:"name"`
	Surname           string `json:"surname" db:"surname"`
	Email             string `json:"email" db:"email"`
	Password          string `json:"password,omitempty" db:"-"`
	EncryptedPassword string `json:"-" db:"password"`
	PhoneNumber       string `json:"phone_number" db:"phone_number"`
	BirthDate         string `json:"birth_date" db:"birth_date"`
	Role              *Role  `json:"role" db:"role"`
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.By(RequiredIf(u.EncryptedPassword == "")),
			validation.Length(6, 50)))
}

func (u *User) BeforeCreating() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}
		u.EncryptedPassword = enc
	}
	u.Sanitize()

	return nil

}

func (u *User) Sanitize() {
	u.Password = ""
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}
