package middleware

import (
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// Admin 要求当前用户具备 admin 角色。
func Admin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if CurrentRole(c) != "admin" {
			response.Fail(c, errs.ErrForbidden)
			c.Abort()
			return
		}
		c.Next()
	}
}
