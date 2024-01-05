package websocket_hub

import (
	"errors"
	"fmt"
	"github.com/gofiber/contrib/websocket"
	"websocket_server/util/logFile"
)

type hub interface {
	Run()
	Broadcast([]byte)
	WsConnect(client *websocket.Conn)
}

type HubManager struct {
	hubs map[string]hub
	l    logFile.LogFile
}

func NewHubManager() *HubManager {
	return &HubManager{
		hubs: make(map[string]hub),
		l:    logFile.NewLogFile("websocket", "hub_manager.log"),
	}
}

func (hm *HubManager) RegisterHub(model string) {
	h := NewHub(model)
	hm.hubs[model] = h
	go h.Run()
}

func (hm *HubManager) Broadcast(model string, msg []byte) {
	hm.l.Info().Printf("model: %s broadcast: %s", model, msg)
	hm.hubs[model].Broadcast(msg)
}

func (hm *HubManager) WsConnect(model string, conn *websocket.Conn) error {
	if h, ok := hm.hubs[model]; !ok {
		eString := fmt.Sprintf("model %s not exist", model)
		hm.l.Error().Println(eString)
		return errors.New(eString)
	} else {
		h.WsConnect(conn)
	}
	return nil
}
