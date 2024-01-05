package api

import (
	"github.com/gofiber/contrib/websocket"
)

type HubManager interface {
	RegisterHub(model string)
	Broadcast(model string, message []byte)
	WsConnect(model string, conn *websocket.Conn) error
}
