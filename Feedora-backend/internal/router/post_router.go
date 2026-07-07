package router

// registerPost 注册帖子模块路由。
func registerPost(x *ctx) {
	x.v1.GET("/posts", x.h.Post.List)
	x.v1.GET("/posts/:postId", x.h.Post.Get)
	x.v1.POST("/posts", x.authMW, x.h.Post.Create)
	x.v1.PUT("/posts/:postId", x.authMW, x.h.Post.Update)
	x.v1.DELETE("/posts/:postId", x.authMW, x.h.Post.Delete)
	x.v1.PUT("/posts/:postId/hide", x.authMW, x.h.Post.Hide)
	x.v1.PUT("/posts/:postId/unhide", x.authMW, x.h.Post.Unhide)
	x.v1.POST("/posts/:postId/share", x.authMW, x.h.Post.Share)
	x.v1.POST("/posts/:postId/repost", x.authMW, x.h.Post.Repost)
}
