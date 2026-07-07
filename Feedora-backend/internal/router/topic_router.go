package router

// registerTopic 注册话题模块路由。
func registerTopic(x *ctx) {
	x.v1.GET("/topics", x.h.Topic.List)
	x.v1.GET("/topics/:topicId", x.h.Topic.Get)
	x.v1.GET("/topics/:topicId/posts", x.h.Topic.Posts)
}
