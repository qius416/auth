package authentication

import (
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

// MakeToken is the method to generate jwt token
func MakeToken(email string, role string) (tokenString string, err error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	token.Claims["email"] = email
	token.Claims["role"] = role
	token.Claims["exp"] = time.Now().Add(time.Second * 30).Unix()
	// Sign and get the complete encoded token as a string
	// Secret is from environment variable which need be set in dockerfile's env
	tokenString, err = token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	return
}
