package router

import (
	"circle/app/pkg/delivery/handler/health"
	orderHandler "circle/app/pkg/delivery/handler/order"
	userHandler "circle/app/pkg/delivery/handler/user"
	"circle/app/pkg/delivery/middleware"
	orderRepo "circle/app/pkg/repository/order"
	userRepo "circle/app/pkg/repository/user"
	orderService "circle/app/pkg/service/order"
	userService "circle/app/pkg/service/user"
	"circle/lib/helper/constant"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *redis.Client) *gin.Engine {
	userRepo := userRepo.NewUserRepo(db)

	orderMainRepo := orderRepo.NewOrderMainRepo(db)
	additionalCostRepo := orderRepo.NewAdditionalCostRepo(db)
	discountRepo := orderRepo.NewDiscountRepo(db)
	orderUserRepo := orderRepo.NewOrderUserRepo(db)
	orderUserItemRepo := orderRepo.NewOrderUserItemRepo(db)
	deleteOrderRepo := orderRepo.NewDeleteOrderRepo(db)

	userService := userService.NewUserService(db, userRepo)
	orderService := orderService.NewOrderService(
		db,
		redis,
		userRepo,
		orderMainRepo,
		orderUserRepo,
		orderUserItemRepo,
		additionalCostRepo,
		discountRepo,
		deleteOrderRepo,
	)

	orderHandler := orderHandler.NewOrderHandler(orderService, userService)
	userHandler := userHandler.NewUserHandler(userService)

	r := gin.Default()
	r.Use(middleware.CaptureRequest(), middleware.Headers(), middleware.LoggingAPIDetail())

	publicRoutes(r, userHandler)
	authenticatedRoutes(r, userHandler)
	orderRoutes(r, orderHandler)

	return r
}

func publicRoutes(r *gin.Engine, userHandler *userHandler.UserHandler) {
	r.GET(constant.HealthCheckPath, health.HealthCheck)
	r.POST(constant.CreateUserPath, userHandler.CreateUser)
	r.POST(constant.LoginUserPath, userHandler.Login)
}

func authenticatedRoutes(r *gin.Engine, userHandler *userHandler.UserHandler) {
	userGroup := r.Group(constant.UserGroupPrefix)
	userGroup.Use(middleware.JWTAuth())
	{
		userGroup.GET(constant.GetUsersPrefix, userHandler.GetUsers)
		userGroup.GET(constant.GetUserByIDPrefix, userHandler.GetUserByID)
		userGroup.PUT(constant.UpdateUserPrefix, userHandler.UpdateUser)
		userGroup.DELETE(constant.DeleteUserPrefix, userHandler.DeleteUser)
	}
}

func orderRoutes(r *gin.Engine, orderHandler *orderHandler.OrderHandler) {
	orderGroup := r.Group(constant.OrderGroupPrefix)
	{
		orderGroup.POST(constant.CreateOrderPrefix, orderHandler.CreateOrder)
		orderGroup.GET(constant.GetOrdersPrefix, orderHandler.GetOrders)
		orderGroup.GET(constant.GetBillPrefix, orderHandler.GetBill)
		orderGroup.POST(constant.GetBillUserPrefix, orderHandler.GetBillUser)
		orderGroup.DELETE(constant.DeleteOrder, orderHandler.DeleteOrder)
	}
}
