package deal_redis

import (
	"context"
	"fmt"
	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	"websocket_server/api"
)

func connectionBaseExecute(
	ctx context.Context, rdb redis.UniversalClient,
	comMap map[string]func(map[string]interface{}) (string, error),
	connectionPayload map[string]interface{}, l api.Logger) {
	com := comMap[connectionPayload["command"].(string)]
	result, err := com(connectionPayload)
	if err != nil {
		l.Errorln(DealCommandErr)
	}
	// isCallBack deal type
	isCallBack := "0"
	switch connectionPayload["is_wait_call_back"].(type) {
	case string:
		isCallBack = connectionPayload["is_wait_call_back"].(string)
	case float64:
		isCallBack = fmt.Sprintf("%v", connectionPayload["is_wait_call_back"])
	}
	fmt.Println(isCallBack)
	if isCallBack == "1" {
		err = connectionBaseCallback(ctx, rdb, connectionPayload["callback_channel"].(string), result, err)
		if err != nil {
			l.Errorln("call back publish error: ", err)
		}
		l.Infoln("return callback success")
	}
}

func connectionBaseCallback(
	ctx context.Context, rdb redis.UniversalClient,
	callBackChannel string, result string, err error) error {
	connectionPayload := make(map[string]interface{})
	if err != nil {
		connectionPayload["data"] = err.Error()
		connectionPayload["status_code"] = "422"
	} else {
		connectionPayload["data"] = result
	}
	cb, _ := json.Marshal(connectionPayload)
	err = rdb.Publish(ctx, callBackChannel, cb).Err()
	return err
}
