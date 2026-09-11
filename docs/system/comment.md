---
title: 评论模块功能文档
date: 2026-09-09
version: v1.2
type: system
module: comment
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/comment_router.go
  - internal/api/comment_api.go
  - internal/dto/comment_dto.go
  - internal/service/comment_service.go
  - internal/repository/comment_repository.go
  - internal/model/comment.go
summary: 评论模块提供帖子评论的发表、楼中楼回复、删除（叶子软删、有子回复保留占位节点）、点赞与树形列表查询，评论计数在事务内对称维护，发评论经 Redis 时间窗防重复提交，并通过 comment.events 事件驱动通知、积分与热度榜单。
---

## 1. 模块概述

评论模块为帖子提供两级评论（根评论 + 楼中楼回复）的发表、回复、删除、点赞与列表查询能力；删除语义区分叶子与分支——叶子评论为软删除，有子回复的评论保留为占位节点（仅置 `status = deleted`、不写 `deleted_at`）以维持树结构。评论创建与删除在数据库事务中对称维护 `posts.comment_count` 与 `users.comment_count` 计数，创建经 Redis 时间窗防重复提交，成功后发布 `CommentCreated` 事件，由 worker 异步生成站内通知、发放积分并累加帖子热度榜单。除本模块路由外，`CommentService` 的"我的评论"查询复用为用户端接口，`CommentRepository` 的全量分页查询复用为后台管理接口。

## 2. 接口清单

路由统一挂载于 `/api/v1` 前缀下，注册于 `internal/router/comment_router.go`。

| 路径 | 方法 | 功能 | 控制器 |
|---|---|---|---|
| `/posts/:postId/comments` | GET | 查询帖子评论树（根评论 + 回复），匿名可访问 | `CommentAPI.ListByPost` |
| `/comments` | POST | 对帖子发表根评论 | `CommentAPI.Create` |
| `/comments/:commentId/replies` | POST | 回复某条评论（楼中楼） | `CommentAPI.Reply` |
| `/comments/:commentId` | DELETE | 删除评论（有子回复时保留占位节点，叶子软删除） | `CommentAPI.Delete` |
| `/comments/:commentId/like` | POST | 点赞评论 | `CommentAPI.Like` |
| `/comments/:commentId/like` | DELETE | 取消点赞评论 | `CommentAPI.Unlike` |

除上表外，另有两条复用本模块服务层/仓储层的接口，注册在其他模块路由中：`GET /users/me/comments`（`UserAPI.MyComments`，用户中心查看自己的评论，调用 `CommentService.ListMine`，注册于 `internal/router/user_router.go`）；`GET /admin/comments`（`AdminAPI.Comments`，后台评论列表分页，使用 `CommentRepository.ListAll`/`Count`）。

## 3. 核心逻辑

### 3.1 楼中楼回复结构

- `parent_id` 表示直接父评论 ID，根评论取 `0`；`root_id` 表示所属根评论（一楼）ID，根评论自身取 `0`。
- `Reply` 时计算根评论：取父评论的 `RootID`，若父评论本身是根评论（`RootID == 0`），则根为其自身 ID。因此无论回复根评论还是回复某条子回复，新评论的 `root_id` 均指向同一根评论。
- `reply_to_user_id` 仅记录被回复人 ID，供前端展示"回复 @ 某人"，不参与树结构组织。
- `ListByPost` 一次性查出帖子全部 `status IN (normal, deleted)` 的评论（`created_at` 升序，无分页）；GORM 仍自动排除 `deleted_at` 非空的软删行，因此查询结果中的 `deleted` 记录均为占位节点。先构建根评论数组并建立 `评论 ID → 数组下标` 索引，再将 `parent_id != 0` 的评论按 `root_id`（为 `0` 时回退取 `parent_id`）挂到对应根评论的 `replies` 数组，回复内部同样按创建时间升序。
- 呈现结构固定为两层：所有子回复平铺在根评论的 `replies` 下，不构建第三层及以上嵌套。
- 占位节点（`status = deleted` 的评论）保留在树中以维持楼中楼结构，返回内容固定为「该评论已删除」，`like_count` 归零、`liked` 恒为 `false`；已删除根评论的子回复照常挂载返回，不再被丢弃。
- 前端列表查询为全量返回，无分页参数；`ListMine` 与后台 `ListAll` 均有分页。

