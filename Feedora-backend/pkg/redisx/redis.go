// Package redisx 封装 Redis 客户端初始化。
// 用于阶段二的业务缓存、计数、榜单 ZSet、幂等与限流。
package redisx

import (
	"context"
	"time"

	"github.com/feedora/backend/pkg/config"
	"github.com/redis/go-redis/v9"
)

// New 根据配置创建 Redis 客户端并校验连通性。
func New(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}
	return client, nil
}
