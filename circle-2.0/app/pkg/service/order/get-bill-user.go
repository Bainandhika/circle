package order

import (
	"fmt"
	"net/http"
	"strings"

	"circle-2.0/lib/helper/constant"
	customError "circle-2.0/lib/helper/custom-error"
	"circle-2.0/lib/helper/tool"
	"circle-2.0/lib/model"
)

func (s *orderService) GetBillUser(req model.GetBillUserRequest) (resp model.GetBillUserResponse, status *model.Status) {
	funcName := "[Service - GetBillUser]"

	user, err := s.UserRepo.GetUserByID(req.UserID)
	if err != nil {
		if strings.EqualFold(err.Error(), customError.NotFoundError("user")) {
			return resp, &model.Status{
				Code:    http.StatusBadRequest,
				Message: customError.MustHaveAccount(),
			}
		} else {
			return resp, &model.Status{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("%s error at s.UserRepo.GetUserByID: %v", funcName, err),
			}
		}
	}

	order, err := s.OrderMainRepo.GetOrderByID(req.OrderID)
	if err != nil {
		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.OrderMainRepo.GetOrderByID: %v", funcName, err),
		}
	}

	additionalCosts, err := s.AdditionalCostRepo.GetAdditionalCostByOrderID(nil, req.OrderID)
	if err != nil {
		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.AdditionalCostRepo.GetAdditionalCostByOrderID: %v", funcName, err),
		}
	}

	var resultAdditionalCosts []model.ResultAdditional
	for _, addition := range additionalCosts {
		resultAdditionalCost := model.ResultAdditional{
			AdditionalID: addition.ID,
			Type:         addition.Type,
			Cost:         addition.Cost,
		}

		resultAdditionalCosts = append(resultAdditionalCosts, resultAdditionalCost)
	}

	discounts, err := s.DiscountRepo.GetDiscountsByOrderID(nil, req.OrderID)
	if err != nil {
		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.DiscountRepo.GetDiscountsByOrderID: %v", funcName, err),
		}
	}

	var resultDiscounts []model.ResultDiscount
	for _, discount := range discounts {
		resultDiscount := model.ResultDiscount{
			DiscountID: discount.ID,
			Type:       discount.Type,
			Cost:       discount.Cost,
		}

		resultDiscounts = append(resultDiscounts, resultDiscount)
	}

	orderUser, err := s.OrderUserRepo.GetOrderUserByUserIDAndOrderID(req.OrderID, req.UserID)
	if err != nil {
		if strings.EqualFold(err.Error(), customError.NotFoundError("order_user")) {
			return resp, &model.Status{
				Code:    http.StatusBadRequest,
				Message: err.Error(),
			}
		} else {
			return resp, &model.Status{
				Code:    http.StatusInternalServerError,
				Message: fmt.Sprintf("%s error at s.OrderUserRepo.GetOrderUserByUserIDAndOrderID: %v", funcName, err),
			}
		}
	}

	orderUserItems, err := s.OrderUserItemRepo.GetOrderUserItemByOrderUserID(nil, orderUser.ID)
	if err != nil {
		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at s.OrderUserItemRepo.GetOrderUserItemByOrderID: %v", funcName, err),
		}
	}

	var resultOrderItems []model.ResultOrderUserItem
	for _, item := range orderUserItems {
		resultOrderItem := model.ResultOrderUserItem{
			OrderUserItemID: item.ID,
			Item:            item.Item,
			Quantity:        int(item.Quantity),
			PricePerItem:    item.PricePerItem,
			Total:           tool.Round(float64(item.Quantity)*item.PricePerItem, constant.DecimalAmount),
		}

		resultOrderItems = append(resultOrderItems, resultOrderItem)
	}

	billUser := model.GetBillUserResponse{
		ResultOrderMain: model.ResultOrderMain{
			OrderMainID: order.ID,
			OrderTitle:  order.OrderTitle,
			Total:       order.Total,
			AdditionalAndDiscounts: model.ResultAdditionalAndDiscounts{
				Additional: resultAdditionalCosts,
				Discounts:  resultDiscounts,
			},
			TotalPayment: order.TotalPayment,
		},
		ResultOrderUser: model.ResultOrderUser{
			OrderUserID: orderUser.ID,
			UserID:      user.ID,
			Name:        user.Name,
			OrderItems:  resultOrderItems,
			Total:       orderUser.TotalPrice,
			PartOfOrder: orderUser.OrderPart,
			PriceToPay:  orderUser.PriceToPay,
		},
	}

	return billUser, nil
}
