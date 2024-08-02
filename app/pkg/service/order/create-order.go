package order

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"circle/lib/helper/constant"
	customError "circle/lib/helper/custom-error"
	"circle/lib/helper/tool"
	"circle/lib/logger"
	"circle/lib/model"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (s *orderService) CreateOrder(req model.OrderRequest, startTime time.Time) (resp model.Bill, status *model.Status) {
	var err error
	funcName := "[Service - CreateOrder]"

	for _, s := range req.OrderDetailUsers {
		s.Name = tool.CapitalizeEachWord(s.Name)
	}

	tx := s.DB.Begin()
	if err = tx.Error; err != nil {
		return resp, &model.Status{
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

	userCreator, err := s.UserRepo.GetUserByID(req.UserID)
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

	userNameAndIDs, err := getUserIDs(s, tx, req, startTime)
	if err != nil {
		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at getUserIDs: %v", funcName, err),
		}
	}

	totalPriceUsers, totalPriceUserItems := calculateTotalPrices(req)
	totalPayment, orderMain, err := createOrderMain(s, tx, req, totalPriceUsers, startTime)
	if err != nil {
		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at createOrderMain: %v", funcName, err),
		}
	}

	resultAdditional, resultDiscounts, err := createAdditionalAndDiscounts(s, tx, req, orderMain, startTime)
	if err != nil {
		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at createAdditionalAndDiscounts: %v", funcName, err),
		}
	}

	partOfOrderUser := calculatePartOfOrder(totalPriceUsers, orderMain.Total)
	priceToPayUsers := calculatePriceToPay(partOfOrderUser, totalPayment)

	resultOrderUsers, err := createOrderUsers(s, tx, req, orderMain, userNameAndIDs, totalPriceUsers, totalPriceUserItems, partOfOrderUser, priceToPayUsers, startTime)
	if err != nil {
		return resp, &model.Status{
			Code:    http.StatusInternalServerError,
			Message: fmt.Sprintf("%s error at createOrderUsers: %v", funcName, err),
		}
	}

	bill := model.Bill{
		CreatedBy: userCreator.Name,
		ResultOrderMain: model.ResultOrderMain{
			OrderMainID: orderMain.ID,
			OrderTitle:  orderMain.OrderTitle,
			Total:       orderMain.Total,
			AdditionalAndDiscounts: model.ResultAdditionalAndDiscounts{
				Additional: resultAdditional,
				Discounts:  resultDiscounts,
			},
			TotalPayment: orderMain.TotalPayment,
		},
		ResultOrderUsers: resultOrderUsers,
	}

	go cacheBill(s.Redis, funcName, bill, orderMain.ID)

	return bill, nil
}

func getUserIDs(s *orderService, tx *gorm.DB, req model.OrderRequest, startTime time.Time) (map[string]string, error) {
	userNameAndIDs := make(map[string]string)
	for _, orderUser := range req.OrderDetailUsers {
		user, err := s.UserRepo.GetUserByName(orderUser.Name)
		if err != nil {
			if err.Error() == customError.NotFoundError("user") {
				user = model.Users{
					ID:        uuid.NewString(),
					Name:      orderUser.Name,
					Email:     orderUser.Name + "@something.com",
					Password:  "123456",
					CreatedAt: startTime,
					UpdatedAt: startTime,
				}

				if err = s.UserRepo.CreateUser(tx, user); err != nil {
					return nil, err
				}
			} else {
				return nil, err
			}
		}
		userNameAndIDs[user.Name] = user.ID
	}
	return userNameAndIDs, nil
}

func calculateTotalPrices(req model.OrderRequest) (map[string]float64, map[string]map[string]float64) {
	totalPriceUserItems := make(map[string]map[string]float64)
	totalPriceUsers := make(map[string]float64)

	for _, orderDetail := range req.OrderDetailUsers {
		totalPriceItems := make(map[string]float64)
		var totalPriceUser float64

		for _, orderItem := range orderDetail.OrderItems {
			total := float64(orderItem.Quantity) * orderItem.PricePerItem
			totalPriceItems[orderItem.Item] = total
			totalPriceUser += total
		}

		totalPriceUserItems[orderDetail.Name] = totalPriceItems
		totalPriceUsers[orderDetail.Name] = totalPriceUser
	}

	return totalPriceUsers, totalPriceUserItems
}

