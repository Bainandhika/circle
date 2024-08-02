package order

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"circle/lib/helper/constant"
	"circle/lib/helper/tool"
	"circle/lib/logger"
	"circle/lib/model"
)

func (s *orderService) UpdateOrder(orderMainID string, req model.UpdateOrderRequest, startTime time.Time) *model.Status {
	funcName := "[Service - Update Order]"

	tx := s.DB.Begin(&sql.TxOptions{
		Isolation: sql.LevelReadUncommitted,
	})
	err := tx.Error
	if err != nil {
		return &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.DB.Begin(): %v", funcName, err),
		}
	}

	defer func(error) {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}(err)

	if len(req.OrderUserItems) > 0 {
		for _, value := range req.OrderUserItems {
			orderUserItem := tool.ToMap(value)
			delete(orderUserItem, "id")
			orderUserItem["updated_at"] = startTime
			orderUserItem["updated_by"] = req.UpdatedBy

			if err := s.OrderUserItemRepo.UpdateOrderUserItem(tx, value.ID, orderUserItem); err != nil {
				return &model.Status{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("%s error at s.OrderUserItemRepo.UpdateOrderUserItem: %v", funcName, err),
				}
			}
		}
	}

	if len(req.OrderAdditionalCosts) > 0 {
		for _, value := range req.OrderAdditionalCosts {
			additionalCost := tool.ToMap(value)
			delete(additionalCost, "id")
			additionalCost["updated_at"] = startTime
			additionalCost["updated_by"] = req.UpdatedBy

			if err := s.AdditionalCostRepo.UpdateAdditionalCost(tx, value.ID, additionalCost); err != nil {
				return &model.Status{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("%s error at s.AdditionalCostRepo.UpdateAdditionalCost: %v", funcName, err),
				}
			}
		}
	}

	if len(req.OrderDiscounts) > 0 {
		for _, value := range req.OrderDiscounts {
			discount := tool.ToMap(value)
			delete(discount, "id")
			discount["updated_at"] = startTime
			discount["updated_by"] = req.UpdatedBy

			if err := s.DiscountRepo.UpdateDiscount(tx, value.ID, discount); err != nil {
				return &model.Status{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("%s error at s.DiscountRepo.UpdateDiscount: %v", funcName, err),
				}
			}
		}
	}

	orderUsers, err := s.OrderUserRepo.GetOrderUserByOrderID(tx, orderMainID)
	if err != nil {
		return &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.OrderMainRepo.GetOrderByID: %v", funcName, err),
		}
	}

	totalPricePerUsers := make(map[string]float64)
	for _, value := range orderUsers {
		orderUserItems, err := s.OrderUserItemRepo.GetOrderUserItemByOrderUserID(tx, value.ID)
		if err != nil {
			return &model.Status{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("%s error at s.OrderUserItemRepo.GetOrderUserItemByOrderID: %v", funcName, err),
			}
		}

		var totalPrice float64
		for _, orderUserItem := range orderUserItems {
			totalPrice += orderUserItem.PricePerItem * float64(orderUserItem.Quantity)
		}
		totalPricePerUsers[value.ID] = totalPrice
	}

	var totalOrderBeforeAdditionalAndDiscounts float64
	for _, totalPricePerUser := range totalPricePerUsers {
		totalOrderBeforeAdditionalAndDiscounts += totalPricePerUser
	}

	partOfOrderPerUsers := make(map[string]float64)
	for orderUserID, totalPrice := range totalPricePerUsers {
		partOfOrderPerUsers[orderUserID] = tool.Round(totalPrice/totalOrderBeforeAdditionalAndDiscounts, constant.DecimalAmount)
	}

	additional, err := s.AdditionalCostRepo.GetAdditionalCostByOrderID(tx, orderMainID)
	if err != nil {
		return &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.AdditionalCostRepo.GetAdditionalCostByOrderID: %v", funcName, err),
		}
	}

	var additionalTotal float64
	for _, value := range additional {
		additionalTotal += value.Cost
	}

	discounts, err := s.DiscountRepo.GetDiscountsByOrderID(tx, orderMainID)
	if err != nil {
		return &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.DiscountRepo.GetDiscountsByOrderID: %v", funcName, err),
		}
	}

	var discountTotal float64
	for _, value := range discounts {
		discountTotal += value.Cost
	}

	totalPayment := totalOrderBeforeAdditionalAndDiscounts + additionalTotal - discountTotal

	for _, orderUser := range orderUsers {
		updateOrderUser := model.UpdateOrderUserRequest{
			ID:         orderUser.ID,
			TotalPrice: totalPricePerUsers[orderUser.ID],
			OrderPart:  partOfOrderPerUsers[orderUser.ID],
			PriceToPay: tool.Round(partOfOrderPerUsers[orderUser.ID]*totalPayment, constant.DecimalAmount),
		}

		updateOrderUserMap := tool.ToMap(updateOrderUser)
		delete(updateOrderUserMap, "id")
		updateOrderUserMap["updated_at"] = startTime
		updateOrderUserMap["updated_by"] = req.UpdatedBy

		if err := s.OrderUserRepo.UpdateOrderUser(tx, orderUser.ID, updateOrderUserMap); err != nil {
			return &model.Status{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("%s error at s.UpdateOrderUser: %v", funcName, err),
			}
		}
	}

	updateOrderMain := model.UpdateOrderMainRequest{
		OrderTitle:      req.NewOrderTitle,
		Total:           totalOrderBeforeAdditionalAndDiscounts,
		AdditionalTotal: additionalTotal,
		DiscountTotal:   discountTotal,
		TotalPayment:    totalPayment,
	}

	updateOrderMainMap := tool.ToMap(updateOrderMain)
	updateOrderMainMap["updated_at"] = startTime
	updateOrderMainMap["updated_by"] = req.UpdatedBy

	if err := s.OrderMainRepo.UpdateOrderMain(tx, orderMainID, updateOrderMainMap); err != nil {
		return &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.OrderMainRepo.UpdateOrderMain: %v", funcName, err),
		}
	}

	go func() {
		setBillCache := s.Redis.Del(s.Redis.Context(), tool.BillRedisKey(orderMainID))
		if err := setBillCache.Err(); err != nil {
			logger.Warning.Printf("%s error creating bill cache: %v", funcName, err)
		}
	}()

	return nil
}
