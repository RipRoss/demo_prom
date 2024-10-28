package redis

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	rbd := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	ctx := context.Background()
	_, err := rbd.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
	}

	return &RedisClient{client: rbd}
}

func (r *RedisClient) Set(key string, value interface{}) {
	ctx := context.Background()
	err := r.client.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Fatal(err)
	}
}

func (r *RedisClient) Get(key string) string {
	ctx := context.Background()
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Fatal(err)
	}

	return val
}

func (r *RedisClient) Close() {
	err := r.client.Close()
	if err != nil {
		log.Fatal(err)
	}
}
