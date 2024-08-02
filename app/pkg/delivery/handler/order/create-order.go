package order

import (
	"net/http"
	"time"

	"circle/app/pkg/delivery/handler/template"
	"circle/lib/model"

	"github.com/gofiber/fiber/v2"
)

func (h *OrderHandler) CreateOrder(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	startTime := c.Locals("start-time").(time.Time)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*fiber.Ctx, model.Response) error {
		return template.SetResponse(c, response)
	}(c, response)

	var request model.OrderRequest
	if err := c.BodyParser(&request); err != nil {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: err.Error()}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	bill, status := h.OrderService.CreateOrder(request, startTime)
	if status != nil {
		response.Status = *status
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	response.Data = fiber.Map{"bill": bill}
	return nil
}
