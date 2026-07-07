package middleware

import (
	"strings"

	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/jwtx"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// Auth 校验 JWT，将 userId、role 写入上下文。未登录直接返回 401。
func Auth(jm *jwtx.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, ok := parseToken(c, jm)
		if !ok {
			response.Fail(c, errs.ErrUnauth)
			c.Abort()
			return
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

// OptionalAuth 尝试解析 JWT，解析成功则写入上下文，失败也放行。
// 用于既支持匿名访问、又能识别登录用户的接口（如帖子列表、详情）。
func OptionalAuth(jm *jwtx.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if claims, ok := parseToken(c, jm); ok {
			c.Set(CtxUserID, claims.UserID)
			c.Set(CtxRole, claims.Role)
		}
		c.Next()
	}
}

func parseToken(c *gin.Context, jm *jwtx.Manager) (*jwtx.Claims, bool) {
	auth := c.GetHeader("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return nil, false
	}
	claims, err := jm.Parse(strings.TrimPrefix(auth, "Bearer "))
	if err != nil {
		return nil, false
	}
	return claims, true
}
