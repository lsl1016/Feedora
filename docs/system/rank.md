---
title: 热榜模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: rank
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/rank_router.go
  - internal/api/rank_api.go
  - internal/dto/rank_dto.go
  - internal/service/rank_service.go
  - internal/worker/rank/doc.go
  - internal/worker/handlers.go
  - internal/cache/keys.go
summary: 热门榜单（帖子/圈子/话题）与用户/圈子排行榜的事件驱动 ZSet 计分、Redis 读取与 DB 回退实现。
---

## 1. 模块概述

rank 模块提供热门榜单与排行榜两类只读查询能力。热门榜单（帖子 / 圈子 / 话题）优先读取 Redis ZSet，ZSet 为空或 Redis 未启用时回退到 DB 排序，保证页面始终有数据。榜单分数由 Worker 消费 Kafka 业务事件后以加权增量方式写入 ZSet，接口本身不触发计算。用户/圈子排行榜由 `RankService.Rankings` 提供，经成长模块路由 `/growth/rankings` 暴露，当前仅走 DB 排序，不读 Redis。

## 2. 接口清单

| 路径 | 方法 | 功能 | 控制器 |
| --- | --- | --- | --- |
| `/hot/ranks` | GET | 热门榜单，按 `rankType`（post/circle/topic）与 `timeRange` 分页查询 | `RankAPI.HotRanks` |
| `/growth/rankings` | GET | 用户/圈子排行榜，按 `type`（active/contribution/creator/circle）分页查询；由 `growth_router.go` 注册，`GrowthService.Rankings` 委托 `RankService.Rankings` | `GrowthAPI.Rankings` |

两个路由均为公开路由，未挂鉴权中间件。完整前缀为 `/api/v1`。

## 3. 核心逻辑

### 3.1 榜单维度

- 热门榜单 `RankService.HotRanks(rankType, timeRange, page, size)`：`rankType` 取 `post` / `circle` / `topic`，缺省与未知取值均落入 `post` 分支；`timeRange` 缺省为 `today`。
- 排行榜 `RankService.Rankings(rankType, timeRange, page, size, currentUserID)`：`rankType` 取 `active`（缺省）/ `contribution` / `creator` / `circle`；`timeRange` 缺省为 `all`，该参数被接收但未参与任何排序逻辑。

### 3.2 分数计算（事件加权增量）

计算逻辑位于 `internal/worker/handlers.go` 的 `Runner.handleRank`，由事件消费触发，非定时任务。每个事件按固定权重对 ZSet 执行 `ZIncrBy`：

| 事件 | 目标榜单 | member | 增量 |
| --- | --- | --- | --- |
| `PostCreated` | post 榜 | 帖子 ID | +1 |
| `PostLiked` | post 榜 | 帖子 ID | +3 |
| `PostFavorited` | post 榜 | 帖子 ID | +4 |
| `CommentCreated` | post 榜 | 评论所属帖子 ID（`comments.FindByID` 反查） | +5 |
| `CircleJoined` | circle 榜 | 圈子 ID | +2 |

帖子与圈子事件均同时写入 `today` / `week` / `all` 三个 timeRange 的 ZSet，写入值完全相同，不存在时间窗衰减或衰减公式。每个事件在处理前通过 `idem.Claim(m.EventID, "rank")` 做消费者级幂等去重；Redis 未启用（`cch.Enabled()` 为 false）时 `handleRank` 直接返回，不写任何 ZSet。

### 3.3 计算周期与触发方式

纯事件驱动：业务方写 Outbox 表 → `Runner.dispatchLoop` 每 `worker.outboxIntervalSeconds`（默认 2 秒）批量投递 Kafka → Worker 消费循环调用 `handle` → `handleRank` 更新 ZSet。模块内不存在定时重算、每日清零或全量回补任务。

### 3.4 结果存储

结果只存 Redis ZSet，不落 DB 表。Key 由 `cache.RankKey(rankType, timeRange)` 生成，格式 `rank:{rankType}:{timeRange}`，实际被写入的 key 为：

- `rank:post:today` / `rank:post:week` / `rank:post:all`
- `rank:circle:today` / `rank:circle:week` / `rank:circle:all`

member 为目标 ID 的十进制字符串，score 为事件加权累计分。ZSet 写入时不设置 TTL，也不存在按日重置逻辑，`today` / `week` 实际为自首次写入以来的累计值，并非真实时间窗口。

### 3.5 接口读取路径

`GET /hot/ranks` → `RankAPI.HotRanks`（绑定 `dto.HotRankQuery`，分页参数补默认值：page 1、pageSize 10、上限 100）→ `RankService.HotRanks`：

