package redis

import (
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/obrr-hhx/simpleDouyin/pkg/constants"
)

var (
	expireTime  = time.Hour * 1
	rdbFollows  *redis.Client
	rdbFavorite *redis.Client
)

func InitRedis() {
	rdbFollows = redis.NewClient(&redis.Options{
		Addr:     constants.RedisDefaultAddr,
		Password: constants.RedisDefaultPwd,
		DB:       0,
	})
	rdbFavorite = redis.NewClient(&redis.Options{
		Addr:     constants.RedisDefaultAddr,
		Password: constants.RedisDefaultPwd,
		DB:       1,
	})
}
