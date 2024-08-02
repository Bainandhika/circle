package template

import (
	"net/http"

	customError "circle/lib/helper/custom-error"
	"circle/lib/model"

	"github.com/gofiber/fiber/v2"
)

func SetResponse(c *fiber.Ctx, response model.Response) error {
	switch response.Status.Code {
	case 0:
		response.Status.Code = http.StatusOK
		response.Status.Message = customError.Success()
	case http.StatusInternalServerError:
		c.Locals("message", response.Status.Message)
		response.Status.Message = customError.InternalServerError()
	default:
		c.Locals("message", response.Status.Message)
	}

	c.Locals("response", response)
	return c.Status(response.Status.Code).JSON(response)
}
