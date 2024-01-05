package ws

import (
	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"websocket_server/app/dbs"
)

type hubManager interface {
	Broadcast(model string, message []byte)
}

type Operate struct {
	cache *cache.Cache
	rdb   *redis.Client
	hm    hubManager
}

func NewOperate(dbs dbs.Dbs, hm hubManager) *Operate {
	o := &Operate{
		cache: dbs.GetCache(),
		rdb:   dbs.GetRdb(),
		hm:    hm,
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
	o.hm.Broadcast("node_object", d)
	return
}

func (o *Operate) streamAlarmBroadcast(rsc map[string]interface{}) (result string, err error) {
	d := []byte(rsc["data"].(string))
	o.hm.Broadcast("alarm", d)
	return
}
