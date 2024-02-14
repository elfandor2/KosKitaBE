package handler

import (
	"KosKita/features/order"
	"time"
)

type OrderResponse struct {
	Id            string  `json:"order_id"`
	UserID        uint    `json:"user_id"`
	StartDate     string  `json:"start_date"`
	PaymentType   string  `json:"payment_method"`
	Total         float64 `json:"total"`
	Status        string  `json:"status"`
	Bank          string  `json:"bank"`
	VirtualNumber string  `json:"virtual_number"`
	ExpiredAt     string  `json:"expired_at"`
}

type OrderHistoryResponse struct {
	Id            string    `json:"order_id"`
	KosId         uint      `json:"kos_id,omitempty"`
	KosName       string    `json:"kos_name,omitempty"`
	KosFasilitas  []string  `json:"kos_fasilitas,omitempty"`
	KosLokasi     string    `json:"kos_lokasi,omitempty"`
	KosRating     float64   `json:"kos_rating,omitempty"`
	StartDate     string    `json:"start_date,omitempty"`
	KosMainFoto   string    `json:"kos_main_foto,omitempty"`
	PaymentStatus string    `json:"payment_status,omitempty"`
	TotalHarga    float64   `json:"total_harga,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	PaidAt        string    `json:"paid_at,omitempty"`
}

func CoreToResponse(o order.OrderCore) OrderResponse {
	return OrderResponse{
		Id:            o.ID,
		UserID:        o.UserID,
		StartDate:     o.StartDate,
		PaymentType:   o.PaymentType,
		Total:         o.Total,
		Status:        o.Status,
		Bank:          o.Bank,
		VirtualNumber: o.VirtualNumber,
		ExpiredAt:     o.ExpiredAt,
	}
}

func CoreToResponseOrderHistory(core *order.OrderCore) OrderHistoryResponse {
	return OrderHistoryResponse{
		Id:        core.ID,
		KosId:     core.BoardingHouse.ID,
		KosName:   core.BoardingHouse.Name,
		KosLokasi: core.BoardingHouse.Address,
		StartDate: core.StartDate,
		// KosRating:     KosRatingResult(core.BoardingHouse.Ratings),
		KosMainFoto:   core.BoardingHouse.PhotoMain,
		PaymentStatus: core.Status,
		TotalHarga:    core.Total,
		CreatedAt:     core.CreatedAt,
		PaidAt:        core.PaidAt,
	}
}
