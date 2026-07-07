# 链路追踪设计

目标：一次请求从 HTTP 入口到 MySQL/Redis/Kafka/Worker/ES/OSS 全链路可串联，凭 `traceId` 定位瓶颈与故障点。

## 1. 追踪覆盖范围

- HTTP 请求链路（Gin）
- Service 内部调用链路
- GORM MySQL 链路（`db.WithContext(ctx)`）
- Redis 链路
- Kafka Producer / Consumer 链路
- Elasticsearch 查询链路
- OSS 上传链路

## 2. Trace 上下文传播规范

HTTP 优先使用 W3C Trace Context 头：`traceparent`、`tracestate`；同时兼容 `X-Trace-Id`、`X-Request-Id`。

传播优先级（入站）：
1. 有 `traceparent` → 解析复用其 trace-id；
2. 否则有 `X-Trace-Id` → 复用；
3. 否则有 `X-Request-Id` → 作为 requestId，并生成新的 traceId；
4. 都没有 → 生成新的 `trace-<uuid>`。

出站：始终在响应头写回 `X-Trace-Id`。

## 3. Gin 中间件设计（现状 + 目标）

现状 `pkg/middleware/trace.go` 生成 `trace-<uuid>` 写入 `gin.Context` 与响应头。目标增强：

```go
func Trace() gin.HandlerFunc {
  return func(c *gin.Context) {
    tid := traceIDFromHeaders(c) // traceparent > X-Trace-Id > 新生成
    c.Set(CtxTrace, tid)
    c.Header("X-Trace-Id", tid)
    c.Next()
  }
}
```

`traceId` 随后写入日志字段、事件 Payload、Kafka Header，实现跨进程传播。

## 4. Service 层传递规范

所有 service 方法应以 `context.Context` 作为首参，禁止在 service 层使用 `gin.Context`：

```go
func (s *PostService) CreatePost(ctx context.Context, userID int64, req dto.CreatePostRequest) (*dto.Post, error)
```

api 层从 `gin.Context` 取标准 `context`（`c.Request.Context()`，并注入 traceId/userId）后传入 service。

> 落地说明：当前代码 service 尚未统一携带 ctx。建议按模块渐进改造（先 post/comment/interaction 等写路径），保持一次一模块、可回归。详见 `observability-implementation-guide.md`。

## 5. Repository 层传递规范

所有 GORM 调用使用 `db.WithContext(ctx)`，使 DB span 挂到当前 trace，并让慢 SQL 日志带上 traceId：

```go
func (r *PostRepository) FindByID(ctx context.Context, id int64) (*model.Post, error) {
  var p model.Post
  err := r.db.WithContext(ctx).First(&p, id).Error
  ...
}
```

## 6. Kafka 链路传播

- **生产**：把 `traceId` 同时写入 Kafka Message Header（key=`traceparent` 或 `X-Trace-Id`）与事件体 `Event.TraceID` 字段。
- **消费**：优先从 Header 恢复 trace，缺失则用 `Event.TraceID`；新建 Consumer Span；日志记录 `eventId`、`eventType`、`traceId`。

事件线格式（`internal/event/outbox.go` 的 `Message`）已预留 `traceId` 字段：

```go
type Message struct {
    EventID     string          `json:"eventId"`
    EventType   string          `json:"eventType"`
    AggregateID int64           `json:"aggregateId"`
    UserID      int64           `json:"userId"`
    TraceID     string          `json:"traceId"`
    Payload     json.RawMessage `json:"payload"`
    CreatedAt   time.Time       `json:"createdAt"`
}
```

## 7. Outbox 链路追踪

`event_outbox` 表新增 `trace_id` 字段（本项目用 GORM AutoMigrate，模型加字段即自动建列）：

```sql
ALTER TABLE event_outbox ADD COLUMN trace_id VARCHAR(64) NOT NULL DEFAULT '';
```

写 Outbox 时保存当前请求 `traceId`；Dispatcher 投递时把 `trace_id` 复制进 `Message.TraceID` 与 Kafka Header，Worker 即可延续同一 trace。

## 8. Span 命名规范

| 场景 | Span 名 |
|---|---|
| HTTP | `HTTP POST /api/v1/posts` |
| Service | `PostService.CreatePost` |
| Repository | `PostRepository.CreatePost` |
| MySQL | `MySQL INSERT posts` |
| Redis | `Redis GET post:detail` |
| Kafka 生产 | `Kafka Produce feedora.post.events` |
| Kafka 消费 | `Kafka Consume feedora.post.events` |
| ES | `ES Search feedora_posts` |
| OSS | `OSS Upload post_image` |

## 9. Span 属性规范

| 属性 | 说明 |
|---|---|
| user.id | 当前用户 |
| module / action / biz.id | 业务维度 |
| db.system / db.statement | `mysql` / 脱敏后 SQL |
| messaging.system / messaging.destination | `kafka` / topic |
| http.method / http.route / http.status_code | HTTP 维度 |

## 10. OpenTelemetry 接入建议

1. 初始化 `TracerProvider`（OTLP exporter，采样率生产可设 0.1）。
2. Gin：`otelgin` 中间件生成 root span，并把已有 `traceId` 作为属性。
3. GORM：`gorm.io/plugin/opentelemetry/tracing` 或自定义 callback 生成 DB span。
4. Redis：`redisotel.InstrumentTracing(rdb)`。
5. Kafka：生产/消费处手动建 span 并注入/提取 Header。
6. ES / OSS：在 `internal/search`、`pkg/ossx` 调用处手动建 span。

接入顺序建议：先打通 traceId 贯通（无需 OTel），再逐步接 OTel span，避免一次性大改。
