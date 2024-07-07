package middleware

import (
	"circle-2.0/lib/model"

	"github.com/gofiber/fiber/v2"
)

func CapturePathParams(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)

	// Capture path parameters
	pathParams := c.AllParams()
	if len(pathParams) > 0 {
		for key, value := range pathParams {
			apiDetail.Request.PathParams = append(apiDetail.Request.PathParams, model.PathParam{
				Key:   key,
				Value: value,
			})
		}
	}

	c.Locals("api-detail", apiDetail)

	return c.Next()
}
