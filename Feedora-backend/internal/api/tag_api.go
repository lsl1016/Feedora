package api

import (
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
// @Success  200  {object}  response.Body
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
// @Param    tagId  path  int  true  "标签ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /tags/{tagId} [get]
func (h *TagAPI) Get(c *gin.Context) {
	res, err := h.svc.Get(paramID(c, "tagId"))
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
// @Param    tagId     path   int     true   "标签ID"
// @Param    sort      query  string  false  "排序方式"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /tags/{tagId}/posts [get]
func (h *TagAPI) Posts(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.Posts(paramID(c, "tagId"), middleware.CurrentUserID(c), c.Query("sort"), page, size)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}
