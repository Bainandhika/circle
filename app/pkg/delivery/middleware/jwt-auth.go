package middleware

import (
	"net/http"
	"strings"

	"circle/app/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware is a sample middleware for authentication
func JWTAuth(c *fiber.Ctx) error {
	// Extract the token from the Authorization header
	tokenHeader := c.Get("Authorization")
	if tokenHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is required"})
	}

	parts := strings.Split(tokenHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
	}

	tokenString := parts[1]

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(config.App.SecretKey), nil
	})
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Check if the token is valid
	if !token.Valid {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Extract claims from the token
	_, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
	}

	// You can now access the claims, including "email" and "iat" (issued at time)
	// For example:
	// userEmail := claims["email"].(string)
	// issuedAt := int64(claims["iat"].(float64))

	// Proceed to the next middleware or handler
	return c.Next()
}
