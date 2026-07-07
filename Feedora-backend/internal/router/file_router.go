package router

// registerFile 注册文件上传模块路由。
func registerFile(x *ctx) {
	x.v1.POST("/files/upload", x.authMW, x.h.File.Upload)
}
