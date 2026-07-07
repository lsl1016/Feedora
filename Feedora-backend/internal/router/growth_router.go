package router

// registerGrowth 注册成长积分模块路由。
func registerGrowth(x *ctx) {
	x.v1.POST("/growth/check-in", x.authMW, x.h.Growth.CheckIn)
	x.v1.GET("/growth/tasks", x.authMW, x.h.Growth.Tasks)
	x.v1.POST("/growth/tasks/:taskId/claim", x.authMW, x.h.Growth.ClaimTask)
	x.v1.GET("/growth/rankings", x.h.Growth.Rankings)
}
