---
title: 基础设施模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: infrastructure
maintainer: Feedora 项目组
status: active
related_code:
  - internal/app/app.go
  - internal/worker/worker.go
  - internal/worker/handlers.go
  - internal/event/outbox.go
  - internal/cache/cache.go
  - internal/cache/keys.go
  - internal/repository/outbox_repository.go
  - pkg/kafkax/producer.go
summary: worker/cache/event/outbox/pkg 横切基础设施总览，覆盖事件可靠投递链路、Worker 消费体系、Redis 缓存键与可插拔组件的降级实现。
---

## 1. 模块概述

本文档描述 Feedora 后端的横切基础设施：`internal/worker`（Outbox Dispatcher 与事件消费）、`internal/cache`（Redis 缓存封装）、`internal/event`（领域事件定义与生产者）、`event_outbox`/`worker_event_records` 两张基础设施表，以及 `pkg/` 下的 kafkax、redisx、esx、ossx、middleware、observability、config 等公共包。各业务模块文档（post/comment/follow 等）只引用本文档描述的事件链路与缓存键，不重复展开。事件链路当前的默认实现为 Outbox 落库加 Kafka 异步投递，Redis、Elasticsearch、Kafka 均可通过 `enabled` 配置开关降级运行。

## 2. 组件清单

| 组件 | 位置 | 职责 |
| --- | --- | --- |
| `App` 应用容器 | `internal/app/app.go` | 加载配置、连接 MySQL/Redis/ES、初始化本地存储与 JWT，按 `kafka.enabled` 选择事件生产者实现并装配全部服务与路由 |
| API 进程入口 | `cmd/api/main.go` | 构建 `App` 并监听 `server.port`；进程启动时执行 `AutoMigrate`（含基础设施表建表） |
| Worker 进程入口 | `cmd/worker/main.go` | 构建 `Runner`，在 `:9091` 暴露 `/metrics`；`kafka.enabled=false` 时直接退出 |
| `Runner` | `internal/worker/worker.go` | Worker 运行器：预创建 Topic、ES 索引全量重建、启动 Outbox Dispatcher 与 Kafka 消费循环 |
| Outbox Dispatcher | `internal/worker/worker.go`（`dispatchLoop`/`dispatchOnce`） | 定时捞取 `event_outbox` 中 `pending` 事件发布到 Kafka，失败按退避序列重试 |
| 事件分发 | `internal/worker/handlers.go` | 将单条 `event.Message` 串行分发给 search、notification、growth、rank 四个处理器，各自独立幂等 |
| ES 索引操作 | `internal/worker/index.go` | 构造帖子/用户/圈子/话题索引文档并写入，按 ID 重建或删除，启动时全量重建 |
| search 消费者 | `internal/worker/search/`（占位包，逻辑在 `handlers.go`/`index.go`） | 消费 `feedora.post.events`、`feedora.user.events`、`feedora.circle.events` 同步 ES 索引 |
| notification 消费者 | `internal/worker/notification/`（占位包，逻辑在 `handlers.go`） | 消费点赞/收藏/评论/加圈/关注事件生成站内通知并累加未读计数 |
| growth 消费者 | `internal/worker/growth/`（占位包，逻辑在 `handlers.go`） | 消费发帖/评论/被赞/被藏/建圈/加圈事件发放积分 |
| rank 消费者 | `internal/worker/rank/`（占位包，逻辑在 `handlers.go`） | 消费内容与圈子事件累加热榜 ZSet 分数 |
| consumer/counter/dispatcher/stat 包 | `internal/worker/{consumer,counter,dispatcher,stat}/` | 仅含 `doc.go` 占位注释，无任何实现（counter 的 Redis 计数回写 MySQL 未实现） |
| `cache.Cache` | `internal/cache/cache.go` | Redis 轻封装：JSON 读写、整型计数、ZSet 榜单，底层 client 为 `nil` 时全部方法安全降级 |
| 缓存键生成 | `internal/cache/keys.go` | 统一生成详情缓存、未读计数、榜单、热门搜索词的 Redis Key |
| `event.Producer`/`NoopProducer` | `internal/event/producer.go` | 事件发布接口与日志实现（Kafka 未启用时使用，仅打印事件 JSON） |
| `OutboxProducer`/`event.Message` | `internal/event/outbox.go` | 事件写入 `event_outbox` 表的可靠投递实现；`Message` 为 Kafka 线格式（Worker 反序列化目标） |
| 事件类型与 Topic 常量 | `internal/event/event_type.go`、`internal/event/topic.go` | 13 个事件类型常量与 5 个 Topic 常量（不含前缀） |
| `OutboxRepository` | `internal/repository/outbox_repository.go` | `event_outbox` 表的写入、可投递捞取、标记投递成功/重试/失败 |
| `WorkerRepository` | `internal/repository/worker_repository.go` | 基于 `worker_event_records` 唯一索引的幂等 `Claim` 与 `Release` |
| `kafkax` | `pkg/kafkax/` | 基于 segmentio/kafka-go 的生产者、消费组消费者与 Topic 预创建 |
| `redisx`/`esx`/`ossx` | `pkg/{redisx,esx,ossx}/` | Redis/ES 客户端初始化（含连通性校验）；文件存储抽象与本地磁盘实现（MinIO/阿里云为占位） |
| `observability` | `pkg/observability/` | Prometheus 指标定义、Gin 指标中间件、`/metrics` 处理器 |
| `middleware` | `pkg/middleware/` | Trace（traceId 传播）、Logger、Recover、CORS、Auth/Admin、上下文键；`RateLimit` 为直接放行的占位实现 |
| `config` | `pkg/config/config.go` | 全量 YAML 配置结构定义与 `Load` 加载 |

