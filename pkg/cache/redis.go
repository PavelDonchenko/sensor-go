package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/PavelDonchenko/sensor-go/config"
	"github.com/go-redis/redis/v8"
)

var (
	ErrorCheckExist = errors.New("error to check if value exist")
	ErrorSetRedis   = errors.New("error to set value to Redis")
	ErrorGetRedis   = errors.New("error to get value from Redis")
)

type CacheRedis interface {
	Set(ctx context.Context, key string, val string) error
	Get(ctx context.Context, key string) (string, error)
	IfExistsInCache(ctx context.Context, key string) (bool, error)
}

type CacheConn struct {
	Client     *redis.Client
	Expiration time.Duration
}

func NewCacheConn(config config.Config) (*CacheConn, error) {
	fmt.Println("config", config)
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Redis.Address,
		DB:   0,
	})
	fmt.Println("redisClient", redisClient)

	ctx := context.Background()

	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err)
		// Sleep for 3 seconds and wait for Redis to initialize
		time.Sleep(3 * time.Second)
		err := redisClient.Ping(ctx).Err()
		if err != nil {
			return nil, err
		}
	}
	fmt.Println(pong)

	return &CacheConn{
		Client:     redisClient,
		Expiration: time.Duration(config.Redis.Expiration) * time.Second,
	}, nil
}

// Set sets a key-value pair
func (cache *CacheConn) Set(ctx context.Context, key string, val string) error {
	if err := cache.Client.Set(ctx, key, val, cache.Expiration).Err(); err != nil {
		return err
	}
	return nil
}

// Get returns true if the key already exists and set dst to the corresponding value
func (cache *CacheConn) Get(ctx context.Context, key string) (string, error) {
	val, err := cache.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func (cache *CacheConn) IfExistsInCache(ctx context.Context, key string) (bool, error) {
	exist, err := cache.Client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	if exist != 1 {
		return false, nil
	}
	return true, nil
}
