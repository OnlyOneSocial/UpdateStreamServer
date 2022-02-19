package apiserver

import (
	"testing"

	"github.com/katelinlis/AuthBackend/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestAuthJWT_Validate(t *testing.T) {
	//jwt := model.TestJWT(t)

	jwtKeys := getJwtKeys("/../../..")

	jwtjson, err := jwtKeys.Create(100)
	assert.NoError(t, err)

	userid, err := jwtKeys.Validate(jwtjson)
	assert.NotEmpty(t, userid)
	assert.NoError(t, err)
}

func TestAuthJWT_Create(t *testing.T) {
	jwt := model.TestJWT(t)
	err := jwt.Validate()
	assert.NoError(t, err)
}
