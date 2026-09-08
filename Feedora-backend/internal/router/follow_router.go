package router

// registerFollow 注册关注体系路由。
func registerFollow(x *ctx) {
	// 关注 / 取关 / 状态。
	x.v1.POST("/users/:userId/follow", x.authMW, x.h.Follow.Follow)
	x.v1.DELETE("/users/:userId/follow", x.authMW, x.h.Follow.Unfollow)
	x.v1.GET("/users/:userId/follow/state", x.h.Follow.State)

	// 关注与粉丝列表（公开）。/users/me/following 为静态路由，优先于 :userId 匹配。
	x.v1.GET("/users/:userId/following", x.h.Follow.Following)
	x.v1.GET("/users/:userId/followers", x.h.Follow.Followers)
	x.v1.GET("/users/me/following", x.authMW, x.h.Follow.MyFollowing)

	// 我的关注聚合（关注中心页面）。
	x.v1.GET("/users/me/following-feed", x.authMW, x.h.Follow.MyFollowingFeed)
	x.v1.GET("/users/me/following-circles", x.authMW, x.h.Follow.MyFollowingCircles)
	x.v1.GET("/users/me/following-topics", x.authMW, x.h.Follow.MyFollowingTopics)
	x.v1.GET("/users/me/following-tags", x.authMW, x.h.Follow.MyFollowingTags)
}
