package main

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"websocket_server/api/group"
	"websocket_server/app/dbs"
	"websocket_server/app/websocket_hub"
	"websocket_server/util/config"
	"websocket_server/util/my_log"
)

// @title           Websocket_Server swagger API
// @version         2.2.8
// @description     This is a websocket server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Wilson
// @contact.url    https://github.com/littlebluewhite
// @contact.email  wwilson008@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      127.0.0.1:5488

func main() {
	mainLog := my_log.NewLog("main")
	// Create context that listens for the interrupt signal from the OS.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// read config
	Config := config.NewConfig[config.Config]("./config", "config", config.Yaml)

	// DBs start includes SQL Cache
	DBS := dbs.NewDbs(mainLog, Config)
	defer func() {
		DBS.GetIdb().Close()
		mainLog.Infoln("influxDB Disconnect")
	}()

	// create websocket manager
	hm := websocket_hub.NewHubManager()

	serverConfig := Config.Server

	var sb strings.Builder
	sb.WriteString(":")
	sb.WriteString(serverConfig.Port)
	//addr := sb.String()
	apiServer := fiber.New(
		fiber.Config{
			ReadTimeout:  serverConfig.ReadTimeout * time.Minute,
			WriteTimeout: serverConfig.WriteTimeout * time.Minute,
			AppName:      "websocket_server",
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		},
	)

	group.Inject(apiServer, DBS, hm)

	// for api server shout down gracefully
	serverShutdown := make(chan struct{})
	go func() {
		_ = <-ctx.Done()
		mainLog.Infoln("Gracefully shutting down api server")
		_ = apiServer.Shutdown()
		serverShutdown <- struct{}{}
	}()

	if err := apiServer.Listen(":5488"); err != nil {
		mainLog.Errorf("listen: %s\n", err)
	}

	// Listen for the interrupt signal.
	<-serverShutdown

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	time.Sleep(1 * time.Second)
	mainLog.Infoln("Server exiting")

}