## 3. 核心逻辑

### 3.1 事件链路

业务服务在写库成功后调用 `producer.Publish(topic, eventType, bizID, userID, payload)` 发布事件。事件链路分两段：

第一段（API 进程内）：`internal/app/app.go` 按 `kafka.enabled` 选择实现。启用时使用 `OutboxProducer`：生成 `utils.UUID()` 作为 `event_id`，将 `topicPrefix + "." + topic`（`topicPrefix` 为空时回退 `"feedora"`）、事件类型、`aggregate_id`、`user_id`、序列化后的 Payload 写入 `event_outbox` 表，`status` 置 `pending`。该写入是与业务写库相互独立的两次数据库操作，不在同一事务内：业务写库成功而 outbox 写入失败时事件丢失。`trace_id` 字段未在写入时赋值，落库与投递日志中的 traceId 恒为空字符串。未启用 Kafka 时使用 `NoopProducer`：仅以 INFO 日志打印 `topicPrefix.topic` 与事件 JSON（前缀为空时回退 `"community"`），事件不落 outbox、不进 Kafka。

第二段（Worker 进程内）：`Runner.dispatchLoop` 按 `worker.outboxIntervalSeconds` 周期执行 `dispatchOnce`，每次最多取 `worker.batchSize` 条。捞取条件为 `status = 'pending' AND (next_retry_at IS NULL OR next_retry_at <= now)`，按 `id ASC` 排序。每条记录组装为 `event.Message`（字段 `eventId`、`eventType`、`aggregateId`、`userId`、`traceId`、`payload`、`createdAt`）序列化后，以 `event_id` 作为消息 Key 发布到 Kafka（`kafkax.Producer`，`AllowAutoTopicCreation=true`，LeastBytes 均衡）。发布成功调用 `MarkDispatched`（`status=dispatched` 并写 `dispatched_at`）并计数 `feedora_kafka_produce_total{result=success}`；失败则 `retry_count+1` 并按退避序列 `[10, 30, 60, 300, 600]` 秒（超出序列取 600）设置 `next_retry_at`，重试次数超过 `worker.maxRetry` 时 `MarkFailed`（`status=failed`）。`failed` 为终态，无恢复路径。

真实 Topic（前缀 `feedora.`）与事件类型映射：

| Topic 常量 | 实际 Topic 名 | 发布的事件类型 |
| --- | --- | --- |
| `TopicUser` | `feedora.user.events` | `UserRegistered`、`UserFollowed`、`UserUnfollowed` |
| `TopicPost` | `feedora.post.events` | `PostCreated`、`PostUpdated`、`PostHidden`、`PostUnhidden`、`PostDeleted` |
| `TopicComment` | `feedora.comment.events` | `CommentCreated` |
| `TopicInteraction` | `feedora.interaction.events` | `PostLiked`、`PostUnliked`、`PostFavorited` |
| `TopicCircle` | `feedora.circle.events` | `CircleCreated`、`CircleJoined` |

