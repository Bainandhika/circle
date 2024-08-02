package middleware

import "github.com/gofiber/fiber/v2"

func RecoveryMiddleware(c *fiber.Ctx) error {
	if err := recover(); err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "message": "Internal Server Error",
            "error":  err,
        })
	}

	return c.Next()
}