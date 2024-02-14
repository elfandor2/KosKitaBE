package data

import (
	"KosKita/features/booking"
	"KosKita/features/kos"
	"KosKita/features/kos/data"
	kd "KosKita/features/kos/data"
	"KosKita/utils/externalapi"
	"errors"
	"fmt"
	"log"

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
	// bookModel.ExpiredAt = time.Time.Local()
	// bookModel.PaidAt = &input.PaidAt
	fmt.Println("book expired at", bookModel.ExpiredAt)

	if err := repo.db.Create(&bookModel).Error; err != nil {
		return nil, err
	}

	input.Code = bookModel.Code

	log.Println("input book", input)
	payment, errPay := repo.paymentMidtrans.NewOrderPayment(input)

	log.Println("input payment", payment)
	if errPay != nil {
		return nil, errPay
	}

	bookModel.Method = payment.Method
	bookModel.Bank = payment.Bank
	bookModel.VirtualNumber = payment.VirtualNumber
	bookModel.Status = payment.Status
	bookModel.ExpiredAt = payment.ExpiredAt
	bookModel.PaidAt = payment.PaidAt

	log.Println("input bookmodel", bookModel)

	if err := repo.db.Updates(&bookModel).Error; err != nil {
		return nil, err
	}

	bookCore := ModelToCoreBook(bookModel)
	// if payment != nil {
	// 	bookCore = *payment
	// }

	return &bookCore, nil
}

// CancelBooking implements booking.BookDataInterface.
func (repo *bookQuery) CancelBooking(userIdLogin int, bookingId string, bookingCore booking.BookingCore) error {
	if bookingCore.Status == "cancelled" {
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

	tx := repo.db.Model(&Booking{}).Where("code = ?", bookingGorm.Code).Updates(bookingGorm)
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
	return result, nil
}

func (repo *bookQuery) GetTotalBooking() (int, error) {
	var count int64
	tx := repo.db.Model(&Booking{}).Count(&count)
	if tx.Error != nil {
		return 0, tx.Error
	}
	return int(count), nil
}

func (repo *bookQuery) GetTotalBookingPerYear(year int) ([]int, error) {
	var counts []int
	rows, err := repo.db.Raw("SELECT COUNT(*) as count FROM bookings WHERE YEAR(created_at) = ? GROUP BY MONTH(created_at) ORDER BY MONTH(created_at)", year).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var count int
		rows.Scan(&count)
		counts = append(counts, count)
	}
	return counts, nil
}
