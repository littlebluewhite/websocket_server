package deal_redis

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"websocket_server/api"
)

type RedisStream struct {
	rdb        redis.UniversalClient
	streamName string
	groupName  string
	l          api.Logger
}

func NewStreamRead(rdb redis.UniversalClient, streamName string, groupName string, l api.Logger) *RedisStream {
	return &RedisStream{
		rdb:        rdb,
		streamName: streamName,
		groupName:  groupName,
		l:          l,
	}
}

func (rs *RedisStream) Start(ctx context.Context, comMap map[string]func(map[string]interface{}) (string, error)) {
	err := rs.streamInit(ctx)
	if err != nil {
		rs.l.Errorln("receiveStream error: ", err)
		return
	}
	for {
		connectionPayload, err := rs.ReadGroup(ctx)
		rs.l.Infoln("get stream")
		if err != nil {
			rs.l.Errorln("receive Stream error: ", err)
			continue
		}
		go func(connectionPayload map[string]interface{}) {
			connectionBaseExecute(ctx, rs.rdb, comMap, connectionPayload, rs.l)
		}(connectionPayload)
	}
}

func (rs *RedisStream) streamInit(ctx context.Context) (err error) {
	r, e := rs.rdb.XInfoGroups(ctx, rs.streamName).Result()
	if e != nil || len(r) == 0 {
		err = rs.rdb.XGroupCreateMkStream(ctx, rs.streamName, rs.groupName, "0").Err()
		if err != nil {
			return
		}
	}
	return nil
}

func (rs *RedisStream) ReadGroup(ctx context.Context) (
	rsr map[string]interface{}, err error) {
	re, err := rs.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    rs.groupName,
		Consumer: "s1",
		Streams:  []string{rs.streamName, ">"},
		Count:    1,
		Block:    0,
		NoAck:    true,
	}).Result()
	if err != nil {
		return
	}
	if len(re) == 0 || len(re[0].Messages) == 0 {
		err = fmt.Errorf("no messages found")
		return
	}

	message := re[0].Messages[0]
	rsr = message.Values

	err = rs.rdb.XDel(ctx, rs.streamName, message.ID).Err()
	if err != nil {
		return
	}
	return
}

func (rs *RedisStream) CallBack(ctx context.Context, callBackChannel string, result string, err error) error {
	connectionPayload := make(map[string]interface{})
	if err != nil {
		connectionPayload["data"] = err.Error()
		connectionPayload["status_code"] = "422"
	} else {
		connectionPayload["data"] = result
	}
	cb, _ := json.Marshal(connectionPayload)
	err = rs.rdb.Publish(ctx, callBackChannel, cb).Err()
	return err
}
