package user

import (
	"circle/lib/model"
	"fmt"
	"net/http"
)

func (s *userService) GetUsers() (users []model.Users, status *model.Status) {
	funcName := "[Service - GetUsers]"

	users, err := s.UserRepo.GetAllUsers()
	if err != nil {
		return nil, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.UserRepo.GetAllUsers(): %v", funcName, err),
		}
	}

	return users, nil
}
