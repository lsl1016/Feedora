package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// FollowAPI 关注体系接口。
type FollowAPI struct {
	svc *service.FollowService
}

func NewFollowAPI(svc *service.FollowService) *FollowAPI {
	return &FollowAPI{svc: svc}
}

// Follow 关注用户
// @Summary  关注用户
// @Tags     关注
// @Produce  json
// @Security BearerAuth
// @Param    req  path  dto.UserIDURI  true  "用户路径参数"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /users/{userId}/follow [post]
func (h *FollowAPI) Follow(c *gin.Context) {
	var req dto.UserIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.Follow(middleware.CurrentUserID(c), req.UserID); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Unfollow 取消关注
// @Summary  取消关注用户
// @Tags     关注
// @Produce  json
// @Security BearerAuth
// @Param    req  path  dto.UserIDURI  true  "用户路径参数"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /users/{userId}/follow [delete]
func (h *FollowAPI) Unfollow(c *gin.Context) {
	var req dto.UserIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.Unfollow(middleware.CurrentUserID(c), req.UserID); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// State 查询关注状态
// @Summary  查询当前用户是否关注了目标用户
// @Tags     关注
// @Produce  json
// @Param    req  path  dto.UserIDURI  true  "用户路径参数"
// @Success  200  {object}  dto.FollowStateResponse
// @Failure  400  {object}  response.Body
// @Router   /users/{userId}/follow/state [get]
func (h *FollowAPI) State(c *gin.Context) {
	var req dto.UserIDURI
	if !bindURI(c, &req) {
		return
	}
	following, err := h.svc.State(middleware.CurrentUserID(c), req.UserID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, dto.FollowState{Following: following})
}

// Following 关注的人列表
// @Summary  获取用户关注的人列表
// @Tags     关注
// @Produce  json
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Param    userId  path  dto.UserIDURI  true  "用户路径参数"
// @Success  200  {object}  dto.FollowingPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/{userId}/following [get]
func (h *FollowAPI) Following(c *gin.Context) {
	var req dto.UserIDURI
	if !bindURI(c, &req) {
		return
	}
	var page dto.PageRequest
	if !bindQuery(c, &page) {
		return
	}
	normalizePageRequest(&page)
	list, total, err := h.svc.Following(req.UserID, page.Page, page.PageSize)
	h.pageOrFail(c, list, total, page.Page, page.PageSize, err)
}

// MyFollowing 我关注的人列表（静态路由，避免 /users/me 命中 :userId 参数校验）。
// @Summary  获取当前用户关注的人列表
// @Tags     关注
// @Produce  json
// @Security BearerAuth
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Success  200  {object}  dto.FollowingPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/following [get]
func (h *FollowAPI) MyFollowing(c *gin.Context) {
	var page dto.PageRequest
	if !bindQuery(c, &page) {
		return
	}
	normalizePageRequest(&page)
	list, total, err := h.svc.Following(middleware.CurrentUserID(c), page.Page, page.PageSize)
	h.pageOrFail(c, list, total, page.Page, page.PageSize, err)
}

// Followers 粉丝列表
// @Summary  获取用户粉丝列表
// @Tags     关注
// @Produce  json
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Param    userId  path  dto.UserIDURI  true  "用户路径参数"
// @Success  200  {object}  dto.FollowingPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/{userId}/followers [get]
func (h *FollowAPI) Followers(c *gin.Context) {
	var req dto.UserIDURI
	if !bindURI(c, &req) {
		return
	}
	var page dto.PageRequest
	if !bindQuery(c, &page) {
		return
	}
	normalizePageRequest(&page)
	list, total, err := h.svc.Followers(req.UserID, page.Page, page.PageSize)
	h.pageOrFail(c, list, total, page.Page, page.PageSize, err)
}

// MyFollowingFeed 我的关注动态
// @Summary  获取当前用户的关注动态
// @Tags     关注
// @Produce  json
// @Security BearerAuth
// @Param    req  query  dto.FollowingFeedQuery  false  "关注动态查询参数"
// @Success  200  {object}  dto.FollowingFeedPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/following-feed [get]
func (h *FollowAPI) MyFollowingFeed(c *gin.Context) {
	var req dto.FollowingFeedQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total, err := h.svc.FollowingFeed(middleware.CurrentUserID(c), req.FeedTab, req.Page, req.PageSize)
	h.pageOrFail(c, list, total, req.Page, req.PageSize, err)
}

// MyFollowingCircles 我关注的圈子（已加入）
// @Summary  获取当前用户加入的圈子列表
// @Tags     关注
// @Produce  json
// @Security BearerAuth
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Success  200  {object}  dto.CirclePageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/following-circles [get]
func (h *FollowAPI) MyFollowingCircles(c *gin.Context) {
	var page dto.PageRequest
	if !bindQuery(c, &page) {
		return
	}
	normalizePageRequest(&page)
	list, total, err := h.svc.FollowingCircles(middleware.CurrentUserID(c), page.Page, page.PageSize)
	h.pageOrFail(c, list, total, page.Page, page.PageSize, err)
}

// MyFollowingTopics 我关注的话题（参与过）
// @Summary  获取当前用户参与的话题列表
// @Tags     关注
// @Produce  json
// @Security BearerAuth
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Success  200  {object}  dto.TopicPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/following-topics [get]
func (h *FollowAPI) MyFollowingTopics(c *gin.Context) {
	var page dto.PageRequest
	if !bindQuery(c, &page) {
		return
	}
	normalizePageRequest(&page)
	list, total, err := h.svc.FollowingTopics(middleware.CurrentUserID(c), page.Page, page.PageSize)
	h.pageOrFail(c, list, total, page.Page, page.PageSize, err)
}

// MyFollowingTags 我关注的标签（使用过）
// @Summary  获取当前用户使用过的标签列表
// @Tags     关注
// @Produce  json
// @Security BearerAuth
// @Param    req  query  dto.PageRequest  false  "分页查询参数"
// @Success  200  {object}  dto.TagPageResponse
// @Failure  400  {object}  response.Body
// @Router   /users/me/following-tags [get]
func (h *FollowAPI) MyFollowingTags(c *gin.Context) {
	var page dto.PageRequest
	if !bindQuery(c, &page) {
		return
	}
	normalizePageRequest(&page)
	list, total, err := h.svc.FollowingTags(middleware.CurrentUserID(c), page.Page, page.PageSize)
	h.pageOrFail(c, list, total, page.Page, page.PageSize, err)
}

func (h *FollowAPI) pageOrFail(c *gin.Context, list any, total int64, page, size int, err error) {
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}
