package main

import (
	"flag"

	"github.com/feedora/backend/internal/model"
	"github.com/feedora/backend/pkg/config"
	"github.com/feedora/backend/pkg/database"
	"github.com/feedora/backend/pkg/logger"
)

// 数据库迁移入口：执行 GORM AutoMigrate 建表。
func main() {
	configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
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
	logger.Infof("数据库迁移完成")
}
