package order

import (
	"errors"

	customError "circle-2.0/lib/helper/custom-error"
	"circle-2.0/lib/model"

	"gorm.io/gorm"
)

type orderUserRepo struct {
	DB *gorm.DB
}

type OrderUserRepo interface {
	CreateOrderUser(tx *gorm.DB, orderUser model.OrderUsers) error
	GetOrderUserByOrderID(tx *gorm.DB, orderID string) ([]model.OrderUsers, error)
	GetOrderUserByUserIDAndOrderID(orderID, userID string) (model.OrderUsers, error)
	UpdateOrderUser(tx *gorm.DB, id string, orderUser map[string]any) error
}

func NewOrderUserRepo(db *gorm.DB) OrderUserRepo {
	return &orderUserRepo{
		DB: db,
	}
}

func (r *orderUserRepo) CreateOrderUser(tx *gorm.DB, orderUser model.OrderUsers) error {
	if tx != nil {
		return tx.Create(&orderUser).Error
	} else {
		return r.DB.Create(&orderUser).Error
	}
}

func (r *orderUserRepo) GetOrderUserByOrderID(tx *gorm.DB, orderID string) ([]model.OrderUsers, error) {
	var orderPeople []model.OrderUsers
	if tx != nil {
		if err := tx.Find(&orderPeople, "order_id =?", orderID).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.DB.Find(&orderPeople, "order_id =?", orderID).Error; err != nil {
			return nil, err
		}
	}

	return orderPeople, nil
}

func (r *orderUserRepo) GetOrderUserByUserIDAndOrderID(orderID, userID string) (model.OrderUsers, error) {
	var orderUser model.OrderUsers
	err := r.DB.First(&orderUser, "order_id = ? AND user_id = ?", orderID, userID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return orderUser, errors.New(customError.NotFoundError("order_user"))
	}

	return orderUser, err
}

func (r *orderUserRepo) UpdateOrderUser(tx *gorm.DB, id string, orderUser map[string]any) error {
	var db *gorm.DB
	if tx == nil {
		db = r.DB
	} else {
		db = tx
	}

	updateActivity := db.Model(&model.OrderUsers{}).Where("id = ?", id).Updates(orderUser)
	if err := updateActivity.Error; err != nil {
		return err
	}

	if updateActivity.RowsAffected == 0 {
		return errors.New(customError.NotFoundError("order_user"))
	}

	return nil
}
