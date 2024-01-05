package group

import (
	"github.com/gofiber/fiber/v2"
	"websocket_server/api"
	"websocket_server/app/dbs"
)

type Group struct {
	app fiber.Router
	dbs dbs.Dbs
	hm  api.HubManager
}

func NewAPIGroup(app fiber.Router, dbs dbs.Dbs, hm api.HubManager) *Group {
	return &Group{
		app: app,
		dbs: dbs,
		hm:  hm,
	}
}

func (g *Group) GetApp() fiber.Router {
	return g.app
}

func (g *Group) GetDbs() dbs.Dbs {
	return g.dbs
}

func (g *Group) GetWebsocketHub() api.HubManager {
	return g.hm
}
