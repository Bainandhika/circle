package user

import (
	"circle-fiber/app/pkg/delivery/handler/template"
	"circle-fiber/lib/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) Login(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*fiber.Ctx, model.Response) error {
		return template.SetResponse(c, response)
	}(c, response)

	var loginData model.LoginUserRequest
	if err := c.BodyParser(&loginData); err != nil {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: err.Error()}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	resp, status := h.UserService.LoginUser(loginData)
	if status != nil {
		response.Status = *status
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	response.Data = resp
	return nil
}
