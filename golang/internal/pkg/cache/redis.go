package cache

import (
	"github.com/go-redis/redis"
	"github.com/kyh0703/template/configs"
)

var client *redis.Client

func NewRedisClient(cfg *configs.Config) (*redis.Client, error) {
	// create new client
	client = redis.NewClient(&redis.Options{Addr: cfg.Infra.Redis.MasterName})
	// ping test
	if _, err := client.Ping().Result(); err != nil {
		return nil, err
	}
	return client, nil
}
