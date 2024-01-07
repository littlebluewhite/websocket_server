package api

import (
	"github.com/gofiber/contrib/websocket"
	"websocket_server/entry/e_module"
)

type HubManager interface {
	RegisterHub(module e_module.Module)
	Broadcast(module e_module.Module, message []byte)
	WsConnect(module e_module.Module, conn *websocket.Conn) error
}
