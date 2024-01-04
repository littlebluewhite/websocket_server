package websocket_manager

import (
	"github.com/gofiber/contrib/websocket"
	"websocket_server/util/logFile"
)

type WebsocketManager struct {
	groups     map[Group]map[*websocket.Conn]struct{}
	register   chan groupConnect
	unregister chan groupConnect
	broadcast  chan groupMessage
	l          logFile.LogFile
}

func NewWebsocketManager() *WebsocketManager {
	return &WebsocketManager{
		groups: map[Group]map[*websocket.Conn]struct{}{
			None:       make(map[*websocket.Conn]struct{}),
			NodeObject: make(map[*websocket.Conn]struct{}),
			Alarm:      make(map[*websocket.Conn]struct{}),
		},
		register:   make(chan groupConnect, 1024),
		unregister: make(chan groupConnect, 1024),
		broadcast:  make(chan groupMessage, 1024),
		l:          logFile.NewLogFile("websocket", "websocket.log"),
	}
}

func (wm *WebsocketManager) Run() {
	defer func() {
		wm.l.Error().Println("websocket manager exit")
	}()
	for {
		select {
		case gc := <-wm.register:
			wm.l.Info().Println("register 2:", gc.group, gc.client)
			wm.groups[gc.group][gc.client] = struct{}{}
		case gc := <-wm.unregister:
			wm.l.Info().Println("unregister 2:", gc.group, gc.client)
			delete(wm.groups[gc.group], gc.client)
			_ = gc.client.Close()
		case gm := <-wm.broadcast:
			wm.l.Info().Println("broadcast 2:", gm.group)
			for client := range wm.groups[gm.group] {
				err := client.WriteMessage(websocket.TextMessage, gm.message)
				if err != nil {
					wm.l.Error().Println("send message error: ", err, "group: ", gm.group)
					delete(wm.groups[gm.group], client)
					_ = client.Close()
				}
			}
		}
	}
}

func (wm *WebsocketManager) Register(d int, client *websocket.Conn) {
	wm.l.Info().Println("register 1:", d, client)
	wm.register <- groupConnect{
		group:  Group(d),
		client: client,
	}
}

func (wm *WebsocketManager) Unregister(d int, client *websocket.Conn) {
	wm.l.Info().Println("unregister 1:", d, client)
	wm.unregister <- groupConnect{
		group:  Group(d),
		client: client,
	}
}

func (wm *WebsocketManager) Broadcast(d int, message []byte) {
	wm.l.Info().Println("broadcast 1:", d, string(message))
	wm.broadcast <- groupMessage{
		group:   Group(d),
		message: message,
	}
}
