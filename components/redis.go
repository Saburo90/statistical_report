package components

import (
	"gitee.com/NotOnlyBooks/statistical_report/conf"
	"github.com/go-redis/redis"
)

var (
	Redis *redis.Client
)

func SetupRedis(redisConf *conf.RedisConfig) error {
	Redis = redis.NewClient(&redis.Options{
		Network: redisConf.Network,
		Addr:    redisConf.Addr,
	})

	return nil
}
