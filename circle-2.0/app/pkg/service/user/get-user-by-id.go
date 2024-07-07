package user

import (
	customError "circle-fiber/lib/helper/custom-error"
	"circle-fiber/lib/model"
	"errors"
	"fmt"
	"net/http"
)

func (s *userService) GetUserByID(id string) (user model.Users, status *model.Status) {
	funcName := "[Service - GetUserByID]"

	user, err := s.UserRepo.GetUserByID(id)
	if err != nil {
		if errors.Is(err, errors.New(customError.NotFoundError("user"))) {
			return user, &model.Status{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
		}

		return user, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.UserRepo.GetUserByID: %v", funcName, err),
		}
	}

	return user, nil
}
