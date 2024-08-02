package order

import (
	"errors"

	customError "circle/lib/helper/custom-error"
	"circle/lib/model"

	"gorm.io/gorm"
)

type additionalCostRepo struct {
	DB *gorm.DB
}

type AdditionalCostRepo interface {
	CreateAdditionalCost(tx *gorm.DB, additionalCost model.AdditionalCosts) error
	GetAdditionalCostByOrderID(tx *gorm.DB, orderID string) ([]model.AdditionalCosts, error)
	UpdateAdditionalCost(tx *gorm.DB, id string, additionalCost map[string]any) error
}

func NewAdditionalCostRepo(db *gorm.DB) AdditionalCostRepo {
	return &additionalCostRepo{
		DB: db,
	}
}

func (r *additionalCostRepo) CreateAdditionalCost(tx *gorm.DB, additionalCost model.AdditionalCosts) error {
	if tx != nil {
		return tx.Create(&additionalCost).Error
	} else {
		return r.DB.Create(&additionalCost).Error
	}
}

func (r *additionalCostRepo) GetAdditionalCostByOrderID(tx *gorm.DB, orderID string) ([]model.AdditionalCosts, error) {
	var additional []model.AdditionalCosts
	if tx != nil {
		if err := tx.Where("order_id =?", orderID).Find(&additional).Error; err != nil {
			return nil, err
		}
	} else {
		if err := r.DB.Where("order_id =?", orderID).Find(&additional).Error; err != nil {
			return nil, err
		}
	}

	return additional, nil
}

func (r *additionalCostRepo) UpdateAdditionalCost(tx *gorm.DB, id string, additionalCost map[string]any) error {
	var db *gorm.DB
	if tx == nil {
		db = r.DB
	} else {
		db = tx
	}

	updateActivity := db.Model(&model.AdditionalCosts{}).Where("id = ?", id).Updates(additionalCost)
	if err := updateActivity.Error; err != nil {
		return err
	}

	if updateActivity.RowsAffected == 0 {
		return errors.New(customError.NotFoundError("additional_cost"))
	}

	return nil
}
