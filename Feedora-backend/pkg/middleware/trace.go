package middleware

import (
	"strings"

	"github.com/feedora/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// Trace 为每个请求确定 traceId 并写入上下文与响应头。
// 传播优先级：W3C traceparent > X-Trace-Id > 新生成（trace-<uuid>）。
// 若仅有 X-Request-Id，则作为 requestId 记录，并生成新的 traceId。
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		tid := traceIDFromHeaders(c)
		c.Set(CtxTrace, tid)
		if rid := c.GetHeader("X-Request-Id"); rid != "" {
			c.Set("requestId", rid)
		}
		c.Header("X-Trace-Id", tid)
		c.Next()
	}
}

// traceIDFromHeaders 依据入站请求头确定 traceId。
func traceIDFromHeaders(c *gin.Context) string {
	// W3C traceparent 格式：version-traceid-spanid-flags，取 trace-id 段。
	if tp := c.GetHeader("traceparent"); tp != "" {
		parts := strings.Split(tp, "-")
		if len(parts) >= 2 && len(parts[1]) == 32 {
			return "trace-" + parts[1]
		}
	}
	if xt := c.GetHeader("X-Trace-Id"); xt != "" {
		return xt
	}
	return "trace-" + utils.UUID()
}
