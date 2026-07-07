package router

// registerCircle 注册圈子模块路由。
func registerCircle(x *ctx) {
	x.v1.GET("/circles", x.h.Circle.List)
	x.v1.POST("/circles", x.authMW, x.h.Circle.Create)
	x.v1.GET("/circles/:circleId", x.h.Circle.Get)
	x.v1.POST("/circles/:circleId/join", x.authMW, x.h.Circle.Join)
	x.v1.POST("/circles/:circleId/leave", x.authMW, x.h.Circle.Leave)
	x.v1.GET("/circles/:circleId/members", x.h.Circle.Members)
	x.v1.GET("/circles/:circleId/posts", x.h.Circle.Posts)
	x.v1.PUT("/circles/:circleId/members/:userId/role", x.authMW, x.h.Circle.SetRole)
	x.v1.PUT("/circles/:circleId/members/:userId/mute", x.authMW, x.h.Circle.Mute)
	x.v1.PUT("/circles/:circleId/members/:userId/unmute", x.authMW, x.h.Circle.Unmute)
	x.v1.DELETE("/circles/:circleId/members/:userId", x.authMW, x.h.Circle.Remove)
}
