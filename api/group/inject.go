package group

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"io"
	"os"
	"websocket_server/api"
	"websocket_server/api/group/ws"
	"websocket_server/app/dbs"
	"websocket_server/util/my_log"
)

func Inject(app *fiber.App, dbs dbs.Dbs, wm api.HubManager) {
	// Middleware
	log := my_log.NewLog("api/inject")
	fiberLog := getFiberLogFile(log)
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Output: fiberLog,
	}))

	// api group add cors middleware
	Api := app.Group("/api", cors.New())

	// create new group
	g := NewAPIGroup(Api, dbs, wm)

	// model registration
	ws.RegisterRouter(g)
}

func getFiberLogFile(log api.Logger) io.Writer {
	fiberFile, err := os.OpenFile("./log/fiber.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Errorln("can not open log file: " + err.Error())
	}
	return io.MultiWriter(fiberFile, os.Stdout)
}
