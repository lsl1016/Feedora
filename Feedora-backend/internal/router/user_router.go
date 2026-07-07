package router

// registerUser 注册用户模块路由。
func registerUser(x *ctx) {
	x.v1.GET("/users", x.h.User.List)
	x.v1.GET("/users/:userId", x.h.User.Get)
	x.v1.PUT("/users/me/profile", x.authMW, x.h.User.UpdateProfile)
	x.v1.GET("/users/me/posts", x.authMW, x.h.User.MyPosts)
	x.v1.GET("/users/me/comments", x.authMW, x.h.User.MyComments)
	x.v1.GET("/users/me/liked-posts", x.authMW, x.h.User.MyLikedPosts)
	x.v1.GET("/users/me/favorite-posts", x.authMW, x.h.User.MyFavoritePosts)
}
