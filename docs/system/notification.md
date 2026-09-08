---
title: 通知模块功能文档
date: 2026-09-09
version: v1.1
type: system
module: notification
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/notification_router.go
  - internal/api/notification_api.go
  - internal/dto/notification_dto.go
  - internal/service/notification_service.go
  - internal/repository/notification_repository.go
  - internal/model/notification.go
  - internal/worker/handlers.go
  - internal/cache/cache.go
summary: 站内通知的列表、已读标记与未读数读路径，写入由 Worker 消费 Kafka 业务事件（点赞、收藏、评论、加圈、关注）异步落库 notifications 表，未读数以 Redis 缓存旁路（cache-aside，10 分钟 TTL）维护。
---

## 1. 模块概述

通知模块提供站内通知的读路径：分页列表、单条已读、全部已读与未读数统计，共 4 个接口，全部需要登录。通知的写入不在 HTTP 请求内完成：业务方（互动、评论、圈子、关注）发布领域事件，经 Outbox 投递 Kafka，由 Worker 的 notification 消费者异步生成通知并落库 `notifications` 表。通知不可删除，已读状态仅支持 `unread` 到 `read` 的单向流转。

## 2. 接口清单

路由实际前缀为 `/api/v1`，下表省略该前缀。4 条路由均注册于 `registerNotification`，均挂认证中间件 `x.authMW`。

| 路径 | 方法 | 功能 | 控制器 |
| --- | --- | --- | --- |
| `/notifications` | GET | 分页查询当前用户通知（可按 category 过滤） | `NotificationAPI.List` |
| `/notifications/unread-count` | GET | 当前用户未读通知数 | `NotificationAPI.UnreadCount` |
| `/notifications/read-all` | PUT | 当前用户全部通知标记已读 | `NotificationAPI.MarkAllRead` |
| `/notifications/:notificationId/read` | PUT | 标记单条通知已读 | `NotificationAPI.MarkRead` |

## 3. 核心逻辑

### 通知的产生路径（异步消费）

通知由 Kafka Worker 异步生成，链路为：业务服务写库成功后调用 `event.Producer.Publish` → `OutboxProducer` 将事件写入 `event_outbox` 表（topic 列已含前缀，如 `feedora.interaction.events`）→ Worker 的 Outbox Dispatcher 定时批量投递 Kafka（失败按 10/30/60/300/600 秒退避重试，超过 `worker.maxRetry` 标记 failed）→ Worker 消费后写入 `notifications` 表。

- `kafka.enabled=false` 时生产者为 `NoopProducer`，事件仅打印日志，不产生任何通知。
- Worker（`worker.Runner`）以消费组订阅 5 个 topic：`{topicPrefix}.user.events`、`{topicPrefix}.post.events`、`{topicPrefix}.comment.events`、`{topicPrefix}.interaction.events`、`{topicPrefix}.circle.events`。消息反序列化为 `event.Message` 后分发给 search / notification / growth / rank 四个消费者，逐条提交 offset。
- notification 消费者的实现是 `internal/worker/handlers.go` 中的 `Runner.handleNotification`；`internal/worker/notification` 包本身仅含占位 `doc.go`（"阶段二占位"），不含任何实现代码。

`handleNotification` 覆盖 5 种事件，其余事件类型一律不处理：

| 事件类型 | 发布方 | topic（不含前缀） | AggregateID / UserID 语义 | 通知接收人 | 自通知抑制 |
| --- | --- | --- | --- | --- | --- |
| `PostLiked` | `InteractionService.SetLike` | `interaction.events` | 帖子 ID / 点赞人 | 帖子作者 | 作者 = 点赞人时跳过 |
| `PostFavorited` | `InteractionService.SetFavorite` | `interaction.events` | 帖子 ID / 收藏人 | 帖子作者 | 作者 = 收藏人时跳过 |
| `CommentCreated` | `CommentService.Create` | `comment.events` | 评论 ID / 评论人 | 帖子作者 | 作者 = 评论人时跳过 |
| `CircleJoined` | `CircleService.Join` | `circle.events` | 圈子 ID / 加入人 | 圈主 `OwnerID` | 圈主 = 加入人时跳过 |
| `UserFollowed` | `FollowService.Follow` | `user.events` | 被关注者 ID / 关注人 | 被关注者 | 被关注者 = 关注人时跳过 |

