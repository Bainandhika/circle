package order

import (
	"errors"

	customError "circle-2.0/lib/helper/custom-error"
	"circle-2.0/lib/model"

	"gorm.io/gorm"
)

type orderUserItemRepo struct {
	DB *gorm.DB
}

type OrderUserItemRepo interface {
	CreateOrderUserItem(tx *gorm.DB, orderUserItem model.OrderUserItems) error
	GetOrderUserItemByOrderUserID(tx *gorm.DB, orderUserID string) ([]model.OrderUserItems, error)
	UpdateOrderUserItem(tx *gorm.DB, id string, orderUserItem map[string]any) error
}

func NewOrderUserItemRepo(db *gorm.DB) OrderUserItemRepo {
	return &orderUserItemRepo{
		DB: db,
	}
}

func (r *orderUserItemRepo) CreateOrderUserItem(tx *gorm.DB, orderUserItem model.OrderUserItems) error {
	if tx != nil {
		return tx.Create(&orderUserItem).Error
	} else {
		return r.DB.Create(&orderUserItem).Error
	}
}

func (r *orderUserItemRepo) GetOrderUserItemByOrderUserID(tx *gorm.DB, orderUserID string) ([]model.OrderUserItems, error) {
	var orderUserItems []model.OrderUserItems
	if tx != nil {
		if err := tx.Where("order_user_id = ?", orderUserID).Find(&orderUserItems).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.DB.Where("order_user_id = ?", orderUserID).Find(&orderUserItems).Error; err != nil {
			return nil, err
		}
	}

	return orderUserItems, nil
}

func (r *orderUserItemRepo) UpdateOrderUserItem(tx *gorm.DB, id string, orderUserItem map[string]any) error {
	var db *gorm.DB
	if tx == nil {
		db = r.DB
	} else {
		db = tx
	}

	updateActivity := db.Model(&model.OrderUserItems{}).Where("id = ?", id).Updates(orderUserItem)
	if err := updateActivity.Error; err != nil {
		return err
	}

	if updateActivity.RowsAffected == 0 {
		return errors.New(customError.NotFoundError("order_user_item"))
	}

	return nil
}
