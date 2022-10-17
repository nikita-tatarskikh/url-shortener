package repository

import (
	"context"

	"github.com/go-redis/redis/v9"
)

type RedisRepository struct {
	conn *redis.Client
}

func NewRedisRepository(conn *redis.Client) *RedisRepository {
	return &RedisRepository{conn: conn}
}

func (r *RedisRepository) Store(shortUrl, url string) error {

	err := r.conn.TxPipeline().Set(context.TODO(), shortUrl, url, 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r *RedisRepository) Retrieve(shortUrl string) string {
	return r.conn.TxPipeline().Get(context.TODO(), shortUrl).Val()
}