func createOrderMain(s *orderService, tx *gorm.DB, req model.OrderRequest, totalPriceUsers map[string]float64, startTime time.Time) (float64, model.OrderMains, error) {
	totalPriceOrder := sum(totalPriceUsers)
	totalAdditional := sumAdditionalCosts(req.OrderAdditionalAndDiscounts.OrderAdditionalCosts)
	totalDiscounts := sumDiscounts(req.OrderAdditionalAndDiscounts.OrderDiscounts)
	totalPayment := totalPriceOrder + totalAdditional - totalDiscounts

	orderMain := model.OrderMains{
		ID:              uuid.NewString(),
		OrderTitle:      req.OrderTitle,
		Total:           totalPriceOrder,
		AdditionalTotal: totalAdditional,
		DiscountTotal:   totalDiscounts,
		TotalPayment:    totalPayment,
		CreatedAt:       startTime,
		CreatedBy:       req.UserID,
		UpdatedAt:       startTime,
		UpdatedBy:       req.UserID,
	}

	if err := s.OrderMainRepo.CreateOrderMain(tx, orderMain); err != nil {
		return 0, orderMain, err
	}

	return totalPayment, orderMain, nil
}

func createAdditionalAndDiscounts(s *orderService, tx *gorm.DB, req model.OrderRequest, orderMain model.OrderMains, startTime time.Time) ([]model.ResultAdditional, []model.ResultDiscount, error) {
	var resultAdditional []model.ResultAdditional
	for _, additional := range req.OrderAdditionalAndDiscounts.OrderAdditionalCosts {
		addition := model.AdditionalCosts{
			ID:        uuid.NewString(),
			OrderID:   orderMain.ID,
			Type:      additional.Type,
			Cost:      additional.Cost,
			CreatedAt: startTime,
			CreatedBy: req.UserID,
			UpdatedAt: startTime,
			UpdatedBy: req.UserID,
		}
		if err := s.AdditionalCostRepo.CreateAdditionalCost(tx, addition); err != nil {
			return nil, nil, err
		}

		resultAdd := model.ResultAdditional{
			AdditionalID: addition.ID,
			Type:         addition.Type,
			Cost:         addition.Cost,
		}
		resultAdditional = append(resultAdditional, resultAdd)
	}

	var resultDiscounts []model.ResultDiscount
	for _, discount := range req.OrderAdditionalAndDiscounts.OrderDiscounts {
		disc := model.Discounts{
			ID:        uuid.NewString(),
			OrderID:   orderMain.ID,
			Type:      discount.Type,
			Cost:      discount.Cost,
			CreatedAt: startTime,
			CreatedBy: req.UserID,
			UpdatedAt: startTime,
			UpdatedBy: req.UserID,
		}
		if err := s.DiscountRepo.CreateDiscounts(tx, disc); err != nil {
			return nil, nil, err
		}

		resultDiscount := model.ResultDiscount{
			DiscountID: disc.ID,
			Type:       disc.Type,
			Cost:       disc.Cost,
		}
		resultDiscounts = append(resultDiscounts, resultDiscount)
	}

	return resultAdditional, resultDiscounts, nil
}

func createOrderUsers(s *orderService, tx *gorm.DB, req model.OrderRequest, orderMain model.OrderMains, userNameAndIDs map[string]string, totalPriceUsers map[string]float64, totalPriceUserItems map[string]map[string]float64, partOfOrderUser, priceToPayUsers map[string]float64, startTime time.Time) ([]model.ResultOrderUser, error) {
	var resultOrderUsers []model.ResultOrderUser
	userNameAndOrderUserIDs := make(map[string]string)
	resultNameAndOrderUserItems := make(map[string][]model.ResultOrderUserItem)

	for name, id := range userNameAndIDs {
		orderUser := model.OrderUsers{
			ID:         uuid.NewString(),
			OrderID:    orderMain.ID,
			UserID:     id,
			TotalPrice: totalPriceUsers[name],
			OrderPart:  partOfOrderUser[name],
			PriceToPay: priceToPayUsers[name],
			CreatedAt:  startTime,
			CreatedBy:  id,
			UpdatedAt:  startTime,
			UpdatedBy:  id,
		}
		if err := s.OrderUserRepo.CreateOrderUser(tx, orderUser); err != nil {
			return nil, err
		}
		userNameAndOrderUserIDs[name] = orderUser.ID

		if len(resultNameAndOrderUserItems) == 0 {
			resultNameAndOrderUserItems = makeResultOrderUserItems(req.OrderDetailUsers, totalPriceUserItems)
		}

		var newResultOrderUserItems []model.ResultOrderUserItem
		for _, resultOrderUserItem := range resultNameAndOrderUserItems[name] {
			item := model.OrderUserItems{
				ID:           uuid.NewString(),
				OrderUserID:  orderUser.ID,
				Item:         resultOrderUserItem.Item,
				Quantity:     int8(resultOrderUserItem.Quantity),
				PricePerItem: resultOrderUserItem.PricePerItem,
				CreatedAt:    startTime,
				CreatedBy:    id,
				UpdatedAt:    startTime,
				UpdatedBy:    id,
			}
			if err := s.OrderUserItemRepo.CreateOrderUserItem(tx, item); err != nil {
				return nil, err
			}

			resultOrderUserItem.OrderUserItemID = item.ID
			newResultOrderUserItems = append(newResultOrderUserItems, resultOrderUserItem)
		}

		resultNameAndOrderUserItems[name] = newResultOrderUserItems
	}

	resultOrderUsers = makeResultOrderUsers(userNameAndIDs, userNameAndOrderUserIDs, req.OrderDetailUsers, resultNameAndOrderUserItems, totalPriceUsers, partOfOrderUser, priceToPayUsers)

	return resultOrderUsers, nil
}

