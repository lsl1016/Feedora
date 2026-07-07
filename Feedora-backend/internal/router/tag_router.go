package router

// registerTag 注册标签模块路由。
func registerTag(x *ctx) {
	x.v1.GET("/tags", x.h.Tag.List)
	x.v1.GET("/tags/:tagId", x.h.Tag.Get)
	x.v1.GET("/tags/:tagId/posts", x.h.Tag.Posts)
}
