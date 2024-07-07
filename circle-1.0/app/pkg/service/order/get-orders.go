package order

import (
	"circle/lib/model"
	"fmt"
	"net/http"
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
