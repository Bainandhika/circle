package order

import (
	customError "circle/lib/helper/custom-error"
	"circle/lib/model"
	"errors"

	"gorm.io/gorm"
)

type orderMainRepo struct {
	DB *gorm.DB
}

type OrderMainRepo interface {
	CreateOrderMain(tx *gorm.DB, orderMain model.OrderMains) error
	GetOrders() ([]model.OrderMains, error)
	GetOrderByID(id string) (model.OrderMains, error)
	UpdateOrderMain(tx *gorm.DB, id string, orderMain map[string]any) error
}

func NewOrderMainRepo(db *gorm.DB) OrderMainRepo {
	return &orderMainRepo{
		DB: db,
	}
}

func (r *orderMainRepo) CreateOrderMain(tx *gorm.DB, orderMain model.OrderMains) error {
	if tx != nil {
		return tx.Create(&orderMain).Error
	} else {
		return r.DB.Create(&orderMain).Error
	}
}

func (r *orderMainRepo) GetOrders() ([]model.OrderMains, error) {
	var orders []model.OrderMains
	err := r.DB.Find(&orders).Order("created_at DESC").Error
	return orders, err
}

func (r *orderMainRepo) GetOrderByID(id string) (model.OrderMains, error) {
	var order model.OrderMains
	err := r.DB.First(&order, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return order, errors.New(customError.NotFoundError("order"))
	}
	return order, err
}

func (r *orderMainRepo) UpdateOrderMain(tx *gorm.DB, id string, orderMain map[string]any) error {
	var db *gorm.DB
	if tx == nil {
		db = r.DB
	} else {
		db = tx
	}

	updateActivity := db.Model(&model.OrderMains{}).Where("id = ?", id).Updates(orderMain)
	if err := updateActivity.Error; err != nil {
		return err
	}

	if updateActivity.RowsAffected == 0 {
		return errors.New(customError.NotFoundError("order_main"))
	}

	return nil
}
