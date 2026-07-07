package router

// registerRank 注册热门榜单路由。
func registerRank(x *ctx) {
	x.v1.GET("/hot/ranks", x.h.Rank.HotRanks)
}
