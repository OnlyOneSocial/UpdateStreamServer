package model

import (
	"testing"
)

//TestJWT ...
func TestJWT(t *testing.T) *JWT {

	return &JWT{
		TokenString: "dfdfsfdfds",
		UserID:      1,
	}
}

//TestGenerateEncrypted ...
func TestGenerateEncrypted(t *testing.T) *Auth {

	return &Auth{
		Login:    "katelinlis",
		Password: "MyPassword",
	}
}
