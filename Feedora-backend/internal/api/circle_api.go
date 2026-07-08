package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// CircleAPI 圈子接口。
type CircleAPI struct {
	svc *service.CircleService
}

func NewCircleAPI(svc *service.CircleService) *CircleAPI {
	return &CircleAPI{svc: svc}
}

// List 获取圈子列表
// @Summary  获取圈子列表
// @Tags     圈子
// @Produce  json
// @Param    req  query  dto.CircleListQuery  false  "圈子列表查询参数"
// @Success  200  {object}  dto.CirclePageResponse
// @Failure  400  {object}  response.Body
// @Router   /circles [get]
func (h *CircleAPI) List(c *gin.Context) {
	var req dto.CircleListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	scope := req.Scope
	if scope == "" {
		scope = "all"
	}
	list, total, err := h.svc.List(scope, req.Keyword, req.Category, req.Sort, middleware.CurrentUserID(c), req.Page, req.PageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// Get 获取圈子详情
// @Summary  获取圈子详情
// @Tags     圈子
// @Produce  json
// @Param    req  path  dto.CircleIDURI  true  "圈子路径参数"
// @Success  200  {object}  dto.CircleResponse
// @Failure  400  {object}  response.Body
// @Router   /circles/{circleId} [get]
func (h *CircleAPI) Get(c *gin.Context) {
	var req dto.CircleIDURI
	if !bindURI(c, &req) {
		return
	}
	res, err := h.svc.Get(req.CircleID, middleware.CurrentUserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Create 创建圈子
// @Summary  创建圈子
// @Tags     圈子
// @Accept   json
// @Produce  json
// @Param    body  body  dto.CreateCircleRequest  true  "创建圈子请求体"
// @Success  200  {object}  dto.CircleResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles [post]
func (h *CircleAPI) Create(c *gin.Context) {
	var in dto.CreateCircleRequest
	if !bindJSON(c, &in) {
		return
	}
	res, err := h.svc.Create(middleware.CurrentUserID(c), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Join 加入圈子
// @Summary  加入圈子
// @Tags     圈子
// @Produce  json
// @Param    req  path  dto.CircleIDURI  true  "圈子路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/join [post]
func (h *CircleAPI) Join(c *gin.Context) {
	var req dto.CircleIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.Join(req.CircleID, middleware.CurrentUserID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Leave 退出圈子
// @Summary  退出圈子
// @Tags     圈子
// @Produce  json
// @Param    req  path  dto.CircleIDURI  true  "圈子路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/leave [post]
func (h *CircleAPI) Leave(c *gin.Context) {
	var req dto.CircleIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.Leave(req.CircleID, middleware.CurrentUserID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Members 获取圈子成员列表
// @Summary  获取圈子成员列表
// @Tags     圈子
// @Produce  json
// @Param    circleId  path   dto.CircleIDURI  true   "圈子路径参数"
// @Param    req       query  dto.PageRequest  false  "分页查询参数"
// @Success  200  {object}  dto.CircleMemberPageResponse
// @Failure  400  {object}  response.Body
// @Router   /circles/{circleId}/members [get]
func (h *CircleAPI) Members(c *gin.Context) {
	var uri dto.CircleIDURI
	if !bindURI(c, &uri) {
		return
	}
	var req dto.PageRequest
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req)
	list, total, err := h.svc.Members(uri.CircleID, req.Page, req.PageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// Posts 获取圈子帖子列表
// @Summary  获取圈子帖子列表
// @Tags     圈子
// @Produce  json
// @Param    circleId  path   dto.CircleIDURI         true   "圈子路径参数"
// @Param    req       query  dto.PostListByTopicQuery false  "帖子列表查询参数"
// @Success  200  {object}  dto.PostPageResponse
// @Failure  400  {object}  response.Body
// @Router   /circles/{circleId}/posts [get]
func (h *CircleAPI) Posts(c *gin.Context) {
	var uri dto.CircleIDURI
	if !bindURI(c, &uri) {
		return
	}
	var req dto.PostListByTopicQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total, err := h.svc.Posts(uri.CircleID, middleware.CurrentUserID(c), req.Sort, req.Page, req.PageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// SetRole 设置成员角色
// @Summary  设置成员角色
// @Tags     圈子
// @Accept   json
// @Produce  json
// @Param    path  path  dto.UserAndCircleURI            true  "成员路径参数"
// @Param    body  body  dto.SetCircleMemberRoleRequest  true  "设置角色请求体"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/members/{userId}/role [put]
func (h *CircleAPI) SetRole(c *gin.Context) {
	var uri dto.UserAndCircleURI
	if !bindURI(c, &uri) {
		return
	}
	var in dto.SetCircleMemberRoleRequest
	if !bindJSON(c, &in) {
		return
	}
	if err := h.svc.SetRole(uri.CircleID, uri.UserID, middleware.CurrentUserID(c), in.Role); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Mute 禁言成员
// @Summary  禁言成员
// @Tags     圈子
// @Accept   json
// @Produce  json
// @Param    path  path  dto.UserAndCircleURI          true  "成员路径参数"
// @Param    body  body  dto.MuteCircleMemberRequest   true  "禁言请求体"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/members/{userId}/mute [put]
func (h *CircleAPI) Mute(c *gin.Context) {
	var uri dto.UserAndCircleURI
	if !bindURI(c, &uri) {
		return
	}
	var in dto.MuteCircleMemberRequest
	if !bindJSON(c, &in) {
		return
	}
	if err := h.svc.Mute(uri.CircleID, uri.UserID, middleware.CurrentUserID(c), in.Duration, in.Reason); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Unmute 解除禁言成员
// @Summary  解除禁言成员
// @Tags     圈子
// @Produce  json
// @Param    req  path  dto.UserAndCircleURI  true  "成员路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/members/{userId}/unmute [put]
func (h *CircleAPI) Unmute(c *gin.Context) {
	var req dto.UserAndCircleURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.Unmute(req.CircleID, req.UserID, middleware.CurrentUserID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Remove 移除圈子成员
// @Summary  移除圈子成员
// @Tags     圈子
// @Produce  json
// @Param    req  path  dto.UserAndCircleURI  true  "成员路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/members/{userId} [delete]
func (h *CircleAPI) Remove(c *gin.Context) {
	var req dto.UserAndCircleURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.Remove(req.CircleID, req.UserID, middleware.CurrentUserID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}
