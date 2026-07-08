package api

import (
	"github.com/feedora/backend/internal/dto"
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
// @Param    req  query  dto.TopicListQuery  false  "话题列表请求参数"
// @Success  200  {object}  dto.TopicPageResponse
// @Failure  400  {object}  response.Body
// @Router   /topics [get]
func (h *TopicAPI) List(c *gin.Context) {
	var req dto.TopicListQuery
	if !bindQuery(c, &req) {
		return
	}
	page, size := req.Page, req.PageSize
	normalizePageRequest(&req.PageRequest)
	page, size = req.Page, req.PageSize
	tab := req.Tab
	if tab == "" {
		tab = "all"
	}
	list, total, err := h.svc.List(tab, page, size)
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
// @Param    req  path  dto.TopicIDURI  true  "话题路径参数"
// @Success  200  {object}  dto.TopicResponse
// @Failure  400  {object}  response.Body
// @Router   /topics/{topicId} [get]
func (h *TopicAPI) Get(c *gin.Context) {
	var req dto.TopicIDURI
	if !bindURI(c, &req) {
		return
	}
	res, err := h.svc.Get(req.TopicID)
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
// @Param    topicId  path   dto.TopicIDURI      true   "话题路径参数"
// @Param    req      query  dto.TopicPostsQuery  false  "帖子列表查询参数"
// @Success  200  {object}  dto.PostPageResponse
// @Failure  400  {object}  response.Body
// @Router   /topics/{topicId}/posts [get]
func (h *TopicAPI) Posts(c *gin.Context) {
	var uri dto.TopicIDURI
	if !bindURI(c, &uri) {
		return
	}
	var req dto.TopicPostsQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total, err := h.svc.Posts(uri.TopicID, middleware.CurrentUserID(c), req.Sort, req.Page, req.PageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}
