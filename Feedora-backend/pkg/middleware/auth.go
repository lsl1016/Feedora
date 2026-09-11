package middleware

import (
	"strings"

	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/jwtx"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AuthChecker 鉴权链的补充校验：token 黑名单、用户级吊销、用户状态。
// 返回非 nil error 表示拒绝该请求；可为 nil（跳过补充校验，仅验签）。
type AuthChecker func(claims *jwtx.Claims, token string) error

// Auth 校验 JWT 与补充规则，将 userId、role 写入上下文。未登录或校验失败返回 401/403。
func Auth(jm *jwtx.Manager, checker AuthChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c)
		claims, ok := parseToken(token, jm)
		if !ok {
			response.Fail(c, errs.ErrUnauth)
			c.Abort()
			return
		}
		if checker != nil {
			if err := checker(claims, token); err != nil {
				response.Fail(c, err)
				c.Abort()
				return
			}
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

// OptionalAuth 尝试解析 JWT，解析并校验成功则写入上下文，否则按匿名放行。
// 用于既支持匿名访问、又能识别登录用户的接口（如帖子列表、详情）。
// 校验失败（封禁 / 已吊销）的 token 按匿名处理，不写入身份。
func OptionalAuth(jm *jwtx.Manager, checker AuthChecker) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := bearerToken(c)
		claims, ok := parseToken(token, jm)
		if !ok {
			c.Next()
			return
		}
		if checker != nil {
			if err := checker(claims, token); err != nil {
				c.Next()
				return
			}
		}
		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxRole, claims.Role)
		c.Next()
	}
}

func bearerToken(c *gin.Context) string {
	return strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
}

func parseToken(token string, jm *jwtx.Manager) (*jwtx.Claims, bool) {
	if token == "" {
		return nil, false
	}
	claims, err := jm.Parse(token)
	if err != nil {
		return nil, false
	}
	return claims, true
}
