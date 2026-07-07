package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// UserAPI 用户接口。
type UserAPI struct {
	svc *service.UserService
}

func NewUserAPI(svc *service.UserService) *UserAPI {
	return &UserAPI{svc: svc}
}

// List 获取用户列表
// @Summary  获取用户列表
// @Tags     用户
// @Produce  json
// @Param    keyword   query  string  false  "搜索关键词"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /users [get]
func (h *UserAPI) List(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.List(c.Query("keyword"), page, size)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// Get 获取用户详情
// @Summary  获取用户详情
// @Tags     用户
// @Produce  json
// @Param    userId  path  int  true  "用户ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /users/{userId} [get]
func (h *UserAPI) Get(c *gin.Context) {
	res, err := h.svc.Get(paramID(c, "userId"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// UpdateProfile 更新当前用户资料
// @Summary  更新当前用户资料
// @Tags     用户
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    body  body  dto.UpdateProfileRequest  true  "资料请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /users/me/profile [put]
func (h *UserAPI) UpdateProfile(c *gin.Context) {
	var in dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, errs.ErrParams)
		return
	}
	res, err := h.svc.UpdateProfile(middleware.CurrentUserID(c), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// MyPosts 获取我的帖子列表
// @Summary  获取我的帖子列表
// @Tags     用户
// @Produce  json
// @Security BearerAuth
// @Param    page      query  int  false  "页码"
// @Param    pageSize  query  int  false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /users/me/posts [get]
func (h *UserAPI) MyPosts(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.MyPosts(middleware.CurrentUserID(c), page, size)
	h.pageOrFail(c, list, total, page, size, err)
}

// MyComments 获取我的评论列表
// @Summary  获取我的评论列表
// @Tags     用户
// @Produce  json
// @Security BearerAuth
// @Param    page      query  int  false  "页码"
// @Param    pageSize  query  int  false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /users/me/comments [get]
func (h *UserAPI) MyComments(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.MyComments(middleware.CurrentUserID(c), page, size)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// MyLikedPosts 获取我点赞的帖子列表
// @Summary  获取我点赞的帖子列表
// @Tags     用户
// @Produce  json
// @Security BearerAuth
// @Param    page      query  int  false  "页码"
// @Param    pageSize  query  int  false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /users/me/liked-posts [get]
func (h *UserAPI) MyLikedPosts(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.MyLikedPosts(middleware.CurrentUserID(c), page, size)
	h.pageOrFail(c, list, total, page, size, err)
}

// MyFavoritePosts 获取我收藏的帖子列表
// @Summary  获取我收藏的帖子列表
// @Tags     用户
// @Produce  json
// @Security BearerAuth
// @Param    page      query  int  false  "页码"
// @Param    pageSize  query  int  false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /users/me/favorite-posts [get]
func (h *UserAPI) MyFavoritePosts(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.MyFavoritePosts(middleware.CurrentUserID(c), page, size)
	h.pageOrFail(c, list, total, page, size, err)
}

func (h *UserAPI) pageOrFail(c *gin.Context, list []dto.Post, total int64, page, size int, err error) {
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}
