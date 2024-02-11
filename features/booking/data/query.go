package data

import (
	"KosKita/features/booking"
	"KosKita/features/kos"
	"KosKita/features/kos/data"
	kd "KosKita/features/kos/data"
	"KosKita/utils/externalapi"
	"errors"

	"gorm.io/gorm"
)

type bookQuery struct {
	db              *gorm.DB
	paymentMidtrans externalapi.MidtransInterface
}

func New(db *gorm.DB, mid externalapi.MidtransInterface) booking.BookDataInterface {
	return &bookQuery{
		db:              db,
		paymentMidtrans: mid,
	}
}

// Insert implements booking.BookDataInterface.
func (repo *bookQuery) Insert(userIdLogin int, input booking.BookingCore) (*booking.BookingCore, error) {
	boardingHouse := kd.BoardingHouse{}
	if err := repo.db.First(&boardingHouse, input.BoardingHouseId).Error; err != nil {
		return nil, err
	}

	input.Total = float64(boardingHouse.Price)

	bookModel := CoreToModelBook(input)
	bookModel.UserId = uint(userIdLogin)
	bookModel.Total = input.Total
	bookModel.Payment.ExpiredAt = nil
	bookModel.Payment.PaidAt = nil

	if err := bookModel.GenerateCode(); err != nil {
		return nil, err
	}

	if err := repo.db.Create(&bookModel).Error; err != nil {
		return nil, err
	}
	tx := repo.db.Preload("User").Where("user_id = ?", userIdLogin).First(&bookModel)
	if tx.Error != nil {
		return nil, tx.Error
	}
	input.Code = bookModel.Code

	payment, errPay := repo.paymentMidtrans.NewOrderPayment(input)

	if errPay != nil {
		return nil, errPay
	}

	bookModel.Payment.Method = payment.Method
	bookModel.Payment.Bank = payment.Bank
	bookModel.Payment.VirtualNumber = payment.VirtualNumber
	bookModel.Payment.BillKey = payment.BillKey
	bookModel.Payment.BillCode = payment.BillCode
	bookModel.Payment.Status = payment.Status
	bookModel.Payment.ExpiredAt = &payment.ExpiredAt
	// bookModel.Payment.PaidAt = &payment.PaidAt
	bookModel.Payment.PaidAt = nil

	if err := repo.db.Save(&bookModel).Error; err != nil {
		return nil, err
	}

	if payment.Status == "settlement" {
		boardingHouse.Rooms -= 1
		if err := repo.db.Save(&boardingHouse).Error; err != nil {
			return nil, err
		}
	}

	bookCore := ModelToCoreBook(bookModel)
	if payment != nil {
		bookCore.Payment = *payment
	}

	return &bookCore, nil
}

// CancelBooking implements booking.BookDataInterface.
func (repo *bookQuery) CancelBooking(userIdLogin int, bookingId string, bookingCore booking.BookingCore) error {
	if bookingCore.Payment.Status == "cancelled" {
		repo.paymentMidtrans.CancelOrderPayment(bookingId)
	}

	booking := Booking{}
	tx := repo.db.Where("code = ? AND user_id = ?", bookingId, userIdLogin).First(&booking)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return errors.New("you do not have permission to edit this product")
		}
		return tx.Error
	}
	bookingInputGorm := CoreToModelBookCancel(bookingCore)

	tx = repo.db.Model(&booking).Updates(&bookingInputGorm)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return errors.New("error record not found")
	}
	return nil
}

// GetBooking implements booking.BookDataInterface.
func (repo *bookQuery) GetBooking(userId uint) ([]booking.BookingCore, error) {
	var bookingGorm []Booking
	tx := repo.db.Preload("BoardingHouse").Preload("User").Find(&bookingGorm, "user_id = ?", userId)
	if tx.Error != nil {
		return nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return nil, errors.New("find booking failed, row affected = 0")
	}
	var bookingCores []booking.BookingCore
	for _, v := range bookingGorm {
		bookingCores = append(bookingCores, ModelToCoreBook(v))
	}

	return bookingCores, nil
}

// WebhoocksData implements booking.BookDataInterface.
func (repo *bookQuery) WebhoocksData(webhoocksReq booking.BookingCore) error {
	bookingGorm := WebhoocksCoreToModel(webhoocksReq)
	tx := repo.db.Model(&Booking{}).Where("code = ?", webhoocksReq.Code).Updates(bookingGorm)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("error record not found ")
	}
	return nil
}

func (repo *bookQuery) GetRatingAndFacility(userId uint) ([]kos.Core, error) {
	var kosData []data.BoardingHouse
	var result []kos.Core

	tx := repo.db.Preload("Ratings").Preload("KosFacilities").Table("boarding_houses").Find(&kosData)

	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, k := range kosData {
		result = append(result, k.ModelToCoreKos())
	}
	// for _, v := range kosData {
	// 	fmt.Println(v.Ratings)

	// }
	return result, nil
}
