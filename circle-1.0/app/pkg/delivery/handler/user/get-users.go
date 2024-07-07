package user

import (
	"circle/app/pkg/delivery/handler/template"
	"circle/lib/model"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) GetUsers(c *gin.Context) {
	a, _ := c.Get("api-detail")
	apiDetail, _ := a.(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*gin.Context, model.Response) {
		template.SetResponse(c, response)
	}(c, response)

	users, status := h.UserService.GetUsers()
	if status != nil {
		response.Status = *status
		return
	}

	response.Data = gin.H{"users": users}
}
