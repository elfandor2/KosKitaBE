package handler

import (
	"crypto/rand"
	"encoding/hex"
)

type CreateRoomReq struct {
	ID string `json:"room_id"`
}

func generateRoomID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}
