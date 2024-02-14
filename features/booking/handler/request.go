package handler

import (
	"KosKita/features/booking"

	"github.com/google/uuid"
)

type BookRequest struct {
	BoardingHouseId uint   `json:"kos_id"`
	Method          string `json:"payment_type" form:"payment_type"`
	Bank            string `json:"bank" form:"bank"`
}

type CancelBookingRequest struct {
	Status string `json:"status"`
}

func RequestToCoreBook(input BookRequest, userIdLogin uint) booking.BookingCore {
	return booking.BookingCore{
		Code:            uuid.New().String(),
		UserId:          userIdLogin,
		BoardingHouseId: input.BoardingHouseId,
		Method:          input.Method,
		Bank:            input.Bank,
		// Payment: booking.PaymentCore{
		// 	Method: input.Method,
		// 	Bank:   input.Bank,
		// },
	}
}

func CancelRequestToCoreBooking(input CancelBookingRequest) booking.BookingCore {
	return booking.BookingCore{
		Status: input.Status,
	}
}

type WebhoocksRequest struct {
	Code   string `json:"order_id"`
	Status string `json:"transaction_status"`
}

func WebhoocksRequestToCore(input WebhoocksRequest) booking.BookingCore {
	return booking.BookingCore{
		Code:   input.Code,
		Status: input.Status,
	}
}
