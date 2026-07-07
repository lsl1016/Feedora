package router

// registerComment 注册评论模块路由。
func registerComment(x *ctx) {
	x.v1.GET("/posts/:postId/comments", x.h.Comment.ListByPost)
	x.v1.POST("/comments", x.authMW, x.h.Comment.Create)
	x.v1.POST("/comments/:commentId/replies", x.authMW, x.h.Comment.Reply)
	x.v1.DELETE("/comments/:commentId", x.authMW, x.h.Comment.Delete)
	x.v1.POST("/comments/:commentId/like", x.authMW, x.h.Comment.Like)
	x.v1.DELETE("/comments/:commentId/like", x.authMW, x.h.Comment.Unlike)
}