生成的通知内容为固定文案：

| 事件 | type（分类） | 标题 | 内容 | targetType / targetURL |
| --- | --- | --- | --- | --- |
| `PostLiked` | `interaction` | 收到新的点赞 | 有人点赞了你的帖子《帖子标题》 | `post` / `/posts/{postID}` |
| `PostFavorited` | `interaction` | 收到新的收藏 | 有人收藏了你的帖子《帖子标题》 | `post` / `/posts/{postID}` |
| `CommentCreated` | `interaction` | 收到新的评论 | 有人评论了你的帖子《帖子标题》 | `post` / `/posts/{postID}` |
| `CircleJoined` | `circle` | 圈子有新成员 | 有新成员加入了你的圈子「圈子名」 | `circle` / `/circles/{circleID}` |
| `UserFollowed` | `follow` | 收到新的关注 | {关注人昵称} 关注了你 | `user` / `/users/{关注人ID}` |

消费者不使用消息 `Payload`，统一按 `AggregateID` / `UserID` 回查数据库（帖子、评论、圈子、用户）取标题与作者后组装通知。写入成功后删除（`Del`）未读数缓存 Key `notify:unread:{接收人ID}`，下次未读数查询回源数据库重建缓存（Redis 未启用时空操作）。

### 幂等与失败处理

`handleNotification` 对每个事件先回查实体并做自通知抑制判断，再调用 `idem.Claim(eventID, "notification")`：向 `worker_event_records` 表插入一条 `(eventID, "notification")` 记录，依赖唯一索引 `uk_event_worker`，重复消费插入失败返回 `false`，直接跳过。首次占用后才写通知。

`notify` 内 `Insert` 失败仅记录错误日志并返回，不删除已占用的幂等记录（未调用 `WorkerRepository.Release`）；该事件即使重投也会被 `Claim` 拦截，对应通知不再生成。

### 已读状态流转

`read_status` 仅在 `unread` 与 `read` 两个值之间单向流转，不存在已读回退为未读的操作。

- `MarkRead`：按 `id + user_id + read_status = 'unread'` 条件更新为 `read`，重复标记或操作他人通知均不生效；仓储返回的 `RowsAffected > 0` 布尔值被 Service 层丢弃，接口恒返回成功。更新执行后无论影响行数如何，均删除未读数缓存 Key（`invalidateUnread`）。
- `MarkAllRead`：将当前用户全部 `unread` 记录一次性置为 `read`，重复调用为空操作；执行后同样删除未读数缓存 Key。

### 列表查询与分类过滤

`List` 按 `user_id` 过滤、`ORDER BY id DESC` 分页返回。查询参数 `category` 映射为 `type` 列的精确匹配过滤；`category` 为空字符串或 `all` 时不过滤；取值不做枚举校验，任意字符串直接作为 SQL 条件。`NotificationRepository.List` 返回 `(rows, total)`，API 层以 `list, _ :=` 丢弃 `total` 与错误，响应体 `data` 仅含数组，不返回总条数。

DTO 转换（`dto.ToNotificationItem`）将模型的 `Type` 字段直接输出为前端 `category`，`ReadStatus`、`TargetURL` 原样透传。

### 未读数维护

未读数采用 Redis cache-aside 模式，Key 为 `notify:unread:{userID}`（`cache.NotifyUnreadKey`），TTL 10 分钟（常量 `unreadCacheTTL`，位于 `internal/service/notification_service.go`）：

