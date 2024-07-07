package order

import (
	"time"

	"circle-2.0/app/pkg/repository/order"
	"circle-2.0/app/pkg/repository/user"
	"circle-2.0/lib/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type orderService struct {
	DB                 *gorm.DB
	Redis              *redis.Client
	UserRepo           user.UserRepo
	OrderMainRepo      order.OrderMainRepo
	OrderUserRepo      order.OrderUserRepo
	OrderUserItemRepo  order.OrderUserItemRepo
	AdditionalCostRepo order.AdditionalCostRepo
	DiscountRepo       order.DiscountRepo
	DeleteOrderRepo    order.DeleteOrderRepo
}

type OrderService interface {
	CreateOrder(req model.OrderRequest, startTime time.Time) (resp model.Bill, status *model.Status)
	GetOrders() (orders []model.OrderMains, status *model.Status)
	GetBill(orderID string) (bill model.Bill, status *model.Status)
	GetBillUser(req model.GetBillUserRequest) (resp model.GetBillUserResponse, status *model.Status)
	DeleteOrder(orderID string) (status *model.Status)
	UpdateOrder(orderMainID string, req model.UpdateOrderRequest, startTime time.Time) *model.Status
}

func NewOrderService(
	db *gorm.DB,
	redis *redis.Client,
	userRepo user.UserRepo,
	orderMainRepo order.OrderMainRepo,
	orderUserRepo order.OrderUserRepo,
	orderUserItemRepo order.OrderUserItemRepo,
	additionalCostRepo order.AdditionalCostRepo,
	discountsRepo order.DiscountRepo,
	deleteOrderRepo order.DeleteOrderRepo,
) OrderService {
	return &orderService{
		DB:                 db,
		Redis:              redis,
		UserRepo:           userRepo,
		OrderMainRepo:      orderMainRepo,
		OrderUserRepo:      orderUserRepo,
		OrderUserItemRepo:  orderUserItemRepo,
		AdditionalCostRepo: additionalCostRepo,
		DiscountRepo:       discountsRepo,
		DeleteOrderRepo:    deleteOrderRepo,
	}
}