事件类型枚举定义于 `internal/event/event_type.go`，共 13 个常量；`PostUnhidden` 无常量定义，在 `post_service.go` 发布处与 `handlers.go` 消费处均为字符串字面量。事件 Payload 阶段一使用通用 `map[string]any`（如发帖事件 `{"title": ...}`、评论事件 `{"postId": ..., "postAuthorId": ...}`），`internal/event/payload.go` 中的结构化 Payload 定义为占位注释。

消费侧幂等由 `WorkerRepository.Claim` 实现：向 `worker_event_records` 插入 `(event_id, worker_name, status="success")`，依赖 `uk_event_worker` 唯一索引，插入成功即首次消费，重复插入返回 false 拦截。

### 3.2 Worker 体系

`cmd/worker/main.go` 启动流程：加载配置 → `kafka.enabled=false` 时记录告警并退出 → 连接 MySQL（不执行 `AutoMigrate`，建表由 API 进程完成）→ 按开关连接 Redis（失败仅记错误并以降级模式运行）与 ES → `worker.New` 构造 `Runner`（内部创建 kafkax 生产者/消费者，订阅带前缀的 5 个 Topic）→ 在 `:9091` 起独立 HTTP 服务暴露 `/metrics` → 监听 SIGINT/SIGTERM 后进入 `Run`。

`Run` 依次执行：`kafkax.EnsureTopics` 预创建 5 个 Topic（单分区单副本，失败仅告警）；ES 启用时 `EnsureIndices` 建索引并 `reindexAll` 全量重建（帖子取 `status <> PostDeleted`，用户取 `status <> UserBanned`，圈子取 `status = "normal"`（复用成员状态常量 `CircleMemberNormal`），话题取 `status = "enabled"`）；随后启动 `dispatchLoop` 协程并进入消费循环。消费循环每次 `Fetch` 一条消息，JSON 反序列化为 `event.Message` 失败则直接提交位点跳过；成功则经 `handle` 串行调用四个处理器后提交位点（手动提交，处理完成后才 commit）。

四个处理器（均在 `internal/worker/handlers.go`，逐事件串行执行）：

| 处理器 | 消费的事件 | 行为 |
| --- | --- | --- |
| `handleSearch` | `PostCreated/PostUpdated/PostHidden/PostDeleted/PostUnhidden`、`UserRegistered`、`CircleCreated` | 帖子事件按 ID 重建索引（帖子不存在或已删除时从索引删除文档）；用户/圈子事件回查后写入 `feedora_posts`、`feedora_users`、`feedora_circles`、`feedora_topics` 索引（前缀 `elasticsearch.indexPrefix` 加下划线拼接） |
| `handleNotification` | `PostLiked`、`PostFavorited`、`CommentCreated`、`CircleJoined`、`UserFollowed` | 通知接收者为帖子作者/圈主/被关注者；操作者即接收者时不通知。写入 `notifications` 记录（`target_url` 形如 `/posts/{id}`、`/circles/{id}`、`/users/{id}`）并对 `notify:unread:{接收者}` 执行 `Incr` |
| `handleGrowth` | `PostCreated`、`CommentCreated`、`PostLiked`、`PostFavorited`、`CircleCreated`、`CircleJoined` | 调用 `AddPointLog` 发放积分：发帖 +10、评论 +3、被点赞 +2、被收藏 +5、建圈 +20、加圈 +1 |
| `handleRank` | `PostCreated`、`PostLiked`、`PostFavorited`、`CommentCreated`、`CircleJoined` | 帖子维度分值 +1/+3/+4/+5（评论作用于其所属帖子），累加 `rank:post:{today|week|all}`；加圈 +2 累加 `rank:circle:{today|week|all}`；Redis 未启用时直接返回 |

局限性：`internal/worker/` 下 consumer、counter、dispatcher、growth、notification、rank、search、stat 八个子包均为仅含 `doc.go` 的占位包，实际逻辑集中在 `worker.go`、`handlers.go`、`index.go` 三个文件；counter（Redis 计数回写 MySQL）与 stat（统计处理）无任何实现。`WorkerRepository.Release` 已定义但无调用点，处理器失败时不删除幂等记录，该事件不会被重新处理。`observability.WorkerProcess` 指标已定义但未被调用。消费循环为单 reader 串行处理，未按 Topic 或处理器并发。`worker.enabled` 配置键无代码读取，Worker 进程是否运行仅由 `kafka.enabled` 决定。

### 3.3 缓存体系

