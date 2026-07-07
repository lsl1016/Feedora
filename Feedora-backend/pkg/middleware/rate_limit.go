package middleware

import "github.com/gin-gonic/gin"

// RateLimit 限流中间件占位。
// 阶段二基于 Redis 实现用户 / IP 维度限流（rate:api:{userId}:{path} 等），
// 阶段一直接放行。
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
