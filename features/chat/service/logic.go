package service

import "KosKita/features/chat"

type chatService struct {
	chatData chat.ChatDataInterface
}

func New(repo chat.ChatDataInterface) chat.ChatServiceInterface {
	return &chatService{
		chatData: repo,
	}
}

// CreateRoom implements chat.ChatServiceInterface.
func (cs *chatService) CreateChat(userIdLogin int, input chat.Core) (chat.Core, error) {
	chatOutput, errInsert := cs.chatData.CreateMessage(userIdLogin, input)
	if errInsert != nil {
		return chat.Core{}, errInsert
	}

	return chatOutput, nil
}


// GetMessage implements chat.ChatServiceInterface.
func (cs *chatService) GetMessage(roomId string) ([]chat.Core, error) {
	chats, errGet := cs.chatData.GetMessage(roomId)
	if errGet != nil {
		return nil, errGet
	}

	return chats, nil
}