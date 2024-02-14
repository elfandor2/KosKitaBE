package order

import (
	kd "KosKita/features/kos"
	ud "KosKita/features/user"
	"time"
)

type OrderCore struct {
	ID              string
	UserID          uint
	BoardingHouseId uint
	StartDate       string
	PaymentType     string
	Total           float64
	Status          string
	Bank            string
	VirtualNumber   string
	ExpiredAt       string
	PaidAt          string
	CreatedAt       time.Time
	User            ud.Core
	BoardingHouse   kd.Core
}

type OrderDataInterface interface {
	PostOrder(userId uint, input OrderCore) (*OrderCore, error)
	GetOrder(userId uint) (*OrderCore, error)
	GetOrders(userId uint) ([]OrderCore, error)
	CancelOrder(userId int, orderId string, orderCore OrderCore) error
	WebhoocksData(webhoocksReq OrderCore) error
}

type OrderServiceInterface interface {
	PostOrder(userId uint, input OrderCore) (*OrderCore, error)
	GetOrders(userId uint) ([]OrderCore, error)
	CancelOrder(userId int, orderId string, orderCore OrderCore) error
	WebhoocksService(webhoocksReq OrderCore) error
}
