package externalapi

import (
	"KosKita/app/config"
	"KosKita/features/booking"
	"errors"
	"fmt"
	"time"

	mid "github.com/midtrans/midtrans-go"

	"github.com/midtrans/midtrans-go/coreapi"
)

type MidtransInterface interface {
	NewOrderPayment(book booking.BookingCore) (*booking.PaymentCore, error)
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
func (pay *midtrans) NewOrderPayment(book booking.BookingCore) (*booking.PaymentCore, error) {
	req := new(coreapi.ChargeReq)
	req.TransactionDetails = mid.TransactionDetails{
		// OrderID:  fmt.Sprintf("%d", book.Code),
		OrderID:  book.Code,
		GrossAmt: int64(book.Total),
	}

	switch book.Payment.Bank {
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
	case "mandiri":
		req.PaymentType = coreapi.PaymentTypeEChannel
		req.EChannel = &coreapi.EChannelDetail{
			BillInfo1: "KosKita Booking",
			BillInfo2: fmt.Sprintf("%d BookedAt", len(book.BookedAt.Format(time.RFC3339))),
			BillKey:   fmt.Sprintf("%d", book.Code),
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

	if res.BillKey != "" {
		book.Payment.BillKey = res.BillKey
	}

	if res.BillerCode != "" {
		book.Payment.BillCode = res.BillerCode
	}

	if len(res.VaNumbers) == 1 {
		book.Payment.VirtualNumber = res.VaNumbers[0].VANumber
	}

	if res.PermataVaNumber != "" {
		book.Payment.VirtualNumber = res.PermataVaNumber
	}

	if res.PaymentType != "" {
		book.Payment.Method = res.PaymentType
	}

	if res.TransactionStatus != "" {
		book.Payment.Status = res.TransactionStatus
	}

	if expiredAt, err := time.Parse("2006-01-02 15:04:05", res.ExpiryTime); err != nil {
		return nil, err
	} else {
		book.Payment.ExpiredAt = expiredAt
	}

	book.Payment.BookingTotal = book.Total

	return &book.Payment, nil
}

func (pay *midtrans) CancelOrderPayment(bookingId string) error {
	res, _ := pay.client.CancelTransaction(bookingId)
	if res.StatusCode != "200" && res.StatusCode != "412" {
		return errors.New(res.StatusMessage)
	}

	return nil
}
