package middleware

import (
	"circle/lib/helper/tool"
	"circle/lib/logger"
	"circle/lib/model"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Headers() gin.HandlerFunc {
	return func(c *gin.Context) {
		apiDetail, _ := c.Get("api-detail")
		apiDetailStruct, _ := apiDetail.(model.APIDetail)

		var status *model.Status
		defer func(*model.Status, model.APIDetail) {
			if status != nil {
				apiDetailStruct.Response = model.Response{
					TransactionID: apiDetailStruct.Headers.TransactionID,
					ChannelID:     apiDetailStruct.Headers.ChannelID,
					Status:        *status,
				}

				apiDetailStruct.Message = status.Message

				startTime, _ := c.Get("start-time")
				startTimeTime, _ := startTime.(time.Time)
				apiDetailStruct.TimeTaken = time.Since(startTimeTime).String()

				logger.Error.Println(tool.ToJSON(apiDetailStruct))
				c.AbortWithStatusJSON(status.Code, apiDetailStruct.Response)
			}

			c.Set("api-detail", apiDetailStruct)
			c.Set("headers", apiDetailStruct.Headers)
			c.Next()
		}(status, apiDetailStruct)

		if apiDetailStruct.Headers.TransactionID == "" {
			apiDetailStruct.Headers.TransactionID = fmt.Sprintf("P-%d-%s", time.Now().Unix(), uuid.New().String())
		}

		// Perform validation (example logic)
		if apiDetailStruct.Headers.APIKey == "" {
			status = &model.Status{Code: http.StatusUnauthorized, Message: "secret-key is needed on header"}
			return
		} else if apiDetailStruct.Headers.APIKey != "AshJ/v@!41" {
			status = &model.Status{Code: http.StatusUnauthorized, Message: "invalid secret-key on header"}
			return
		}

		if apiDetailStruct.Headers.ChannelID == "" {
			status = &model.Status{Code: http.StatusUnauthorized, Message: "channel-id is needed header"}
			return
		}
	}
}
