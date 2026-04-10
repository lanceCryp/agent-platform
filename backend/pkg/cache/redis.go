package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type RedisCache struct {
	Client *redis.Client
}

func NewRedisCache(ctx context.Context, addr, password string, db int) (*RedisCache, error) {
	logrus.Info("Connecting to Redis...")

	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     100,
		MinIdleConns: 10,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	// Test connection
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logrus.Info("Redis connection established successfully")
	return &RedisCache{Client: client}, nil
}

func (c *RedisCache) Close() error {
	if c.Client != nil {
		logrus.Info("Closing Redis connection")
		return c.Client.Close()
	}
	return nil
}

// Get retrieves a value by key
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.Client.Get(ctx, key).Result()
}

// Set stores a value with optional expiration
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		data = string(jsonData)
	}
	return c.Client.Set(ctx, key, data, expiration).Err()
}

// Delete removes a key
func (c *RedisCache) Delete(ctx context.Context, keys ...string) error {
	return c.Client.Del(ctx, keys...).Err()
}

// Exists checks if a key exists
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.Client.Exists(ctx, key).Result()
	return result > 0, err
}

// SetNX sets a value only if it doesn't exist (for distributed locks)
func (c *RedisCache) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	var data string
	switch v := value.(type) {
	case string:
		data = v
	default:
		jsonData, err := json.Marshal(v)
		if err != nil {
			return false, fmt.Errorf("failed to marshal value: %w", err)
		}
		data = string(jsonData)
	}
	return c.Client.SetNX(ctx, key, data, expiration).Result()
}

// Incr increments a counter
func (c *RedisCache) Incr(ctx context.Context, key string) (int64, error) {
	return c.Client.Incr(ctx, key).Result()
}

// Expire sets a key's expiration
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.Client.Expire(ctx, key, expiration).Err()
}

// TTL returns the remaining time to live of a key
func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.Client.TTL(ctx, key).Result()
}

// Rate limiting helper
type RateLimiter struct {
	cache *RedisCache
	limit int
	window time.Duration
}

func NewRateLimiter(cache *RedisCache, limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		cache:  cache,
		limit:  limit,
		window: window,
	}
}

// Allow checks if a request should be allowed based on rate limit
func (rl *RateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	redisKey := fmt.Sprintf("rate_limit:%s", key)
	
	// Get current count
	count, err := rl.cache.Incr(ctx, redisKey)
	if err != nil {
		return false, err
	}
	
	// Set expiration on first request
	if count == 1 {
		if err := rl.cache.Expire(ctx, redisKey, rl.window); err != nil {
			return false, err
		}
	}
	
	return count <= int64(rl.limit), nil
}

// Health check
func (c *RedisCache) HealthCheck(ctx context.Context) map[string]interface{} {
	status := "healthy"
	latency := time.Duration(0)

	start := time.Now()
	if err := c.Client.Ping(ctx).Err(); err != nil {
		status = "unhealthy"
	}
	latency = time.Since(start)

	info := make(map[string]interface{})
	if client := c.Client; client != nil {
		if dbSize, err := client.DBSize(ctx).Result(); err == nil {
			info["db_size"] = dbSize
		}
	}

	return map[string]interface{}{
		"status":     status,
		"latency_ms": latency.Milliseconds(),
		"info":       info,
	}
}
