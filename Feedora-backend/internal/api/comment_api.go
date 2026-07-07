package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	errs "github.com/feedora/backend/pkg/errors"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// CommentAPI 评论接口。
type CommentAPI struct {
	svc *service.CommentService
}

func NewCommentAPI(svc *service.CommentService) *CommentAPI {
	return &CommentAPI{svc: svc}
}

// ListByPost 获取帖子评论列表
// @Summary  获取帖子评论列表
// @Tags     评论
// @Produce  json
// @Param    postId  path  int  true  "帖子ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/comments [get]
func (h *CommentAPI) ListByPost(c *gin.Context) {
	list, err := h.svc.ListByPost(paramID(c, "postId"), middleware.CurrentUserID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, list)
}

// Create 发表评论
// @Summary  发表评论
// @Tags     评论
// @Accept   json
// @Produce  json
// @Param    body  body  dto.CreateCommentRequest  true  "评论请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments [post]
func (h *CommentAPI) Create(c *gin.Context) {
	var in dto.CreateCommentRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, errs.ErrParams)
		return
	}
	res, err := h.svc.Create(in.PostID, middleware.CurrentUserID(c), in.Content)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Reply 回复评论
// @Summary  回复评论
// @Tags     评论
// @Accept   json
// @Produce  json
// @Param    commentId  path  int  true  "评论ID"
// @Param    body  body  dto.ReplyCommentRequest  true  "回复请求体"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments/{commentId}/replies [post]
func (h *CommentAPI) Reply(c *gin.Context) {
	var in dto.ReplyCommentRequest
	if err := c.ShouldBindJSON(&in); err != nil {
		response.Fail(c, errs.ErrParams)
		return
	}
	res, err := h.svc.Reply(paramID(c, "commentId"), middleware.CurrentUserID(c), in.Content, in.ReplyToUserID)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}

// Delete 删除评论
// @Summary  删除评论
// @Tags     评论
// @Produce  json
// @Param    commentId  path  int  true  "评论ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments/{commentId} [delete]
func (h *CommentAPI) Delete(c *gin.Context) {
	isAdmin := middleware.CurrentRole(c) == "admin"
	if err := h.svc.Delete(paramID(c, "commentId"), middleware.CurrentUserID(c), isAdmin); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Like 点赞评论
// @Summary  点赞评论
// @Tags     评论
// @Produce  json
// @Param    commentId  path  int  true  "评论ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments/{commentId}/like [post]
func (h *CommentAPI) Like(c *gin.Context) {
	if err := h.svc.SetLike(paramID(c, "commentId"), middleware.CurrentUserID(c), true); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Unlike 取消点赞评论
// @Summary  取消点赞评论
// @Tags     评论
// @Produce  json
// @Param    commentId  path  int  true  "评论ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments/{commentId}/like [delete]
func (h *CommentAPI) Unlike(c *gin.Context) {
	if err := h.svc.SetLike(paramID(c, "commentId"), middleware.CurrentUserID(c), false); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}
