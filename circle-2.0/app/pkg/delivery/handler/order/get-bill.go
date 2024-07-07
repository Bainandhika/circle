package order

import (
	"circle-fiber/app/pkg/delivery/handler/template"
	customError "circle-fiber/lib/helper/custom-error"
	"circle-fiber/lib/model"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *OrderHandler) GetBill(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*fiber.Ctx, model.Response) error {
		return template.SetResponse(c, response)
	}(c, response)

	pathParam := apiDetail.Request.PathParams
	if len(pathParam) == 0 {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: customError.NeedIDasPathParam()}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	id := pathParam[0].Value
	bill, status := h.OrderService.GetBill(id)
	if status != nil {
		response.Status = *status
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	response.Data = fiber.Map{"bill": bill}
	return nil
}
