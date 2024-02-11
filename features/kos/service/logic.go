package service

import (
	"KosKita/features/kos"
	"KosKita/features/user"
	"errors"

	"github.com/go-playground/validator/v10"
)

type kosService struct {
	kosData     kos.KosDataInterface
	userService user.UserServiceInterface
	validate    *validator.Validate
}

func New(repo kos.KosDataInterface, us user.UserServiceInterface) kos.KosServiceInterface {
	return &kosService{
		kosData:     repo,
		validate:    validator.New(),
		userService: us,
	}
}

// Create implements kos.KosServiceInterface.
func (ks *kosService) Create(userIdLogin int, input kos.CoreInput) (uint, error) {
	user, err := ks.userService.GetById(userIdLogin)
	if err != nil {
		return 0, err
	}

	if user.Role != "owner" {
		return 0, errors.New("anda bukan owner")
	}

	if input.Name == "" {
		return 0, errors.New("name anda kosong")
	}

	if input.Category == "" {
		return 0, errors.New("category anda kosong")
	}

	if input.Price <= 0 {
		return 0, errors.New("harga anda kosong")
	}

	if input.Rooms <= 0 {
		return 0, errors.New("rooms anda kosong")
	}

	if input.Address == "" {
		return 0, errors.New("alamat anda kosong")
	}

	errValidate := ks.validate.Struct(input)
	if errValidate != nil {
		return 0, errValidate
	}

	kosId, errInsert := ks.kosData.Insert(userIdLogin, input)
	if errInsert != nil {
		return 0, errInsert
	}

	return kosId, nil
}

// Put implements kos.KosServiceInterface.
func (ks *kosService) Put(userIdLogin int, input kos.Core) error {
	err := ks.kosData.Update(userIdLogin, input)
	if err != nil {
		return err
	}
	return nil
}

// CreateRating implements kos.KosServiceInterface.
func (ks *kosService) CreateRating(userIdLogin int, kosId int, input kos.RatingCore) error {
	if input.Score < 1 || input.Score > 5 {
		return errors.New("skor hanya bisa 1 sampai 5")
	}

	existingRating, err := ks.kosData.CekRating(userIdLogin, kosId)
	if err != nil {
		return err
	}

	if existingRating != nil {
		return errors.New("anda sudah pernah memberikan rating untuk kos ini")
	}

	errInsert := ks.kosData.InsertRating(userIdLogin, kosId, input)
	if errInsert != nil {
		return errInsert
	}

	return nil
}

// GetByRating implements kos.KosServiceInterface.
func (ks *kosService) GetByRating() ([]kos.Core, error) {
	result, err := ks.kosData.SelectByRating()
	return result, err
}

// Delete implements kos.KosServiceInterface.
func (ks *kosService) Delete(userIdLogin int, kosId int) error {
	if userIdLogin <= 0 {
		return errors.New("perlu login dulu")
	}
	if kosId <= 0 {
		return errors.New("kos id tidak valid")
	}

	kos, err := ks.kosData.SelectById(kosId)
	if err != nil {
		if err.Error() == "record not found" {
			return errors.New("kos id tidak ada")
		}
		return err
	}

	if kos.User.ID != uint(userIdLogin) {
		return errors.New("kos ini bukan milik Anda")
	}

	err = ks.kosData.Delete(userIdLogin, kosId)
	if err != nil {
		return err
	}

	return nil
}

// GetById implements kos.KosServiceInterface.
func (ks *kosService) GetById(kosId int) (*kos.Core, error) {
	result, err := ks.kosData.SelectById(kosId)
	return result, err
}

// GetByUserId implements kos.KosServiceInterface.
func (ks *kosService) GetByUserId(userIdLogin int) ([]kos.Core, error) {
	kos, err := ks.kosData.SelectByUserId(userIdLogin)
	if err != nil {
		return nil, err
	}
	return kos, nil
}

// SearchKos implements kos.KosServiceInterface.
func (ks *kosService) SearchKos(query string, category string, minPrice int, maxPrice int) ([]kos.Core, error) {
	kos, err := ks.kosData.SearchKos(query, category, minPrice, maxPrice)
	if err != nil {
		return nil, err
	}
	if len(kos) == 0 {
		return nil, errors.New("tidak ada kos yang ditemukan dengan filter yang dipilih")
	}
	return kos, nil
}

// CreateImage implements kos.KosServiceInterface.
func (ks *kosService) CreateImage(userIdLogin int, kosId int, input kos.CoreFoto) error {
	user, err := ks.userService.GetById(userIdLogin)
	if err != nil {
		return err
	}

	if user.Role != "owner" {
		return errors.New("anda bukan owner")
	}

	errInsert := ks.kosData.InsertImage(userIdLogin, kosId, input)
	if errInsert != nil {
		return errInsert
	}

	return nil
}
