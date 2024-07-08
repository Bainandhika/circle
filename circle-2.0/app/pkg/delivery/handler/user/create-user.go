package user

import (
	"net/http"
	"sync"
	"time"

	"circle-2.0/app/pkg/delivery/handler/template"
	"circle-2.0/lib/model"

	"github.com/gofiber/fiber/v2"
)

var wg sync.WaitGroup

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	startTime := c.Locals("start-time").(time.Time)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	var request model.CreateUserRequest
	if err := c.BodyParser(&request); err != nil {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: err.Error()}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		newUser, status := h.UserService.CreateUser(request, startTime)
		if status != nil {
			response.Status = *status
		}

		response.Data = newUser
	}()
	wg.Wait()

	return template.SetResponse(c, response)
}
