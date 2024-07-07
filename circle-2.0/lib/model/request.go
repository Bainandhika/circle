package model

// user request
type (
	CreateUserRequest struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,ascii" validation:"PasswordValidation"`
	}

	LoginUserRequest struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	UpdateUserRequest struct {
		Name     string `json:"name"`
		Email    string `json:"email" binding:"email"`
		Password string `json:"password" binding:"ascii" validation:"PasswordValidation"`
	}
)

// create order request
type (
	OrderItem struct {
		Item         string  `json:"item"`
		Quantity     int     `json:"quantity"`
		PricePerItem float64 `json:"price-per-item"`
	}

	OrderDetailUser struct {
		Name       string      `json:"name" binding:"required"`
		OrderItems []OrderItem `json:"order-items"`
	}

	OrderAdditionalCost struct {
		Type string  `json:"type"`
		Cost float64 `json:"cost"`
	}

	OrderDiscount struct {
		Type string  `json:"type"`
		Cost float64 `json:"cost"`
	}

	OrderAdditionalAndDiscounts struct {
		OrderAdditionalCosts []OrderAdditionalCost `json:"additional"`
		OrderDiscounts       []OrderDiscount       `json:"discounts"`
	}

	OrderRequest struct {
		UserID                      string                      `json:"user-id" binding:"required"`
		OrderTitle                  string                      `json:"order-title" binding:"required"`
		OrderDetailUsers            []OrderDetailUser           `json:"order-detail-users" binding:"required"`
		OrderAdditionalAndDiscounts OrderAdditionalAndDiscounts `json:"order-additional-and-discounts" binding:"required"`
	}
)

type (
	GetBillUserRequest struct {
		UserID  string `json:"user-id" binding:"required"`
		OrderID string `json:"order-id" binding:"required"`
	}
)

type (
	UpdateOrderUserItemRequest struct {
		ID           string  `json:"id"`
		Item         string  `json:"item"`
		Quantity     int8    `json:"quantity"`
		PricePerItem float64 `json:"price_per_item"`
	}

	UpdateOrderUserRequest struct {
		ID         string  `json:"id"`
		TotalPrice float64 `json:"total_price"`
		OrderPart  float64 `json:"order_part"`
		PriceToPay float64 `json:"price_to_pay"`
	}

	UpdateOrderAdditionalCostRequest struct {
		ID   string  `json:"id"`
		Type string  `json:"type"`
		Cost float64 `json:"cost"`
	}

	UpdateOrderDiscountRequest struct {
		ID   string  `json:"id"`
		Type string  `json:"type"`
		Cost float64 `json:"cost"`
	}

	UpdateOrderMainRequest struct {
		OrderTitle      string  `json:"order_title"`
		Total           float64 `json:"total"`
		AdditionalTotal float64 `json:"additional_total"`
		DiscountTotal   float64 `json:"discount_total"`
		TotalPayment    float64 `json:"total_payment"`
	}

	UpdateOrderRequest struct {
		OrderUserItems       []UpdateOrderUserItemRequest       `json:"order_user_items"`
		OrderAdditionalCosts []UpdateOrderAdditionalCostRequest `json:"order_additional_costs"`
		OrderDiscounts       []UpdateOrderDiscountRequest       `json:"order_discounts"`
		NewOrderTitle        string                             `json:"new_order_title"`
		UpdatedBy            string                             `json:"updated_by"`
	}
)
