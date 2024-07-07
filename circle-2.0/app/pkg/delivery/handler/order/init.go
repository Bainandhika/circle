package order

import (
	"circle-fiber/app/pkg/service/order"
	"circle-fiber/app/pkg/service/user"
)

type OrderHandler struct {
	OrderService order.OrderService
	UserService  user.UserService
}

func NewOrderHandler(
	orderService order.OrderService,
	userService user.UserService,
) *OrderHandler {
	return &OrderHandler{
		OrderService: orderService,
		UserService:  userService,
	}
}
