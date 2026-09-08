---
title: 成长体系模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: growth
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/growth_router.go
  - internal/api/growth_api.go
  - internal/dto/growth_dto.go
  - internal/service/growth_service.go
  - internal/repository/growth_repository.go
  - internal/model/growth.go
  - internal/worker/handlers.go
summary: 提供每日签到、基于 Kafka 事件的积分异步发放、静态成长任务列表与领取占位接口，以及基于用户统计字段排序的成长排行榜。
---

## 1. 模块概述

成长体系模块提供每日签到、积分发放、成长任务与成长排行榜四类能力。签到在请求路径上同步发放固定积分；发帖、评论、点赞、收藏、建圈、加圈六类行为的积分由 Worker 消费 Kafka 事件后异步发放，依赖数据库唯一索引保证幂等。任务列表为服务层硬编码的静态任务加用户计数近似进度，任务领取接口为占位实现，不发放积分。模块不含等级成长逻辑，用户等级注册后固定为 1。

## 2. 接口清单

路由注册于 `internal/router/growth_router.go`，全部挂载在 `/api/v1` 分组下。下表路径省略 `/api/v1` 前缀。

| 路径 | 方法 | 功能 | 控制器 |
|------|------|------|--------|
| `/growth/check-in` | POST | 每日签到，同日重复签到不重复发放积分 | `GrowthAPI.CheckIn` |
| `/growth/tasks` | GET | 成长任务列表，支持 `type` 类型过滤 | `GrowthAPI.Tasks` |
| `/growth/tasks/:taskId/claim` | POST | 领取任务奖励（占位实现，恒返回成功） | `GrowthAPI.ClaimTask` |
| `/growth/rankings` | GET | 成长排行榜（用户榜与圈子榜） | `GrowthAPI.Rankings` |

共 4 个接口。前 3 个接口经 `authMW` 认证，用户 ID 取自 `middleware.CurrentUserID(c)`；`/growth/rankings` 为公开接口，未挂认证中间件，匿名访问时无法标记 `isCurrentUser`。

## 3. 核心逻辑

### 3.1 每日签到

`GrowthService.CheckIn` 执行以下流程：

1. 计算当日日期键：`year*10000 + month*100 + day`（如 2026-09-09 得 `20260909`）。
2. 调用 `GrowthRepository.AddPointLog(userID, "check_in", 5, "checkin", 日期键, "每日签到")` 写入积分流水。表 `user_point_logs` 的唯一索引 `uk_user_action_biz(user_id, action, biz_type, biz_id)` 使同一用户同一天只能写入一条签到流水；插入冲突（或失败）时返回 `false`，本次发放积分为 0；插入成功后以 `point_count + 5` 原子累加 `users.point_count`。
3. 返回当前积分（`users.point_count`）、本次发放积分（0 或 5）、连续签到天数。

签到奖励为固定值 5 分（常量 `checkInPoints`），不随连续天数递增，无连签额外奖励，无补签功能。

连续签到天数为近似值：`CountRecentCheckIns(userID, 7)` 统计最近 7 天（`created_at >= now-7d`）内 `action = "check_in"` 的流水条数，上限 7。该统计不校验日期是否连续，漏签一天不清零，返回的是"最近 7 天签到次数"而非严格连续天数。

### 3.2 积分获取途径与事件消费

除签到外，积分由 Worker 消费 Kafka 事件后异步发放。消费逻辑位于 `internal/worker/handlers.go` 的 `Runner.handleGrowth`（worker 名 `growth`）。Worker 以消费组 `feedora-worker` 订阅 5 个 topic，实际名称由配置 `kafka.topicPrefix`（当前值 `feedora`）拼接：`feedora.user.events`、`feedora.post.events`、`feedora.comment.events`、`feedora.interaction.events`、`feedora.circle.events`。

各事件的积分发放规则：

| 事件类型 | action | biz_type | 积分 | 受益用户 |
|----------|--------|----------|------|----------|
| `PostCreated` | `create_post` | `post` | +10 | 帖子作者 |
| `CommentCreated` | `create_comment` | `comment` | +3 | 评论者 |
| `PostLiked` | `post_liked` | `post` | +2 | 帖子作者 |
| `PostFavorited` | `post_favorited` | `post` | +5 | 帖子作者 |
| `CircleCreated` | `create_circle` | `circle` | +20 | 圈子创建者（事件 `UserID`） |
| `CircleJoined` | `join_circle` | `circle` | +1 | 加入者（事件 `UserID`） |

幂等保障为两层：先用 `WorkerRepository.Claim(eventID, "growth")` 以 `uk_event_worker` 唯一索引拦截同一事件的重复消费，再由 `user_point_logs` 的 `uk_user_action_biz` 唯一索引兜底，同一 `(user_id, action, biz_type, biz_id)` 组合只发放一次。

事件消费不区分操作者与受益者是否为同一用户：用户给自己的帖子点赞、收藏同样给作者（本人）加分。`PostUnliked`、`UserUnfollowed` 等逆向事件不触发积分回收，`PostUpdated`、`PostDeleted`、`PostHidden`、`UserRegistered`、`UserFollowed` 不触发积分发放。当前实现不存在任何积分扣减路径，所有发放值均为正数。

### 3.3 成长任务

`GrowthService.Tasks` 返回服务层硬编码的 3 个静态任务，进度用用户表冗余计数近似，不查任务表：

