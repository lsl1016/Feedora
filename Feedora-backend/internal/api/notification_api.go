package api

import (
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
// @Param    category  query  string  false  "通知分类"
// @Param    page      query  int     false  "页码"
// @Param    pageSize  query  int     false  "每页数量"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /notifications [get]
func (h *NotificationAPI) List(c *gin.Context) {
	page, size := pageParams(c)
	list, _ := h.svc.List(middleware.CurrentUserID(c), c.Query("category"), page, size)
	response.OK(c, list)
}

// MarkRead 标记单条通知已读
// @Summary  标记单条通知已读
// @Tags     通知
// @Produce  json
// @Param    notificationId  path  int  true  "通知ID"
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /notifications/{notificationId}/read [put]
func (h *NotificationAPI) MarkRead(c *gin.Context) {
	h.svc.MarkRead(paramID(c, "notificationId"), middleware.CurrentUserID(c))
	response.OK(c, gin.H{})
}

// MarkAllRead 标记全部通知已读
// @Summary  标记全部通知已读
// @Tags     通知
// @Produce  json
// @Success  200  {object}  response.Body
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
// @Success  200  {object}  response.Body
// @Failure  400  {object}  response.Body
// @Security BearerAuth
// @Router   /notifications/unread-count [get]
func (h *NotificationAPI) UnreadCount(c *gin.Context) {
	response.OK(c, gin.H{"count": h.svc.UnreadCount(middleware.CurrentUserID(c))})
}
