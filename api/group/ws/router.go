package ws

import (
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"websocket_server/api"
	"websocket_server/app/dbs"
	"websocket_server/util/logFile"
)

func RegisterRouter(g group) {
	log := logFile.NewLogFile("router", "websocket.log")
	app := g.GetApp()

	hm := g.GetWebsocketHub()
	hm.RegisterHub("node_object")
	hm.RegisterHub("alarm")

	o := NewOperate(g.GetDbs(), hm)

	go func() {
		receiveNodeObjectStream(o, log)
	}()
	go func() {
		receiveAlarmStream(o, log)
	}()

	ws := app.Group("/ws")
	ws.Use("/", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	ws.Get("/node_object", websocket.New(func(c *websocket.Conn) {
		err := hm.WsConnect("node_object", c)
		if err != nil {
			log.Error().Println(err)
		}
	}))
	ws.Get("/alarm", websocket.New(func(c *websocket.Conn) {
		err := hm.WsConnect("alarm", c)
		if err != nil {
			log.Error().Println(err)
		}
	}))
}

type group interface {
	GetApp() fiber.Router
	GetDbs() dbs.Dbs
	GetWebsocketHub() api.HubManager
}
