---
title: 帖子模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: post
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/post_router.go
  - internal/api/post_api.go
  - internal/dto/post_dto.go
  - internal/service/post_service.go
  - internal/repository/post_repository.go
  - internal/model/post.go
  - internal/model/constants.go
  - internal/event/producer.go
summary: 帖子模块提供发帖（立即/草稿/预约）、编辑、隐藏、软删除、分享与转发，以及按圈子/话题/标签/关注流等多维度的分页列表和基于 Redis 缓存的详情查询，变更通过 post.events 事件对外发布。
---

## 1. 模块概述

帖子模块提供帖子的发布（立即发布、草稿、预约三种状态）、编辑、隐藏/取消隐藏、软删除、分享计数与转发能力，并提供多维度分页列表与单帖详情查询。`PostService` 同时承担帖子 DTO 的批量聚合装配（`Assemble`，一次性加载作者、图片、标签、话题、圈子与点赞/收藏状态），被话题、标签、圈子、用户、搜索等模块复用。详情查询采用 Redis Cache-Aside 缓存（键 `post:detail:{postId}`，TTL 10 分钟），帖子变更通过 `post.events` 主题以事件形式对外发布。

## 2. 接口清单

路由统一挂载于 `/api/v1` 前缀下，注册于 `internal/router/post_router.go`。`List` 与 `Get` 为匿名可访问，其余接口经过 `authMW` 认证中间件。

| 路径 | 方法 | 功能 | 控制器 |
|---|---|---|---|
| `/posts` | GET | 分页查询帖子列表（支持信息流/排序/状态/关键词/标签/圈子/话题/作者等筛选） | `PostAPI.List` |
| `/posts/:postId` | GET | 查询帖子详情（含阅读计数与当前用户点赞/收藏状态） | `PostAPI.Get` |
| `/posts` | POST | 创建帖子 | `PostAPI.Create` |
| `/posts/:postId` | PUT | 编辑帖子（仅作者本人） | `PostAPI.Update` |
| `/posts/:postId` | DELETE | 删除帖子（软删除，作者本人或 admin 角色） | `PostAPI.Delete` |
| `/posts/:postId/hide` | PUT | 隐藏帖子（仅作者本人） | `PostAPI.Hide` |
| `/posts/:postId/unhide` | PUT | 取消隐藏帖子（仅作者本人） | `PostAPI.Unhide` |
| `/posts/:postId/share` | POST | 分享计数 +1 | `PostAPI.Share` |
| `/posts/:postId/repost` | POST | 转发帖子（生成 repost 类型新帖） | `PostAPI.Repost` |

除上表外，`PostService.List` / `Assemble` / `AssembleByIDs` 被其他模块路由复用：`GET /topics/:topicId/posts`（话题下帖子）、`GET /tags/:tagId/posts`（标签下帖子）、`GET /circles/:circleId/posts`（圈子下帖子）、用户主页帖子与搜索等。后台管理端另有 `GET /admin/posts`（`AdminAPI.Posts`），走 `AdminRepository.ListPosts` 独立查询路径，不复用本模块仓储层。

## 3. 核心逻辑

### 3.1 发帖事务与圈子权限

`Create` 按以下顺序处理：

1. 标题去首尾空白，标题或正文为空返回参数错误；`visibility` 为空时默认 `public`。
2. `publishMode` 映射状态：`draft` → `draft`，`schedule` → `scheduled`，其余值（含空）→ `published`；仅 `published` 状态写入 `published_at`。
3. 携带 `circleId` 时执行 `checkCirclePostPermission`：圈子不存在返回圈子不存在错误；成员记录不存在或状态为 `removed` 返回未加入错误；状态为 `muted` 返回禁言错误；圈子 `post_permission = "admin_only"` 时要求成员角色为 `owner` 或 `moderator`。
4. 标签/话题 ID 先经 `dedup` 去重并过滤非正整数，再调用 `CreateWithRelations` 落库。

`PostRepository.CreateWithRelations` 在单个数据库事务内完成：

- 插入 `posts` 主记录；`summary` 取正文前 120 个字符（超出追加 `...`），`images` 非空时首图写入 `cover_url`。
- 逐条插入 `post_images`（`sort_order` 为图片下标）。
- 逐条插入 `post_tags` 并对每个标签执行 `use_count + 1`；逐条插入 `post_topics` 并对每个话题执行 `post_count + 1`。
- 对作者执行 `users.post_count + 1`；帖子归属圈子时对 `circles.post_count + 1`。

