package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Handle struct {
	DB *redis.Client
}

func (h Handle) CheckConnection() error {
	if h.DB == nil {
		panic("Redis not initialized !")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := h.DB.Ping(ctx).Result()
	if err != nil {
		return err
	}
	// logrus.Debugf("Redis ping: %s", result)
	return nil
}
