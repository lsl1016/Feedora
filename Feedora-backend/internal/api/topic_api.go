package api

import (
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// TopicAPI 话题接口。
type TopicAPI struct {
	svc *service.TopicService
}

func NewTopicAPI(svc *service.TopicService) *TopicAPI {
	return &TopicAPI{svc: svc}
}

// List 话题列表
// @Summary  话题列表
// @Tags     话题
// @Produce  json
// @Param    tab       query  string  false  "分类标签（默认 all）"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /topics [get]
func (h *TopicAPI) List(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.List(c.DefaultQuery("tab", "all"), page, size)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// Get 话题详情
// @Summary  话题详情
// @Tags     话题
// @Produce  json
// @Param    topicId  path  int  true  "话题ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /topics/{topicId} [get]
func (h *TopicAPI) Get(c *gin.Context) {
	res, err := h.svc.Get(paramID(c, "topicId"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Posts 话题下的帖子列表
// @Summary  话题下的帖子列表
// @Tags     话题
// @Produce  json
// @Param    topicId   path   int     true   "话题ID"
// @Param    sort      query  string  false  "排序方式"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /topics/{topicId}/posts [get]
func (h *TopicAPI) Posts(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.Posts(paramID(c, "topicId"), middleware.CurrentUserID(c), c.Query("sort"), page, size)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}
