package user

import (
	"circle/app/pkg/delivery/handler/template"
	"circle/lib/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) CreateUser(c *gin.Context) {
	a, _ := c.Get("api-detail")
	apiDetail, _ := a.(model.APIDetail)
	t, _ := c.Get("start-time")
	startTime, _ := t.(time.Time)
	response := model.Response{TransactionID: apiDetail.Headers.TransactionID, ChannelID: apiDetail.Headers.ChannelID}

	defer func(*gin.Context, model.Response) {
		template.SetResponse(c, response)
	}(c, response)

	var request model.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Status = model.Status{Code: http.StatusBadRequest, Message: err.Error()}
		return
	}

	newUser, status := h.UserService.CreateUser(request, startTime)
	if status != nil {
		response.Status = *status
		return
	}

	response.Data = newUser
}
