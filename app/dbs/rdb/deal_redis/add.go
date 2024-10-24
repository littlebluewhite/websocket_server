package deal_redis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func StreamAdd(ctx context.Context, rdb redis.UniversalClient,
	streamName string, values map[string]interface{}) error {
	_, err := rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		Values: values,
	}).Result()
	return err
}