事务提交后，`Create` 对 `published` 状态再次调用 `users.IncColumn(authorID, "post_count", 1)`，即一次成功发布中作者 `post_count` 被递增两次；`draft`/`scheduled` 帖子仅事务内递增一次。

`Repost` 转发生成新帖：`post_type = repost`，标题为 `转发：` + 源帖标题，正文与 `repost_comment` 为转发附言，`source_post_id` 指向源帖，可见性固定 `public`、状态固定 `published`；随后源帖 `repost_count + 1`、转发者 `users.post_count + 1`。转发走 `PostRepository.Create` 单条插入，不建立标签/话题关系。

### 3.2 列表筛选与排序

`PostRepository.List` 的过滤维度：

- 状态：默认只取 `published`；显式传入 `status` 时精确匹配；`status = all` 且 `includeHidden = true` 且已登录时放宽为 `status = published OR (author_id = 当前用户 AND status <> deleted)`（作者可在列表中看到自己的草稿/隐藏帖），否则仍只取 `published`。
- 关注流：`feedType = following` 且已登录时，以 `user_follows` 子查询限定 `author_id IN (该用户关注的 followee_id 集合)`；未登录时该条件不生效，退化为普通信息流。`feedType = hot` 等价于热度排序。
- 圈子：传入单个 `circleId` 或 `circleIds` 集合时按圈子过滤且不叠加可见性限制（圈内帖对圈子维度可见）；未指定圈子时限定 `visibility = public OR circle_id IS NULL`，即 `circle_only` 帖子仅在按其所属圈子过滤时出现。
- 关键词：`title LIKE %kw% OR content_md LIKE %kw%`。
- 标签：`tagId` 经 `post_tags` 子查询转换为帖子 ID 集合；话题：`topicId` 或 `topicIds` 集合经 `post_topics` 子查询。
- 作者：单个 `authorId` 或 `authorIds` 集合（供关注体系使用）。

排序规则（均以 `created_at DESC` 兜底）：

| 条件 | 排序 |
|---|---|
| 默认 / `sort = latest` 及其它未识别值 | `created_at DESC` |
| `feedType = hot` 或 `sort = hot` | `hot_score DESC, created_at DESC` |
| `sort = comment` | `comment_count DESC, created_at DESC` |
| `sort = favorite` | `favorite_count DESC, created_at DESC` |
| `sort = view` | `view_count DESC, created_at DESC` |

先 `Count` 总数，再按排序与 `OFFSET/LIMIT` 取当前页。查询结果交由 `Assemble` 批量装配：按帖子 ID 集合一次性加载作者、圈子、图片、标签、话题及当前用户的点赞/收藏集合，避免逐帖查询。

### 3.3 详情查询与阅读计数

`Get` 采用 Cache-Aside：

- 缓存命中：直接对 `posts.view_count` 执行一次 DB 增量（`IncViewCount`），将内存中的 `ViewCount + 1` 后叠加当前用户状态返回；不校验帖子当前状态。
- 缓存未命中：`FindByID` 查询（GORM 软删除自动排除已删行），记录不存在或状态为 `deleted` 返回帖子不存在错误；状态非 `published` 且访问者不是作者时返回帖子不可见错误；随后递增阅读数、以 `viewerID = 0` 装配可共享的详情主体，仅 `published` 状态写入缓存，最后叠加当前用户状态。

阅读计数为每次读取同步写库（含缓存命中路径），无异步合并；`liked`/`favorited` 字段由 `overlayViewer` 按 `viewerID` 实时查询 `InteractionRepository` 叠加，不进入缓存。

### 3.4 事件发送

帖子事件统一发布到 `event.TopicPost`（常量值 `post.events`）。生产者由全局配置决定：`kafka.enabled = true` 时为 `OutboxProducer`，事件写入 `event_outbox` 表且 topic 列为 `{kafka.topicPrefix}.post.events`（配置默认前缀 `feedora`）；否则为 `NoopProducer`，仅打印日志。将 outbox 投递到 Kafka 的 Dispatcher（`internal/worker/dispatcher`）为阶段二占位空实现，当前无下游消费者。

| 动作 | 事件类型 | Payload |
|---|---|---|
| `Create`（任意状态） | `PostCreated` | `{"title": 帖子标题}` |
| `Update` | `PostUpdated` | `nil` |
| `SetHidden(true)` | `PostHidden` | `nil` |
| `SetHidden(false)` | `PostUnhidden`（字面量字符串，`event_type.go` 中无常量定义） | `nil` |
| `Delete` | `PostDeleted` | `nil` |

