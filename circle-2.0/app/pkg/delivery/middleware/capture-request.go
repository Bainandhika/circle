package middleware

import (
	"os"
	"time"

	"circle-2.0/lib/model"

	"github.com/gofiber/fiber/v2"
)

func CaptureRequest(c *fiber.Ctx) error {
	startTime := time.Now()

	var apiDetail model.APIDetail

	apiDetail.Hostname, _ = os.Hostname()
	apiDetail.URL = string(c.Path())
	apiDetail.Method = c.Method()

	apiDetail.Headers = model.Headers{
		TransactionID: c.Get("transaction-id"),
		SecretKey:     c.Get("secret-key"),
		ChannelID:     c.Get("channel-id"),
	}

	// Capture query parameters
	queryParams := c.Queries()
	if len(queryParams) > 0 {
		for key, value := range queryParams {
			apiDetail.Request.QueryParams = append(apiDetail.Request.QueryParams, model.QueryParam{
				Key:   key,
				Value: value,
			})
		}
	}

	var body fiber.Map
	c.App().Config().JSONDecoder(c.Body(), &body)
	if body != nil {
		apiDetail.Request.Body = body
	}

	c.Locals("start-time", startTime)
	c.Locals("api-detail", apiDetail)

	return c.Next()
}
