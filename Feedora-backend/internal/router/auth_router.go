package router

// registerAuth 注册认证模块路由。
func registerAuth(x *ctx) {
	x.v1.POST("/auth/register", x.h.Auth.Register)
	x.v1.POST("/auth/login", x.h.Auth.Login)
	x.v1.GET("/auth/me", x.authMW, x.h.Auth.Me)
	x.v1.POST("/auth/logout", x.authMW, x.h.Auth.Logout)
}
