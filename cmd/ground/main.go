package main

import (
	"context"
	"websocket_server/app/dbs"
	"websocket_server/app/dbs/rdb/deal_redis"
	"websocket_server/util/config"
	"websocket_server/util/my_log"
)

func main() {
	testLog := my_log.NewLog("test")
	Config := config.NewConfig[config.Config]("./config", "config", config.Yaml)

	DBS := dbs.NewDbs(testLog, Config)
	rps := deal_redis.NewRedisPubSub(DBS.GetRdb(), "NodeObjectWebsocket", testLog)
	var comMap map[string]func(rsc map[string]interface{}) (string, error)
	rps.Subscribe(context.Background(), comMap)
}
