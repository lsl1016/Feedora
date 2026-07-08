package api

import (
	"github.com/feedora/backend/internal/dto"
	errspkg "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// normalizePageRequest 对分页参数补默认值并限制上限。
func normalizePageRequest(req *dto.PageRequest) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}
}

// bindJSON 绑定 JSON 请求体，失败时统一返回参数错误。
func bindJSON(c *gin.Context, req any) bool {
	if err := c.ShouldBindJSON(req); err != nil {
		response.Fail(c, errspkg.ErrParams)
		return false
	}
	return true
}

// bindQuery 绑定 query 参数，失败时统一返回参数错误。
func bindQuery(c *gin.Context, req any) bool {
	if err := c.ShouldBindQuery(req); err != nil {
		response.Fail(c, errspkg.ErrParams)
		return false
	}
	return true
}

// bindURI 绑定路径参数，失败时统一返回参数错误。
func bindURI(c *gin.Context, req any) bool {
	if err := c.ShouldBindUri(req); err != nil {
		response.Fail(c, errspkg.ErrParams)
		return false
	}
	return true
}
