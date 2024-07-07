package user

import (
	"net/http"
	"time"

	"circle-2.0/app/pkg/delivery/handler/template"
	"circle-2.0/lib/model"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	startTime := c.Locals("start-time").(time.Time)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*fiber.Ctx, model.Response) error {
		return template.SetResponse(c, response)
	}(c, response)

	var request model.CreateUserRequest
	if err := c.BodyParser(&request); err != nil {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: err.Error()}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	newUser, status := h.UserService.CreateUser(request, startTime)
	if status != nil {
		response.Status = *status
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	response.Data = newUser
	return nil
}
