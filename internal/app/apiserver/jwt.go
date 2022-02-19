package apiserver

import (
	"crypto/rsa"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

const (
	privKeyPath = "/configs/app.rsa"     // rsa 512 private
	pubKeyPath  = "/configs/app.rsa.pub" // rsa 512 public
)

//JwtKeys object
type JwtKeys struct {
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
}

func getJwtKeys(localPath string) JwtKeys {
	path, err := os.Getwd()
	fmt.Println(path)
	if err != nil {
		log.Fatal(err)
	}

	signBytes, err := ioutil.ReadFile(path + localPath + privKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	verifyBytes, err := ioutil.ReadFile(path + localPath + pubKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal(err)
	}
	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal(err)
	}
	return JwtKeys{
		signKey:   signKey,
		verifyKey: verifyKey,
	}
}

//Validate jwt token and get userid
func (w *JwtKeys) Validate(TokenString string) (int, error) {
	token, err := jwt.Parse(TokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return w.verifyKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return int(claims["userid"].(float64)), nil
	}

	return 0, err
}

//Create ...
func (w *JwtKeys) Create(UserID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"userid": UserID,
		"exp":    time.Now().Unix() + (60 * 60 * 6),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(w.signKey)

	fmt.Println(tokenString, err)
	return tokenString, nil
}

//CreateRT создание Refresh токена
func (w *JwtKeys) CreateRT(UserID int) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims{
		"userid": UserID,
		"exp":    time.Now().Unix() + (60 * 60 * 48),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(w.signKey)

	fmt.Println(tokenString, err)
	return tokenString, nil
}
