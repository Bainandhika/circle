package order

import (
	customError "circle-fiber/lib/helper/custom-error"
	"circle-fiber/lib/model"
	"errors"

	"gorm.io/gorm"
)

type discountRepo struct {
	DB *gorm.DB
}

type DiscountRepo interface {
	CreateDiscounts(tx *gorm.DB, discount model.Discounts) error
	GetDiscountsByOrderID(tx *gorm.DB, orderID string) ([]model.Discounts, error)
	UpdateDiscount(tx *gorm.DB, id string, discount map[string]any) error
}

func NewDiscountRepo(db *gorm.DB) DiscountRepo {
	return &discountRepo{
		DB: db,
	}
}

func (r *discountRepo) CreateDiscounts(tx *gorm.DB, discount model.Discounts) error {
	if tx != nil {
		return tx.Create(&discount).Error
	} else {
		return r.DB.Create(&discount).Error
	}
}

func (r *discountRepo) GetDiscountsByOrderID(tx *gorm.DB, orderID string) ([]model.Discounts, error) {
	var discounts []model.Discounts
	if tx != nil {
		if err := tx.Find(&discounts, "order_id = ?", orderID).Error; err != nil {
			return nil, err
        }
	} else {
		if err := r.DB.Find(&discounts, "order_id = ?", orderID).Error; err!= nil {
            return nil, err
        }
	}

	return discounts, nil
}

func (r *discountRepo) UpdateDiscount(tx *gorm.DB, id string, discount map[string]any) error {
	var db *gorm.DB
	if tx == nil {
		db = r.DB
	} else {
		db = tx
	}

	updateActivity := db.Model(&model.Discounts{}).Where("id = ?", id).Updates(discount)
	if err := updateActivity.Error; err != nil {
		return err
	}

	if updateActivity.RowsAffected == 0 {
		return errors.New(customError.NotFoundError("discount"))
	}

	return nil
}
