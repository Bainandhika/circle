package order

import (
	"circle-2.0/app/pkg/delivery/handler/template"
	"circle-2.0/lib/model"

	"github.com/gofiber/fiber/v2"
)

func (h *OrderHandler) GetOrders(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*fiber.Ctx, model.Response) error {
		return template.SetResponse(c, response)
	}(c, response)

	orders, status := h.OrderService.GetOrders()
	if status != nil {
		response.Status = *status
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	response.Data = fiber.Map{"orders": orders}
	return nil
}
