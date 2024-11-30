package credis

import "github.com/redis/go-redis/v9"

type RedisDB struct {
	Client *redis.Client
}

func NewRedisClient(addr, password string, db int) *RedisDB {
	return &RedisDB{
		Client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}
