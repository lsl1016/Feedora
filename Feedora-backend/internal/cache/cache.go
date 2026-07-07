package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache 是对 Redis 的轻量封装。所有方法在底层 client 为 nil 时安全降级（视为未命中 / 空操作），
// 使得未启用 Redis 时业务仍可运行。
type Cache struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) *Cache {
	return &Cache{rdb: rdb}
}

// Enabled 是否已接入 Redis。
func (c *Cache) Enabled() bool { return c != nil && c.rdb != nil }

// GetJSON 读取 JSON 缓存，命中返回 true。
func (c *Cache) GetJSON(ctx context.Context, key string, dest any) bool {
	if !c.Enabled() {
		return false
	}
	b, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return false
	}
	return json.Unmarshal(b, dest) == nil
}

// SetJSON 写入 JSON 缓存。
func (c *Cache) SetJSON(ctx context.Context, key string, val any, ttl time.Duration) {
	if !c.Enabled() {
		return
	}
	if b, err := json.Marshal(val); err == nil {
		c.rdb.Set(ctx, key, b, ttl)
	}
}

// Del 删除缓存 Key。
func (c *Cache) Del(ctx context.Context, keys ...string) {
	if !c.Enabled() || len(keys) == 0 {
		return
	}
	c.rdb.Del(ctx, keys...)
}

// ZIncr 榜单 ZSet 增加分数。
func (c *Cache) ZIncr(ctx context.Context, key string, member string, delta float64) {
	if !c.Enabled() {
		return
	}
	c.rdb.ZIncrBy(ctx, key, delta, member)
}

// ZTop 按分数倒序取排名，返回 member 与 score。
func (c *Cache) ZTop(ctx context.Context, key string, offset, count int) []redis.Z {
	if !c.Enabled() {
		return nil
	}
	res, err := c.rdb.ZRevRangeWithScores(ctx, key, int64(offset), int64(offset+count-1)).Result()
	if err != nil {
		return nil
	}
	return res
}

// ZCard 榜单元素数量。
func (c *Cache) ZCard(ctx context.Context, key string) int64 {
	if !c.Enabled() {
		return 0
	}
	n, _ := c.rdb.ZCard(ctx, key).Result()
	return n
}

// Incr 计数 +1 并返回最新值。
func (c *Cache) Incr(ctx context.Context, key string) int64 {
	if !c.Enabled() {
		return 0
	}
	n, _ := c.rdb.Incr(ctx, key).Result()
	return n
}

// Decr 计数 -1（不低于 0）。
func (c *Cache) Decr(ctx context.Context, key string) {
	if !c.Enabled() {
		return
	}
	if n, _ := c.rdb.Get(ctx, key).Int(); n <= 0 {
		return
	}
	c.rdb.Decr(ctx, key)
}

// GetInt 读取整型计数。
func (c *Cache) GetInt(ctx context.Context, key string) int64 {
	if !c.Enabled() {
		return 0
	}
	n, _ := c.rdb.Get(ctx, key).Int64()
	return n
}

// SetInt 设置整型值。
func (c *Cache) SetInt(ctx context.Context, key string, val int64, ttl time.Duration) {
	if !c.Enabled() {
		return
	}
	c.rdb.Set(ctx, key, val, ttl)
}
