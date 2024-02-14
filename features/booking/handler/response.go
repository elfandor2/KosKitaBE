package handler

import (
	"KosKita/features/booking"
	"KosKita/features/kos"
	"time"
)

type BookingResponse struct {
	Code                 string     `json:"booking_code,omitempty"`
	Status               string     `json:"status,omitempty"`
	Total                float64    `json:"total,omitempty"`
	PaymentMethod        string     `json:"payment_method,omitempty"`
	PaymentBank          string     `json:"bank,omitempty"`
	PaymentVirtualNumber string     `json:"virtual_number,omitempty"`
	PaymentBillKey       string     `json:"key_bill,omitempty"`
	PaymentBillCode      string     `json:"code_bill,omitempty"`
	PaymentExpiredAt     *time.Time `json:"payment_expired,omitempty"`
}

type BookingHistoryResponse struct {
	KosId         uint     `json:"kos_id,omitempty"`
	KosName       string   `json:"kos_name,omitempty"`
	KosFasilitas  []string `json:"kos_fasilitas,omitempty"`
	KosLokasi     string   `json:"kos_lokasi,omitempty"`
	KosRating     float64  `json:"kos_rating,omitempty"`
	KosMainFoto   string   `json:"kos_main_foto,omitempty"`
	PaymentStatus string   `json:"payment_status,omitempty"`
	TotalHarga    float64  `json:"total_harga,omitempty"`
}

func CoreToResponseBook(core *booking.BookingCore) BookingResponse {
	return BookingResponse{
		Code:                 core.Code,
		Status:               core.Status,
		Total:                core.Total,
		PaymentMethod:        core.Method,
		PaymentBank:          core.Bank,
		PaymentVirtualNumber: core.VirtualNumber,
		PaymentExpiredAt:     &core.ExpiredAt,
	}
}

func CoreToResponseBookHistory(core *booking.BookingCore) BookingHistoryResponse {
	return BookingHistoryResponse{
		KosId:         core.BoardingHouse.ID,
		KosName:       core.BoardingHouse.Name,
		KosLokasi:     core.BoardingHouse.Address,
		KosRating:     KosRatingResult(core.BoardingHouse.Ratings),
		KosMainFoto:   core.BoardingHouse.PhotoMain,
		PaymentStatus: core.Status,
		TotalHarga:    core.Total,
	}
}

func KosFasilitasList(kf []kos.KosFacilityCore) []string {
	var results []string
	for _, v := range kf {
		results = append(results, v.Facility)
	}
	return results
}

func KosRatingResult(numbers []kos.RatingCore) float64 {
	var results float64
	if len(numbers) > 0 {
		for _, num := range numbers {
			results += float64(num.Score)
		}
		return float64(results) / float64(len(numbers))
	}
	return 0
}
