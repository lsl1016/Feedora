// Package observability 提供 Prometheus 指标定义、Gin 指标中间件与 /metrics 暴露。
// 指标 label 保持低基数：使用注册路由模板作为 route，禁止 userId/traceId 作为 label。
package observability

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// apiRequestTotal HTTP 请求总数。
	apiRequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "feedora_api_request_total",
		Help: "HTTP 请求总数",
	}, []string{"method", "route", "status"})

	// apiRequestDuration HTTP 请求耗时（秒）。
	apiRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "feedora_api_request_duration_seconds",
		Help:    "HTTP 请求耗时（秒）",
		Buckets: []float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
	}, []string{"method", "route"})

	// businessActionTotal 业务动作计数（发帖/评论/点赞等）。
	businessActionTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "feedora_business_action_total",
		Help: "业务动作计数",
	}, []string{"action"})

	// cacheAccessTotal 缓存命中/未命中计数。
	cacheAccessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "feedora_cache_access_total",
		Help: "缓存访问计数（命中/未命中）",
	}, []string{"cache", "result"})

	// kafkaProduceTotal Kafka 生产结果计数。
	kafkaProduceTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "feedora_kafka_produce_total",
		Help: "Kafka 生产结果计数",
	}, []string{"topic", "result"})

	// kafkaConsumeTotal Kafka 消费结果计数。
	kafkaConsumeTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "feedora_kafka_consume_total",
		Help: "Kafka 消费结果计数",
	}, []string{"topic", "result"})

	// workerProcessTotal 各 Worker 处理结果计数。
	workerProcessTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "feedora_worker_process_total",
		Help: "Worker 处理结果计数",
	}, []string{"worker", "result"})
)

// HTTPMiddleware 记录每个请求的计数与耗时。route 使用注册路由模板避免高基数。
func HTTPMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}
		status := strconv.Itoa(c.Writer.Status())
		apiRequestTotal.WithLabelValues(c.Request.Method, route, status).Inc()
		apiRequestDuration.WithLabelValues(c.Request.Method, route).Observe(time.Since(start).Seconds())
	}
}

// MetricsHandler 返回 Prometheus /metrics 处理器。
func MetricsHandler() http.Handler { return promhttp.Handler() }

// BizInc 业务动作计数 +1。
func BizInc(action string) { businessActionTotal.WithLabelValues(action).Inc() }

// CacheHit / CacheMiss 缓存命中统计。
func CacheHit(cache string)  { cacheAccessTotal.WithLabelValues(cache, "hit").Inc() }
func CacheMiss(cache string) { cacheAccessTotal.WithLabelValues(cache, "miss").Inc() }

// KafkaProduce 记录一次 Kafka 生产结果（result: success/fail）。
func KafkaProduce(topic, result string) { kafkaProduceTotal.WithLabelValues(topic, result).Inc() }

// KafkaConsume 记录一次 Kafka 消费结果。
func KafkaConsume(topic, result string) { kafkaConsumeTotal.WithLabelValues(topic, result).Inc() }

// WorkerProcess 记录一次 Worker 处理结果。
func WorkerProcess(worker, result string) { workerProcessTotal.WithLabelValues(worker, result).Inc() }
