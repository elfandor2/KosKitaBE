package handler

import (
	"KosKita/features/kos"
	"mime/multipart"
)

type KosRequest struct {
	Name          string   `json:"kos_name"`
	Description   string   `json:"description" form:"description"`
	Category      string   `json:"category"`
	Price         int      `json:"price"`
	Rooms         int      `json:"rooms"`
	Address       string   `json:"address"`
	Longitude     string   `json:"longitude"`
	Latitude      string   `json:"latitude"`
	KosFacilities []string `json:"kos_facilities" `
	KosRules      []string `json:"kos_rules"`
	UserID        uint
}

type KosFotoRequest struct {
	PhotoMain       *multipart.FileHeader `form:"main_kos_photo"`
	PhotoFront      *multipart.FileHeader ` form:"front_kos_photo"`
	PhotoBack       *multipart.FileHeader `form:"back_kos_photo"`
	PhotoRoomFront  *multipart.FileHeader ` form:"front_room_photo"`
	PhotoRoomInside *multipart.FileHeader ` form:"inside_room_photo"`
	UserID          uint
}

type RatingRequest struct {
	Score int `json:"score" form:"score"`
}

func RequestToCore(input KosRequest, userIdLogin uint) kos.CoreInput {
	var kosFacilities []kos.KosFacilityCore
	for _, facility := range input.KosFacilities {
		kosFacilities = append(kosFacilities, kos.KosFacilityCore{
			Facility: facility,
		})
	}

	var kosRules []kos.KosRuleCore
	for _, rule := range input.KosRules {
		kosRules = append(kosRules, kos.KosRuleCore{
			Rule: rule,
		})
	}

	return kos.CoreInput{
		UserID:        userIdLogin,
		Name:          input.Name,
		Description:   input.Description,
		Category:      input.Category,
		Price:         input.Price,
		Rooms:         input.Rooms,
		Address:       input.Address,
		Longitude:     input.Longitude,
		Latitude:      input.Latitude,
		KosFacilities: kosFacilities,
		KosRules:      kosRules,
	}
}

func RequestToCoreFoto(imageURLs []string, userIdLogin uint) kos.CoreFoto {
	return kos.CoreFoto{
		UserID:          userIdLogin,
		PhotoMain:       imageURLs[0],
		PhotoFront:      imageURLs[1],
		PhotoBack:       imageURLs[2],
		PhotoRoomFront:  imageURLs[3],
		PhotoRoomInside: imageURLs[4],
	}
}

func RequestToCoreFotoPut(imageURLs []string, userIdLogin uint) kos.CoreFoto {
    coreFoto := kos.CoreFoto{
        UserID: userIdLogin,
    }

    numImages := len(imageURLs)
    if numImages > 0 {
        coreFoto.PhotoMain = imageURLs[0]
    }
    if numImages > 1 {
        coreFoto.PhotoFront = imageURLs[1]
    }
    if numImages > 2 {
        coreFoto.PhotoBack = imageURLs[2]
    }
    if numImages > 3 {
        coreFoto.PhotoRoomFront = imageURLs[3]
    }
    if numImages > 4 {
        coreFoto.PhotoRoomInside = imageURLs[4]
    }

    return coreFoto
}


func RequestToCorePut(input KosRequest, userIdLogin uint) kos.Core {
	var kosFacilities []kos.KosFacilityCore
	for _, facility := range input.KosFacilities {
		kosFacilities = append(kosFacilities, kos.KosFacilityCore{
			Facility: facility,
		})
	}

	var kosRules []kos.KosRuleCore
	for _, rule := range input.KosRules {
		kosRules = append(kosRules, kos.KosRuleCore{
			Rule: rule,
		})
	}

	kos := kos.Core{
		UserID:      userIdLogin,
		Name:        input.Name,
		Description: input.Description,
		Category:    input.Category,
		Price:       input.Price,
		Rooms:       input.Rooms,
		Address:     input.Address,
		Longitude:   input.Longitude,
		Latitude:    input.Latitude,
		KosFacilities: kosFacilities,
		KosRules:      kosRules,
	}

	return kos
}

func RequestToCoreRating(input RatingRequest, kosId uint, userIdLogin uint) kos.RatingCore {
	return kos.RatingCore{
		UserID:          userIdLogin,
		BoardingHouseID: kosId,
		Score:           input.Score,
	}
}
