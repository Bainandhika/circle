package model

import "time"

type Users struct {
	ID        string    `gorm:"type:varchar(255);primaryKey" json:"id"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex:uni_user_name" json:"name"`
	Email     string    `gorm:"type:varchar(255);not null;uniqueIndex:uni_user_email" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null;default:'123456'" json:"-"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"updated_at"`
	UpdatedBy string    `gorm:"type:varchar(255);null" json:"updated_by"`
}

type OrderUsers struct {
	ID         string    `gorm:"type:varchar(255);primaryKey" json:"id"`
	OrderID    string    `gorm:"type:varchar(255);not null;index:fk1_order_person_order_main_id" json:"order_id"`
	UserID     string    `gorm:"type:varchar(255);not null;index:fk2_order_person_user_id" json:"user_id"`
	TotalPrice float64   `gorm:"type:decimal(10,2);not null;default:0.00" json:"total_price"`
	OrderPart  float64   `gorm:"type:decimal(2,2);not null" json:"order_part"`
	PriceToPay float64   `gorm:"type:decimal(10,2);not null" json:"price_to_pay"`
	CreatedAt  time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"created_at"`
	CreatedBy  string    `gorm:"type:varchar(255);null" json:"created_by"`
	UpdatedAt  time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"updated_at"`
	UpdatedBy  string    `gorm:"type:varchar(255);null" json:"updated_by"`

	OrderMain OrderMains `gorm:"foreignKey:OrderID;references:ID" json:"-"`
	User      Users      `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

type OrderUserItems struct {
	ID           string    `gorm:"type:varchar(255);primaryKey" json:"id"`
	OrderUserID  string    `gorm:"type:varchar(255);not null;index:fk1_order_user_item_order_user_id" json:"order_user_id"`
	Item         string    `gorm:"type:varchar(255);not null" json:"item"`
	Quantity     int8      `gorm:"type:tinyint(4);not null;default:0" json:"quantity"`
	PricePerItem float64   `gorm:"type:decimal(10,2);not null;default:0.00" json:"price_per_item"`
	CreatedAt    time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"created_at"`
	CreatedBy    string    `gorm:"type:varchar(255);null" json:"created_by"`
	UpdatedAt    time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"updated_at"`
	UpdatedBy    string    `gorm:"type:varchar(255);null" json:"updated_by"`

	OrderUser OrderUsers `gorm:"foreignKey:OrderUserID;references:ID" json:"-"`
}

type OrderMains struct {
	ID              string    `gorm:"type:varchar(255);primaryKey" json:"id"`
	OrderTitle      string    `gorm:"type:varchar(500);uniqueIndex:uni_order_title;not null" json:"order_title"`
	Total           float64   `gorm:"type:decimal(20,2);not null" json:"total"`
	AdditionalTotal float64   `gorm:"type:decimal(10,2);not null" json:"additional_total"`
	DiscountTotal   float64   `gorm:"type:decimal(10,2);not null" json:"discount_total"`
	TotalPayment    float64   `gorm:"type:decimal(20,2);not null" json:"total_payment"`
	CreatedAt       time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"created_at"`
	CreatedBy       string    `gorm:"type:varchar(255);null" json:"created_by"`
	UpdatedAt       time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"updated_at"`
	UpdatedBy       string    `gorm:"type:varchar(255);null" json:"updated_by"`
}

type AdditionalCosts struct {
	ID        string    `gorm:"type:varchar(255);primaryKey" json:"id"`
	OrderID   string    `gorm:"type:varchar(255);not null;index:fk1_additional_cost_order_main_id" json:"order_id"`
	Type      string    `gorm:"type:varchar(255);not null" json:"type"`
	Cost      float64   `gorm:"type:decimal(10,2);not null" json:"cost"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(255);null" json:"created_by"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"updated_at"`
	UpdatedBy string    `gorm:"type:varchar(255);null" json:"updated_by"`

	OrderMain OrderMains `gorm:"foreignKey:OrderID;references:ID" json:"-"`
}

type Discounts struct {
	ID        string    `gorm:"type:varchar(255);primaryKey" json:"id"`
	OrderID   string    `gorm:"type:varchar(255);not null;index:fk1_discount_order_main_id" json:"order_id"`
	Type      string    `gorm:"type:varchar(255);not null" json:"type"`
	Cost      float64   `gorm:"type:decimal(10,2);not null" json:"cost"`
	CreatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"created_at"`
	CreatedBy string    `gorm:"type:varchar(255);null" json:"created_by"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;default:current_timestamp" json:"updated_at"`
	UpdatedBy string    `gorm:"type:varchar(255);null" json:"updated_by"`

	OrderMain OrderMains `gorm:"foreignKey:OrderID;references:ID" json:"-"`
}
