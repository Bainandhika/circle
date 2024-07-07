package order

import (
	"circle/app/pkg/delivery/handler/template"
	"circle/lib/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) GetBillUser(c *gin.Context) {
	a, _ := c.Get("api-detail")
	apiDetail, _ := a.(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*gin.Context, model.Response) {
		template.SetResponse(c, response)
	}(c, response)

	var request model.GetBillUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: err.Error()}
		return
	}

	billUser, status := h.OrderService.GetBillUser(request)
	if status != nil {
		response.Status = *status
		return
	}

	response.Data = gin.H{"bill-user": billUser}
}
