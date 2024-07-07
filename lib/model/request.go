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
		OrderID string `json:"order-main-id" binding:"required"`
	}
)
