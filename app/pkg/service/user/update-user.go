package user

import (
	customError "circle/lib/helper/custom-error"
	"circle/lib/model"
	"fmt"
	"net/http"
	"reflect"
	"strings"
)

func (s *userService) UpdateUser(userID string, req model.UpdateUserRequest) *model.Status {
	funcName := "[Service - UpdateUser]"

	updates := make(map[string]any)
	// Use reflection to populate the updates map
	v := reflect.ValueOf(req)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		if field.Kind() == reflect.String && field.String() != "" {
			jsonTag := fieldType.Tag.Get("json")
			if jsonTag != "" {
				updates[jsonTag] = field.String()
			}
		}
	}

	if err := s.UserRepo.UpdateUser(userID, updates); err != nil {
		if strings.EqualFold(err.Error(), customError.NotFoundError("user")) {
			return &model.Status{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
		}

		return &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.UserRepo.UpdateUser: %v", funcName, err),
		}
	}

	return nil
}
