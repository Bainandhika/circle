package order

import (
	"circle/app/pkg/service/order"
	"circle/app/pkg/service/user"
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
