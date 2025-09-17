package adaptor

import (
	"fmt"

	"github.com/michaelyusak/go-helper/entity"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(config entity.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Host, config.Port),
		Password: config.Password,
		DB:       config.Db,
	})
}
