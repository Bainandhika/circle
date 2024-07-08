package user

import (
	"fmt"
	"net/http"
	"time"

	"circle-2.0/lib/helper/tool"
	"circle-2.0/lib/model"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

func (s *userService) CreateUser(req model.CreateUserRequest, startTime time.Time) (newUser model.Users, status *model.Status) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
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
