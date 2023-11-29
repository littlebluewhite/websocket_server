package ws

import (
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"websocket_server/app/dbs"
)

type websocketManager interface {
	Broadcast(d int, message []byte)
}

type Operate struct {
	cache *cache.Cache
	rdb   *redis.Client
	wm    websocketManager
}

func NewOperate(dbs dbs.Dbs, wm websocketManager) *Operate {
	o := &Operate{
		cache: dbs.GetCache(),
		rdb:   dbs.GetRdb(),
		wm:    wm,
	}
	return o
}

var StreamNodeObjectMap = make(map[string]func(rsc map[string]interface{}) (string, error))
var StreamAlarmMap = make(map[string]func(rsc map[string]interface{}) (string, error))

func (o *Operate) getStreamNodeObjectMap() map[string]func(rsc map[string]interface{}) (string, error) {
	StreamNodeObjectMap["send_to_websocket"] = o.streamNodeObjectBroadcast
	return StreamNodeObjectMap
}

func (o *Operate) getStreamAlarmMap() map[string]func(rsc map[string]interface{}) (string, error) {
	StreamAlarmMap["send_to_websocket"] = o.streamAlarmBroadcast
	return StreamAlarmMap
}

func (o *Operate) streamNodeObjectBroadcast(rsc map[string]interface{}) (result string, err error) {
	d := []byte(rsc["data"].(string))
	o.wm.Broadcast(1, d)
	return
}

func (o *Operate) streamAlarmBroadcast(rsc map[string]interface{}) (result string, err error) {
	d := []byte(rsc["data"].(string))
	o.wm.Broadcast(2, d)
	return
}
