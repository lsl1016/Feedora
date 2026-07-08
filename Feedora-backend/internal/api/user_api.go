package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
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
// @Param    req  query  dto.UserListQuery  false  "用户列表查询参数"
// @Success  200  {object}  dto.UserListResponse
// @Failure  400  {object}  response.Body
// @Router   /users [get]
func (h *UserAPI) List(c *gin.Context) {
	var req dto.UserListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total, err := h.svc.List(req.Keyword, req.Page, req.PageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// Get 获取用户详情
// @Summary  获取用户详情
// @Tags     用户
// @Produce  json
// @Param    req  path  dto.UserIDURI  true  "用户路径参数"
// @Success  200  {object}  dto.UserResponse
// @Failure  400  {object}  response.Body
// @Router   /users/{userId} [get]
func (h *UserAPI) Get(c *gin.Context) {
	var req dto.UserIDURI
	if !bindURI(c, &req) {
		return
	}
	res, err := h.svc.Get(req.UserID)
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
// @Success  200  {object}  dto.UserResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/profile [put]
func (h *UserAPI) UpdateProfile(c *gin.Context) {
	var in dto.UpdateProfileRequest
	if !bindJSON(c, &in) {
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
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Success  200  {object}  dto.PostPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/posts [get]
func (h *UserAPI) MyPosts(c *gin.Context) {
	var req dto.PageRequest
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req)
	list, total, err := h.svc.MyPosts(middleware.CurrentUserID(c), req.Page, req.PageSize)
	h.pageOrFail(c, list, total, req.Page, req.PageSize, err)
}

// MyComments 获取我的评论列表
// @Summary  获取我的评论列表
// @Tags     用户
// @Produce  json
// @Security BearerAuth
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Success  200  {object}  dto.MyCommentPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/comments [get]
func (h *UserAPI) MyComments(c *gin.Context) {
	var req dto.PageRequest
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req)
	list, total, err := h.svc.MyComments(middleware.CurrentUserID(c), req.Page, req.PageSize)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// MyLikedPosts 获取我点赞的帖子列表
// @Summary  获取我点赞的帖子列表
// @Tags     用户
// @Produce  json
// @Security BearerAuth
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Success  200  {object}  dto.PostPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/liked-posts [get]
func (h *UserAPI) MyLikedPosts(c *gin.Context) {
	var req dto.PageRequest
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req)
	list, total, err := h.svc.MyLikedPosts(middleware.CurrentUserID(c), req.Page, req.PageSize)
	h.pageOrFail(c, list, total, req.Page, req.PageSize, err)
}

// MyFavoritePosts 获取我收藏的帖子列表
// @Summary  获取我收藏的帖子列表
// @Tags     用户
// @Produce  json
// @Security BearerAuth
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Success  200  {object}  dto.PostPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/favorite-posts [get]
func (h *UserAPI) MyFavoritePosts(c *gin.Context) {
	var req dto.PageRequest
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req)
	list, total, err := h.svc.MyFavoritePosts(middleware.CurrentUserID(c), req.Page, req.PageSize)
	h.pageOrFail(c, list, total, req.Page, req.PageSize, err)
}

func (h *UserAPI) pageOrFail(c *gin.Context, list []dto.Post, total int64, page, size int, err error) {
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}
