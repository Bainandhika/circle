package user

import (
	"circle-fiber/lib/helper/tool"
	"circle-fiber/lib/model"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

func (s *userService) CreateUser(req model.CreateUserRequest, startTime time.Time) (newUser model.Users, status *model.Status) {
	funcName := "[Service - CreateUser]"

	req.Name = tool.CapitalizeEachWord(req.Name)

	if err := copier.Copy(&newUser, &req); err != nil {
		return newUser, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at copier.Copy(&newUser, &req): %v", funcName, err),
		}
	}

	newUser.ID = uuid.NewString()
	newUser.CreatedAt = startTime
	newUser.UpdatedAt = startTime
	newUser.UpdatedBy = newUser.ID

	if err := s.UserRepo.CreateUser(nil, newUser); err != nil {
		return newUser, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.UserRepo.CreateUser: %v", funcName, err),
		}
	}

	return newUser, nil
}
