package api

import (
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// SearchAPI 搜索接口。
type SearchAPI struct {
	svc *service.SearchService
}

func NewSearchAPI(svc *service.SearchService) *SearchAPI {
	return &SearchAPI{svc: svc}
}

// Search 综合搜索
// @Summary  综合搜索
// @Tags     搜索
// @Produce  json
// @Param    keyword   query  string  false  "搜索关键词"
// @Param    type      query  string  false  "搜索类型"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /search [get]
func (h *SearchAPI) Search(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.Search(c.Query("keyword"), c.Query("type"), page, size, middleware.CurrentUserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// Suggest 搜索联想建议
// @Summary  搜索联想建议
// @Tags     搜索
// @Produce  json
// @Param    keyword  query  string  false  "搜索关键词"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /search/suggest [get]
func (h *SearchAPI) Suggest(c *gin.Context) {
	res, err := h.svc.Suggest(c.Query("keyword"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// HotKeywords 热门搜索词
// @Summary  热门搜索词
// @Tags     搜索
// @Produce  json
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /search/hot-keywords [get]
func (h *SearchAPI) HotKeywords(c *gin.Context) {
	response.OK(c, h.svc.HotKeywords())
}
