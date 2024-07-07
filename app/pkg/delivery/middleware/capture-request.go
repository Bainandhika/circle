package middleware

import (
	"bytes"
	"circle/lib/model"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func CaptureRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var apiDetail model.APIDetail

		apiDetail.Hostname, _ = os.Hostname()
		apiDetail.URL = c.Request.URL.Path
		apiDetail.Method = c.Request.Method

		apiDetail.Headers = model.Headers{
			TransactionID: c.GetHeader("transaction-id"),
			APIKey:        c.GetHeader("api-key"),
			ChannelID:     c.GetHeader("channel-id"),
		}
		c.Set("headers", apiDetail.Headers)

		// Capture path parameters
		pathParams := c.Params
		if len(pathParams) > 0 {
			for _, param := range pathParams {
				apiDetail.Request.PathParams = append(apiDetail.Request.PathParams, model.PathParam{
					Key:   param.Key,
					Value: param.Value,
				})
			}
		}

		// Capture query parameters
		queryParams := c.Request.URL.Query()
		if len(queryParams) > 0 {
			for key, values := range queryParams {
				apiDetail.Request.QueryParams = append(apiDetail.Request.QueryParams, model.QueryParam{
					Key:   key,
					Value: values,
				})
			}
		}

		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			c.Abort()
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body so it can be read again

		if len(bodyBytes) > 0 {
			// requestBody := tool.RequestBodyFactory[apiDetail.URL]
			// if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
			// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal request body"})
			// } else {
			// 	createUserRequest, okCreateUserRequest := requestBody.(*model.CreateUserRequest)
			// 	if okCreateUserRequest {
			// 		createUserRequest.Password = ""
			// 		apiDetail.Request.Body = requestBody
			// 	}

			// 	loginUserRequest, okLoginUserRequest := requestBody.(*model.LoginUserRequest)
			// 	if okLoginUserRequest {
			// 		loginUserRequest.Password = ""
			// 		apiDetail.Request.Body = requestBody
			// 	}

			// 	apiDetail.Request.Body = requestBody
			// }

			var requestBody any
			if err := json.Unmarshal(bodyBytes, &requestBody); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unmarshal request body"})
				c.Abort()
				return
			}

		}

		c.Set("start-time", startTime)
		c.Set("api-detail", apiDetail)

		c.Next()
	}
}
