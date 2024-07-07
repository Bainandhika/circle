package order

import (
	"circle/app/pkg/delivery/handler/template"
	"circle/lib/model"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) GetOrders(c *gin.Context) {
	a, _ := c.Get("api-detail")
	apiDetail, _ := a.(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*gin.Context, model.Response) {
		template.SetResponse(c, response)
	}(c, response)

	orders, status := h.OrderService.GetOrders()
	if status != nil {
		response.Status = *status
		return
	}

	response.Data = gin.H{"orders": orders}
}
