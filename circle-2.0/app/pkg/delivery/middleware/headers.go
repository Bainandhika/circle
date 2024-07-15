package middleware

import (
	"fmt"
	"net/http"
	"time"

	"circle-2.0/app/config"
	"circle-2.0/lib/helper/tool"
	"circle-2.0/lib/logger"
	"circle-2.0/lib/model"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func Headers(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)

	var status *model.Status
	defer func(*model.Status, model.APIDetail) error {
		if status != nil {
			apiDetail.Response = model.Response{
				TransactionID: apiDetail.Headers.TransactionID,
				ChannelID:     apiDetail.Headers.ChannelID,
				Status:        *status,
			}

			apiDetail.Message = status.Message
			apiDetail.TimeTaken = time.Since(c.Locals("start-time").(time.Time)).String()

			logger.Error.Println(tool.ToJSONString(apiDetail))
			return c.Status(status.Code).JSON(apiDetail.Response)
		}

		c.Locals("api-detail", apiDetail)
		return c.Next()
	}(status, apiDetail)

	if apiDetail.Headers.TransactionID == "" {
		apiDetail.Headers.TransactionID = fmt.Sprintf("P-%d-%s", time.Now().Unix(), uuid.New().String())
	}

	// Perform validation (example logic)
	if apiDetail.Headers.SecretKey == "" {
		status = &model.Status{Code: http.StatusUnauthorized, Message: "secret-key is needed on header"}
		return fiber.NewError(status.Code, status.Message)
	} else if apiDetail.Headers.SecretKey != config.App.SecretKey {
		status = &model.Status{Code: http.StatusUnauthorized, Message: "invalid secret-key on header"}
		return fiber.NewError(status.Code, status.Message)
	}

	if apiDetail.Headers.ChannelID == "" {
		status = &model.Status{Code: http.StatusUnauthorized, Message: "channel-id is needed header"}
		return fiber.NewError(status.Code, status.Message)
	}

	return nil
}
