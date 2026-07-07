package router

// registerNotification 注册通知模块路由。
func registerNotification(x *ctx) {
	x.v1.GET("/notifications", x.authMW, x.h.Notification.List)
	x.v1.GET("/notifications/unread-count", x.authMW, x.h.Notification.UnreadCount)
	x.v1.PUT("/notifications/read-all", x.authMW, x.h.Notification.MarkAllRead)
	x.v1.PUT("/notifications/:notificationId/read", x.authMW, x.h.Notification.MarkRead)
}
