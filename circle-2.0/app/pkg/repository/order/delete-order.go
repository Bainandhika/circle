package order

import (
	"circle-fiber/lib/model"

	"gorm.io/gorm"
)

type deleteOrderRepo struct {
	DB *gorm.DB
}

type DeleteOrderRepo interface {
	DeleteOrderByOrderID(tx *gorm.DB, orderID string) error
}

func NewDeleteOrderRepo(db *gorm.DB) DeleteOrderRepo {
	return &deleteOrderRepo{
		DB: db,
	}
}

func (r *deleteOrderRepo) DeleteOrderByOrderID(tx *gorm.DB, orderID string) error {
	if err := tx.Where("order_user_id IN (SELECT id FROM order_users WHERE order_id = ?)", orderID).Delete(&model.OrderUserItems{}).Error; err != nil {
		return err
	}
	if err := tx.Where("order_id = ?", orderID).Delete(&model.OrderUsers{}).Error; err != nil {
		return err
	}
	if err := tx.Where("order_id = ?", orderID).Delete(&model.AdditionalCosts{}).Error; err != nil {
		return err
	}
	if err := tx.Where("order_id = ?", orderID).Delete(&model.Discounts{}).Error; err != nil {
		return err
	}
	if err := tx.Where("id = ?", orderID).Delete(&model.OrderMains{}).Error; err != nil {
		return err
	}
	return nil
}
