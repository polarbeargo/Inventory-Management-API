package database

import (
	"context"
	"encoding/json"
	"inventory_management/models"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

func GetItemFromCache(id string, item *models.Item) bool {
	if RedisClient != nil && RedisCtx != nil {
		cached, err := RedisClient.Get(RedisCtx, "item:"+id).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cached), item); err == nil {
				return true
			}
		}
	}
	return false
}

func SetItemToCache(id string, item models.Item) {
	if RedisClient != nil && RedisCtx != nil {
		b, _ := json.Marshal(item)
		_ = RedisClient.Set(RedisCtx, "item:"+id, b, 0).Err()
	}
}

func DeleteItemFromCache(id string) {
	if RedisClient != nil && RedisCtx != nil {
		_ = RedisClient.Del(RedisCtx, "item:"+id).Err()
	}
}