`internal/cache/cache.go` 的 `Cache` 封装 Redis 客户端，所有方法在底层 client 为 `nil`（未启用 Redis）时安全降级：读操作视为未命中、写操作为空操作、计数返回 0，业务逻辑自动回源数据库。缓存键统一在 `internal/cache/keys.go` 生成：

| Redis Key | 生成函数 | 数据形态 | TTL | 失效与更新方式 |
| --- | --- | --- | --- | --- |
| `post:detail:{postID}` | `PostDetailKey` | JSON 字符串 | 10 分钟（`postDetailTTL`） | 帖子更新、点赞/取消收藏、评论新增或删除时 `Del`，下次读回源重建 |
| `user:profile:{userID}` | `UserProfileKey` | JSON 字符串 | 30 分钟（`userProfileTTL`） | 用户资料更新时 `Del`；关注/取关时双向 `Del` 两个用户 |
| `circle:detail:{circleID}` | `CircleDetailKey` | JSON 字符串 | 15 分钟（`circleDetailTTL`） | 圈子更新时 `Del` |
| `notify:unread:{userID}` | `NotifyUnreadKey` | 整型计数 | 无 TTL | Worker 生成通知时 `Incr`；当前无任何读取方，未读数接口 `NotificationService.UnreadCount` 直接查询数据库 |
| `rank:{post\|circle}:{today\|week\|all}` | `RankKey` | ZSet | 无 TTL | Worker 事件驱动 `ZIncr`；读侧 `ZTop` 取分页 Top，ZSet 为空时回退数据库按 `hot_score` 排序 |
| `rank:user:{type}:{range}` | `UserRankKey` | ZSet | — | 已定义，当前无任何调用方 |
| `search:hot_keywords` | `HotKeywordsKey` 常量 | ZSet | 无 TTL | 每次执行搜索时 `ZIncr`；热门搜索词读取 Top 10 |

防缓存击穿的 singleflight 未实现：`pkg/singleflightx/singleflight.go` 为仅含包注释的占位文件。`observability.CacheHit`/`CacheMiss` 指标已定义但无调用点，缓存命中率不产生指标。`Cache.Decr` 提供计数递减（当前值不大于 0 时不递减），无业务调用方。

### 3.4 可插拔实现

| 组件 | 降级实现 | 启用真实实现的方式 | 当前状态 |
| --- | --- | --- | --- |
| 事件生产者 | `NoopProducer`：打印事件日志 | `kafka.enabled=true` 时切换为 `OutboxProducer`，接口签名一致，业务代码无感知 | 两实现均可用，默认配置启用 Outbox |
| Kafka 客户端 | API 进程不创建 Kafka 连接（Noop 模式）；`cmd/worker` 在 `kafka.enabled=false` 时退出 | `kafka.enabled=true` | `pkg/kafkax` 基于 segmentio/kafka-go，完整实现 |
| Elasticsearch | `search.Client` 为 `nil`：API 侧搜索服务以 nil 客户端运行，Worker 侧跳过全部索引逻辑（`r.sc == nil` 短路） | `elasticsearch.enabled=true`，`esx.New` 内以 `client.Info()` 校验连通性 | 完整实现，默认启用 |
| OSS 存储 | `LocalStorage` 本地磁盘实现（写 `oss.basePath`，经 `oss.publicBaseUrl` 对外提供静态访问） | `pkg/ossx/minio.go`、`pkg/ossx/aliyun.go` 为占位注释包，接口已抽象但实现未编写 | 阶段一固定使用 `LocalStorage`；`oss.type` 配置键存在但 `app.go` 不读取 |
| Redis | `cache.New(nil)`，全部方法空操作/未命中 | `redis.enabled=true`，`redisx.New` 内以 `Ping` 校验连通性 | 完整实现，默认启用 |

`pkg/middleware/RateLimit` 同为占位实现：直接 `c.Next()` 放行，未基于 Redis 做用户/IP 维度限流。

## 4. 数据模型

两张基础设施表经 `model.AllModels()` 注册，由 API 进程启动时的 `AutoMigrate` 建表。

`event_outbox`（`model.EventOutbox`）：

