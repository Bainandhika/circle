package user

import (
	"fmt"
	"net/http"

	"circle/lib/model"
)

func (s *userService) DeleteUser(userID string) *model.Status {
	var err error
	funcName := "[Service - DeleteUser]"

	if err = s.UserRepo.DeleteUser(userID); err != nil {
		return &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.UserRepo.DeleteUser: %v", funcName, err),
		}
	}

	return nil
}
