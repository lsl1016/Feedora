package app

import (
	"fmt"

	"github.com/feedora/backend/internal/api"
	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/event"
	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/internal/repository"
	"github.com/feedora/backend/internal/router"
	"github.com/feedora/backend/internal/search"
	"github.com/feedora/backend/internal/service"
	"github.com/feedora/backend/pkg/config"
	"github.com/feedora/backend/pkg/database"
	"github.com/feedora/backend/pkg/esx"
	"github.com/feedora/backend/pkg/jwtx"
	"github.com/feedora/backend/pkg/logger"
	"github.com/feedora/backend/pkg/ossx"
	"github.com/feedora/backend/pkg/redisx"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// App 应用容器，组装配置、数据库、缓存、搜索、事件与 HTTP 引擎。
type App struct {
	Cfg    *config.Config
	DB     *gorm.DB
	Engine *gin.Engine
}

// New 加载配置、连接依赖并装配全部服务，返回可运行的应用。
func New(configPath string) (*App, error) {
	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}

	db, err := database.New(cfg.MySQL)
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}
	if err := db.AutoMigrate(model.AllModels()...); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}
	logger.Infof("数据库连接与迁移完成")

	// Redis（可选）。
	cch := cache.New(nil)
	if cfg.Redis.Enabled {
		rdb, err := redisx.New(cfg.Redis)
		if err != nil {
			return nil, fmt.Errorf("连接 Redis 失败: %w", err)
		}
		cch = cache.New(rdb)
		logger.Infof("Redis 连接成功")
	}

	// Elasticsearch（可选）。
	var searchClient *search.Client
	if cfg.Elasticsearch.Enabled {
		es, err := esx.New(cfg.Elasticsearch)
		if err != nil {
			return nil, fmt.Errorf("连接 Elasticsearch 失败: %w", err)
		}
		searchClient = search.NewClient(es, cfg.Elasticsearch.IndexPrefix)
		logger.Infof("Elasticsearch 连接成功")
	}

	var storage *ossx.LocalStorage
	switch cfg.OSS.Type {
	case "", "local":
		storage, err = ossx.NewLocalStorage(cfg.OSS.BasePath, cfg.OSS.PublicBaseUrl)
		if err != nil {
			return nil, fmt.Errorf("初始化存储失败: %w", err)
		}
	default:
		// minio / aliyun 实现尚未落地，显式失败优于静默降级为本地磁盘。
		return nil, fmt.Errorf("不支持的 oss.type=%q（minio / aliyun 尚未实现，请使用 local）", cfg.OSS.Type)
	}
	jm := jwtx.NewManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// 仓储层。
	userRepo := repository.NewUserRepository(db)
	postRepo := repository.NewPostRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	interRepo := repository.NewInteractionRepository(db)
	tagRepo := repository.NewTagRepository(db)
	topicRepo := repository.NewTopicRepository(db)
	circleRepo := repository.NewCircleRepository(db)
	fileRepo := repository.NewFileRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	outboxRepo := repository.NewOutboxRepository(db)
	notifRepo := repository.NewNotificationRepository(db)
	growthRepo := repository.NewGrowthRepository(db)
	followRepo := repository.NewFollowRepository(db)

	// 事件生产者：启用 Kafka 时走 Outbox 可靠投递，否则打印日志。
	var producer event.Producer
	if cfg.Kafka.Enabled {
		producer = event.NewOutboxProducer(outboxRepo, cfg.Kafka.TopicPrefix)
		logger.Infof("事件生产者：Outbox 模式")
	} else {
		producer = event.NewNoopProducer(cfg.Kafka.TopicPrefix)
		logger.Infof("事件生产者：Noop 模式")
	}

	// 服务层。
	postSvc := service.NewPostService(postRepo, userRepo, tagRepo, topicRepo, circleRepo, interRepo, producer, cch)
	authSvc := service.NewAuthService(userRepo, jm, producer)
	commentSvc := service.NewCommentService(commentRepo, postRepo, userRepo, interRepo, producer, cch)
	interSvc := service.NewInteractionService(postRepo, userRepo, interRepo, producer, cch)
	tagSvc := service.NewTagService(tagRepo, postSvc)
	topicSvc := service.NewTopicService(topicRepo, postSvc)
	circleSvc := service.NewCircleService(circleRepo, userRepo, postSvc, producer, cch)
	userSvc := service.NewUserService(userRepo, interRepo, postSvc, commentSvc, cch)
	fileSvc := service.NewFileService(fileRepo, storage)
	adminSvc := service.NewAdminService(adminRepo, userRepo, tagRepo, topicRepo, circleRepo, commentRepo, producer)
	searchSvc := service.NewSearchService(searchClient, postRepo, userRepo, topicRepo, circleRepo, postSvc, cch)
	rankSvc := service.NewRankService(db, postRepo, userRepo, topicRepo, circleRepo, cch)
	notifSvc := service.NewNotificationService(notifRepo, cch)
	growthSvc := service.NewGrowthService(growthRepo, userRepo, rankSvc)
	followSvc := service.NewFollowService(followRepo, userRepo, postSvc, producer, cch)

	// 接口层。
	handlers := router.Handlers{
		Auth:         api.NewAuthAPI(authSvc),
		User:         api.NewUserAPI(userSvc),
		Post:         api.NewPostAPI(postSvc),
		Comment:      api.NewCommentAPI(commentSvc),
		Interaction:  api.NewInteractionAPI(interSvc),
		Tag:          api.NewTagAPI(tagSvc),
		Topic:        api.NewTopicAPI(topicSvc),
		Circle:       api.NewCircleAPI(circleSvc),
		File:         api.NewFileAPI(fileSvc),
		Admin:        api.NewAdminAPI(adminSvc),
		Search:       api.NewSearchAPI(searchSvc),
		Rank:         api.NewRankAPI(rankSvc),
		Notification: api.NewNotificationAPI(notifSvc),
		Growth:       api.NewGrowthAPI(growthSvc),
		Follow:       api.NewFollowAPI(followSvc),
	}

	engine := router.New(router.Options{
		CORS:      cfg.CORS,
		JWT:       jm,
		StaticDir: storage.BasePath(),
		Handlers:  handlers,
	})

	return &App{Cfg: cfg, DB: db, Engine: engine}, nil
}

// Run 启动 HTTP 服务。
func (a *App) Run() error {
	addr := fmt.Sprintf(":%d", a.Cfg.Server.Port)
	logger.Infof("%s 服务启动，监听 %s", a.Cfg.Server.Name, addr)
	return a.Engine.Run(addr)
}
