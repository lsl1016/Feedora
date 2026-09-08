package router

// registerAdmin 注册后台管理路由（需 admin 权限）。
func registerAdmin(x *ctx) {
	g := x.v1.Group("/admin", x.authMW, x.adminMW)
	g.GET("/users", x.h.Admin.Users)
	g.PUT("/users/:userId/status", x.h.Admin.UserStatus)
	g.GET("/posts", x.h.Admin.Posts)
	g.PUT("/posts/:postId/status", x.h.Admin.PostStatus)
	g.GET("/comments", x.h.Admin.Comments)
	g.GET("/tags", x.h.Admin.Tags)
	g.POST("/tags", x.h.Admin.CreateTag)
	g.PUT("/tags/:tagId", x.h.Admin.UpdateTag)
	g.GET("/topics", x.h.Admin.Topics)
	g.POST("/topics", x.h.Admin.CreateTopic)
	g.PUT("/topics/:topicId", x.h.Admin.UpdateTopic)
	g.GET("/circles", x.h.Admin.Circles)
	g.GET("/dashboard/stats", x.h.Admin.Stats)
	g.GET("/operation-logs", x.h.Admin.Logs)
}
