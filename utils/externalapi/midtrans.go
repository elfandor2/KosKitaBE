package externalapi

import (
	"KosKita/app/config"
	"KosKita/features/booking"
	"errors"
	"time"

	mid "github.com/midtrans/midtrans-go"

	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransInterface interface {
	NewOrderPayment(book booking.BookingCore) (*booking.BookingCore, error)
	CancelOrderPayment(bookingId string) error
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

	if expiredAt, err := time.Parse("2006-01-02 15:04:05", res.ExpiryTime); err != nil {
		return nil, err
	} else {
		book.ExpiredAt = expiredAt
	}

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
