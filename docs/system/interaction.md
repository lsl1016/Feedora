---
title: 互动模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: interaction
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/interaction_router.go
  - internal/api/interaction_api.go
  - internal/dto/interaction_dto.go
  - internal/service/interaction_service.go
  - internal/repository/interaction_repository.go
  - internal/model/interaction.go
summary: 帖子与评论的点赞、收藏关系维护，含唯一索引幂等写、计数同步增减、interaction.events 事件发布与我的点赞/收藏列表查询。
---

## 1. 模块概述

互动模块提供帖子点赞、帖子收藏与评论点赞三类互动关系的建立与取消。帖子点赞/收藏由本模块的 4 个接口承接；评论点赞的 HTTP 接口注册在评论模块，但关系表的读写复用本模块仓储。目标类型不存在运行时枚举：关系按目标类型落在 `post_likes`、`post_favorites`、`comment_likes` 三张独立表，不存在统一 `interactions` 表。每次操作同步返回目标帖子的最新点赞数、收藏数与当前用户的互动状态。

## 2. 接口清单

路由实际前缀为 `/api/v1`，下表省略该前缀。前 4 条注册于 `registerInteraction`，均挂认证中间件；后 4 条操作同一份互动数据，注册于 `user_router` 与 `comment_router`。

| 路径 | 方法 | 功能 | 控制器 |
| --- | --- | --- | --- |
| `/posts/:postId/like` | POST | 点赞帖子（需登录） | `InteractionAPI.Like` |
| `/posts/:postId/like` | DELETE | 取消点赞帖子（需登录） | `InteractionAPI.Unlike` |
| `/posts/:postId/favorite` | POST | 收藏帖子（需登录） | `InteractionAPI.Favorite` |
| `/posts/:postId/favorite` | DELETE | 取消收藏帖子（需登录） | `InteractionAPI.Unfavorite` |
| `/users/me/liked-posts` | GET | 我点赞的帖子列表（需登录） | `UserAPI.MyLikedPosts` |
| `/users/me/favorite-posts` | GET | 我收藏的帖子列表（需登录） | `UserAPI.MyFavoritePosts` |
| `/comments/:commentId/like` | POST | 点赞评论（需登录，业务在评论模块） | `CommentAPI.Like` |
| `/comments/:commentId/like` | DELETE | 取消点赞评论（需登录，业务在评论模块） | `CommentAPI.Unlike` |

## 3. 核心逻辑

### 幂等写：唯一索引 + RowsAffected

幂等由数据库唯一索引保证，不是先查后插。`InteractionRepository.AddPostLike` 直接 `Create`，重复点赞命中唯一索引 `uk_post_like_user` 后 `RowsAffected=0`，返回 `false`；仅返回 `true` 时才递增计数并发布事件。取消按 `post_id + user_id` 物理删除，以 `RowsAffected > 0` 判断是否实际删除，重复取消不扣减计数。收藏与评论点赞同构，复用同一套 `Add/Remove` + 布尔返回模式。

前置校验：`SetLike` 先 `FindByID` 加载帖子，不存在返回 `ErrPostNotFound`；`SetFavorite` 用 `Exists` 仅做存在性判断，不加载帖子。

### 计数同步维护

计数在请求内同步更新 MySQL，各动作用的列如下（`UpdateColumn` 原子表达式，扣减侧用 `GREATEST(col - ?, 0)` 兜底不为负）：

| 动作 | `posts.like_count` | `users.like_count`（作者总量） | `posts.favorite_count` | `comments.like_count` |
| --- | --- | --- | --- | --- |
| 点赞 / 取消点赞帖子 | +1 / -1 | +1 / -1 | — | — |
| 收藏 / 取消收藏帖子 | — | — | +1 / -1 | — |
| 点赞 / 取消点赞评论 | — | — | — | +1 / -1 |

关系插入与计数更新不在同一事务中执行。每次操作后删除缓存 Key `post:detail:{postID}`（`cache.PostDetailKey`），未启用 Redis 时空操作；点赞数/收藏数本身不写 Redis，不存在 Redis 计数缓存。

每次操作返回 `InteractionResult`：重新 `FindByID` 读取帖子最新计数，另执行 `HasLiked`/`HasFavorited` 两个 `COUNT` 查询取当前用户状态。

### 事件发布与消费

事件发布到 topic `interaction.events`（`event.TopicInteraction`），实际 topic 为 `{topicPrefix}.interaction.events`；`bizID` 为帖子 ID，`userID` 为操作人，`payload` 为 `nil`。

