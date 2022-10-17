package redis

import (
	"context"
	"url-shortener/internal/configuration"
	"url-shortener/internal/logger"

	"github.com/go-redis/redis/v9"
)

func NewConnection(redisCfg configuration.Redis, logger *logger.Logger) (*redis.Client, func(), error) {

	rdb := redis.NewClient(&redis.Options{
		Addr: redisCfg.Host,
	})

	if _, err := rdb.Ping(context.TODO()).Result(); err != nil {
		return nil, nil, err
	}

	return rdb, func() {
		if err := rdb.Close(); err != nil {
			logger.LogError("fail to close redis connection", err)
		}
	}, nil
}