### 3.2 评论计数维护（同步）

- 创建评论通过 `CreateWithCounters` 在单个数据库事务内完成三步：插入 `comments` 记录、`posts.comment_count + 1`、`users.comment_count + 1`。计数为同步维护，不存在异步补偿。
- 删除评论时先经 `HasChildren(commentID)` 判断是否存在 `parent_id` 指向该评论的子回复，再调用 `DeleteWithCounters(c, placeholder)` 在单个数据库事务内执行：叶子评论（`placeholder = false`）置 `status = deleted` 并执行 GORM 软删写 `deleted_at`；有子回复的评论（`placeholder = true`）仅置 `status = deleted`（同时刷新 `updated_at`），不写 `deleted_at`，保留为占位节点。两种路径均对称回减 `posts.comment_count` 与 `users.comment_count`（以 `GREATEST(col - 1, 0)` 为下限 0），与创建时 `CreateWithCounters` 的 +1/+1 对称。
- 删除有子回复的评论不级联删除子回复：该评论以占位节点形式留在列表树中，子回复照常可见；`comment_count` 仅回减被删除评论自身，子回复计数保留，占位节点不占评论计数。
- 帖子被删除时不存在对本模块评论数的级联处理逻辑。

### 3.3 发表与删除校验

- 发表前校验：内容去除首尾空白后非空（仅 `required` 校验，无应用层长度上限，长度上限由 `content` 列 `size:2000` 约束）；目标帖子存在且 `status` 为 `PostPublished`，否则返回 `ErrPostInvisible`。
- 发表防重复提交：创建前以 `IdemCommentKey` 生成键 `idem:comment:{userId}:{hash(postId|content)}` 执行 Redis `SetNX`，10 秒窗口内同用户同帖同内容的重复提交返回 429（`errs.ErrRateLimit`）；窗口占用后任何校验失败（帖子不存在 / 不可见、参数错误等）会删除该键释放窗口允许立即重试；Redis 未启用时 `SetNX` 恒成功，不拦截。
- 回复前校验父评论存在（`FindByID` 受 GORM 软删过滤：叶子软删评论不可查返回 `ErrCommentNotFound`；占位节点 `deleted_at` 为空、可查到，回复已删除的占位父评论仍被允许，以维持楼中楼线索）。
- 删除权限：评论作者本人（`m.UserID == userID`）或角色为 `admin` 的用户可删除；帖子作者无权删除他人对自己帖子的评论。

### 3.4 事件发送

- 仅创建评论发送事件：topic 为 `comment.events`（常量 `event.TopicComment`），事件类型为 `CommentCreated`，`aggregateId` 为评论 ID，`userId` 为评论者 ID，payload 为 `{postId, postAuthorId}`。
- 实际投递取决于全局装配：`kafka.enabled` 为真时经 `OutboxProducer` 写入 `event_outbox` 表，由 Outbox Dispatcher 异步投递 Kafka，实际 topic 为 `<topicPrefix>.comment.events`（前缀默认 `feedora`）；未启用 Kafka 时使用 `NoopProducer`，仅打印日志。
- worker 侧消费 `CommentCreated` 的下游（每个消费者独立幂等去重）：

| 消费者 | 行为 |
|---|---|
| notification | 向帖子作者生成"收到新的评论"站内通知并累加未读数；评论者即帖子作者时跳过 |
| growth | 为评论者发放 `create_comment` 积分 +3 |
| rank | 将所属帖子热度 ZSet（today/week/all 三个周期）+5 |

