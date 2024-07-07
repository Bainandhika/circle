package constant

const (
	ParadisePrefix = "/circle"

	HealthCheckPath = ParadisePrefix + "/health"
	CreateUserPath  = ParadisePrefix + "/new/user"
	LoginUserPath   = ParadisePrefix + "/login"

	UserGroupPrefix   = ParadisePrefix + "/users"
	GetUsersPrefix    = "/get/all"
	GetUserByIDPrefix = "/get/:id"
	UpdateUserPrefix  = "/update/:id"
	DeleteUserPrefix  = "/delete/:id"

	OrderGroupPrefix  = ParadisePrefix + "/order"
	CreateOrderPrefix = "/new"
	GetOrdersPrefix   = "/get/all"
	GetBillPrefix     = "/get/bill/:id"
	GetBillUserPrefix = "/get/bill/user"
	DeleteOrderPrefix = "/delete/:id"
	UpdateOrderPrefix = "/update/:id"
)
