package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
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
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.CommentListResponse
// @Failure  400  {object}  response.Body
// @Router   /posts/{postId}/comments [get]
func (h *CommentAPI) ListByPost(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	list, err := h.svc.ListByPost(req.PostID, middleware.CurrentUserID(c))
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
// @Success  200  {object}  dto.CommentResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments [post]
func (h *CommentAPI) Create(c *gin.Context) {
	var in dto.CreateCommentRequest
	if !bindJSON(c, &in) {
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
// @Param    commentId  path  dto.CommentIDURI          true  "评论路径参数"
// @Param    body       body  dto.ReplyCommentRequest   true  "回复请求体"
// @Success  200  {object}  dto.CommentResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments/{commentId}/replies [post]
func (h *CommentAPI) Reply(c *gin.Context) {
	var uri dto.CommentIDURI
	if !bindURI(c, &uri) {
		return
	}
	var in dto.ReplyCommentRequest
	if !bindJSON(c, &in) {
		return
	}
	res, err := h.svc.Reply(uri.CommentID, middleware.CurrentUserID(c), in.Content, in.ReplyToUserID)
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
// @Param    req  path  dto.CommentIDURI  true  "评论路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments/{commentId} [delete]
func (h *CommentAPI) Delete(c *gin.Context) {
	var req dto.CommentIDURI
	if !bindURI(c, &req) {
		return
	}
	isAdmin := middleware.CurrentRole(c) == "admin"
	if err := h.svc.Delete(req.CommentID, middleware.CurrentUserID(c), isAdmin); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Like 点赞评论
// @Summary  点赞评论
// @Tags     评论
// @Produce  json
// @Param    req  path  dto.CommentIDURI  true  "评论路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments/{commentId}/like [post]
func (h *CommentAPI) Like(c *gin.Context) {
	var req dto.CommentIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.SetLike(req.CommentID, middleware.CurrentUserID(c), true); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}

// Unlike 取消点赞评论
// @Summary  取消点赞评论
// @Tags     评论
// @Produce  json
// @Param    req  path  dto.CommentIDURI  true  "评论路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /comments/{commentId}/like [delete]
func (h *CommentAPI) Unlike(c *gin.Context) {
	var req dto.CommentIDURI
	if !bindURI(c, &req) {
		return
	}
	if err := h.svc.SetLike(req.CommentID, middleware.CurrentUserID(c), false); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{})
}
