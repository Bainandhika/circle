package template

import (
	customError "circle/lib/helper/custom-error"
	"circle/lib/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetResponse(c *gin.Context, response model.Response) {
	switch response.Status.Code {
	case 0:
		response.Status.Code = http.StatusOK
		response.Status.Message = customError.Success()
	case http.StatusInternalServerError:
		c.Set("message", response.Status.Message)
		response.Status.Message = customError.InternalServerError()
	default:
		c.Set("message", response.Status.Message)
	}

	c.Set("response", response)
	c.JSON(response.Status.Code, response)
}
