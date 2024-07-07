package user

import (
	"circle/app/pkg/repository/user"
	"circle/lib/model"
	"time"

	"gorm.io/gorm"
)

type userService struct {
	DB       *gorm.DB
	UserRepo user.UserRepo
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
