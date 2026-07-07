package main

import (
	"flag"
	"time"

	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/config"
	"github.com/feedora/backend/pkg/database"
	"github.com/feedora/backend/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func main() {
	configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
	force := flag.Bool("force", false, "已存在数据时仍然写入")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Errorf("加载配置失败: %v", err)
		return
	}
	db, err := database.New(cfg.MySQL)
	if err != nil {
		logger.Errorf("连接数据库失败: %v", err)
		return
	}
	if err := db.AutoMigrate(model.AllModels()...); err != nil {
		logger.Errorf("数据库迁移失败: %v", err)
		return
	}

	var userCount int64
	db.Model(&model.User{}).Count(&userCount)
	if userCount > 0 && !*force {
		logger.Infof("已存在用户数据，跳过种子写入（使用 -force 强制写入）")
		return
	}

	now := time.Now()
	hash, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	pwd := string(hash)

	// 用户。
	users := []model.User{
		{Account: "admin", PasswordHash: pwd, Nickname: "平台管理员", Role: "admin", Status: "normal", Level: 10, Bio: "社区管理员", CreatedAt: now, UpdatedAt: now},
		{Account: "zhangsan", PasswordHash: pwd, Nickname: "张三", Role: "user", Status: "normal", Level: 5, Bio: "后端开发，喜欢 Go 和云原生", CreatedAt: now, UpdatedAt: now},
		{Account: "lisi", PasswordHash: pwd, Nickname: "李四", Role: "user", Status: "normal", Level: 3, Bio: "前端工程师，React 爱好者", CreatedAt: now, UpdatedAt: now},
		{Account: "wangwu", PasswordHash: pwd, Nickname: "王五", Role: "user", Status: "normal", Level: 2, Bio: "全栈开发者", CreatedAt: now, UpdatedAt: now},
	}
	mustCreate(db, &users)

	// 标签。
	tags := []model.Tag{
		{Name: "Go语言", Description: "Go 后端开发", Status: "enabled", UseCount: 0, CreatedAt: now, UpdatedAt: now},
		{Name: "React", Description: "React 前端", Status: "enabled", CreatedAt: now, UpdatedAt: now},
		{Name: "项目实战", Description: "真实项目经验", Status: "enabled", CreatedAt: now, UpdatedAt: now},
		{Name: "云原生", Description: "Docker/K8s", Status: "enabled", CreatedAt: now, UpdatedAt: now},
		{Name: "数据库", Description: "MySQL/Redis", Status: "enabled", CreatedAt: now, UpdatedAt: now},
	}
	mustCreate(db, &tags)

	// 话题。
	topics := []model.Topic{
		{Name: "个人项目部署", Description: "分享项目部署经验", IsOfficial: true, IsRecommended: true, Status: "enabled", CreatedAt: now, UpdatedAt: now},
		{Name: "AI编程助手实践", Description: "AI 辅助编程", IsOfficial: true, IsRecommended: true, Status: "enabled", CreatedAt: now, UpdatedAt: now},
		{Name: "求职面经", Description: "面试经验分享", IsRecommended: false, Status: "enabled", CreatedAt: now, UpdatedAt: now},
	}
	mustCreate(db, &topics)

	// 圈子。
	circles := []model.Circle{
		{OwnerID: users[1].ID, Name: "Go 后端开发圈", Description: "专注 Go、微服务、性能优化", Category: "后端开发", JoinType: "direct", PostPermission: "all", Status: "normal", MemberCount: 1, IsRecommended: true, Rules: "请遵守圈子规则", CreatedAt: now, UpdatedAt: now},
		{OwnerID: users[2].ID, Name: "前端成长圈", Description: "React / Vue / 工程化", Category: "前端开发", JoinType: "direct", PostPermission: "all", Status: "normal", MemberCount: 1, IsRecommended: true, Rules: "友善交流", CreatedAt: now, UpdatedAt: now},
	}
	mustCreate(db, &circles)

	// 圈主成员记录。
	members := []model.CircleMember{
		{CircleID: circles[0].ID, UserID: users[1].ID, Role: "owner", Status: "normal", JoinedAt: now, UpdatedAt: now},
		{CircleID: circles[1].ID, UserID: users[2].ID, Role: "owner", Status: "normal", JoinedAt: now, UpdatedAt: now},
	}
	mustCreate(db, &members)

	// 帖子。
	circleID0 := circles[0].ID
	posts := []model.Post{
		{AuthorID: users[1].ID, Title: "Go 项目 Docker Compose 部署复盘", ContentMD: "# 部署复盘\n\n本文记录了使用 Docker Compose 部署 Go 项目的完整流程，包括镜像构建、多服务编排与常见踩坑。", Summary: "记录 Go 项目部署流程", PostType: "original", Visibility: "public", Status: "published", LikeCount: 12, CommentCount: 2, ViewCount: 120, HotScore: 200, PublishedAt: &now, CreatedAt: now, UpdatedAt: now},
		{AuthorID: users[2].ID, Title: "React 18 并发特性实践", ContentMD: "# React 18\n\n聊聊 useTransition、Suspense 在真实项目中的应用与坑点。", Summary: "React 18 并发特性", PostType: "original", Visibility: "public", Status: "published", LikeCount: 8, CommentCount: 1, ViewCount: 88, HotScore: 150, PublishedAt: &now, CreatedAt: now, UpdatedAt: now},
		{AuthorID: users[1].ID, Title: "圈内分享：Go 性能优化清单", ContentMD: "# 性能优化\n\n圈内专属分享，pprof、内存逃逸、GC 调优实战。", Summary: "Go 性能优化清单", PostType: "original", CircleID: &circleID0, Visibility: "public", Status: "published", LikeCount: 5, ViewCount: 40, HotScore: 90, PublishedAt: &now, CreatedAt: now, UpdatedAt: now},
	}
	mustCreate(db, &posts)

	// 帖子标签 / 话题关系。
	db.Create(&model.PostTag{PostID: posts[0].ID, TagID: tags[0].ID, CreatedAt: now})
	db.Create(&model.PostTag{PostID: posts[0].ID, TagID: tags[3].ID, CreatedAt: now})
	db.Create(&model.PostTag{PostID: posts[1].ID, TagID: tags[1].ID, CreatedAt: now})
	db.Create(&model.PostTag{PostID: posts[2].ID, TagID: tags[0].ID, CreatedAt: now})
	db.Create(&model.PostTopic{PostID: posts[0].ID, TopicID: topics[0].ID, CreatedAt: now})
	db.Create(&model.PostTopic{PostID: posts[1].ID, TopicID: topics[1].ID, CreatedAt: now})

	// 评论。
	comments := []model.Comment{
		{PostID: posts[0].ID, UserID: users[2].ID, Content: "写得很详细，学到了！", Status: "normal", LikeCount: 3, CreatedAt: now, UpdatedAt: now},
		{PostID: posts[0].ID, UserID: users[3].ID, Content: "请问镜像体积怎么优化？", Status: "normal", CreatedAt: now, UpdatedAt: now},
		{PostID: posts[1].ID, UserID: users[1].ID, Content: "Suspense 真香。", Status: "normal", LikeCount: 1, CreatedAt: now, UpdatedAt: now},
	}
	mustCreate(db, &comments)

	logger.Infof("种子数据写入完成：用户 %d，标签 %d，话题 %d，圈子 %d，帖子 %d", len(users), len(tags), len(topics), len(circles), len(posts))
	logger.Infof("默认账号：admin / zhangsan / lisi / wangwu，密码均为 123456")
}

// mustCreate 批量插入，失败则记录日志。
func mustCreate(db *gorm.DB, value any) {
	if err := db.Create(value).Error; err != nil {
		logger.Errorf("写入种子数据失败: %v", err)
	}
}
