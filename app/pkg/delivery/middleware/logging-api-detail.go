package middleware

import (
	"net/http"
	"time"

	"circle/lib/helper/tool"
	"circle/lib/logger"
	"circle/lib/model"

	"github.com/gofiber/fiber/v2"
)

func LoggingAPIDetail(c *fiber.Ctx) error {
	c.Next()

	apiDetail := c.Locals("api-detail").(model.APIDetail)

	if response := c.Locals("response"); response != nil {
		apiDetail.Response = response.(model.Response)
	} else {
		apiDetail.Response = response
	}

	if message := c.Locals("message"); message != nil {
		apiDetail.Message = message.(string)
	} else {
		apiDetail.Message = message
	}

	startTime := c.Locals("start-time").(time.Time)
	apiDetail.TimeTaken = time.Since(startTime).String()

	apiDetailJson := tool.ToJSONString(apiDetail)
	if c.Response().StatusCode() != http.StatusOK {
		logger.Error.Println(apiDetailJson)
	} else {
		logger.Info.Println(apiDetailJson)
	}

	return nil
}