| TaskID | 标题 | 类型 | 奖励积分 | 目标值 | 进度来源 |
|--------|------|------|----------|--------|----------|
| 1 | 发布首篇帖子 | `newbie` | 10 | 1 | `min(users.post_count, 1)` |
| 2 | 参与评论（发表 3 条） | `daily` | 3 | 3 | `min(users.comment_count, 3)` |
| 3 | 每日签到 | `daily` | 5 | 1 | 恒为 0，状态恒为 `todo` |

任务状态仅有 `done` / `todo` 两种，任务 1、2 依据计数是否达到目标值判定。`type` 参数为空或 `all` 时返回全部任务，否则按类型精确匹配过滤；接口无分页。

`POST /growth/tasks/:taskId/claim` 为占位实现：仅校验路径参数 `taskId`（必填，最小 1），随后恒返回 `{"claimed": true}`，不校验任务完成状态，不写入任何数据，不发放积分。

### 3.4 等级

用户表含 `level` 字段（`int`，默认 1），注册时写入 1。当前代码不存在任何更新 `level` 的写入路径，等级不随积分或行为增长，恒为注册值 1。等级名称映射（`dto.LevelNameOf`）为：`level >= 10` 资深专家、`>= 6` 活跃达人、`>= 3` 进阶用户、其余新手上路；该映射仅在排行榜与用户信息展示时使用。

经验值为近似：`experience` 直接取 `users.point_count`，`nextLevelExperience` 取 `level * 100`。不存在独立经验曲线公式，也不存在升级触发逻辑。

### 3.5 排行榜

`GrowthService.Rankings` 全权委托 `RankService.Rankings`，数据来源为数据库直查（`users`、`circles` 表），不使用 Redis ZSet（ZSet 仅用于热门榜单，与本接口无关）。默认参数：`type = "active"`、`range = "all"`。

| type | 数据源与排序 | 得分 |
|------|--------------|------|
| `active`（默认，含未知类型） | `users`（`status = normal`）按 `point_count DESC, post_count DESC` | `point_count` |
| `creator` / `contribution` | `users`（`status = normal`）按 `post_count DESC, like_count DESC` | `post_count` |
| `circle` | `circles`（正常状态）按 `member_count DESC` | `member_count` |

用户榜项包含昵称、头像、等级、等级名、积分、发帖数、获赞数，并以 `isCurrentUser` 标记当前登录用户。`range` 参数被接受但未参与任何查询条件，所有时间范围返回相同结果。

## 4. 数据模型

模块仅一张自有表 `user_point_logs`（模型 `model.UserPointLog`），同时读写 `users` 表的计数字段。

`user_point_logs` 字段：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | bigint | 主键 |
| `user_id` | bigint | 所属用户 |
| `action` | varchar(64) | 积分动作：`check_in`、`create_post`、`create_comment`、`post_liked`、`post_favorited`、`create_circle`、`join_circle` |
| `point` | int | 积分值，当前实现均为正数 |
| `biz_type` | varchar(64) | 业务类型：`checkin`、`post`、`comment`、`circle` |
| `biz_id` | bigint | 业务 ID；签到为日期键（`YYYYMMDD`），其余为帖子/评论/圈子 ID |
| `remark` | varchar(255) | 备注 |
| `created_at` | datetime | 创建时间 |

索引：`uk_user_action_biz(user_id, action, biz_type, biz_id)` 唯一索引防重复发放；`idx_user_created(user_id, created_at)` 支撑签到次数统计。

`users` 表关联字段：`point_count`（积分余额，默认 0，发放时原子累加）、`post_count`、`comment_count`、`like_count`（任务进度与排行依据）、`level`（默认 1）。

不存在独立的签到记录表（签到以 `action = "check_in"` 的积分流水表达）、任务表（任务硬编码于服务层）和勋章表（无勋章功能）。

## 5. 配置项与限制

模块无独立配置文件段，全部奖励数值为代码硬编码：

| 数值 | 值 | 定义位置 |
|------|-----|----------|
| 每日签到积分 | 5 | `internal/service/growth_service.go` 常量 `checkInPoints` |
| 发帖奖励 | 10 | `internal/worker/handlers.go` `handleGrowth` |
| 评论奖励 | 3 | 同上 |
| 被点赞奖励 | 2 | 同上 |
| 被收藏奖励 | 5 | 同上 |
| 创建圈子奖励 | 20 | 同上 |
| 加入圈子奖励 | 1 | 同上 |

依赖的 Kafka 与 Worker 配置（`configs/config.yaml`）：`kafka.topicPrefix = feedora`（实际 topic 加前缀）、`kafka.consumerGroup = feedora-worker`、`worker.batchSize = 100`、`worker.outboxIntervalSeconds = 2`、`worker.maxRetry = 5`（投递失败退避 10/30/60/300/600 秒）。

当前实现的局限：

- `internal/worker/growth/` 目录仅含 `doc.go` 占位文件，积分发放逻辑实际位于 `internal/worker/handlers.go`。
- 任务领取接口为占位实现，不校验完成状态、不发放积分。
- 无补签、无连签递增奖励，签到固定 5 分。
- 连续签到天数为最近 7 天签到次数近似，非严格连续判定。
- 无积分扣减与逆向事件回收；自交互（如给本人帖子点赞）同样发放积分。
- 用户等级注册后恒为 1，无升级路径；经验值用积分近似。
- 排行榜 `range` 参数不生效。
- `GrowthRepository.HasPointLog` 当前无任何调用方。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
|------|------|------|------|
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
