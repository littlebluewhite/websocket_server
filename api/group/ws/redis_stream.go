package ws

import (
	"context"
	"websocket_server/api"
	"websocket_server/app/dbs/rdb/deal_redis"
)

func receiveNodeObjectStream(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start node_object receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := deal_redis.NewStreamRead(o.rdb, "NodeObjectWebsocket", "server", l)
	rs.Start(ctx, o.getStreamNodeObjectMap())
}

func receiveAlarmStream(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start alarm receiveStream --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rs := deal_redis.NewStreamRead(o.rdb, "AlarmWebsocket", "server", l)
	rs.Start(ctx, o.getStreamAlarmMap())
}

func subscribeNodeObject(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start node_object subscribe --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rpb := deal_redis.NewRedisPubSub(o.rdb, "NodeObjectWebsocket", l)
	rpb.Subscribe(ctx, o.getStreamNodeObjectMap())
}

func subscribeAlarm(o *Operate, l api.Logger) {
	l.Infoln("----------------------------------- start alarm subscribe --------------------------------")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rpb := deal_redis.NewRedisPubSub(o.rdb, "AlarmWebsocket", l)
	rpb.Subscribe(ctx, o.getStreamAlarmMap())
}