1. 调用 `cache.ZTop`（底层 `ZRevRangeWithScores`）按分数倒序取当前页 member 与 score，member 解析为 ID；
2. ZSet 命中时按 `rankType` 走 `buildHot` 组装：post 批量查 `posts.FindByIDs` + `users.FindByIDs` 取作者昵称；circle 批量查 `circles.FindByIDs`；topic 逐 ID 调用 `topics.FindByID`（逐条查询）；`Rank = offset + 序号 + 1`，`Score` 取 ZSet 分数；
3. ZSet 为空（含 Redis 未启用、`timeRange` 为非既有关键字的情况）时回退 `hotFromDB` 直接查 DB：
   - post：`status = PostPublished`，`ORDER BY hot_score DESC, like_count DESC, created_at DESC`，`Score = hot_score`；
   - circle：`ORDER BY member_count DESC, post_count DESC`，`Score = member_count`；
   - topic：`status = StatusEnabled`，`ORDER BY participant_count DESC, post_count DESC`，`Score = participant_count`。

`GET /growth/rankings` → `GrowthAPI.Rankings` → `GrowthService.Rankings` → `RankService.Rankings`，不读任何 ZSet，直接 DB 排序：

| type | 过滤 | 排序 | Score |
| --- | --- | --- | --- |
| `active`（缺省） | `users.status = UserNormal` | `point_count DESC, post_count DESC` | `point_count` |
| `creator` / `contribution` | `users.status = UserNormal` | `post_count DESC, like_count DESC` | `post_count` |
| `circle` | 圈子 status 过滤 | `member_count DESC` | `member_count` |

用户榜返回等级、积分、发帖数、获赞数，并以 `currentUserID` 标记 `isCurrentUser`。

## 4. 数据模型

本模块不新增数据表。存储结构如下：

- Redis ZSet：`rank:post:{today|week|all}`、`rank:circle:{today|week|all}`，member 为目标 ID 十进制字符串，score 为事件加权累计分，无 TTL。
- `internal/cache/keys.go` 另定义 `UserRankKey` 生成 `rank:user:{type}:{range}`（如 `rank:user:creator:all`），当前无任何调用方，属于未使用的预留 key。
- DB 回退读取依赖既有表与列：`posts`（`hot_score`、`like_count`、`comment_count`、`status`）、`users`（`point_count`、`post_count`、`like_count`、`level`、`status`）、`circles`（`member_count`、`post_count`、`featured_count`、`status`）、`topics`（`participant_count`、`post_count`、`status`）。
- `posts.hot_score` 列带索引，业务代码中不存在写入路径，仅 `cmd/seed` 种子数据赋值，线上数据保持默认值 0。

## 5. 配置项与限制

配置读取自 `configs/config.yaml`，模块无专属配置节，依赖以下全局配置：

| 配置 | 默认值（config.yaml） | 对本模块的影响 |
| --- | --- | --- |
| `redis.enabled` | `true` | 关闭后 ZSet 不写不读，热门榜单与排行榜全部走 DB 回退 |
| `kafka.enabled`、`kafka.brokers`、`kafka.consumerGroup` | `true`、`localhost:9092`、`feedora-worker` | 控制榜单事件的消费通道 |
| `worker.enabled`、`worker.batchSize`、`worker.outboxIntervalSeconds`、`worker.maxRetry` | `true`、`100`、`2`、`5` | 决定事件投递节奏与榜单更新时效 |

不存在榜单计算周期类配置，因为不存在定时计算。当前实现的限制与否定事实：

- `internal/worker/rank/` 目录仅含 `doc.go` 占位文件，实际的榜单计算逻辑在 `internal/worker/handlers.go` 的 `handleRank` 中。
- 话题榜单无任何写入方：不存在 `rank:topic:*` 的 ZSet 写入，`rankType=topic` 恒走 DB 回退。
- 接口接受任意 `timeRange` 字符串拼入 key；非 `today` / `week` / `all` 的取值对应空 ZSet，结果同样来自 DB 回退，回退排序不区分时间范围。
- `today` / `week` / `all` 三个 ZSet 写入值相同且无 TTL、无重置机制，时间范围参数对实际数据不产生差异。
- `posts.hot_score` 无业务写入路径，post 榜 DB 回退中 `hot_score` 恒为 0（种子数据除外），实际排序主要由 `like_count`、`created_at` 决定。
- 用户/圈子排行榜不使用 `UserRankKey` ZSet，无 Redis 加速，每次请求直接查 DB。
- ZSet 命中路径下 topic 榜单逐 ID 查询 `topics.FindByID`，存在 N+1 查询。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
| --- | --- | --- | --- |
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
