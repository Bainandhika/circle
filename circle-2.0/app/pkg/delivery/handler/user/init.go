package user

import (
	"circle-fiber/app/pkg/service/user"
)

type UserHandler struct {
	UserService user.UserService
}

func NewUserHandler(
	userService user.UserService,
) *UserHandler {
	return &UserHandler{
		UserService: userService,
	}
}
