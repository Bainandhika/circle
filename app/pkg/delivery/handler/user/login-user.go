package user

import (
	"circle/app/pkg/delivery/handler/template"
	"circle/lib/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) Login(c *gin.Context) {
	a, _ := c.Get("api-detail")
	apiDetail, _ := a.(model.APIDetail)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*gin.Context, model.Response) {
		template.SetResponse(c, response)
	}(c, response)

	var loginData model.LoginUserRequest
	if err := c.ShouldBindJSON(&loginData); err != nil {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: err.Error()}
		return
	}

	resp, status := h.UserService.LoginUser(loginData)
	if status != nil {
		response.Status = *status
		return
	}

	response.Data = resp
}
