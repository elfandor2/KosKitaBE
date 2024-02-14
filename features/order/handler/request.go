package handler

import (
	"KosKita/features/order"

	"github.com/google/uuid"
)

type OrderRequest struct {
	ID              string
	BoardingHouseId uint   `json:"kos_id"`
	Method          string `json:"payment_type" form:"payment_type"`
	Bank            string `json:"bank" form:"bank"`
	StartDate       string `json:"start_date" form:"start_date"`
}

type CancelOrderRequest struct {
	Status string `json:"status"`
}

func RequestToCoreOrder(input OrderRequest) order.OrderCore {
	return order.OrderCore{
		ID:              uuid.New().String(),
		BoardingHouseId: input.BoardingHouseId,
		StartDate:       input.StartDate,
		PaymentType:     input.Method,
		Bank:            input.Bank,
	}
}

func CancelRequestToCoreOrder(input CancelOrderRequest) order.OrderCore {
	return order.OrderCore{
		Status: input.Status,
	}
}
