package order

import (
	"circle-fiber/app/pkg/delivery/handler/template"
	customError "circle-fiber/lib/helper/custom-error"
	"circle-fiber/lib/model"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *OrderHandler) UpdateOrder(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	startTime := c.Locals("start-time").(time.Time)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*fiber.Ctx, model.Response) error {
		return template.SetResponse(c, response)
	}(c, response)

	pathParam := apiDetail.Request.PathParams
	if len(pathParam) == 0 {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: customError.NeedIDasPathParam()}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	orderMainID := pathParam[0].Value

	var request model.UpdateOrderRequest
	if err := c.BodyParser(&request); err != nil {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: err.Error()}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	if status := h.OrderService.UpdateOrder(orderMainID, request, startTime); status != nil {
		response.Status = *status
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	return nil
}