| 字段 | 说明 |
| --- | --- |
| `id` | 主键，投递顺序按此列 ASC |
| `event_id` | 事件唯一 ID（UUID），唯一索引 `uk_event_id`，同时作为 Kafka 消息 Key |
| `topic` | 含前缀的完整 Topic 名，如 `feedora.post.events` |
| `event_type` | 事件类型，如 `PostCreated` |
| `aggregate_id` | 聚合根 ID（帖子/评论/圈子/用户 ID，语义随事件类型变化） |
| `user_id` | 触发者用户 ID，默认 0 |
| `trace_id` | 链路追踪 ID，当前恒为空字符串 |
| `payload` | JSON 文本 |
| `status` | 状态，默认 `pending`，复合索引 `idx_status_next_retry` |
| `retry_count` | 已重试次数，默认 0 |
| `next_retry_at` | 下次可投递时间，可为 NULL，复合索引 `idx_status_next_retry` |
| `dispatched_at` | 投递成功时间，可为 NULL |
| `created_at`/`updated_at` | 创建/更新时间，`created_at` 单列索引 |

状态机：`pending` →（投递成功）`dispatched`（终态，写 `dispatched_at`）；`pending` →（重试次数超过 `worker.maxRetry`）`failed`（终态）。重试期间保持 `pending`，仅递增 `retry_count` 并推迟 `next_retry_at`。无 `failed` 重置或人工重放路径。

`worker_event_records`（`model.WorkerEventRecord`）：

| 字段 | 说明 |
| --- | --- |
| `id` | 主键 |
| `event_id` + `worker_name` | 联合唯一索引 `uk_event_worker`，幂等键 |
| `status` | 状态，默认 `success` |
| `error_message` | 错误信息，当前无写入方 |
| `created_at` | 创建时间，复合索引 `idx_worker_created` |

状态机：`Claim` 插入即 `success`（单状态）；`Release` 物理删除记录（无调用方）。不存在 `failed` 状态的写入路径。

模型注释将两表标注为"阶段一预留、业务尚未写入"，实际在 `kafka.enabled=true` 时 `OutboxProducer` 与幂等 `Claim` 均已持续写入。

## 5. 配置项

配置文件为 `configs/config.yaml`，结构定义在 `pkg/config/config.go`，以下列出与基础设施相关的真实配置键及当前值：

| 配置键 | 当前值 | 用途 |
| --- | --- | --- |
| `server.name` / `server.env` / `server.port` | `feedora-api` / `local` / `8090` | API 进程监听端口与静态文件基地址 |
| `mysql.dsn` | `root:root@tcp(localhost:3307)/feedora...` | 数据库连接（Docker 映射 3307） |
| `mysql.maxOpenConns` / `mysql.maxIdleConns` / `mysql.connMaxLifetimeSeconds` | `50` / `10` / `3600` | 连接池参数 |
| `kafka.enabled` | `true` | 开关：API 进程选择 Outbox/Noop 生产者，Worker 进程是否运行 |
| `kafka.brokers` | `localhost:9092` | Kafka broker 列表 |
| `kafka.topicPrefix` | `feedora` | Topic 前缀，拼接为 `feedora.{name}`；空值时生产者回退 `feedora`、Noop 回退 `community` |
| `kafka.consumerGroup` | `feedora-worker` | Worker 消费组 ID |
| `redis.enabled` | `true` | 缓存开关 |
| `redis.addr` / `redis.password` / `redis.db` | `localhost:6379` / 空 / `0` | Redis 连接 |
| `redis.poolSize` / `redis.minIdleConns` | `20` / `5` | 连接池参数 |
| `elasticsearch.enabled` | `true` | 搜索开关 |
| `elasticsearch.addresses` | `http://localhost:9200` | ES 地址 |
| `elasticsearch.username` / `elasticsearch.password` | 空 | ES 认证 |
| `elasticsearch.indexPrefix` | `feedora` | 索引前缀，生成 `feedora_posts` 等 4 个索引 |
| `worker.enabled` | `true` | 已定义但无代码读取 |
| `worker.batchSize` | `100` | Dispatcher 单批捞取上限（`<=0` 时回退 100） |
| `worker.outboxIntervalSeconds` | `2` | Dispatcher 轮询间隔秒数（`<=0` 时回退 2） |
| `worker.maxRetry` | `5` | 投递最大重试次数，超过标记 `failed` |
| `oss.type` | `local` | 已定义但无代码读取，固定实例化 `LocalStorage` |
| `oss.basePath` / `oss.publicBaseUrl` | `./uploads` / `http://localhost:8090/static` | 本地存储根目录与公开访问基地址 |

Worker 指标端口 `:9091` 与消费组订阅的 5 个 Topic 名为代码内固定值，不在配置文件中。API 进程的 `/metrics` 由路由注册在主服务端口 `/metrics` 路径。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
| --- | --- | --- | --- |
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
