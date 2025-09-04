package db

import (
	"context"

	redisv9 "github.com/redis/go-redis/v9"
)

type RedisOptions struct {
	Addr     string
	Password string
	DB       int
}

func NewRedisClient(opts RedisOptions) *redisv9.Client {
	return redisv9.NewClient(&redisv9.Options{
		Addr:     opts.Addr,
		Password: opts.Password,
		DB:       opts.DB,
	})
}

func RedisPing(ctx context.Context, client *redisv9.Client) error {
	return client.Ping(ctx).Err()
} 