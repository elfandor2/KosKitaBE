package kos

import (
	"KosKita/features/user"
	"time"
)

type Core struct {
	ID              uint
	Name            string `validate:"required"`
	Description     string
	Category        string
	Price           int    `validate:"required"`
	Rooms           int    `validate:"required"`
	Address         string `validate:"required"`
	Longitude       string
	Latitude        string
	PhotoMain       string
	PhotoFront      string
	PhotoBack       string
	PhotoRoomFront  string
	PhotoRoomInside string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UserID          uint
	Ratings         []RatingCore
	KosFacilities   []KosFacilityCore
	KosRules        []KosRuleCore
	User            user.Core
}


type CoreInput struct {
	ID              uint
	Name            string `validate:"required"`
	Description     string
	Category        string
	Price           int    `validate:"required"`
	Rooms           int    `validate:"required"`
	Address         string `validate:"required"`
	Longitude       string
	Latitude        string
	PhotoMain       string
	PhotoFront      string
	PhotoBack       string
	PhotoRoomFront  string
	PhotoRoomInside string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	UserID          uint
	Ratings         []RatingCore
	KosFacilities   []KosFacilityCore
	KosRules        []KosRuleCore
}

type CoreFoto struct {
	PhotoMain       string
	PhotoFront      string
	PhotoBack       string
	PhotoRoomFront  string
	PhotoRoomInside string
	UserID          uint
}

type KosFacilityCore struct {
	ID              uint
	Facility        string
	BoardingHouseID uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type KosRuleCore struct {
	ID              uint
	Rule            string
	BoardingHouseID uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type RatingCore struct {
	ID              uint
	Score           int
	UserID          uint
	BoardingHouseID uint
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type KosDataInterface interface {
	Insert(userIdLogin int, input CoreInput) (uint, error)
	Update(userIdLogin int, input CoreInput) error
	CekRating(userId, kosId int) (*RatingCore, error)
	InsertRating(userIdLogin, kosId int, score RatingCore) error
	SelectByRating() ([]Core, error)
	Delete(userIdLogin, kosId int) error
	SelectById(kosId int) (*Core, error)
	SelectByUserId(userIdLogin int) ([]Core, error)
	SearchKos(query, category string, minPrice, maxPrice int) ([]Core, error)
	InsertImage(userIdLogin, kosId int, input CoreFoto) error
	GetTotalKos() (int, error) 
}

// interface untuk Service Layer
type KosServiceInterface interface {
	Create(userIdLogin int, input CoreInput) (uint, error)
	Put(userIdLogin int, input CoreInput) error
	CreateRating(userIdLogin, kosId int, score RatingCore) error
	GetByRating() ([]Core, error)
	Delete(userIdLogin, kosId int) error
	GetById(kosId int) (*Core, error)
	GetByUserId(userIdLogin int) ([]Core, error)
	SearchKos(query, category string, minPrice, maxPrice int) ([]Core, error)
	CreateImage(userIdLogin, kosId int, input CoreFoto) error
}
