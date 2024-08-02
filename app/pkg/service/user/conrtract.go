package user

import (
	"sync"
	"time"

	"circle/app/pkg/repository/user"
	"circle/lib/model"

	"gorm.io/gorm"
)

type userService struct {
	DB       *gorm.DB
	UserRepo user.UserRepo
	mu       sync.Mutex
}

type UserService interface {
	CreateUser(req model.CreateUserRequest, startTime time.Time) (newUser model.Users, status *model.Status)
	LoginUser(req model.LoginUserRequest) (resp model.LoginUserResponse, status *model.Status)
	GetUserByID(id string) (user model.Users, status *model.Status)
	GetUsers() (users []model.Users, status *model.Status)
	UpdateUser(userID string, req model.UpdateUserRequest) *model.Status
	DeleteUser(userID string) *model.Status
}

func NewUserService(
	db *gorm.DB,
	userRepo user.UserRepo,
) UserService {
	return &userService{
		DB:       db,
		UserRepo: userRepo,
	}
}
