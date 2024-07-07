package middleware

import (
	"circle/lib/helper/tool"
	"circle/lib/logger"
	"circle/lib/model"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func LoggingAPIDetail() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		apiDetail, _ := c.Get("api-detail")
		apiDetailStruct, _ := apiDetail.(model.APIDetail)

		// Handle request body based on content type
		var responseStruct model.Response
		response, ok := c.Get("response")
		if ok {
			apiDetailStruct.Response = response

			responseByte, _ := json.Marshal(response)
			_ = json.Unmarshal(responseByte, &responseStruct)
		}

		// Store details in the context
		status, exists := c.Get("message")
		if exists {
			apiDetailStruct.Message = status.(string)
		}

		startTime, _ := c.Get("start-time")
		startTimeTime, _ := startTime.(time.Time)
		apiDetailStruct.TimeTaken = time.Since(startTimeTime).String()

		if responseStruct.Status.Code != http.StatusOK {
			logger.Error.Println(tool.ToJSON(apiDetailStruct))
		} else {
			logger.Info.Println(tool.ToJSON(apiDetailStruct))
		}
	}
}
