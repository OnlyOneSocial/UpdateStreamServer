package model

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

//JWT ...
type JWT struct {
	UserID      int
	TokenString string
}

var hmacSampleSecret []byte

//Validate ...
func (w *JWT) Validate() error {
	token, err := jwt.Parse(w.TokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return hmacSampleSecret, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims["userid"], claims["nbf"])
	} else {
		fmt.Println(err)
	}
	return nil
}

//Create ...
func (w *JWT) Create() error {

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"userid": w.UserID,
		"nbf":    time.Now().Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(hmacSampleSecret)

	fmt.Println(tokenString, err)
	return nil
}
