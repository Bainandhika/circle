package router

import (
	"circle-fiber/app/pkg/delivery/handler/health"
	orderHandler "circle-fiber/app/pkg/delivery/handler/order"
	userHandler "circle-fiber/app/pkg/delivery/handler/user"
	"circle-fiber/app/pkg/delivery/middleware"
	orderRepo "circle-fiber/app/pkg/repository/order"
	userRepo "circle-fiber/app/pkg/repository/user"
	orderService "circle-fiber/app/pkg/service/order"
	userService "circle-fiber/app/pkg/service/user"
	"circle-fiber/lib/helper/constant"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, redis *redis.Client) *fiber.App {
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

	r := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})

	r.Use(middleware.CaptureRequest, middleware.Headers, middleware.LoggingAPIDetail)

	publicRoutes(r, userHandler)
	authenticatedRoutes(r, userHandler)
	orderRoutes(r, orderHandler)

	return r
}

func publicRoutes(r *fiber.App, userHandler *userHandler.UserHandler) {
	r.Get(constant.HealthCheckPath, health.HealthCheck)
	r.Post(constant.CreateUserPath, userHandler.CreateUser)
	r.Post(constant.LoginUserPath, userHandler.Login)
}

func authenticatedRoutes(r *fiber.App, userHandler *userHandler.UserHandler) {
	userGroup := r.Group(constant.UserGroupPrefix)
	userGroup.Use(middleware.JWTAuth)
	{
		userGroup.Get(constant.GetUsersPrefix, userHandler.GetUsers)
		userGroup.Get(constant.GetUserByIDPrefix, middleware.CapturePathParams, userHandler.GetUserByID)
		userGroup.Put(constant.UpdateUserPrefix, middleware.CapturePathParams, userHandler.UpdateUser)
		userGroup.Delete(constant.DeleteUserPrefix, middleware.CapturePathParams, userHandler.DeleteUser)
	}
}

func orderRoutes(r *fiber.App, orderHandler *orderHandler.OrderHandler) {
	orderGroup := r.Group(constant.OrderGroupPrefix)
	{
		orderGroup.Post(constant.CreateOrderPrefix, orderHandler.CreateOrder)
		orderGroup.Get(constant.GetOrdersPrefix, orderHandler.GetOrders)
		orderGroup.Get(constant.GetBillPrefix, middleware.CapturePathParams, orderHandler.GetBill)
		orderGroup.Post(constant.GetBillUserPrefix, orderHandler.GetBillUser)
		orderGroup.Delete(constant.DeleteOrderPrefix, middleware.CapturePathParams, orderHandler.DeleteOrder)
		orderGroup.Put(constant.UpdateOrderPrefix, middleware.CapturePathParams, orderHandler.UpdateOrder)
	}
}
