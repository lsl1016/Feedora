package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// PostAPI 帖子接口。
type PostAPI struct {
	svc *service.PostService
}

func NewPostAPI(svc *service.PostService) *PostAPI {
	return &PostAPI{svc: svc}
}

// List 获取帖子列表
// @Summary  获取帖子列表
// @Tags     帖子
// @Produce  json
// @Param    req  query  dto.PostListQuery  false  "帖子列表查询参数"
// @Success  200  {object}  dto.PostPageResponse
// @Failure  400  {object}  response.Body
// @Router   /posts [get]
func (h *PostAPI) List(c *gin.Context) {
	var req dto.PostListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, total, err := h.svc.List(service.ListFilter{
		FeedType:      req.FeedType,
		Sort:          req.Sort,
		Status:        req.Status,
		Keyword:       req.Keyword,
		TagID:         req.TagID,
		CircleID:      req.CircleID,
		TopicID:       req.TopicID,
		AuthorID:      req.AuthorID,
		IncludeHidden: req.IncludeHidden,
		ViewerID:      middleware.CurrentUserID(c),
		Page:          req.Page,
		PageSize:      req.PageSize,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, req.Page, req.PageSize)
}

// Get 获取帖子详情
// @Summary  获取帖子详情
// @Tags     帖子
// @Produce  json
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.PostResponse
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId} [get]
func (h *PostAPI) Get(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	res, err := h.svc.Get(req.PostID, middleware.CurrentUserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Create 创建帖子
// @Summary  创建帖子
// @Tags     帖子
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    body  body  dto.CreatePostRequest  true  "创建帖子请求体"
// @Success  200  {object}  dto.PostResponse
// @Failure  400  {object}  response.Body
// @Router   /posts [post]
func (h *PostAPI) Create(c *gin.Context) {
	var in dto.CreatePostRequest
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

// Update 更新帖子
// @Summary  更新帖子
// @Tags     帖子
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    postId  path  dto.PostIDURI          true  "帖子路径参数"
// @Param    body    body  dto.UpdatePostRequest  true  "更新帖子请求体"
// @Success  200  {object}  dto.PostResponse
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId} [put]
func (h *PostAPI) Update(c *gin.Context) {
	var uri dto.PostIDURI
	if !bindURI(c, &uri) {
		return
	}
	var in dto.UpdatePostRequest
	if !bindJSON(c, &in) {
		return
	}
	res, err := h.svc.Update(uri.PostID, middleware.CurrentUserID(c), in)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Hide 隐藏帖子
// @Summary  隐藏帖子
// @Tags     帖子
// @Produce  json
// @Security BearerAuth
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/hide [put]
func (h *PostAPI) Hide(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.SetHidden(req.PostID, middleware.CurrentUserID(c), true); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Unhide 取消隐藏帖子
// @Summary  取消隐藏帖子
// @Tags     帖子
// @Produce  json
// @Security BearerAuth
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/unhide [put]
func (h *PostAPI) Unhide(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.SetHidden(req.PostID, middleware.CurrentUserID(c), false); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Delete 删除帖子
// @Summary  删除帖子
// @Tags     帖子
// @Produce  json
// @Security BearerAuth
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId} [delete]
func (h *PostAPI) Delete(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	isAdmin := middleware.CurrentRole(c) == "admin"
	if err := h.svc.Delete(req.PostID, middleware.CurrentUserID(c), isAdmin); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Share 分享帖子
// @Summary  分享帖子
// @Tags     帖子
// @Produce  json
// @Security BearerAuth
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/share [post]
func (h *PostAPI) Share(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.Share(req.PostID); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Repost 转发帖子
// @Summary  转发帖子
// @Tags     帖子
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    postId  path  dto.PostIDURI           true  "帖子路径参数"
// @Param    body    body  dto.RepostPostRequest   true  "转发请求体"
// @Success  200  {object}  dto.PostResponse
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/repost [post]
func (h *PostAPI) Repost(c *gin.Context) {
	var uri dto.PostIDURI
	if !bindURI(c, &uri) {
		return
	}
	var body dto.RepostPostRequest
	if !bindJSON(c, &body) {
		return
	}
	res, err := h.svc.Repost(uri.PostID, middleware.CurrentUserID(c), body.RepostComment)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
