package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// AdminAPI 后台管理接口。
type AdminAPI struct {
	svc *service.AdminService
}

func NewAdminAPI(svc *service.AdminService) *AdminAPI {
	return &AdminAPI{svc: svc}
}

// Users 用户列表
// @Summary  用户列表
// @Tags     后台管理
// @Produce  json
// @Param    req  query  dto.AdminUserListQuery  false  "后台用户列表查询参数"
// @Success  200  {object}  dto.UserListResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/users [get]
func (h *AdminAPI) Users(c *gin.Context) {
	var req dto.AdminUserListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total, err := h.svc.Users(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// Posts 帖子列表
// @Summary  帖子列表
// @Tags     后台管理
// @Produce  json
// @Param    req  query  dto.AdminPostListQuery  false  "后台帖子列表查询参数"
// @Success  200  {object}  dto.AdminPostPageResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/posts [get]
func (h *AdminAPI) Posts(c *gin.Context) {
	var req dto.AdminPostListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total := h.svc.Posts(req.Status, req.Page, req.PageSize)
	response.Page(c, list, total, req.Page, req.PageSize)
}

// Comments 评论列表
// @Summary  评论列表
// @Tags     后台管理
// @Produce  json
// @Param    req  query  dto.AdminCommentListQuery  false  "后台评论列表查询参数"
// @Success  200  {object}  dto.AdminCommentPageResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/comments [get]
func (h *AdminAPI) Comments(c *gin.Context) {
	var req dto.AdminCommentListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total := h.svc.Comments(req.Page, req.PageSize)
	response.Page(c, list, total, req.Page, req.PageSize)
}

// Tags 标签列表
// @Summary  标签列表
// @Tags     后台管理
// @Produce  json
// @Success  200  {object}  dto.TagListResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/tags [get]
func (h *AdminAPI) Tags(c *gin.Context) {
	list, err := h.svc.Tags()
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

// CreateTag 创建标签
// @Summary  创建标签
// @Tags     后台管理
// @Accept   json
// @Produce  json
// @Param    body  body  dto.CreateTagRequest  true  "创建标签请求体"
// @Success  200  {object}  dto.TagResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/tags [post]
func (h *AdminAPI) CreateTag(c *gin.Context) {
	var in dto.CreateTagRequest
	if !bindJSON(c, &in) {
		return
	}
	res, err := h.svc.CreateTag(in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// UpdateTag 更新标签
// @Summary  更新标签
// @Tags     后台管理
// @Accept   json
// @Produce  json
// @Param    tagId  path  dto.TagIDURI           true  "标签路径参数"
// @Param    body   body  dto.UpdateTagRequest   true  "更新标签请求体"
// @Success  200  {object}  dto.TagResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/tags/{tagId} [put]
func (h *AdminAPI) UpdateTag(c *gin.Context) {
	var uri dto.TagIDURI
	if !bindURI(c, &uri) {
		return
	}
	var in dto.UpdateTagRequest
	if !bindJSON(c, &in) {
		return
	}
	res, err := h.svc.UpdateTag(uri.TagID, in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Topics 话题列表
// @Summary  话题列表
// @Tags     后台管理
// @Produce  json
// @Param    req  query  dto.AdminTopicListQuery  false  "后台话题列表查询参数"
// @Success  200  {object}  dto.TopicPageResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/topics [get]
func (h *AdminAPI) Topics(c *gin.Context) {
	var req dto.AdminTopicListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total, err := h.svc.Topics(req.Page, req.PageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// CreateTopic 创建话题
// @Summary  创建话题
// @Tags     后台管理
// @Accept   json
// @Produce  json
// @Param    body  body  dto.CreateTopicRequest  true  "创建话题请求体"
// @Success  200  {object}  dto.TopicResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/topics [post]
func (h *AdminAPI) CreateTopic(c *gin.Context) {
	var in dto.CreateTopicRequest
	if !bindJSON(c, &in) {
		return
	}
	res, err := h.svc.CreateTopic(in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// UpdateTopic 更新话题
// @Summary  更新话题
// @Tags     后台管理
// @Accept   json
// @Produce  json
// @Param    topicId  path  dto.TopicIDURI           true  "话题路径参数"
// @Param    body     body  dto.UpdateTopicRequest   true  "更新话题请求体"
// @Success  200  {object}  dto.TopicResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/topics/{topicId} [put]
func (h *AdminAPI) UpdateTopic(c *gin.Context) {
	var uri dto.TopicIDURI
	if !bindURI(c, &uri) {
		return
	}
	var in dto.UpdateTopicRequest
	if !bindJSON(c, &in) {
		return
	}
	res, err := h.svc.UpdateTopic(uri.TopicID, in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Circles 圈子列表
// @Summary  圈子列表
// @Tags     后台管理
// @Produce  json
// @Param    req  query  dto.AdminCircleListQuery  false  "后台圈子列表查询参数"
// @Success  200  {object}  dto.AdminCirclePageResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/circles [get]
func (h *AdminAPI) Circles(c *gin.Context) {
	var req dto.AdminCircleListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total := h.svc.Circles(req.Page, req.PageSize)
	response.Page(c, list, total, req.Page, req.PageSize)
}

// Stats 仪表盘统计
// @Summary  仪表盘统计
// @Tags     后台管理
// @Produce  json
// @Success  200  {object}  dto.AdminStatsResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/dashboard/stats [get]
func (h *AdminAPI) Stats(c *gin.Context) {
	response.OK(c, h.svc.Stats())
}

// Logs 操作日志列表
// @Summary  操作日志列表
// @Tags     后台管理
// @Produce  json
// @Param    req  query  dto.AdminLogListQuery  false  "后台操作日志列表查询参数"
// @Success  200  {object}  dto.AdminOperationLogPageResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/operation-logs [get]
func (h *AdminAPI) Logs(c *gin.Context) {
	var req dto.AdminLogListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total := h.svc.Logs(req.Page, req.PageSize)
	response.Page(c, list, total, req.Page, req.PageSize)
}
