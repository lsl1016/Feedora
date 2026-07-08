package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// InteractionAPI 点赞 / 收藏接口。
type InteractionAPI struct {
	svc *service.InteractionService
}

func NewInteractionAPI(svc *service.InteractionService) *InteractionAPI {
	return &InteractionAPI{svc: svc}
}

// Like 点赞帖子
// @Summary  点赞帖子
// @Tags     互动
// @Produce  json
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.InteractionResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /posts/{postId}/like [post]
func (h *InteractionAPI) Like(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	res, err := h.svc.SetLike(req.PostID, middleware.CurrentUserID(c), true)
	h.reply(c, res, err)
}

// Unlike 取消点赞帖子
// @Summary  取消点赞帖子
// @Tags     互动
// @Produce  json
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.InteractionResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /posts/{postId}/like [delete]
func (h *InteractionAPI) Unlike(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	res, err := h.svc.SetLike(req.PostID, middleware.CurrentUserID(c), false)
	h.reply(c, res, err)
}

// Favorite 收藏帖子
// @Summary  收藏帖子
// @Tags     互动
// @Produce  json
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.InteractionResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /posts/{postId}/favorite [post]
func (h *InteractionAPI) Favorite(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	res, err := h.svc.SetFavorite(req.PostID, middleware.CurrentUserID(c), true)
	h.reply(c, res, err)
}

// Unfavorite 取消收藏帖子
// @Summary  取消收藏帖子
// @Tags     互动
// @Produce  json
// @Param    req  path  dto.PostIDURI  true  "帖子路径参数"
// @Success  200  {object}  dto.InteractionResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /posts/{postId}/favorite [delete]
func (h *InteractionAPI) Unfavorite(c *gin.Context) {
	var req dto.PostIDURI
	if !bindURI(c, &req) {
		return
	}
	res, err := h.svc.SetFavorite(req.PostID, middleware.CurrentUserID(c), false)
	h.reply(c, res, err)
}

func (h *InteractionAPI) reply(c *gin.Context, res *dto.InteractionResult, err error) {
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, res)
}
