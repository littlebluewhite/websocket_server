package ws

import (
	"context"
	"websocket_server/util/logFile"
	"websocket_server/util/redis_stream"
)

func receiveNodeObjectStream(o *Operate, l logFile.LogFile) {
	l.Info().Println("----------------------------------- start node_object receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := redis_stream.NewStreamRead(o.rdb, "NodeObjectWebsocket", "server", l)
	rs.Start(ctx, o.getStreamNodeObjectMap())
}

func receiveAlarmStream(o *Operate, l logFile.LogFile) {
	l.Info().Println("----------------------------------- start alarm receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := redis_stream.NewStreamRead(o.rdb, "AlarmWebsocket", "server", l)
	rs.Start(ctx, o.getStreamAlarmMap())
}
