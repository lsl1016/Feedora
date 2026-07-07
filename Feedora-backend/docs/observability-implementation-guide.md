# 可观测性代码落地指南

本指南说明在 Feedora 后端落地日志、指标、链路追踪的最小改动，并标注哪些已在本轮接入、哪些为后续渐进项。原则：不改变业务功能与路由、不大规模重构、不把 SQL 写进 service、不把业务写进 api。

## 1. 新增包

```
pkg/observability
├── metrics.go     # Prometheus 指标定义 + Gin 指标中间件 + Handler
└── context.go     # 从 gin.Context 取 traceId/userId 的小工具
```

（logger/tracing/middleware 已分别复用现有 `pkg/logger`、`pkg/middleware/trace.go`、`pkg/middleware/logger.go`，无需重复造。）

## 2. 涉及改动的位置

| 文件 | 改动 | 状态 |
|---|---|---|
| `pkg/observability/metrics.go` | 新增指标与中间件 | 本轮 |
| `internal/router/router.go` | 挂载 `/metrics` 与指标中间件 | 本轮 |
| `pkg/middleware/trace.go` | 复用入站 traceparent/X-Trace-Id/X-Request-Id | 本轮 |
| `pkg/middleware/logger.go` | 慢请求 WARN + 结构化字段 | 已具备/微调 |
| `internal/event/outbox.go` | `Message` 增加 `TraceID` | 本轮 |
| `internal/model/outbox.go` | `event_outbox` 增加 `trace_id` 列 | 本轮 |
| `internal/worker/*` | 消费日志含 eventId/eventType/traceId/workerName/durationMs | 本轮 |
| `pkg/database/mysql.go` | GORM 慢 SQL 阈值日志（Logger 配置） | 本轮（阈值）/ ctx 渐进 |
| `pkg/redisx/redis.go` | 命中率/耗时 hook | 后续 |
| `pkg/kafkax/producer.go` | 生产结果指标 | 后续 |
| `internal/search` `pkg/ossx` | 查询/上传耗时与错误 | 后续 |

## 3. Gin 请求日志中间件

见 `logging-spec.md` 第 12 节；`pkg/middleware/logger.go` 已输出 traceId/method/path/status/durationMs，并对 >500ms 打 WARN。

## 4. Trace 中间件

`pkg/middleware/trace.go` 增强为：入站若含 `traceparent` 则解析其 trace-id 复用；否则用 `X-Trace-Id`；否则生成 `trace-<uuid>`。始终写 `gin.Context` 与响应头 `X-Trace-Id`。

## 5. Prometheus metrics 接入

`pkg/observability/metrics.go` 定义 `feedora_api_request_total{method,route,status}` 与 `feedora_api_request_duration_seconds{method,route}`，并提供 `HTTPMiddleware()`（用注册路由模板 `c.FullPath()` 作为 route，避免高基数）与 `MetricsHandler()`。router 中：

```go
r.GET("/metrics", gin.WrapH(observability.MetricsHandler()))
v1.Use(observability.HTTPMiddleware()) // 或在根路由挂载
```

业务计数：`observability.BizInc("post_create")` 在各 service 成功路径调用（可渐进补充）。

## 6. GORM 慢查询

`pkg/database/mysql.go` 配置 GORM Logger：`SlowThreshold: 200ms`、`LogLevel: Warn`，即可自动打印慢 SQL。接入 ctx 后（`db.WithContext(ctx)`）慢 SQL 日志可带 traceId。

## 7. Redis 命中率

在 `internal/cache` 的 `GetJSON` 命中/未命中处调用 `observability` 计数器（`cache_hit_total{cache}` / `cache_miss_total{cache}`）。命令耗时可用 go-redis Hook 统一采集（后续）。

## 8. Kafka Producer 记录投递结果

`internal/worker` 的 Outbox Dispatcher 投递成功/失败时 `feedora_kafka_produce_total{topic,result}` +1，并已有日志（见 `logging-spec.md` 第 10 节）。

## 9. Kafka Consumer 记录消费结果

Worker 每条消息处理完 `feedora_kafka_consume_total{topic,result}` +1、`feedora_worker_process_total{worker,result}` +1、记录耗时直方图。

## 10. Worker 记录 eventId/eventType/traceId

`internal/worker/handlers.go`、`runner.go` 消费入口打印：

```go
logger.Infof("consume event, workerName:%s, eventType:%s, eventId:%s, traceId:%s, durationMs:%d",
    workerName, m.EventType, m.EventID, m.TraceID, ms)
```

## 11. ES 查询耗时

`internal/search/client.go` 的 `SearchIDs`/`IndexDoc` 前后计时，超 300ms WARN，并计入 `feedora_es_query_duration_seconds`（后续）。

## 12. OSS 上传耗时与大小

`pkg/ossx` 或 `internal/service/file_service.go` 上传前后计时，记录 `size`、`durationMs`，超 1000ms WARN（后续）。

## 13. 渐进项：Service/Repository 携带 context.Context

为实现 DB span 与跨进程 trace 贯通，目标是让 service/repository 方法首参为 `ctx context.Context` 并 `db.WithContext(ctx)`。因涉及面广，建议**按写路径模块逐个改造**（post → comment → interaction → circle …），每次改完回归，避免一次性重构。当前仓库尚未统一携带 ctx，属已知待办。

## 14. 验收标准

- [x] 每个 HTTP 请求都有 traceId，响应头返回 `X-Trace-Id`
- [x] 请求日志包含 traceId/method/path/statusCode/durationMs/userId，慢请求 WARN
- [x] `event_outbox` 与事件 `Message` 含 `traceId` 字段（列/字段/投递/消费日志已打通；跨进程取值待 ctx 线程化填充，见第 13 节）
- [x] Worker 日志含 eventId/eventType/traceId/durationMs 字段
- [x] `/metrics` 可访问：API `:8090/metrics`（请求数/耗时/状态码）、Worker `:9091/metrics`（Kafka 生产/消费计数）
- [ ] Grafana 展示 API QPS/错误率/P95（部署 Prometheus+Grafana 后）
- [x] 慢 SQL 记录：GORM SlowThreshold=200ms 已配；ES 慢查询 / OSS 失败为后续
- [ ] Outbox pending/failed 指标（导出器/定时查库，后续）
- [~] 出错时凭 traceId 串联：traceId 生成/响应头/日志字段已具备；跨进程贯通与 OTel span 为渐进项（第 13 节）

## 15. 运行与访问

```bash
cd deployments && docker compose up -d      # 中间件
go run ./cmd/api                            # :8090，/metrics 暴露
go run ./cmd/worker                         # 事件消费
curl -i http://localhost:8090/api/v1/posts  # 响应头含 X-Trace-Id
curl http://localhost:8090/metrics          # Prometheus 指标
```

凭 traceId 排查：拿响应头 `X-Trace-Id` → 在日志（Loki/终端）检索该 traceId → 依次看 API/Service/Repository/Outbox/Worker/ES 日志。
