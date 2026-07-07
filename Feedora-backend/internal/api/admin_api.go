package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	errs "github.com/feedora/backend/pkg/errors"
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
// @Param    keyword   query  string  false  "关键词"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/users [get]
func (h *AdminAPI) Users(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.Users(c.Query("keyword"), page, size)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// Posts 帖子列表
// @Summary  帖子列表
// @Tags     后台管理
// @Produce  json
// @Param    status    query  string  false  "帖子状态"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/posts [get]
func (h *AdminAPI) Posts(c *gin.Context) {
	page, size := pageParams(c)
	list, total := h.svc.Posts(c.Query("status"), page, size)
	response.Page(c, list, total, page, size)
}

// Comments 评论列表
// @Summary  评论列表
// @Tags     后台管理
// @Produce  json
// @Param    page      query  int  false  "页码"
// @Param    pageSize  query  int  false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/comments [get]
func (h *AdminAPI) Comments(c *gin.Context) {
	page, size := pageParams(c)
	list, total := h.svc.Comments(page, size)
	response.Page(c, list, total, page, size)
}

// Tags 标签列表
// @Summary  标签列表
// @Tags     后台管理
// @Produce  json
// @Success  200  {object}  response.Body
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
// @Param    body  body  dto.CreateTagRequest  true  "请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/tags [post]
func (h *AdminAPI) CreateTag(c *gin.Context) {
	var in dto.CreateTagRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, errs.ErrParams)
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
// @Param    tagId  path  int  true  "标签ID"
// @Param    body   body  dto.UpdateTagRequest  true  "请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/tags/{tagId} [put]
func (h *AdminAPI) UpdateTag(c *gin.Context) {
	var in dto.UpdateTagRequest
	_ = c.ShouldBindJSON(&in)
	res, err := h.svc.UpdateTag(paramID(c, "tagId"), in)
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
// @Param    page      query  int  false  "页码"
// @Param    pageSize  query  int  false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/topics [get]
func (h *AdminAPI) Topics(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.Topics(page, size)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// CreateTopic 创建话题
// @Summary  创建话题
// @Tags     后台管理
// @Accept   json
// @Produce  json
// @Param    body  body  dto.CreateTopicRequest  true  "请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/topics [post]
func (h *AdminAPI) CreateTopic(c *gin.Context) {
	var in dto.CreateTopicRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, errs.ErrParams)
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
// @Param    topicId  path  int     true  "话题ID"
// @Param    body     body  object  true  "请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/topics/{topicId} [put]
func (h *AdminAPI) UpdateTopic(c *gin.Context) {
	var in map[string]any
	_ = c.ShouldBindJSON(&in)
	res, err := h.svc.UpdateTopic(paramID(c, "topicId"), in)
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
// @Param    page      query  int  false  "页码"
// @Param    pageSize  query  int  false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/circles [get]
func (h *AdminAPI) Circles(c *gin.Context) {
	page, size := pageParams(c)
	list, total := h.svc.Circles(page, size)
	response.Page(c, list, total, page, size)
}

// Stats 仪表盘统计
// @Summary  仪表盘统计
// @Tags     后台管理
// @Produce  json
// @Success  200  {object}  response.Body
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
// @Param    page      query  int  false  "页码"
// @Param    pageSize  query  int  false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /admin/operation-logs [get]
func (h *AdminAPI) Logs(c *gin.Context) {
	page, size := pageParams(c)
	list, total := h.svc.Logs(page, size)
	response.Page(c, list, total, page, size)
}
