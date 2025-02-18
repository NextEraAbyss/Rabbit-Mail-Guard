package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func NewRedisClient(host, port, password string, db int) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	return &Client{rdb: rdb}
}

func (c *Client) SetVerificationCode(ctx context.Context, email, code string) error {
	return c.rdb.Set(ctx, "verification:"+email, code, 5*time.Minute).Err()
}

func (c *Client) VerifyCode(ctx context.Context, email, code string) bool {
	storedCode, err := c.rdb.Get(ctx, "verification:"+email).Result()
	if err != nil {
		return false
	}
	return storedCode == code
} 