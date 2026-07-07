# 日志记录规范

适用于 Feedora（开发者知识社区 V2.1）Go 后端。技术栈：Gin + GORM + MySQL + Redis + Kafka + Elasticsearch + OSS/MinIO + JWT。分层：`router → api → service → repository(model) → cache/event/worker/search`，基础设施在 `pkg`。

日志打印统一使用格式化占位符，禁止字符串拼接：

```go
logger.Infof("create post success, postId:%d", postID)
```

## 1. 日志设计目标

- 问题排查：任意请求可凭 `traceId` 串起 API → Service → Repository → Kafka → Worker → ES 全链路。
- 接口耗时分析：每个 HTTP 请求记录 `durationMs`、`statusCode`。
- 用户行为追踪：关键业务动作记录 `userId`、`module`、`action`、`bizId`。
- Worker 消费排查：记录 `eventId`、`eventType`、`workerName`、`retryCount`、`status`。
- Kafka 事件追踪：Outbox 投递与消费均携带 `eventId`、`traceId`。
- 慢查询定位：MySQL/Redis/ES/OSS/Kafka 慢操作按阈值打 WARN。

## 2. 日志级别规范

| 级别 | 使用场景 | 示例 |
|---|---|---|
| DEBUG | 本地调试细节，生产默认关闭 | `logger.Debugf("assemble post relations, postId:%d, tags:%d", id, n)` |
| INFO | 正常业务完成、请求完成、事件投递/消费成功 | `logger.Infof("http request completed, path:%s, status:%d, durationMs:%d", path, code, ms)` |
| WARN | 可继续但需关注：慢操作、幂等拦截、重试、降级、参数校验失败 | `logger.Warnf("slow sql, durationMs:%d, sql:%s", ms, sql)` |
| ERROR | 业务失败、外部依赖失败、panic 恢复 | `logger.Errorf("consume event failed, eventId:%s, err:%v", id, err)` |
| PANIC/FATAL | 启动失败、不可恢复错误（仅启动阶段） | 配置加载失败、DB 无法连接直接退出 |

原则：能被 Recover 的运行期错误一律用 ERROR，不要 panic；只有进程无法启动才 FATAL。

## 3. 统一日志字段规范

| 字段 | 说明 | 必填 |
|---|---|---|
| timestamp | RFC3339 时间 | 是 |
| level | 日志级别 | 是 |
| message | 事件描述（英文短句） | 是 |
| traceId | 全链路追踪 ID | HTTP/Worker 是 |
| spanId | 跨度 ID（接入 OTel 后） | 否 |
| userId | 当前用户，未登录为 0 | HTTP 尽量 |
| requestId | 客户端请求 ID（X-Request-Id） | 否 |
| method / path | HTTP 方法与路由 | HTTP 是 |
| statusCode | HTTP 状态码 | HTTP 是 |
| durationMs | 处理耗时（毫秒） | HTTP/慢操作 是 |
| clientIp / userAgent | 客户端信息 | HTTP 尽量 |
| module | 业务模块：auth/user/post/... | 业务日志是 |
| action | 业务动作：create_post/login/... | 业务日志是 |
| bizId | 业务对象 ID（postId/circleId 等） | 业务日志尽量 |
| errorCode / errorMessage | 业务错误码与信息 | 错误日志是 |

## 4. HTTP 请求日志规范

由 `pkg/middleware/logger.go` 统一输出，每个请求一条 INFO（慢请求升级为 WARN）：

```json
{
  "level": "info",
  "traceId": "trace-xxx",
  "userId": 10001,
  "method": "POST",
  "path": "/api/v1/posts",
  "statusCode": 200,
  "durationMs": 35,
  "clientIp": "127.0.0.1",
  "module": "post",
  "action": "create_post",
  "message": "http request completed",
  "timestamp": "2026-07-07T10:00:00+08:00"
}
```

## 5. 业务日志规范（按模块）

只记录关键节点，避免每行都打日志。

- **auth**：注册成功/失败、登录成功/失败（失败记账号，不记密码）、封禁用户登录拦截。
- **user**：资料更新成功。
- **post**：发帖成功（postId）、编辑、隐藏/取消隐藏、删除、转发、圈子发帖权限拦截（WARN）。
- **comment**：评论/回复成功（commentId, postId）、删除。
- **interaction**：点赞/收藏状态变更（幂等未变更时可 DEBUG）。
- **circle**：创建、加入/退出、成员禁言/解禁/移除、权限拦截（WARN）。
- **search**：搜索关键词与命中数、ES 查询失败（ERROR）。
- **notification**：通知生成成功（userId, type, targetId）。
- **rank**：榜单 ZSet 更新（DEBUG/INFO）。
- **admin**：所有后台写操作记 INFO（adminId, action, targetType, targetId），并落 `operation_logs` 表。
- **file**：上传成功（size, mimeType, url）、校验失败（WARN）、上传失败（ERROR）。
- **worker**：见第 9 节。

示例：

```go
logger.Infof("create post success, module:post, action:create_post, userId:%d, postId:%d", userID, postID)
logger.Warnf("circle post denied, module:circle, action:post, userId:%d, circleId:%d", userID, circleID)
```

## 6. 错误日志规范

