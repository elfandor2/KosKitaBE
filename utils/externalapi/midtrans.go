package externalapi

import (
	"KosKita/app/config"
	"KosKita/features/booking"
	"KosKita/features/order"
	"errors"
	"fmt"

	mid "github.com/midtrans/midtrans-go"

	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransInterface interface {
	NewOrderPayment(book booking.BookingCore) (*booking.BookingCore, error)
	CancelOrderPayment(bookingId string) error
	NewOrderPaymentOrder(order order.OrderCore) (*order.OrderCore, error)
	CancelOrderPaymentOrder(orderId string) error
}

type midtrans struct {
	client      coreapi.Client
	environment mid.EnvironmentType
}

func New() MidtransInterface {
	environment := mid.Sandbox
	var client coreapi.Client
	client.New(config.MID_KEY, environment)

	return &midtrans{
		client: client,
	}
}

// NewOrderPayment implements Midtrans.
func (pay *midtrans) NewOrderPayment(book booking.BookingCore) (*booking.BookingCore, error) {
	req := new(coreapi.ChargeReq)
	req.TransactionDetails = mid.TransactionDetails{
		OrderID:  book.Code,
		GrossAmt: int64(book.Total),
	}

	switch book.Bank {
	case "bca":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankBca,
		}
	case "bni":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankBni,
		}
	case "bri":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankBri,
		}
	case "permata":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankPermata,
		}
	default:
		return nil, errors.New("unsupported payment")
	}

	res, err := pay.client.ChargeTransaction(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != "201" {
		return nil, errors.New(res.StatusMessage)
	}

	if len(res.VaNumbers) == 1 {
		book.VirtualNumber = res.VaNumbers[0].VANumber
	}

	if res.PermataVaNumber != "" {
		book.VirtualNumber = res.PermataVaNumber
	}

	if res.PaymentType != "" {
		book.Method = res.PaymentType
	}

	if res.TransactionStatus != "" {
		book.Status = res.TransactionStatus
	}

	if res.ExpiryTime != "" {
		book.Status = res.TransactionStatus
	}

	// if expiredAt, err := time.Parse("2006-01-02 15:04:05", res.ExpiryTime); err != nil {
	// 	return nil, err
	// } else {
	// 	book.ExpiredAt = expiredAt
	// }

	book.BookingTotal = book.Total

	return &book, nil
}

func (pay *midtrans) CancelOrderPayment(bookingId string) error {
	res, _ := pay.client.CancelTransaction(bookingId)
	if res.StatusCode != "200" && res.StatusCode != "412" {
		return errors.New(res.StatusMessage)
	}

	return nil
}

// NewOrderPayment implements Midtrans.
func (pay *midtrans) NewOrderPaymentOrder(order order.OrderCore) (*order.OrderCore, error) {
	req := new(coreapi.ChargeReq)
	req.TransactionDetails = mid.TransactionDetails{
		OrderID:  order.ID,
		GrossAmt: int64(order.Total),
	}

	switch order.Bank {
	case "bca":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankBca,
		}
	case "bni":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankBni,
		}
	case "bri":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankBri,
		}
	case "permata":
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: mid.BankPermata,
		}
	default:
		return nil, errors.New("unsupported payment")
	}

	res, err := pay.client.ChargeTransaction(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != "201" {
		return nil, errors.New(res.StatusMessage)
	}

	if len(res.VaNumbers) == 1 {
		order.VirtualNumber = res.VaNumbers[0].VANumber
	}

	if res.PermataVaNumber != "" {
		order.VirtualNumber = res.PermataVaNumber
	}

	if res.PaymentType != "" {
		order.PaymentType = res.PaymentType
	}

	if res.TransactionStatus != "" {
		order.Status = res.TransactionStatus
	}

	// if expiredAt, err := time.Parse("2006-01-02 15:04:05", res.ExpiryTime); err != nil {
	// 	return nil, err
	// } else {
	// 	order.ExpiredAt = expiredAt
	// }
	fmt.Println("expired time from Midtrans", res.ExpiryTime)
	if res.ExpiryTime != "" {
		order.ExpiredAt = res.ExpiryTime
	}

	return &order, nil
}

func (pay *midtrans) CancelOrderPaymentOrder(orderId string) error {
	res, _ := pay.client.CancelTransaction(orderId)
	if res.StatusCode != "200" && res.StatusCode != "412" {
		return errors.New(res.StatusMessage)
	}
	fmt.Println("transaction time from Midtrans", res.TransactionTime)

	return nil
}
