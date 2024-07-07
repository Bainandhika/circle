package order

import (
	"circle/app/pkg/delivery/handler/template"
	customError "circle/lib/helper/custom-error"
	"circle/lib/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	a, _ := c.Get("api-detail")
	apiDetail, _ := a.(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*gin.Context, model.Response) {
		template.SetResponse(c, response)
	}(c, response)

	pathParam := apiDetail.Request.PathParams
	if len(pathParam) == 0 {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: customError.NeedIDasPathParam()}
		return
	}

	id := pathParam[0].Value
	if status := h.OrderService.DeleteOrder(id); status != nil {
		response.Status = *status
		return
	}
}
