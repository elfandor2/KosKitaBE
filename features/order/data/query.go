package data

import (
	kd "KosKita/features/kos/data"
	"KosKita/features/order"
	"KosKita/utils/externalapi"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type orderQuery struct {
	db              *gorm.DB
	paymentMidtrans externalapi.MidtransInterface
}

func NewOrder(db *gorm.DB, mid externalapi.MidtransInterface) order.OrderDataInterface {
	return &orderQuery{
		db:              db,
		paymentMidtrans: mid,
	}
}

// PostOrder implements order.OrderDataInterface.
func (repo *orderQuery) PostOrder(userId uint, input order.OrderCore) (*order.OrderCore, error) {
	var orderGorm Order
	// var order order.OrderCore

	boardingHouse := kd.BoardingHouse{}
	if err := repo.db.First(&boardingHouse, input.BoardingHouseId).Error; err != nil {
		return nil, err
	}
	var amount = boardingHouse.Price

	input.Total = float64(amount)
	input.UserID = userId

	payment, errPay := repo.paymentMidtrans.NewOrderPaymentOrder(input)
	if errPay != nil {
		return nil, errPay
	}

	fmt.Println(payment.ExpiredAt)
	// repo.db.Transaction
	repo.db.Transaction(func(tx *gorm.DB) error {
		// Create Data Order
		orderGorm = OrderCoreToModel(input)
		orderGorm.PaymentType = payment.PaymentType
		orderGorm.Status = payment.Status
		orderGorm.VirtualNumber = payment.VirtualNumber
		orderGorm.PaidAt = payment.PaidAt
		orderGorm.ExpiredAt = payment.ExpiredAt
		orderGorm.Total = float64(amount)
		if errOrder := tx.Create(&orderGorm).Error; errOrder != nil {
			return errOrder
		}

		return nil
	})
	var orderCores = ModelToCore(orderGorm)

	return &orderCores, nil
}

// GetOrders implements order.OrderDataInterface.
func (repo *orderQuery) GetOrders(userId uint) ([]order.OrderCore, error) {
	var orderGorm []Order
	// tx := repo.db.Preload("ItemOrders").Preload("User").Find(&orderGorm, "user_id = ?", userId)
	tx := repo.db.Preload("BoardingHouse").Preload("User").Find(&orderGorm, "user_id = ?", userId)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("find order failed, row affected = 0")
	}
	var orderCores []order.OrderCore
	for _, v := range orderGorm {
		orderCores = append(orderCores, ModelToCore(v))
	}

	return orderCores, nil
}

// CancelOrder implements order.OrderDataInterface.
func (repo *orderQuery) CancelOrder(userId int, orderId string, orderCore order.OrderCore) error {
	if orderCore.Status == "cancelled" {
		repo.paymentMidtrans.CancelOrderPaymentOrder(orderId)
	}

	dataGorm := Order{
		Status: orderCore.Status,
	}
	fmt.Println("order id::", orderId)
	tx := repo.db.Model(&Order{}).Where("id = ? AND user_id = ?", orderId, userId).Updates(dataGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found ")
	}
	return nil
}

// GetOrder implements order.OrderDataInterface.
func (repo *orderQuery) GetOrder(userId uint) (*order.OrderCore, error) {
	panic("unimplemented")
}

// WebhoocksData implements order.OrderDataInterface.
func (repo *orderQuery) WebhoocksData(webhoocksReq order.OrderCore) error {
	dataGorm := WebhoocksCoreToModel(webhoocksReq)
	tx := repo.db.Model(&Order{}).Where("id = ?", webhoocksReq.ID).Updates(dataGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found ")
	}
	return nil
}
