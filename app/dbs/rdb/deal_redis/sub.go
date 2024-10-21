package deal_redis

import (
	"context"
	"errors"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"websocket_server/api"
)

type RedisPubSub struct {
	rdb   redis.UniversalClient
	topic string
	l     api.Logger
}

func NewRedisPubSub(rdb redis.UniversalClient, topic string, l api.Logger) *RedisPubSub {
	return &RedisPubSub{
		rdb:   rdb,
		topic: topic,
		l:     l,
	}
}

func (t *RedisPubSub) Subscribe(ctx context.Context, comMap map[string]func(rsc map[string]interface{}) (string, error)) {
	pubsub := t.rdb.Subscribe(ctx, t.topic)
	defer func(pubsub *redis.PubSub) {
		err := pubsub.Close()
		if err != nil {
			t.l.Errorln(err)
		}
	}(pubsub)
	for {
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				t.l.Infoln("rdbSub receive canceled")
				return
			}
			t.l.Errorln(err)
			if errors.Is(err, redis.ErrClosed) {
				return
			}
			continue
		}

		// deal with message
		go func(payload string) {
			var connectionPayload map[string]interface{}
			if e := json.Unmarshal([]byte(payload), &connectionPayload); e != nil {
				t.l.Errorln(SubscribeErr)
			}
			connectionBaseExecute(ctx, t.rdb, comMap, connectionPayload, t.l)
		}(msg.Payload)
	}
}