- 读取（`GET /notifications/unread-count`）：`Cache.GetIntOK` 命中（Key 存在且可解析为整数）时直接返回缓存值；未命中时执行数据库 `COUNT(user_id = ? AND read_status = 'unread')` 并 `SetInt` 写回缓存后返回。
- 失效：Worker `notify` 落库成功后、`MarkRead` / `MarkAllRead` 执行后，均 `Del` 该 Key，下次读取回源重建。
- `redis.enabled=false` 时 `GetIntOK` 恒未命中、`SetInt` / `Del` 为空操作，接口每次请求执行数据库 COUNT，与无缓存行为一致。
- 升级前由旧 `Incr` 逻辑写入的存量 Key 无 TTL：命中后原样返回历史累加值，直到该用户下一次触发 `Del`（收到通知或标记已读）才回源重建。

## 4. 数据模型

单表 `notifications`（`model.Notification`，`TableName()` 固定为 `notifications`）：

| 字段 | 类型 | 说明 |
| --- | --- | --- |
| `id` | int64 | 主键 |
| `user_id` | int64 | 接收人 ID |
| `actor_id` | int64 | 触发人 ID，默认 0 |
| `type` | varchar(64) | 通知分类，取值 `interaction` / `circle` / `follow` |
| `title` | varchar(200) | 通知标题 |
| `content` | varchar(1000) | 通知内容 |
| `target_type` | varchar(64) | 关联目标类型，取值 `post` / `circle` / `user` |
| `target_id` | int64 | 关联目标 ID |
| `target_url` | varchar(512) | 前端跳转链接 |
| `read_status` | varchar(32) | 已读状态，默认 `unread`，取值 `unread` / `read` |
| `created_at` / `updated_at` | datetime | 创建 / 更新时间 |

索引：复合索引 `idx_user_read_created(user_id, read_status, created_at)`、`idx_user_created(user_id, created_at)`。

类型无 Go 枚举常量：`type`、`target_type`、`read_status` 的取值均以字符串字面量散落在 Worker 的 `handleNotification`、`notify` 与仓储 SQL 条件中（如 `"interaction"`、`"unread"`），数据库层无枚举约束。模型注释仍写着"阶段一预留表结构，暂不产生数据"，与当前 Worker 已实际写入的现状不符，属过时注释。

## 5. 配置项与限制

本模块无独立配置，行为受以下全局配置影响（当前值取自 `configs/config.yaml`）：

| 配置 | 当前值 | 影响 |
| --- | --- | --- |
| `kafka.enabled` | `true` | `true` 走 Outbox + Kafka，Worker 生成通知；`false` 为 `NoopProducer`，不产生任何通知 |
| `kafka.topicPrefix` | `feedora` | 实际 topic 前缀，如 `feedora.interaction.events`；`OutboxProducer` 缺省值 `feedora` |
| `kafka.consumerGroup` | `feedora-worker` | Worker 消费组 ID |
| `worker.batchSize` / `worker.outboxIntervalSeconds` / `worker.maxRetry` | `100` / `2` / `5` | Outbox 批量投递与重试 |
| `redis.enabled` | `true` | 控制未读数缓存（cache-aside）；`false` 时未读数每次请求数据库 COUNT，不影响通知落库 |

当前实现的边界与否定事实：

- `PostUnliked`、`UserUnfollowed`、`PostCreated`、`CircleCreated`、`UserRegistered` 等事件虽被发布，notification 消费者均不处理，不产生通知。
- 回复评论不通知被回复人：`CommentCreated` 仅通知帖子作者，评论的 `ReplyToUserID` 不参与通知生成。
- 通知无删除接口、无清理任务，数据只增不减；已读不可回退。
- 未读数缓存最多滞后 10 分钟（TTL 兜底），缓存删除失败时以 TTL 过期为准；`Cache.Incr` / `Decr` 在通知链路不再被调用。
- 通知 `Insert` 失败时幂等记录不释放，该事件对应通知永久丢失，仅留错误日志。
- `internal/worker/notification` 包为占位，无实现代码；实际消费逻辑位于 `internal/worker/handlers.go`。
- 列表接口丢弃 `total` 与查询错误，前端无法获得总条数。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
| --- | --- | --- | --- |
| v1.1 | 2026-09-09 | Feedora 项目组 | 未读数改为 Redis 缓存旁路（10 分钟 TTL），落库与已读后失效缓存 |
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
