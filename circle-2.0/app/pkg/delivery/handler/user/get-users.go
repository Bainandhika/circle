package user

import (
	"circle-fiber/app/pkg/delivery/handler/template"
	"circle-fiber/lib/model"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*fiber.Ctx, model.Response) error {
		return template.SetResponse(c, response)
	}(c, response)

	users, status := h.UserService.GetUsers()
	if status != nil {
		response.Status = *status
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	response.Data = fiber.Map{"users": users}
	return nil
}