| 动作 | 事件类型 |
| --- | --- |
| 点赞帖子 | `PostLiked` |
| 取消点赞 | `PostUnliked` |
| 收藏帖子 | `PostFavorited` |
| 取消收藏 | 不发布（不存在 `PostUnfavorited` 事件类型） |
| 评论点赞/取消 | 不发布事件 |

生产者由 `kafka.enabled` 决定：启用时为 `OutboxProducer`，事件先写 `event_outbox` 表，由 Worker 的 Outbox Dispatcher 定时批量投递 Kafka，失败按 10/30/60/300/600 秒退避重试，超过 `Worker.MaxRetry` 标记失败；未启用时为 `NoopProducer`，仅打印事件日志。

Worker 消费 `{prefix}.interaction.events`，逐条分发给四个消费者，各自以 `idem.Claim(eventID, worker)` 独立幂等。互动事件的真实消费行为：

| 消费者 | `PostLiked` | `PostFavorited` | `PostUnliked` |
| --- | --- | --- | --- |
| notification | 向帖子作者写"收到新的点赞"通知（作者与操作人相同时跳过），`Incr` `notify:unread:{authorID}` | "收到新的收藏"，同左 | 不处理 |
| growth | 作者 +2 积分，动作 `post_liked` | 作者 +5 积分，动作 `post_favorited` | 不处理 |
| rank | `rank:post:{today,week,all}` ZSet +3 分 | 同左 +4 分 | 不处理（不回扣分数） |
| search | 不处理 | 不处理 | 不处理 |

### counter 消费者现状

`internal/worker/counter` 包仅含占位说明（"阶段二占位"），`worker.go` 包注释将"counter：Redis 计数回写 MySQL"列为规划项。当前不存在计数异步回写消费者，点赞/收藏计数的全部维护路径均为上述同步 MySQL 更新。

### 我的点赞/收藏列表

`GET /users/me/liked-posts` 与 `GET /users/me/favorite-posts` 经 `UserService.postsByRelation` 查询：`InteractionRepository.PostIDsByUser` 一次性取出该用户在 `post_likes` 或 `post_favorites` 中的全部 `post_id`（按关系行 `id DESC`，SQL 层无 LIMIT），`total` 取 ID 数量，分页在内存中切片完成，再由 `AssembleByIDs` 装配帖子。

帖子列表与详情中"当前用户是否点赞/收藏"的标记由帖子模块调用本模块 `LikedSet`/`FavoritedSet` 以 `IN` 批量查询；`userID <= 0`（未登录）时直接返回空集合。

## 4. 数据模型

按目标类型分三张独立表，结构同构，字段均为 `id` 主键、目标 ID、`user_id`、`created_at`：

| 表 | 模型 | 唯一索引 | 单列索引 | 覆盖动作 |
| --- | --- | --- | --- | --- |
| `post_likes` | `model.PostLike` | `uk_post_like_user(post_id, user_id)` | `user_id` | 帖子点赞 |
| `post_favorites` | `model.PostFavorite` | `uk_post_fav_user(post_id, user_id)` | `user_id` | 帖子收藏 |
| `comment_likes` | `model.CommentLike` | `uk_comment_like_user(comment_id, user_id)` | `user_id` | 评论点赞（评论无收藏） |

表内无类型列、无状态列，语义为"存在即关系"：建立是插入一行，取消是物理删除该行，无软删除。不存在统一 `interactions` 表，也不存在目标类型枚举字段，目标类型由表名静态区分。

## 5. 配置项与限制

无新增配置项，模块不读取独立配置。行为受以下全局配置影响：

| 配置 | 影响 |
| --- | --- |
| `kafka.enabled` | `true` 走 Outbox + Kafka 可靠投递；`false` 事件仅打印日志（Noop 模式），Worker 消费链路无互动事件进入 |
| `kafka.topicPrefix` | 实际 topic 前缀；`OutboxProducer` 缺省 `feedora`，`NoopProducer` 缺省 `community` |
| Redis 是否启用 | 未启用时缓存删除为空操作，rank 消费者整体跳过，通知未读数不累加 |

当前实现的边界与否定事实：

- 取消收藏不发布事件；`PostUnliked` 虽被发布但四个消费者均不处理，取消点赞不回扣积分与榜单分数。
- 评论点赞不发布事件、不影响作者 `users.like_count`、不失效任何缓存。
- 关系写入与计数更新无事务保证，两者可能不一致；不存在对账或修复任务。
- 我的点赞/收藏列表每次请求全量加载关系 ID 后内存分页，关系量大时开销线性增长。
- counter 消费者未实现，不存在 Redis 计数回写 MySQL 的路径。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
| --- | --- | --- | --- |
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