func cacheBill(redis *redis.Client, funcName string, bill model.Bill, orderMainID string) {
	billBytes, _ := json.Marshal(bill)
	setBillCache := redis.Set(redis.Context(), tool.BillRedisKey(orderMainID), string(billBytes), constant.RedisExpiredDuration)
	if err := setBillCache.Err(); err != nil {
		logger.Warning.Printf("%s error at creating bill cache: %v", funcName, err)
	}

	if err := tool.CreateBillTextFile(bill); err != nil {
		logger.Warning.Printf("%s error creating bill text file: %v", funcName, err)
	}
}

func sum(m map[string]float64) (total float64) {
	for _, value := range m {
		total += value
	}
	return
}

func sumAdditionalCosts(costs []model.OrderAdditionalCost) (total float64) {
	for _, cost := range costs {
		total += cost.Cost
	}
	return
}

func sumDiscounts(discounts []model.OrderDiscount) (total float64) {
	for _, discount := range discounts {
		total += discount.Cost
	}
	return
}

func calculatePartOfOrder(totalPriceUsers map[string]float64, totalPriceOrder float64) map[string]float64 {
	partOfOrderUser := make(map[string]float64)
	for name, totalPrice := range totalPriceUsers {
		partOfOrderUser[name] = tool.Round(totalPrice/totalPriceOrder, constant.DecimalAmount)
	}
	return partOfOrderUser
}

func calculatePriceToPay(partOfOrderUser map[string]float64, totalPayment float64) map[string]float64 {
	priceToPayUsers := make(map[string]float64)
	for name, part := range partOfOrderUser {
		priceToPayUsers[name] = tool.Round(part*totalPayment, constant.DecimalAmount)
	}
	return priceToPayUsers
}

func makeResultOrderUserItems(orderDetailUsers []model.OrderDetailUser, totalPriceUserItems map[string]map[string]float64) map[string][]model.ResultOrderUserItem {
	result := make(map[string][]model.ResultOrderUserItem)
	for _, orderDetail := range orderDetailUsers {
		var resultOrderItems []model.ResultOrderUserItem
		for _, orderItem := range orderDetail.OrderItems {
			resultOrderItems = append(resultOrderItems, model.ResultOrderUserItem{
				Item:         orderItem.Item,
				Quantity:     orderItem.Quantity,
				PricePerItem: orderItem.PricePerItem,
				Total:        totalPriceUserItems[orderDetail.Name][orderItem.Item],
			})
		}
		result[orderDetail.Name] = resultOrderItems
	}
	return result
}

func makeResultOrderUsers(userNameAndIDs map[string]string, userNameAndOrderUserIDs map[string]string, orderDetailUsers []model.OrderDetailUser, resultOrderUserItems map[string][]model.ResultOrderUserItem, totalPriceUsers, partOfOrderUser, priceToPayUsers map[string]float64) []model.ResultOrderUser {
	var result []model.ResultOrderUser
	for _, orderDetailUser := range orderDetailUsers {
		result = append(result, model.ResultOrderUser{
			OrderUserID: userNameAndOrderUserIDs[orderDetailUser.Name],
			UserID:      userNameAndIDs[orderDetailUser.Name],
			Name:        orderDetailUser.Name,
			OrderItems:  resultOrderUserItems[orderDetailUser.Name],
			Total:       totalPriceUsers[orderDetailUser.Name],
			PartOfOrder: partOfOrderUser[orderDetailUser.Name],
			PriceToPay:  priceToPayUsers[orderDetailUser.Name],
		})
	}

	return result
}
