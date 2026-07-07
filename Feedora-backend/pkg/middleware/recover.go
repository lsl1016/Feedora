package middleware

import (
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/logger"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// Recover 捕获 panic，返回统一的 500 错误。
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("panic recovered: %v", err)
				response.Fail(c, errs.ErrInternal)
				c.Abort()
			}
		}()
		c.Next()
	}
}
