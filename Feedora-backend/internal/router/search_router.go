package router

// registerSearch 注册搜索模块路由。
func registerSearch(x *ctx) {
	x.v1.GET("/search", x.h.Search.Search)
	x.v1.GET("/search/suggest", x.h.Search.Suggest)
	x.v1.GET("/search/hot-keywords", x.h.Search.HotKeywords)
}
