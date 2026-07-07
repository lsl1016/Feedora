package middleware

import (
	"time"

	"github.com/feedora/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

// slowRequestMs 慢请求阈值（毫秒），超过则以 WARN 级别记录。
const slowRequestMs = 500

// Logger 记录请求方法、路径、状态码与耗时。
// 分级：5xx 记 ERROR，慢请求(>500ms)记 WARN，其余记 INFO，便于线上排查。
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		cost := time.Since(start).Milliseconds()
		status := c.Writer.Status()
		format := "request completed, traceId:%s method:%s path:%s status:%d durationMs:%d"
		args := []interface{}{c.GetString(CtxTrace), c.Request.Method, c.Request.URL.Path, status, cost}
		switch {
		case status >= 500:
			logger.Errorf(format, args...)
		case cost >= slowRequestMs:
			logger.Warnf("slow "+format, args...)
		default:
			logger.Infof(format, args...)
		}
	}
}
