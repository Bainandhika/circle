package user

import (
	"net/http"

	"circle/app/pkg/delivery/handler/template"
	customError "circle/lib/helper/custom-error"
	"circle/lib/model"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*fiber.Ctx, model.Response) error {
		return template.SetResponse(c, response)
	}(c, response)

	var updateUser model.UpdateUserRequest
	if err := c.BodyParser(&updateUser); err != nil {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: err.Error()}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	pathParam := apiDetail.Request.PathParams
	if len(pathParam) == 0 {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: customError.NeedIDasPathParam()}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	id := pathParam[0].Value
	if status := h.UserService.UpdateUser(id, updateUser); status != nil {
		response.Status = *status
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	return nil
}
