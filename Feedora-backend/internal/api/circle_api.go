package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	errs "github.com/feedora/backend/pkg/errors"
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
// @Param    scope     query  string  false  "范围"
// @Param    keyword   query  string  false  "关键词"
// @Param    category  query  string  false  "分类"
// @Param    sort      query  string  false  "排序"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /circles [get]
func (h *CircleAPI) List(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.List(
		c.DefaultQuery("scope", "all"), c.Query("keyword"), c.Query("category"), c.Query("sort"),
		middleware.CurrentUserID(c), page, size,
	)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// Get 获取圈子详情
// @Summary  获取圈子详情
// @Tags     圈子
// @Produce  json
// @Param    circleId  path  int  true  "圈子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /circles/{circleId} [get]
func (h *CircleAPI) Get(c *gin.Context) {
	res, err := h.svc.Get(paramID(c, "circleId"), middleware.CurrentUserID(c))
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
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles [post]
func (h *CircleAPI) Create(c *gin.Context) {
	var in dto.CreateCircleRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, errs.ErrParams)
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
// @Param    circleId  path  int  true  "圈子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/join [post]
func (h *CircleAPI) Join(c *gin.Context) {
	if err := h.svc.Join(paramID(c, "circleId"), middleware.CurrentUserID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Leave 退出圈子
// @Summary  退出圈子
// @Tags     圈子
// @Produce  json
// @Param    circleId  path  int  true  "圈子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/leave [post]
func (h *CircleAPI) Leave(c *gin.Context) {
	if err := h.svc.Leave(paramID(c, "circleId"), middleware.CurrentUserID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Members 获取圈子成员列表
// @Summary  获取圈子成员列表
// @Tags     圈子
// @Produce  json
// @Param    circleId  path   int  true   "圈子ID"
// @Param    page      query  int  false  "页码"
// @Param    pageSize  query  int  false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /circles/{circleId}/members [get]
func (h *CircleAPI) Members(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.Members(paramID(c, "circleId"), page, size)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// Posts 获取圈子帖子列表
// @Summary  获取圈子帖子列表
// @Tags     圈子
// @Produce  json
// @Param    circleId  path   int     true   "圈子ID"
// @Param    sort      query  string  false  "排序"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /circles/{circleId}/posts [get]
func (h *CircleAPI) Posts(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.Posts(paramID(c, "circleId"), middleware.CurrentUserID(c), c.Query("sort"), page, size)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// SetRole 设置成员角色
// @Summary  设置成员角色
// @Tags     圈子
// @Accept   json
// @Produce  json
// @Param    circleId  path  int  true  "圈子ID"
// @Param    userId    path  int  true  "用户ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/members/{userId}/role [put]
func (h *CircleAPI) SetRole(c *gin.Context) {
	var in struct {
		Role string `json:"role"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := h.svc.SetRole(paramID(c, "circleId"), paramID(c, "userId"), middleware.CurrentUserID(c), in.Role); err != nil {
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
// @Param    circleId  path  int  true  "圈子ID"
// @Param    userId    path  int  true  "用户ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/members/{userId}/mute [put]
func (h *CircleAPI) Mute(c *gin.Context) {
	var in struct {
		Duration int    `json:"duration"`
		Reason   string `json:"reason"`
	}
	_ = c.ShouldBindJSON(&in)
	if err := h.svc.Mute(paramID(c, "circleId"), paramID(c, "userId"), middleware.CurrentUserID(c), in.Duration, in.Reason); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Unmute 解除禁言成员
// @Summary  解除禁言成员
// @Tags     圈子
// @Produce  json
// @Param    circleId  path  int  true  "圈子ID"
// @Param    userId    path  int  true  "用户ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/members/{userId}/unmute [put]
func (h *CircleAPI) Unmute(c *gin.Context) {
	if err := h.svc.Unmute(paramID(c, "circleId"), paramID(c, "userId"), middleware.CurrentUserID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Remove 移除圈子成员
// @Summary  移除圈子成员
// @Tags     圈子
// @Produce  json
// @Param    circleId  path  int  true  "圈子ID"
// @Param    userId    path  int  true  "用户ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /circles/{circleId}/members/{userId} [delete]
func (h *CircleAPI) Remove(c *gin.Context) {
	if err := h.svc.Remove(paramID(c, "circleId"), paramID(c, "userId"), middleware.CurrentUserID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}
