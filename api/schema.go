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

type Logger interface {
	Infoln(args ...interface{})
	Infof(s string, args ...interface{})
	Errorln(args ...interface{})
	Errorf(s string, args ...interface{})
	Warnln(args ...interface{})
	Warnf(s string, args ...interface{})
}
