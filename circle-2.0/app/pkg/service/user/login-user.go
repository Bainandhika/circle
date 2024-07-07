package user

import (
	"circle-fiber/app/pkg/auth"
	customError "circle-fiber/lib/helper/custom-error"
	"circle-fiber/lib/model"
	"fmt"
	"net/http"
	"strings"
)

func (s *userService) LoginUser(req model.LoginUserRequest) (resp model.LoginUserResponse, status *model.Status) {
	funcName := "[Service - LoginUser]"

	user, err := s.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		if strings.EqualFold(err.Error(), customError.NotFoundError("user")) {
			return resp, &model.Status{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
		}

		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.UserRepo.GetUserByEmail: %v", funcName, err),
		}
	}

	if req.Password != user.Password {
		return resp, &model.Status{
			Code:    http.StatusBadRequest,
			Message: customError.IncorrectPassword(),
		}
	}

	token, err := auth.GenerateToken(req.Email)
	if err != nil {
		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at auth.GenerateToken: %v", funcName, err),
		}
	}

	resp = model.LoginUserResponse{
		UserID: user.ID,
		Token:  token,
	}

	return resp, status
}