错误日志必须包含：`traceId`、`userId`、`module`、`action`、`bizId`、`errorCode`、`errorMessage`、`error`（原始 err）；可选 `stack`（仅 panic）、脱敏后的请求摘要。

```go
logger.Errorf("create post failed, module:post, action:create_post, userId:%d, errorCode:%d, err:%v", userID, code, err)
```

约定：Service 返回 `pkg/errors.Error`（含 code/message），api 层在 `response.Fail` 时由中间件统一记录 errorCode/errorMessage，业务层只在“非预期错误”处补充 ERROR 日志，避免重复打印。

## 7. 敏感信息脱敏规范

以下字段禁止明文打印，需脱敏或省略：

`password`、`passwordHash`、`token`、`Authorization`、`apiKey`、`secretKey`、OSS secret、JWT secret；`email`/`phone` 部分脱敏（`z***@x.com`、`138****0000`）。

- 请求日志不打印 body；确需打印时先脱敏并截断（≤256 字符）。
- 日志中间件不记录 `Authorization` 头。
- 打印用户对象时使用只含公开字段的摘要，切勿直接 `%+v` 整个 model.User。

## 8. 慢操作日志规范（WARN 阈值）

| 操作 | 阈值 | 级别 |
|---|---|---|
| HTTP 接口 | > 500ms | WARN |
| MySQL 查询 | > 200ms | WARN |
| Redis 命令 | > 50ms | WARN |
| ES 查询 | > 300ms | WARN |
| OSS 上传 | > 1000ms | WARN |
| Kafka 生产/消费 | > 1000ms | WARN |

```go
logger.Warnf("slow http, path:%s, durationMs:%d", path, ms)
logger.Warnf("slow sql, durationMs:%d, rows:%d, sql:%s", ms, rows, sql)
```

## 9. Worker 日志规范

每条事件消费一条日志，必须包含：`traceId`、`eventId`、`eventType`、`topic`、`partition`、`offset`、`workerName`、`retryCount`、`durationMs`、`status`、`errorMessage`。

```go
logger.Infof("consume event success, workerName:%s, eventType:%s, eventId:%s, traceId:%s, durationMs:%d",
    workerName, eventType, eventID, traceID, ms)
logger.Errorf("consume event failed, workerName:%s, eventId:%s, retryCount:%d, err:%v",
    workerName, eventID, retry, err)
```

幂等拦截记 DEBUG/INFO：`logger.Infof("event skipped by idempotency, workerName:%s, eventId:%s", w, id)`。

## 10. Kafka Outbox 日志规范

Outbox Dispatcher 每轮扫描输出汇总，单条投递输出明细：

```go
logger.Infof("outbox dispatch batch, scanned:%d, success:%d, failed:%d", scanned, ok, fail)
logger.Infof("outbox dispatched, eventId:%s, topic:%s, retryCount:%d", id, topic, retry)
logger.Warnf("outbox dispatch retry, eventId:%s, retryCount:%d, nextRetryAt:%s", id, retry, next)
logger.Errorf("outbox dispatch failed permanently, eventId:%s", id)
```

## 11. 日志落盘建议

| 环境 | 输出 | 采集 |
|---|---|---|
| 本地 | console（可读文本） | 无 |
| 测试 | console + file | 可选 |
| 生产 | stdout（JSON） | Docker/K8s → Loki 或 ELK 采集 |

生产使用 JSON 结构化输出，容器只写 stdout，由采集侧统一处理，应用不自行管理文件轮转。

## 12. 日志示例代码（Zap 封装目标）

当前 `pkg/logger` 为轻量封装（`Infof/Warnf/Errorf`）。生产建议升级为 Zap，接口保持不变，调用点无需改动：

```go
// pkg/logger/logger.go（目标：Zap 封装）
package logger

import "go.uber.org/zap"

var sugar *zap.SugaredLogger

func Init(env string) {
    var l *zap.Logger
    if env == "local" {
        l, _ = zap.NewDevelopment()
    } else {
        l, _ = zap.NewProduction() // JSON 输出到 stdout
    }
    sugar = l.Sugar()
}

func Infof(format string, args ...any)  { sugar.Infof(format, args...) }
func Warnf(format string, args ...any)  { sugar.Warnf(format, args...) }
func Errorf(format string, args ...any) { sugar.Errorf(format, args...) }

// WithTrace 返回带 traceId 的 logger（接入结构化字段时使用）
func WithTrace(traceID string) *zap.SugaredLogger { return sugar.With("traceId", traceID) }
```

Gin 请求日志中间件（现有 `pkg/middleware/logger.go`，含慢请求 WARN）：

```go
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        ms := time.Since(start).Milliseconds()
        f := logger.Infof
        if ms > 500 { f = logger.Warnf }
        f("http request completed, traceId:%s, userId:%d, method:%s, path:%s, status:%d, durationMs:%d, clientIp:%s",
            c.GetString(CtxTrace), CurrentUserID(c), c.Request.Method, c.Request.URL.Path,
            c.Writer.Status(), ms, c.ClientIP())
    }
}
```

Worker 消费日志：

```go
start := time.Now()
// ... 处理 ...
logger.Infof("consume event success, workerName:%s, eventType:%s, eventId:%s, traceId:%s, durationMs:%d",
    "search", m.EventType, m.EventID, m.TraceID, time.Since(start).Milliseconds())
```
