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
	PaymentBank          string     `json:"payment_method,omitempty"`
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
		Status:               core.Payment.Status,
		Total:                core.Total,
		PaymentBank:          core.Payment.Bank,
		PaymentVirtualNumber: core.Payment.VirtualNumber,
		PaymentBillKey:       core.Payment.BillKey,
		PaymentBillCode:      core.Payment.BillCode,
		PaymentExpiredAt:     &core.Payment.ExpiredAt,
	}
}

func CoreToResponseBookHistory(core *booking.BookingCore) BookingHistoryResponse {
	return BookingHistoryResponse{
		KosId:   core.BoardingHouse.ID,
		KosName: core.BoardingHouse.Name,
		// KosFasilitas: KosFasilitasList(core.BoardingHouse.KosFacilities),
		KosLokasi:     core.BoardingHouse.Address,
		KosRating:     KosRatingResult(core.BoardingHouse.Ratings),
		KosMainFoto:   core.BoardingHouse.PhotoMain,
		PaymentStatus: core.Payment.Status,
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
	// menjumlahkan semua angka dalam slice
	if len(numbers) > 0 {
		for _, num := range numbers {
			results += float64(num.Score)
		}
		// mengembalikan rata-rata
		return float64(results) / float64(len(numbers))
	}
	return 0
}
