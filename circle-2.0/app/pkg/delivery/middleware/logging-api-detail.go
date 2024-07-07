package middleware

import (
	"circle-fiber/lib/helper/tool"
	"circle-fiber/lib/logger"
	"circle-fiber/lib/model"
	"net/http"
	"time"

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
