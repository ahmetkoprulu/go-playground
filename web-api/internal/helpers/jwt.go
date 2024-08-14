package helpers

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateJwtToken(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expiration time
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJwtToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}
	return claims, nil
}
