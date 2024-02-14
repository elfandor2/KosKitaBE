package handler

import (
	ch "KosKita/features/chat"
	cd "KosKita/features/chat/data"
	hub "KosKita/features/chat/service"
	"KosKita/utils/middlewares"
	"KosKita/utils/responses"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
	chatService ch.ChatServiceInterface
	hub         *hub.Hub
}

func New(cs ch.ChatServiceInterface, h *hub.Hub) *ChatHandler {
	return &ChatHandler{
		chatService: cs,
		hub:         h,
	}
}

func (ch *ChatHandler) CreateRoom(c echo.Context) error {
	var req CreateRoomReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	roomID, err := generateRoomID()
	if err != nil {
		return err
	}

	ch.hub.Rooms[req.ID] = &hub.Room{
		ID:      roomID,
		Clients: make(map[string]*hub.Client),
	}

	return c.JSON(http.StatusOK, req)
}

func (ch *ChatHandler) JoinRoom(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	roomID := c.Param("roomId")
	clientID := c.QueryParam("userId")

	cl := &hub.Client{
		Conn:    conn,
		Message: make(chan *cd.Chat, 10),
		ID:      clientID,
		RoomID:  roomID,
	}

	m := &cd.Chat{
		// Model:        gorm.Model{},
		Message:    "",
		RoomID:     roomID,
		ReceiverID: 0,
		SenderID:   uint(userIdLogin),
		// UserReceiver: data.User{},
		// UserSender:   data.User{},
	}

	ch.hub.Register <- cl
	ch.hub.Broadcast <- m

	go cl.WriteMessage()
	cl.ReadMessage(ch.hub, ch.chatService)

	return nil
}

func (ch *ChatHandler) GetMessages(c echo.Context) error {
	roomID := c.Param("roomId")

	chats, errGet := ch.chatService.GetMessage(roomID)
	if errGet != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": errGet.Error()})
	}

	chatResult := CoreToGetChats(chats)

	return c.JSON(http.StatusOK, responses.WebResponse("success get message.", chatResult))

}

func (ch *ChatHandler) GetRooms(c echo.Context) error {
	rooms := make([]RoomRes, 0)

	for _, r := range ch.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID: r.ID,
		})
	}

	return c.JSON(http.StatusOK, rooms)
}
