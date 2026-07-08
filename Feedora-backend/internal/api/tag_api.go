package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// TagAPI 标签接口。
type TagAPI struct {
	svc *service.TagService
}

func NewTagAPI(svc *service.TagService) *TagAPI {
	return &TagAPI{svc: svc}
}

// List 标签列表
// @Summary  标签列表
// @Tags     标签
// @Produce  json
// @Success  200  {object}  dto.TagListResponse
// @Failure  400  {object}  response.Body
// @Router   /tags [get]
func (h *TagAPI) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

// Get 标签详情
// @Summary  标签详情
// @Tags     标签
// @Produce  json
// @Param    req  path  dto.TagIDURI  true  "标签路径参数"
// @Success  200  {object}  dto.TagResponse
// @Failure  400  {object}  response.Body
// @Router   /tags/{tagId} [get]
func (h *TagAPI) Get(c *gin.Context) {
	var req dto.TagIDURI
	if !bindURI(c, &req) {
		return
	}
	res, err := h.svc.Get(req.TagID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Posts 标签下的帖子列表
// @Summary  标签下的帖子列表
// @Tags     标签
// @Produce  json
// @Param    tagId  path   dto.TagIDURI      true   "标签路径参数"
// @Param    req    query  dto.TagPostsQuery false  "帖子列表查询参数"
// @Success  200  {object}  dto.PostPageResponse
// @Failure  400  {object}  response.Body
// @Router   /tags/{tagId}/posts [get]
func (h *TagAPI) Posts(c *gin.Context) {
	var uri dto.TagIDURI
	if !bindURI(c, &uri) {
		return
	}
	var req dto.TagPostsQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total, err := h.svc.Posts(uri.TagID, middleware.CurrentUserID(c), req.Sort, req.Page, req.PageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}
