package router

import (
	_ "github.com/feedora/backend/docs" // swag init 生成的接口文档。
	"github.com/feedora/backend/internal/api"
	"github.com/feedora/backend/pkg/config"
	"github.com/feedora/backend/pkg/jwtx"
	"github.com/feedora/backend/pkg/middleware"
	"github.com/feedora/backend/pkg/observability"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Handlers 汇总各模块 API 处理器。
type Handlers struct {
	Auth         *api.AuthAPI
	User         *api.UserAPI
	Post         *api.PostAPI
	Comment      *api.CommentAPI
	Interaction  *api.InteractionAPI
	Tag          *api.TagAPI
	Topic        *api.TopicAPI
	Circle       *api.CircleAPI
	File         *api.FileAPI
	Admin        *api.AdminAPI
	Search       *api.SearchAPI
	Rank         *api.RankAPI
	Notification *api.NotificationAPI
	Growth       *api.GrowthAPI
	Follow       *api.FollowAPI
}

// Options 路由装配所需的配置与依赖。
type Options struct {
	CORS      config.CORSConfig
	JWT       *jwtx.Manager
	Checker   middleware.AuthChecker // 鉴权链补充校验（token 黑名单 / 用户状态），可为 nil
	StaticDir string                // 非空时对外提供本地静态文件（模拟 OSS）
	Handlers  Handlers
}

// ctx 在各模块路由注册函数间共享的上下文。
type ctx struct {
	v1      *gin.RouterGroup
	authMW  gin.HandlerFunc
	adminMW gin.HandlerFunc
	h       Handlers
}

// New 构建 gin 引擎并注册全部路由。
func New(opts Options) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Recover(), middleware.Trace(), middleware.Logger(), observability.HTTPMiddleware(), middleware.CORS(opts.CORS))

	// Prometheus 指标暴露（不鉴权，生产环境用网络策略限制访问）。
	r.GET("/metrics", gin.WrapH(observability.MetricsHandler()))

	// Swagger 接口文档，访问 /swagger/index.html。
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if opts.StaticDir != "" {
		r.Static("/static", opts.StaticDir)
	}

	v1 := r.Group("/api/v1")
	v1.Use(middleware.OptionalAuth(opts.JWT, opts.Checker)) // 全局尝试识别登录用户，校验失败按匿名放行。

	x := &ctx{
		v1:      v1,
		authMW:  middleware.Auth(opts.JWT, opts.Checker),
		adminMW: middleware.Admin(),
		h:       opts.Handlers,
	}

	registerAuth(x)
	registerUser(x)
	registerPost(x)
	registerComment(x)
	registerInteraction(x)
	registerTag(x)
	registerTopic(x)
	registerCircle(x)
	registerFile(x)
	registerAdmin(x)
	registerSearch(x)
	registerRank(x)
	registerNotification(x)
	registerGrowth(x)
	registerFollow(x)

	return r
}
