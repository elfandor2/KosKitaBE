package booking

import (
	"time"

	kd "KosKita/features/kos"
	ud "KosKita/features/user"

	"gorm.io/gorm"
)

type BookingCore struct {
	Code            string
	Total           float64
	BookedAt        time.Time
	DeletedAt       gorm.DeletedAt
	UserId          uint
	User            ud.Core
	BoardingHouseId uint
	BoardingHouse   kd.Core
	Method          string
	Bank            string
	VirtualNumber   string
	BookingCode     int
	BookingTotal    float64
	Status          string
	CreatedAt       time.Time
	ExpiredAt       time.Time
	PaidAt          time.Time
	// Payment         PaymentCore
}

type MonthCount struct {
	Month int
	Count int
}

// type PaymentCore struct {
// 	Method        string
// 	Bank          string
// 	VirtualNumber string
// 	BillKey       string
// 	BillCode      string
// 	BookingCode   int
// 	BookingTotal  float64
// 	Status        string
// 	CreatedAt     time.Time
// 	ExpiredAt     time.Time
// 	PaidAt        time.Time
// }

type BookDataInterface interface {
	Insert(userIdLogin int, input BookingCore) (*BookingCore, error)
	CancelBooking(userIdLogin int, bookingId string, bookingCore BookingCore) error
	GetBooking(userId uint) ([]BookingCore, error)
	WebhoocksData(webhoocksReq BookingCore) error
	GetTotalBooking() (int, error)
	GetTotalBookingPerYear(year int) ([]int, error)
}

// interface untuk Service Layer
type BookServiceInterface interface {
	Create(userIdLogin int, input BookingCore) (*BookingCore, error)
	CancelBooking(userIdLogin int, bookingId string, bookingCore BookingCore) error
	GetBooking(userId uint) ([]BookingCore, error)
	WebhoocksData(webhoocksReq BookingCore) error
}
