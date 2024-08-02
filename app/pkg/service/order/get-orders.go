package order

import (
	"fmt"
	"net/http"

	"circle/lib/model"
)

func (s *orderService) GetOrders() (orders []model.OrderMains, status *model.Status) {
	funcName := "[Service - GetOrders]"

	orders, err := s.OrderMainRepo.GetOrders()
	if err != nil {
		return nil, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.OrderMainRepo.GetOrders(): %v", funcName, err),
		}
	}

	return orders, nil
}
