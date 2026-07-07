package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	errs "github.com/feedora/backend/pkg/errors"
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
// @Param    feedType       query  string  false  "信息流类型"
// @Param    sort           query  string  false  "排序方式"
// @Param    status         query  string  false  "状态"
// @Param    keyword        query  string  false  "搜索关键词"
// @Param    tagId          query  int     false  "标签ID"
// @Param    circleId       query  int     false  "圈子ID"
// @Param    topicId        query  int     false  "话题ID"
// @Param    authorId       query  int     false  "作者ID"
// @Param    includeHidden  query  string  false  "是否包含隐藏帖子"
// @Param    page           query  int     false  "页码"
// @Param    pageSize       query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts [get]
func (h *PostAPI) List(c *gin.Context) {
	page, size := pageParams(c)
	list, total, err := h.svc.List(service.ListFilter{
		FeedType:      c.Query("feedType"),
		Sort:          c.Query("sort"),
		Status:        c.Query("status"),
		Keyword:       c.Query("keyword"),
		TagID:         queryID(c, "tagId"),
		CircleID:      queryID(c, "circleId"),
		TopicID:       queryID(c, "topicId"),
		AuthorID:      queryID(c, "authorId"),
		IncludeHidden: c.Query("includeHidden") == "true",
		ViewerID:      middleware.CurrentUserID(c),
		Page:          page,
		PageSize:      size,
	})
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Page(c, list, total, page, size)
}

// Get 获取帖子详情
// @Summary  获取帖子详情
// @Tags     帖子
// @Produce  json
// @Param    postId  path  int  true  "帖子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId} [get]
func (h *PostAPI) Get(c *gin.Context) {
	res, err := h.svc.Get(paramID(c, "postId"), middleware.CurrentUserID(c))
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
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts [post]
func (h *PostAPI) Create(c *gin.Context) {
	var in dto.CreatePostRequest
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

// Update 更新帖子
// @Summary  更新帖子
// @Tags     帖子
// @Accept   json
// @Produce  json
// @Security BearerAuth
// @Param    postId  path  int  true  "帖子ID"
// @Param    body    body  dto.UpdatePostRequest  true  "更新帖子请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId} [put]
func (h *PostAPI) Update(c *gin.Context) {
	var in dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, errs.ErrParams)
		return
	}
	res, err := h.svc.Update(paramID(c, "postId"), middleware.CurrentUserID(c), in)
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
// @Param    postId  path  int  true  "帖子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/hide [put]
func (h *PostAPI) Hide(c *gin.Context) {
	if err := h.svc.SetHidden(paramID(c, "postId"), middleware.CurrentUserID(c), true); err != nil {
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
// @Param    postId  path  int  true  "帖子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/unhide [put]
func (h *PostAPI) Unhide(c *gin.Context) {
	if err := h.svc.SetHidden(paramID(c, "postId"), middleware.CurrentUserID(c), false); err != nil {
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
// @Param    postId  path  int  true  "帖子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId} [delete]
func (h *PostAPI) Delete(c *gin.Context) {
	isAdmin := middleware.CurrentRole(c) == "admin"
	if err := h.svc.Delete(paramID(c, "postId"), middleware.CurrentUserID(c), isAdmin); err != nil {
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
// @Param    postId  path  int  true  "帖子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/share [post]
func (h *PostAPI) Share(c *gin.Context) {
	if err := h.svc.Share(paramID(c, "postId")); err != nil {
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
// @Param    postId  path  int  true  "帖子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/repost [post]
func (h *PostAPI) Repost(c *gin.Context) {
	var body struct {
		RepostComment string `json:"repostComment"`
	}
	_ = c.ShouldBindJSON(&body)
	res, err := h.svc.Repost(paramID(c, "postId"), middleware.CurrentUserID(c), body.RepostComment)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
