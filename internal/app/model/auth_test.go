package model_test

import (
	"testing"

	"github.com/katelinlis/AuthBackend/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestAuthJWT_Validate(t *testing.T) {
	jwt := model.TestJWT(t)
	assert.NoError(t, jwt.Validate())
	assert.NotEmpty(t, jwt.UserID)
}

func TestAuthJWT_Create(t *testing.T) {
	jwt := model.TestJWT(t)
	err := jwt.Validate()
	assert.NoError(t, err)
}

func TestAuthCreate(t *testing.T) {
	auth := model.TestGenerateEncrypted(t)
	err := auth.Validate()
	assert.NoError(t, err)
	err = auth.BeforeCreate()
	assert.NoError(t, err)

	assert.NotEmpty(t, auth.EncryptedPassword)
}