- 搜索消费者不处理 `CommentCreated`，评论不写入 ES 索引。
- 删除评论、点赞、取消点赞均不发布事件。

### 3.5 点赞与缓存

- 点赞记录写入 `comment_likes` 表，`comment_id + user_id` 唯一索引保证幂等：插入实际生效（`RowsAffected > 0`）时才执行 `like_count + 1`；取消点赞对称处理，负增量以 `GREATEST(like_count - 1, 0)` 为下限 0。
- 点赞前校验评论存在：叶子软删评论被 GORM 过滤、返回 `ErrCommentNotFound`；占位节点可被查到，点赞照常执行（列表返回时其 `like_count` 归零、`liked` 为 `false`）。
- 列表查询通过 `CommentLikedSet` 批量回填当前查看者的点赞状态；匿名访问（`GET /posts/:postId/comments` 无鉴权，`CurrentUserID` 返回 0）时所有评论的 `liked` 恒为 `false`。
- 创建与删除评论成功后删除该帖子的详情缓存（`cache.PostDetailKey(postID)`）；点赞不清理缓存。

### 3.6 我的评论

`ListMine` 分页返回当前用户 `status = normal` 的评论（`created_at` 倒序），批量回查帖子标题；帖子不存在（已软删）时标题显示为"原帖已删除"。

## 4. 数据模型

表 `comments`（`internal/model/comment.go`）：

| 字段 | 类型/约束 | 说明 |
|---|---|---|
| `id` | 主键 | 评论 ID |
| `post_id` | 索引 `idx_post_parent`、`idx_post_created` | 所属帖子 ID |
| `user_id` | 索引 `idx_user_created` | 评论者 ID |
| `parent_id` | 默认 `0`，索引 `idx_post_parent` | 直接父评论 ID，根评论为 `0` |
| `root_id` | 默认 `0`，索引 | 所属根评论 ID，根评论自身为 `0` |
| `reply_to_user_id` | 可空 | 被回复用户 ID |
| `content` | `size:2000` | 评论内容 |
| `status` | `size:32`，默认 `normal` | `normal` / `deleted` |
| `like_count` | 默认 `0` | 点赞数冗余计数 |
| `created_at` / `updated_at` | 时间戳 | 创建/更新时间 |
| `deleted_at` | GORM 软删标记，索引 | 非空即已软删 |

删除语义：叶子评论删除时同时写 `status = deleted` 与 `deleted_at` 两套标记；有子回复的评论仅置 `status = deleted`、不写 `deleted_at`，作为占位节点保留（树形结构下双标志的唯一例外）。查询侧：GORM 自动排除 `deleted_at` 非空记录，`ListByPost` 以 `status IN (normal, deleted)` 查询，两者叠加后返回「正常评论 + 占位节点」集合；`ListMine` 仍只取 `status = normal`。点赞关系表 `comment_likes` 不随评论删除而清理。

## 5. 配置项与限制

本模块无自有配置项。受以下全局配置影响：`kafka.enabled` 决定事件走 Outbox 投递还是 Noop 日志，`kafka.topicPrefix` 决定实际 topic 前缀（默认 `feedora`）。

当前实现的限制：

- 帖子评论列表一次性返回全部评论，无分页参数，评论量大时响应体积随之增长。
- 评论内容仅在应用层校验非空，长度上限依赖数据库列约束（`size:2000`）。
- 删除评论不发送评论删除事件，下游通知与热度榜单不回退。
- 不支持评论编辑。

## 6. 历史版本

| 版本 | 日期 | 维护者 | 说明 |
|---|---|---|---|
| v1.2 | 2026-09-09 | Feedora 项目组 | 删除语义改为：有子回复的评论保留占位节点（仅置 status、不设 deleted_at），列表含占位并显示「该评论已删除」；新增发评论 10 秒时间窗防重复提交 |
| v1.1 | 2026-09-09 | Feedora 项目组 | 修正删除路径描述：软删除事务化并补齐 `users.comment_count` 回减 |
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