`Share` 与 `Repost` 不发送任何事件。

### 3.5 缓存策略

- 缓存键：`PostDetailKey` 生成 `post:detail:{postId}`，TTL 为常量 `postDetailTTL = 10 * time.Minute`。
- 失效方式：`Update`、`SetHidden`、`Delete` 调用 `PostService.InvalidateDetail` 删除该键；互动与评论模块在计数变化后直接对同一键执行 `cache.Del`（`interaction_service.go`、`comment_service.go`），不经由 `InvalidateDetail`。
- 阅读数增量不回写缓存，缓存内 `viewCount` 与库值的偏差在 TTL 到期前不刷新。
- `Cache` 在未启用 Redis 时所有操作安全降级为空操作，详情每次回源数据库。

## 4. 数据模型

主表 `posts`（`model.Post`）：

| 字段 | 说明 |
|---|---|
| `id` | 主键 |
| `author_id` / `status` | 组合索引 `idx_author_status` |
| `title` / `content_md` / `summary` | 标题（200）、Markdown 正文（longtext）、摘要（500） |
| `post_type` | `original`（默认）或 `repost`，转发帖携带 `source_post_id` 与 `repost_comment`（1000） |
| `circle_id` / `status` | `circle_id` 与 `status` 组合索引 `idx_circle_status` |
| `visibility` | `public`（默认）或 `circle_only` |
| `view_count` / `like_count` / `comment_count` / `favorite_count` / `share_count` / `repost_count` | 计数列，默认 0 |
| `hot_score` | 热度分，默认 0，带索引 |
| `scheduled_at` / `published_at` / `created_at` / `updated_at` / `deleted_at` | 时间列，`deleted_at` 为 GORM 软删除标记（带索引） |

关联表：`post_images`（`post_id`、`image_url`、`sort_order`）；`post_tags`（唯一索引 `uk_post_tag (post_id, tag_id)`）；`post_topics`（唯一索引 `uk_post_topic (post_id, topic_id)`）。

状态枚举（`model/constants.go`）共 8 个：`draft`、`scheduled`、`reviewing`、`published`、`hidden`、`rejected`、`deleted`、`takedown`。其中 `reviewing`、`rejected`、`takedown` 三个常量当前无任何写入点。

软删与隐藏语义：

- 删除（`SoftDelete`）：先将 `status` 更新为 `deleted`，再执行 GORM 软删除写入 `deleted_at`，两条语句不在同一事务；软删后的帖子对所有查询不可见，列表中作者自己的帖子也通过 `status <> deleted` 排除。
- 隐藏（`SetHidden`）：`status` 在 `hidden` 与 `published` 之间切换，不删数据；隐藏帖不出现在默认列表，仅作者本人可通过详情接口访问。

## 5. 配置项与限制

本模块无新增配置项；事件 topic 前缀依赖全局 `kafka.topicPrefix`（默认 `feedora`），详情缓存依赖全局 Redis 配置。

限制与约定：

- 分页：`page` 默认 1，`pageSize` 默认 10、上限 100（DTO 绑定 `max=100` 与 `normalizePageRequest` 双重限制）。
- 摘要截取上限 120 个字符；标题长度上限 200。
- 详情缓存 TTL 10 分钟。

当前实现的局限：

- `CreatePostRequest.scheduledAt` 参数被接受但未映射到 `posts.scheduled_at`，预约时间不持久化；`scheduled` 状态帖子无定时发布任务，不会自动转为 `published`。
- 已发布帖创建时 `users.post_count` 在事务内与 Service 层各递增一次（共 +2），而 `Delete` 仅递减 1；草稿/预约帖也在事务内递增作者 `post_count`。
- `posts.hot_score` 列在全部代码中无写入点，`hot` 排序实际按默认值 0 排序后退化为 `created_at DESC`。
- `Repost` 不校验源帖状态与可见性，源帖为 `hidden` 或 `draft` 时仍可转发。
- `Share` 仅递增计数，不记录分享者；`Update` 不能修改标签、话题与圈子归属。
- DTO 中 `isTop`、`isFeatured`、`isSelected`、`followedAuthor` 字段无数据来源，`ToPost` 不赋值，响应中恒为 `false`。

## 6. 历史版本

| 版本 | 日期 | 维护者 | 说明 |
|---|---|---|---|
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
