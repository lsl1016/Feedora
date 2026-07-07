package middleware

import "github.com/gin-gonic/gin"

// 存入 gin.Context 的键名。
const (
	CtxUserID = "userId"
	CtxRole   = "role"
	CtxTrace  = "traceId"
)

// CurrentUserID 从上下文读取当前登录用户 ID，未登录返回 0。
func CurrentUserID(c *gin.Context) int64 {
	if v, ok := c.Get(CtxUserID); ok {
		if id, ok := v.(int64); ok {
			return id
		}
	}
	return 0
}

// CurrentRole 从上下文读取当前用户角色。
func CurrentRole(c *gin.Context) string {
	if v, ok := c.Get(CtxRole); ok {
		if r, ok := v.(string); ok {
			return r
		}
	}
	return ""
}
