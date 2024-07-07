package model

// order response
type (
	ResultOrderUserItem struct {
		OrderUserItemID string  `json:"order-user-item-id"`
		Item            string  `json:"item"`
		Quantity        int     `json:"quantity"`
		PricePerItem    float64 `json:"price-per-item"`
		Total           float64 `json:"total"`
	}

	ResultOrderUser struct {
		OrderUserID string                `json:"order-user-id"`
		UserID      string                `json:"user-id"`
		Name        string                `json:"name"`
		OrderItems  []ResultOrderUserItem `json:"order-items"`
		Total       float64               `json:"total"`
		PartOfOrder float64               `json:"part-of-order"`
		PriceToPay  float64               `json:"price-to-pay"`
	}

	ResultAdditional struct {
		AdditionalID string  `json:"additional-id"`
		Type         string  `json:"type"`
		Cost         float64 `json:"cost"`
	}

	ResultDiscount struct {
		DiscountID string  `json:"discount-id"`
		Type       string  `json:"type"`
		Cost       float64 `json:"cost"`
	}

	ResultAdditionalAndDiscounts struct {
		Additional []ResultAdditional `json:"additional"`
		Discounts  []ResultDiscount   `json:"discounts"`
	}

	ResultOrderMain struct {
		OrderMainID            string                       `json:"order-main-id"`
		OrderTitle             string                       `json:"order-title"`
		Total                  float64                      `json:"total"`
		AdditionalAndDiscounts ResultAdditionalAndDiscounts `json:"additional-and-discounts"`
		TotalPayment           float64                      `json:"total-payment"`
	}

	Bill struct {
		CreatedBy        string            `json:"created-by"`
		ResultOrderMain  ResultOrderMain   `json:"order-main"`
		ResultOrderUsers []ResultOrderUser `json:"order-users"`
	}
)

type (
	GetBillUserResponse struct {
		ResultOrderMain ResultOrderMain `json:"order-main"`
		ResultOrderUser ResultOrderUser `json:"order-user"`
	}
)

type (
	LoginUserResponse struct {
		UserID string `json:"user-id"`
		Token  string `json:"token"`
	}
)
