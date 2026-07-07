package main

import (
	"flag"

	"github.com/feedora/backend/internal/app"
	"github.com/feedora/backend/pkg/logger"
)

// @title                      Feedora API
// @version                    1.0
// @description                Feedora 后端接口文档。
// @BasePath                   /api/v1
// @securityDefinitions.apikey BearerAuth
// @in                         header
// @name                       Authorization
func main() {
	configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
	flag.Parse()

	a, err := app.New(*configPath)
	if err != nil {
		logger.Errorf("%v", err)
		return
	}
	if err := a.Run(); err != nil {
		logger.Errorf("服务启动失败: %v", err)
	}
}
