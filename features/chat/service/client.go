package service

import (
	cc "KosKita/features/chat"
	cd "KosKita/features/chat/data"
	"log"
	"strconv"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn    *websocket.Conn
	Message chan *cd.Chat
	ID      string `json:"id"`
	RoomID  string `json:"roomId"`
}

type ChatRes struct {
	Message    string `json:"message"`
	SenderID   uint   `json:"sender_id"`
	ReceiverID uint   `json:"receiver_id"`
	RoomID     string `json:"room_id"`
}

func (c *Client) WriteMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}

		result := ChatRes{
			Message:    message.Message,
			RoomID:     message.RoomID,
			ReceiverID: message.ReceiverID,
			SenderID:   message.SenderID,
		}

		c.Conn.WriteJSON(result)
	}
}

func (c *Client) ReadMessage(hub *Hub, chatService cc.ChatServiceInterface) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		msg := &cd.Chat{
			Message: string(m),
			RoomID:  c.RoomID,
		}

		coreMsg := cc.Core{
			Message:    msg.Message,
			RoomID:     msg.RoomID,
			ReceiverID: msg.ReceiverID,
			SenderID:   msg.SenderID,
		}

		userID, err := strconv.Atoi(c.ID)
		if err != nil {
			log.Printf("Error converting ID to integer: %v", err)
			continue
		}
		
		_, err = chatService.CreateChat(userID, coreMsg)
		if err != nil {
			log.Printf("Error saving message: %v", err)
			continue
		}
		hub.Broadcast <- msg
	}
}
