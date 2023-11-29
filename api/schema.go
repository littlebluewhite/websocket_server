package api

import (
	"github.com/gofiber/contrib/websocket"
)

type WebsocketManager interface {
	Register(d int, client *websocket.Conn)
	Unregister(d int, client *websocket.Conn)
	Broadcast(d int, message []byte)
}
