package api

import (
	"github.com/feedora/backend/internal/dto"
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/response"
	"github.com/gin-gonic/gin"
)

// NotificationAPI 通知接口。
type NotificationAPI struct {
	svc *service.NotificationService
}

func NewNotificationAPI(svc *service.NotificationService) *NotificationAPI {
	return &NotificationAPI{svc: svc}
}

// List 通知列表
// @Summary  通知列表
// @Tags     通知
// @Produce  json
// @Param    req  query  dto.NotificationListQuery  false  "通知列表查询参数"
// @Success  200  {object}  dto.NotificationListResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /notifications [get]
func (h *NotificationAPI) List(c *gin.Context) {
	var req dto.NotificationListQuery
	if !bindQuery(c, &req) {
		return
	}
	normalizePageRequest(&req.PageRequest)
	list, _ := h.svc.List(middleware.CurrentUserID(c), req.Category, req.Page, req.PageSize)
	response.OK(c, list)
}

// MarkRead 标记单条通知已读
// @Summary  标记单条通知已读
// @Tags     通知
// @Produce  json
// @Param    req  path  dto.NotificationIDURI  true  "通知路径参数"
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /notifications/{notificationId}/read [put]
func (h *NotificationAPI) MarkRead(c *gin.Context) {
	var req dto.NotificationIDURI
	if !bindURI(c, &req) {
		return
	}
	h.svc.MarkRead(req.NotificationID, middleware.CurrentUserID(c))
	response.OK(c, gin.H{})
}

// MarkAllRead 标记全部通知已读
// @Summary  标记全部通知已读
// @Tags     通知
// @Produce  json
// @Success  200  {object}  dto.EmptyResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /notifications/read-all [put]
func (h *NotificationAPI) MarkAllRead(c *gin.Context) {
	h.svc.MarkAllRead(middleware.CurrentUserID(c))
	response.OK(c, gin.H{})
}

// UnreadCount 未读通知数量
// @Summary  未读通知数量
// @Tags     通知
// @Produce  json
// @Success  200  {object}  dto.CountResponse
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /notifications/unread-count [get]
func (h *NotificationAPI) UnreadCount(c *gin.Context) {
	response.OK(c, gin.H{"count": h.svc.UnreadCount(middleware.CurrentUserID(c))})
}
