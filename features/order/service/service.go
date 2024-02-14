package service

import (
	"KosKita/features/order"
	"errors"
	// _midtransService "BE-REPO-20/features/midtrans/service"
)

type orderService struct {
	orderData order.OrderDataInterface
}

func NewOrder(repo order.OrderDataInterface) order.OrderServiceInterface {
	return &orderService{
		orderData: repo,
	}
}

// PostOrder implements order.OrderServiceInterface.
func (service *orderService) PostOrder(userId uint, input order.OrderCore) (*order.OrderCore, error) {
	if userId <= 0 {
		return nil, errors.New("invalid id")
	}
	res, err := service.orderData.PostOrder(userId, input)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetOrders implements order.OrderServiceInterface.
func (service *orderService) GetOrders(userId uint) ([]order.OrderCore, error) {
	results, err := service.orderData.GetOrders(userId)
	if err != nil {
		return nil, err
	}
	return results, nil
}

// CancelOrder implements order.OrderServiceInterface.
func (service *orderService) CancelOrder(userId int, orderId string, orderCore order.OrderCore) error {
	if orderCore.Status == "" {
		orderCore.Status = "cancelled"
	}

	err := service.orderData.CancelOrder(userId, orderId, orderCore)
	return err
}

// WebhoocksService implements order.OrderServiceInterface.
func (service *orderService) WebhoocksService(webhoocksReq order.OrderCore) error {
	if webhoocksReq.ID == "" {
		return errors.New("invalid order id")
	}

	err := service.orderData.WebhoocksData(webhoocksReq)
	if err != nil {
		return err
	}

	return nil
}
