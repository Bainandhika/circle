package auth

import (
	"time"

	"circle-2.0/app/config"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken generates a JWT token for the provided user email
func GenerateToken(email string) (string, error) {
	// Create a new token object, specifying the signing method and the claims.
	token := jwt.New(jwt.SigningMethodHS256)

	// Create a map to store our claims.
	timeNow := time.Now()
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["iat"] = timeNow.Unix()
	claims["exp"] = timeNow.Add(24 * time.Hour).Unix()

	// Sign and get the complete encoded token as a string using the secret
	secretKey := []byte(config.App.SecretKey)
	return token.SignedString(secretKey)
}
