package order

import (
	"circle-2.0/app/pkg/service/order"
	"circle-2.0/app/pkg/service/user"
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
