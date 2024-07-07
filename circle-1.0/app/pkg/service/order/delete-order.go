package order

import (
	"circle/lib/model"
	"fmt"
	"net/http"
)

func (s *orderService) DeleteOrder(orderID string) (status *model.Status) {
	var err error
	funcName := "[Service - DeleteOrder]"

	tx := s.DB.Begin()
	if err = tx.Error; err != nil {
		return &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at o.DB.Begin(): %v", funcName, err),
		}
	}

	defer func(error) {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}(err)

	if err := s.DeleteOrderRepo.DeleteOrderByOrderID(tx, orderID); err != nil {
		return &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.DeleteOrderRepo.DeleteOrderByOrderID: %v", funcName, err),
		}
	}

	return nil
}
