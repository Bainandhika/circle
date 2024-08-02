package order

import (
	"net/http"

	"circle/app/pkg/delivery/handler/template"
	customError "circle/lib/helper/custom-error"
	"circle/lib/model"

	"github.com/gofiber/fiber/v2"
)

func (h *OrderHandler) GetBillUser(c *fiber.Ctx) error {
	apiDetail := c.Locals("api-detail").(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*fiber.Ctx, model.Response) error {
		return template.SetResponse(c, response)
	}(c, response)

	request := model.GetBillUserRequest{
		UserID:  c.Query("user_id"),
		OrderID: c.Query("order-id"),
	}

	if request.UserID == "" || request.OrderID == "" {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: customError.NeedQueryParams("user-id", "order-id")}
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	billUser, status := h.OrderService.GetBillUser(request)
	if status != nil {
		response.Status = *status
		return fiber.NewError(response.Status.Code, response.Status.Message)
	}

	response.Data = fiber.Map{"bill-user": billUser}
	return nil
}
