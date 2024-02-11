package data

import (
	"KosKita/features/kos"

	"gorm.io/gorm"
)

type kosQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) kos.KosDataInterface {
	return &kosQuery{
		db: db,
	}
}

// Insert implements kos.KosDataInterface.
func (repo *kosQuery) Insert(userIdLogin int, input kos.CoreInput) (uint, error) {
	kosInput := CoreToModel(input)
	kosInput.UserID = uint(userIdLogin)

	tx := repo.db.Create(&kosInput)
	if tx.Error != nil {
		return 0, tx.Error
	}

	for _, facility := range input.KosFacilities {
		facilityModel := KosFacility{
			Facility:        facility.Facility,
			BoardingHouseID: kosInput.ID,
		}
		tx = repo.db.Create(&facilityModel)
		if tx.Error != nil {
			return 0, tx.Error
		}
	}

	for _, rule := range input.KosRules {
		ruleModel := KosRule{
			Rule:            rule.Rule,
			BoardingHouseID: kosInput.ID,
		}
		tx = repo.db.Create(&ruleModel)
		if tx.Error != nil {
			return 0, tx.Error
		}
	}

	return kosInput.ID, nil
}

// InsertImage implements kos.KosDataInterface.
func (repo *kosQuery) InsertImage(userIdLogin int, kosId int, input kos.CoreFoto) error {
	kos := BoardingHouse{}
	tx := repo.db.Where("id = ? AND user_id = ?", kosId, userIdLogin).First(&kos)
	if tx.Error != nil {
		return tx.Error
	}

	kosInput := CoreToModelFoto(input)

	tx = repo.db.Model(&kos).Updates(&kosInput)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// Update implements kos.KosDataInterface.
func (repo *kosQuery) Update(userIdLogin int, input kos.Core) error {
	kos := BoardingHouse{}
	tx := repo.db.Where("id = ? AND user_id = ?", input.ID, userIdLogin).First(&kos)
	if tx.Error != nil {
		return tx.Error
	}

	kosInput := CoreToModelPut(input)

	tx = repo.db.Model(&kos).Updates(&kosInput)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

func (repo *kosQuery) CekRating(userId, kosId int) (*kos.RatingCore, error) {
	var ratingData Rating

	tx := repo.db.Where("user_id = ? AND boarding_house_id = ?", userId, kosId).First(&ratingData)

	if tx.Error != nil {
		return nil, tx.Error
	}

	ratingCore := ratingData.ModelToCoreRating()
	return &ratingCore, nil
}

// InsertRating implements kos.KosDataInterface.
func (repo *kosQuery) InsertRating(userIdLogin int, kosId int, score kos.RatingCore) error {
	ratingInput := CoreToModelRating(score)
	ratingInput.UserID = uint(userIdLogin)
	ratingInput.BoardingHouseID = uint(kosId)

	tx := repo.db.Create(&ratingInput)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}

// SelectByRating implements kos.KosDataInterface.
func (repo *kosQuery) SelectByRating() ([]kos.Core, error) {
	var kosData []BoardingHouse
	var result []kos.Core

	tx := repo.db.Preload("User").Preload("Ratings").Preload("KosFacilities").Preload("KosRules").Table("boarding_houses").
		Joins("left join ratings on boarding_houses.id = ratings.boarding_house_id").
		Group("boarding_houses.id").
		Select("boarding_houses.*, AVG(ratings.score) as average_rating").
		Order("average_rating desc").
		Find(&kosData)

	if tx.Error != nil {
		return nil, tx.Error
	}

	for _, k := range kosData {
		result = append(result, k.ModelToCoreKos())
	}

	return result, nil
}

// Delete implements kos.KosDataInterface.
func (repo *kosQuery) Delete(userIdLogin int, kosId int) error {
	tx := repo.db.Where("id = ? AND user_id = ?", kosId, userIdLogin).Delete(&BoardingHouse{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// SelectById implements kos.KosDataInterface.
func (repo *kosQuery) SelectById(kosId int) (*kos.Core, error) {
	var kosData BoardingHouse
	tx := repo.db.Preload("User").Preload("Ratings").Preload("KosFacilities").Preload("KosRules").Where("id = ?", kosId).First(&kosData)
	if tx.Error != nil {
		return nil, tx.Error
	}

	result := kosData.ModelToCoreKos()
	return &result, nil
}

// SelectByUserId implements kos.KosDataInterface.
func (repo *kosQuery) SelectByUserId(userIdLogin int) ([]kos.Core, error) {
	var kosData []BoardingHouse
	tx := repo.db.Preload("User").Preload("Ratings").Preload("KosFacilities").Preload("KosRules").Where("user_id = ?", userIdLogin).Find(&kosData)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var result []kos.Core
	for _, k := range kosData {
		result = append(result, k.ModelToCoreKos())
	}

	return result, nil
}

// SearchKos implements kos.KosDataInterface.
func (repo *kosQuery) SearchKos(query, category string, minPrice, maxPrice int) ([]kos.Core, error) {
	var kosDataGorms []BoardingHouse
	tx := repo.db.Preload("User").Preload("Ratings").Preload("KosFacilities").Preload("KosRules")

	if query != "" {
		tx = tx.Where("address LIKE ?", "%"+query+"%")
	}
	if category != "" {
		tx = tx.Where("category = ?", category)
	}
	if minPrice != 0 {
		tx = tx.Where("price >= ?", minPrice)
	}
	if maxPrice != 0 {
		tx = tx.Where("price <= ?", maxPrice)
	}

	tx = tx.Find(&kosDataGorms)
	if tx.Error != nil {
		return nil, tx.Error
	}

	var result []kos.Core
	for _, k := range kosDataGorms {
		result = append(result, k.ModelToCoreKos())
	}

	return result, nil
}
