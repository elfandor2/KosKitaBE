package data

import (
	"KosKita/features/kos"
	"KosKita/features/user/data"

	"gorm.io/gorm"
)

type BoardingHouse struct {
	gorm.Model
	Name            string
	Description     string
	Category        string 
	Price           int
	Rooms           int
	Address         string
	Longitude       string
	Latitude        string
	PhotoMain       string
	PhotoFront      string
	PhotoBack       string
	PhotoRoomFront  string
	PhotoRoomInside string
	UserID          uint
	User            data.User
	Ratings         []Rating
	KosFacilities   []KosFacility `gorm:"foreignKey:BoardingHouseID"`
	KosRules        []KosRule `gorm:"foreignKey:BoardingHouseID"`
}

type KosFacility struct {
	gorm.Model
	Facility        string
	BoardingHouseID uint
}

type KosRule struct {
	gorm.Model
	Rule            string
	BoardingHouseID uint
}

type Rating struct {
	gorm.Model
	Score           int
	UserID          uint
	User            data.User
	BoardingHouseID uint
	BoardingHouse   BoardingHouse
}

func CoreToModel(input kos.CoreInput) BoardingHouse {
	return BoardingHouse{
		UserID:          input.UserID,
		Name:            input.Name,
		Description:     input.Description,
		Category:        input.Category,
		Price:           input.Price,
		Rooms:           input.Rooms,
		Address:         input.Address,
		Longitude:       input.Longitude,
		Latitude:        input.Latitude,
	}
}

func CoreToModelPut(input kos.Core) BoardingHouse {
	return BoardingHouse{
		UserID:          input.UserID,
		Name:            input.Name,
		Description:     input.Description,
		Category:        input.Category,
		Price:           input.Price,
		Rooms:           input.Rooms,
		Address:         input.Address,
		Longitude:       input.Longitude,
		Latitude:        input.Latitude,
	}
}

func CoreToModelFoto(input kos.CoreFoto) BoardingHouse {
	return BoardingHouse{
		UserID:          input.UserID,
		PhotoMain:       input.PhotoMain,
		PhotoFront:      input.PhotoFront,
		PhotoBack:       input.PhotoBack,
		PhotoRoomFront:  input.PhotoRoomFront,
		PhotoRoomInside: input.PhotoRoomInside,
	}
}

func CoreToModelFacility(facility kos.KosFacilityCore) KosFacility {
	return KosFacility{
		Model:           gorm.Model{ID: facility.ID},
		Facility:        facility.Facility,
		BoardingHouseID: facility.BoardingHouseID,
	}
}

func CoreToModelRule(rule kos.KosRuleCore) KosRule {
	return KosRule{
		Model:           gorm.Model{ID: rule.ID},
		Rule:        rule.Rule,
		BoardingHouseID: rule.BoardingHouseID,
	}
}

func (bh BoardingHouse) ModelToCoreKos() kos.Core {
	var ratings []kos.RatingCore
	for _, r := range bh.Ratings {
		ratings = append(ratings, r.ModelToCoreRating())
	}

	var kosFacilities []kos.KosFacilityCore
	for _, f := range bh.KosFacilities {
		kosFacilities = append(kosFacilities, f.ModelToCoreFacility())
	}

	var kosRules []kos.KosRuleCore
	for _, r := range bh.KosRules {
		kosRules = append(kosRules, r.ModelToCoreRule())
	}

	return kos.Core{
		UserID:          bh.UserID,
		ID:              bh.ID,
		Name:            bh.Name,
		Description:     bh.Description,
		Category:        bh.Category,
		Price:           bh.Price,
		Rooms:           bh.Rooms,
		Address:         bh.Address,
		Longitude:       bh.Longitude,
		Latitude:        bh.Latitude,
		KosFacilities:   kosFacilities,
		KosRules:        kosRules,
		PhotoMain:       bh.PhotoMain,
		PhotoFront:      bh.PhotoFront,
		PhotoBack:       bh.PhotoBack,
		PhotoRoomFront:  bh.PhotoRoomFront,
		PhotoRoomInside: bh.PhotoRoomInside,
		CreatedAt:       bh.CreatedAt,
		UpdatedAt:       bh.UpdatedAt,
		Ratings:         ratings,
		// User:            bh.User.ModelToCore(),
		// User: kos.User{
		// 	ID:           bh.User.ID,
		// 	Name:         bh.User.Name,
		// 	UserName:     bh.User.UserName,
		// 	Email:        bh.User.Email,
		// 	Password:     bh.User.Password,
		// 	Gender:       bh.User.Gender,
		// 	Role:         bh.User.Role,
		// 	PhotoProfile: bh.User.PhotoProfile,
		// 	CreatedAt:    bh.User.CreatedAt,
		// 	UpdatedAt:    bh.User.UpdatedAt,
		// },
	}
}


func CoreToModelRating(input kos.RatingCore) Rating {
	return Rating{
		Score:           input.Score,
		UserID:          input.UserID,
		BoardingHouseID: input.BoardingHouseID,
	}
}

func (r Rating) ModelToCoreRating() kos.RatingCore {
	return kos.RatingCore{
		ID:              r.ID,
		Score:           r.Score,
		UserID:          r.UserID,
		BoardingHouseID: r.BoardingHouseID,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}

func (f KosFacility) ModelToCoreFacility() kos.KosFacilityCore {
	return kos.KosFacilityCore{
		ID:              f.ID,
		Facility:        f.Facility,
		BoardingHouseID: f.BoardingHouseID,
		CreatedAt:       f.CreatedAt,
		UpdatedAt:       f.UpdatedAt,
	}
}

func (r KosRule) ModelToCoreRule() kos.KosRuleCore {
	return kos.KosRuleCore{
		ID:              r.ID,
		Rule:            r.Rule,
		BoardingHouseID: r.BoardingHouseID,
		CreatedAt:       r.CreatedAt,
		UpdatedAt:       r.UpdatedAt,
	}
}
