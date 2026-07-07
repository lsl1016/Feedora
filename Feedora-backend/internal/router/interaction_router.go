package router

// registerInteraction 注册点赞收藏模块路由。
func registerInteraction(x *ctx) {
	x.v1.POST("/posts/:postId/like", x.authMW, x.h.Interaction.Like)
	x.v1.DELETE("/posts/:postId/like", x.authMW, x.h.Interaction.Unlike)
	x.v1.POST("/posts/:postId/favorite", x.authMW, x.h.Interaction.Favorite)
	x.v1.DELETE("/posts/:postId/favorite", x.authMW, x.h.Interaction.Unfavorite)
}
