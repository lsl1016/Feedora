package api

import (
	"github.com/feedora/backend/internal/dto"
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
// @Param    req  query  dto.SearchQuery  false  "综合搜索查询参数"
// @Success  200  {object}  dto.SearchResultPageResponse
// @Failure  400  {object}  response.Body
// @Router   /search [get]
func (h *SearchAPI) Search(c *gin.Context) {
	var req dto.SearchQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total, err := h.svc.Search(req.Keyword, req.Type, req.Page, req.PageSize, middleware.CurrentUserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// Suggest 搜索联想建议
// @Summary  搜索联想建议
// @Tags     搜索
// @Produce  json
// @Param    req  query  dto.SearchSuggestQuery  false  "搜索联想查询参数"
// @Success  200  {object}  dto.SearchSuggestResponse
// @Failure  400  {object}  response.Body
// @Router   /search/suggest [get]
func (h *SearchAPI) Suggest(c *gin.Context) {
	var req dto.SearchSuggestQuery
	if !bindQuery(c, &req) {
		return
	}
	res, err := h.svc.Suggest(req.Keyword)
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
// @Success  200  {object}  dto.HotKeywordResponse
// @Failure  400  {object}  response.Body
// @Router   /search/hot-keywords [get]
func (h *SearchAPI) HotKeywords(c *gin.Context) {
	response.OK(c, h.svc.HotKeywords())
}
