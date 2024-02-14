package data

import (
	"KosKita/features/chat"
	"KosKita/features/user/data"

	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	Message      string
	RoomID       string
	ReceiverID   uint
	SenderID     uint
	UserReceiver data.User `gorm:"foreignKey:ReceiverID;references:ID"`
	UserSender   data.User `gorm:"foreignKey:SenderID;references:ID"`
}

func CoreToModelChat(input chat.Core) Chat {
	return Chat{
		Message:    input.Message,
		RoomID:     input.RoomID,
		ReceiverID: input.ReceiverID,
		SenderID:   input.SenderID,
	}
}

func (m Chat) ModelToCoreChat() chat.Core {
	return chat.Core{
		Message:    m.Message,
		RoomID:     m.RoomID,
		SenderID:   m.SenderID,
		ReceiverID: m.ReceiverID,
	}
}
