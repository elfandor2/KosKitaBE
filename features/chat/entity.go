package chat

import (
	"KosKita/features/user"
	"time"
)

type Core struct {
	ID         uint
	Message    string
	RoomID     string
	ReceiverID uint
	SenderID   uint
	User       user.Core
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// interface untuk Data Layer
type ChatDataInterface interface {
	CreateMessage(userIdLogin int, input Core) (Core, error)
	GetMessage(roomId string) ([]Core, error)
}

// interface untuk Service Layer
type ChatServiceInterface interface {
	CreateChat(userIdLogin int, input Core) (Core, error)
	GetMessage(roomId string) ([]Core, error)
}
