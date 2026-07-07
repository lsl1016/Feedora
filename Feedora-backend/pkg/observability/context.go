package observability

import "github.com/gin-gonic/gin"

// TraceID 从 gin.Context 读取 traceId（由 trace 中间件写入）。
func TraceID(c *gin.Context) string {
	if v, ok := c.Get("traceId"); ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}
