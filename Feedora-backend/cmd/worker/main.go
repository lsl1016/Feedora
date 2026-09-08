package main

import (
	"context"
	"flag"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/feedora/backend/internal/cache"
	"github.com/feedora/backend/internal/search"
	"github.com/feedora/backend/internal/worker"
	"github.com/feedora/backend/pkg/config"
	"github.com/feedora/backend/pkg/database"
	"github.com/feedora/backend/pkg/esx"
	"github.com/feedora/backend/pkg/logger"
	"github.com/feedora/backend/pkg/observability"
	"github.com/feedora/backend/pkg/redisx"
)

// Kafka Worker 启动入口：Outbox Dispatcher + 各事件消费者（搜索索引 / 通知 / 积分 / 榜单）。
func main() {
	configPath := flag.String("config", "configs/config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Errorf("加载配置失败: %v", err)
		return
	}
	if !cfg.Worker.Enabled {
		logger.Warnf("worker.enabled=false，worker 退出")
		return
	}
	if !cfg.Kafka.Enabled {
		logger.Warnf("kafka 未启用，worker 退出")
		return
	}

	db, err := database.New(cfg.MySQL)
	if err != nil {
		logger.Errorf("连接数据库失败: %v", err)
		return
	}

	cch := cache.New(nil)
	if cfg.Redis.Enabled {
		if rdb, err := redisx.New(cfg.Redis); err != nil {
			logger.Errorf("连接 Redis 失败: %v", err)
		} else {
			cch = cache.New(rdb)
			logger.Infof("Redis 连接成功")
		}
	}

	var sc *search.Client
	if cfg.Elasticsearch.Enabled {
		if es, err := esx.New(cfg.Elasticsearch); err != nil {
			logger.Errorf("连接 Elasticsearch 失败: %v", err)
		} else {
			sc = search.NewClient(es, cfg.Elasticsearch.IndexPrefix)
			logger.Infof("Elasticsearch 连接成功")
		}
	}

	runner := worker.New(cfg, db, cch, sc)
	defer runner.Close()

	// 暴露 Worker 自身的 Prometheus 指标（Kafka 消费、投递等）。
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", observability.MetricsHandler())
		logger.Infof("worker metrics 暴露于 :9091/metrics")
		if err := http.ListenAndServe(":9091", mux); err != nil {
			logger.Errorf("worker metrics 服务退出: %v", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Infof("feedora worker 启动")
	runner.Run(ctx)
	logger.Infof("feedora worker 退出")
}
