package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

//Auth obj
type Auth struct {
	ID                int    `json:"id"`
	Login             string `json:"login"`
	EncryptedPassword string `json:"-"`
	Password          string `json:"password,omitempty"`
}

//BeforeCreate криптография пароля
func (auth *Auth) BeforeCreate() error {
	if len(auth.Password) > 0 {
		enc, err := encryptString(auth.Password)
		if err != nil {
			return err
		}

		auth.EncryptedPassword = enc
	}
	return nil
}

//Sanitize удаление пароля
func (auth *Auth) Sanitize() {
	auth.Password = ""
}

//Validate ...
func (auth *Auth) Validate() error {
	return validation.ValidateStruct(
		auth,
		validation.Field(&auth.Login, validation.Required, validation.Length(1, 200)),
		validation.Field(&auth.Password, validation.Required),
	)
}

//ComparePassword ...
func (auth *Auth) ComparePassword() bool {
	return bcrypt.CompareHashAndPassword([]byte(auth.EncryptedPassword), []byte(auth.Password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(b), nil
}
