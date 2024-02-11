package handler

import (
	"KosKita/features/kos"
	"KosKita/features/user/handler"
)

type KosFacilityResponse struct {
	ID       uint   `json:"id"`
	Facility string `json:"facility"`
}

type KosRuleResponse struct {
	ID   uint   `json:"id"`
	Rule string `json:"rule"`
}

type KosResponseRating struct {
	ID            uint                  `json:"id"`
	Name          string                `json:"kos_name"`
	Rating        float64               `json:"rating"`
	Category      string                `json:"category"`
	Price         int                   `json:"price"`
	Rooms         int                   `json:"rooms"`
	Address       string                `json:"address"`
	KosFacilities []KosFacilityResponse `json:"kos_facilities"`
	PhotoKos      PhotoMainResponse     `json:"photo_kos"`
}

type KosResponseUser struct {
	ID            uint                  `json:"id"`
	Name          string                `json:"kos_name"`
	Rating        float64               `json:"rating"`
	Price         int                   `json:"price"`
	Rooms         int                   `json:"rooms"`
	Address       string                `json:"address"`
	KosFacilities []KosFacilityResponse `json:"kos_facilities"`
	PhotoKos      PhotoMainResponse     `json:"photo_kos"`
}

type PhotoMainResponse struct {
	PhotoMain string `json:"main_kos_photo"`
}

type KosIdResponse struct {
	ID string `json:"kos_id"`
}

type KosResponseDetail struct {
	ID            uint                          `json:"id"`
	Name          string                         `json:"kos_name"`
	Description   string                        `json:"description"`
	Rooms         int                           `json:"stock"`
	Rating        float64                       `json:"rating"`
	Category      string                        `json:"category"`
	Price         int                           `json:"price"`
	Address       string                        `json:"address"`
	Longitude     string                        `json:"longitude"`
	Latitude      string                        `json:"latitude"`
	KosFacilities []KosFacilityResponse         `json:"kos_facilities"`
	KosRules      []KosRuleResponse             `json:"kos_rules"`
	PhotoKos      PhotoDetailResponse           `json:"photo_kos"`
	User          handler.UserKosDetailResponse `json:"user"`
}

type PhotoDetailResponse struct {
	PhotoMain       string `json:"main_kos_photo"`
	PhotoFront      string `json:"front_kos_photo"`
	PhotoBack       string `json:"back_kos_photo"`
	PhotoRoomFront  string `json:"front_room_photo"`
	PhotoRoomInside string `json:"inside_room_photo"`
}

func CoreToGetRating(kos kos.Core) KosResponseRating {
	var totalRating float64
	for _, rating := range kos.Ratings {
		totalRating += float64(rating.Score)
	}

	var averageRating float64
	if len(kos.Ratings) > 0 {
		averageRating = totalRating / float64(len(kos.Ratings))
	}

	var kosFacilities []KosFacilityResponse
	for _, f := range kos.KosFacilities {
		kosFacilities = append(kosFacilities, KosFacilityResponse{
			ID:       f.ID,
			Facility: f.Facility,
		})
	}

	return KosResponseRating{
		ID:            kos.ID,
		Name:          kos.Name,
		Rating:        averageRating,
		Category:      kos.Category,
		Rooms:         kos.Rooms,
		Price:         kos.Price,
		Address:       kos.Address,
		KosFacilities: kosFacilities,
		PhotoKos:      PhotoMainResponse{PhotoMain: kos.PhotoMain},
	}
}

func CoreToGetDetail(kos kos.Core) KosResponseDetail {
	var totalRating float64
	for _, rating := range kos.Ratings {
		totalRating += float64(rating.Score)
	}

	var averageRating float64
	if len(kos.Ratings) > 0 {
		averageRating = totalRating / float64(len(kos.Ratings))
	}

	var kosFacilities []KosFacilityResponse
	for _, f := range kos.KosFacilities {
		kosFacilities = append(kosFacilities, KosFacilityResponse{
			ID:       f.ID,
			Facility: f.Facility,
		})
	}

	var kosRules []KosRuleResponse
	for _, r := range kos.KosRules {
		kosRules = append(kosRules, KosRuleResponse{
			ID:   r.ID,
			Rule: r.Rule,
		})
	}

	return KosResponseDetail{
		ID:            kos.ID,
		Name:          kos.Name,
		Description:   kos.Description,
		Rooms:         kos.Rooms,
		Rating:        averageRating,
		Category:      kos.Category,
		Price:         kos.Price,
		Address:       kos.Address,
		Longitude:     kos.Longitude,
		Latitude:      kos.Latitude,
		KosFacilities: kosFacilities,
		KosRules:      kosRules,
		PhotoKos: PhotoDetailResponse{
			PhotoMain:       kos.PhotoMain,
			PhotoFront:      kos.PhotoFront,
			PhotoBack:       kos.PhotoBack,
			PhotoRoomFront:  kos.PhotoRoomFront,
			PhotoRoomInside: kos.PhotoRoomInside,
		},
		User: handler.UserKosDetailResponse{
			Name:         kos.User.Name,
			UserName:     kos.User.UserName,
			PhotoProfile: kos.User.PhotoProfile,
		},
	}
}

func CoreToGetUser(kos kos.Core) KosResponseUser {
	var totalRating float64
	for _, rating := range kos.Ratings {
		totalRating += float64(rating.Score)
	}

	var averageRating float64
	if len(kos.Ratings) > 0 {
		averageRating = totalRating / float64(len(kos.Ratings))
	}

	var kosFacilities []KosFacilityResponse
	for _, f := range kos.KosFacilities {
		kosFacilities = append(kosFacilities, KosFacilityResponse{
			ID:       f.ID,
			Facility: f.Facility,
		})
	}

	return KosResponseUser{
		ID:            kos.ID,
		Name:          kos.Name,
		Rating:        averageRating,
		Rooms:         kos.Rooms,
		Price:         kos.Price,
		Address:       kos.Address,
		KosFacilities: kosFacilities,
		PhotoKos:      PhotoMainResponse{PhotoMain: kos.PhotoMain},
	}
}
