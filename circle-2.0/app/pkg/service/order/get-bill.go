package order

import (
	"circle-fiber/lib/helper/constant"
	customError "circle-fiber/lib/helper/custom-error"
	"circle-fiber/lib/helper/tool"
	"circle-fiber/lib/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-redis/redis/v8"
)

func (s *orderService) GetBill(orderID string) (bill model.Bill, status *model.Status) {
	funcName := "[Service - GetBill]"

	redisKey := tool.BillRedisKey(orderID)
	billString, err := s.Redis.Get(s.Redis.Context(), redisKey).Result()
	if err != nil && err != redis.Nil {
		return bill, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.Redis.Get: %v", funcName, err),
		}
	}

	if billString != "" {
		if err := json.Unmarshal([]byte(billString), &bill); err != nil {
			return bill, &model.Status{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("%s error at json.Unmarshal: %v", funcName, err),
			}
		}

		return bill, nil
	}

	order, err := s.OrderMainRepo.GetOrderByID(orderID)
	if err != nil {
		return bill, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.OrderMainRepo.GetOrderByID: %v", funcName, err),
		}
	}

	userCreator, err := s.UserRepo.GetUserByID(order.CreatedBy)
	if err != nil {
		if !strings.EqualFold(err.Error(), customError.NotFoundError("user")) {
			return bill, &model.Status{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("%s error at s.UserRepo.GetUserByID: %v", funcName, err),
			}
		}
	}

	var billCreatedBy string
	if userCreator.Name == "" {
		billCreatedBy = order.CreatedBy
	} else {
		billCreatedBy = userCreator.Name
	}

	additional, err := s.AdditionalCostRepo.GetAdditionalCostByOrderID(nil, order.ID)
	if err != nil {
		return bill, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.AdditionalCostRepo.GetAdditionalCostByOrderID: %v", funcName, err),
		}
	}

	var orderAdditionalCost []model.ResultAdditional
	for _, addition := range additional {
		orderAddition := model.ResultAdditional{
			AdditionalID: addition.ID,
			Type:         addition.Type,
			Cost:         addition.Cost,
		}

		orderAdditionalCost = append(orderAdditionalCost, orderAddition)
	}

	discounts, err := s.DiscountRepo.GetDiscountsByOrderID(nil, order.ID)
	if err != nil {
		return bill, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.DiscountRepo.GetDiscountsByOrderID: %v", funcName, err),
		}
	}

	var orderDiscounts []model.ResultDiscount
	for _, disc := range discounts {
		discount := model.ResultDiscount{
			DiscountID: disc.ID,
			Type:       disc.Type,
			Cost:       disc.Cost,
		}

		orderDiscounts = append(orderDiscounts, discount)
	}

	orderUsers, err := s.OrderUserRepo.GetOrderUserByOrderID(nil, order.ID)
	if err != nil {
		return bill, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.OrderUserRepo.GetOrderUserByOrderID: %v", funcName, err),
		}
	}

	var resultOrderUsers []model.ResultOrderUser
	for _, orderUser := range orderUsers {
		user, err := s.UserRepo.GetUserByID(orderUser.UserID)
		if err != nil {
			if !strings.EqualFold(err.Error(), customError.NotFoundError("user")) {
				return bill, &model.Status{
					Code:    http.StatusInternalServerError,
					Message: fmt.Sprintf("%s error at s.UserRepo.GetUserByID: %v", funcName, err),
				}
			}
		}

		orderUserItems, err := s.OrderUserItemRepo.GetOrderUserItemByOrderUserID(nil, orderUser.ID)
		if err != nil {
			return bill, &model.Status{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("%s error at s.OrderUserItemRepo.GetOrderUserItemByOrderID: %v", funcName, err),
			}
		}

		var resultOrderUserItems []model.ResultOrderUserItem
		for _, item := range orderUserItems {
			resultOrderUserItem := model.ResultOrderUserItem{
				OrderUserItemID: item.ID,
				Item:            item.Item,
				Quantity:        int(item.Quantity),
				PricePerItem:    item.PricePerItem,
				Total:           tool.Round(float64(item.Quantity)*item.PricePerItem, constant.DecimalAmount),
			}

			resultOrderUserItems = append(resultOrderUserItems, resultOrderUserItem)
		}

		resultOrderUser := model.ResultOrderUser{
			OrderUserID: orderUser.ID,
			UserID:      user.ID,
			Name:        user.Name,
			OrderItems:  resultOrderUserItems,
			Total:       orderUser.TotalPrice,
			PartOfOrder: orderUser.OrderPart,
			PriceToPay:  orderUser.PriceToPay,
		}

		resultOrderUsers = append(resultOrderUsers, resultOrderUser)
	}

	bill = model.Bill{
		CreatedBy: billCreatedBy,
		ResultOrderMain: model.ResultOrderMain{
			OrderMainID: order.ID,
			OrderTitle:  order.OrderTitle,
			Total:       order.Total,
			AdditionalAndDiscounts: model.ResultAdditionalAndDiscounts{
				Additional: orderAdditionalCost,
				Discounts:  orderDiscounts,
			},
			TotalPayment: order.TotalPayment,
		},
		ResultOrderUsers: resultOrderUsers,
	}

	go cacheBill(s.Redis, funcName, bill, orderID)

	return bill, nil
}
