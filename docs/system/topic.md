---
title: 话题模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: topic
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/topic_router.go
  - internal/api/topic_api.go
  - internal/dto/topic_dto.go
  - internal/service/topic_service.go
  - internal/repository/topic_repository.go
  - internal/model/topic.go
  - internal/model/post.go
summary: 话题模块提供话题广场（all/official/hot 分类）、话题详情与话题下帖子列表的只读查询能力，话题的创建与维护由后台接口承担，帖子发帖时在事务内绑定话题并同步维护 topics.post_count 计数。
---

## 1. 模块概述

话题模块为帖子提供按话题聚合的能力：话题广场支持 `all`/`official`/`hot` 分类浏览，话题详情返回基础信息与计数，话题下帖子列表复用帖子模块的列表装配逻辑。前台三个接口均为只读查询；话题的创建、编辑与启停由后台管理接口（`/admin/topics` 系列，注册于 `internal/router/admin_router.go`）承担。帖子发布时在数据库事务内写入 `post_topics` 关联并同步递增 `topics.post_count`。`TopicRepository` 另被搜索、榜单、关注与 worker 索引模块复用。

## 2. 接口清单

路由统一挂载于 `/api/v1` 前缀下，注册于 `internal/router/topic_router.go`，位于 `OptionalAuth` 组（匿名可访问，登录时识别用户）。

| 路径 | 方法 | 功能 | 控制器 |
|---|---|---|---|
| `/topics` | GET | 话题广场分页列表，`tab` 参数支持 `all`/`official`/`hot`/`latest` | `TopicAPI.List` |
| `/topics/:topicId` | GET | 话题详情 | `TopicAPI.Get` |
| `/topics/:topicId/posts` | GET | 话题下的帖子分页列表，`sort` 参数透传帖子列表排序 | `TopicAPI.Posts` |

除上表外，另有三条复用本模块仓储层/服务层的接口，注册于 `internal/router/admin_router.go`（需登录且角色为 `admin`）：`GET /admin/topics`（`AdminAPI.Topics`，后台话题分页，使用 `TopicRepository.ListAll`）；`POST /admin/topics`（`AdminAPI.CreateTopic`，创建话题）；`PUT /admin/topics/:topicId`（`AdminAPI.UpdateTopic`，更新话题）。

## 3. 核心逻辑

### 3.1 话题广场列表与排序

`TopicRepository.List` 以 `status = enabled` 为基础过滤，按 `tab` 分支：

| tab | 附加条件 | 排序 |
|---|---|---|
| `official` | `is_official = true` | `created_at DESC` |
| `hot` | 无 | `participant_count DESC, post_count DESC` |
| `all` / `latest` / 其他任意值 | 无 | `created_at DESC` |

- 控制器将空 `tab` 归一为 `all`；service/repository 层不校验 `tab` 取值，未知值按默认分支处理。
- `latest` 与 `all` 行为完全相同，均为创建时间倒序，无独立实现。
- `hot` 排序无时间衰减因子，仅依赖两个字段的静态值。
- 分页规则（service 层 `normPage`）：`page <= 0` 取 1，`size <= 0` 取 10，`size > 100` 收敛为 100；控制器侧 `normalizePageRequest` 执行同一套规整后再调用 service。

### 3.2 话题详情

- `FindByID` 按 ID 单查，记录不存在时返回 `(nil, nil)`，service 层映射为 `ErrNotFound`。
- 详情查询不过滤 `status`：已置为 `disabled` 的话题不出现在广场列表，但仍可通过详情接口按 ID 查询。

### 3.3 话题下的帖子

`TopicService.Posts` 直接委托 `PostService.List`，携带 `TopicID` 与当前查看者 ID：

- 帖子过滤条件为 `status = published` 叠加 `(visibility = 'public' OR circle_id IS NULL)`：公开帖与未入圈的帖子可见，圈子内的非公开帖子不出现在话题帖列表。
- `TopicID` 通过 `post_topics` 子查询（`id IN (SELECT post_id FROM post_topics WHERE topic_id = ?)`）过滤。
- `sort` 透传帖子模块排序，支持 `latest`（默认，`created_at DESC`）、`hot`（`hot_score DESC`）、`comment`、`favorite`、`view`。
- `ViewerID` 用于回填每条帖子的点赞/收藏状态；匿名访问时该 ID 为 0，状态恒为 `false`。

### 3.4 话题创建与更新（后台）

- 创建：`name` 为必填（DTO `binding:"required"` 且 service 层再判空，不做 trim）；新话题固定 `status = enabled`，`is_official` 由请求指定，`is_recommended` 创建时不可指定（默认 `false`，仅能通过更新接口设置）。
- 创建不预先查询重名，直接 `INSERT`，唯一索引 `uk_topic_name` 兜底去重；`Create` 返回任意错误（含非重名的数据库错误）时统一返回 `409 话题已存在`。
- 更新：白名单字段为 `name`、`description`、`cover_url`、`is_official`、`is_recommended`、`status`，全部为可选指针字段，仅更新请求中出现的字段并刷新 `updated_at`；`status` 取值不做枚举校验。更新语句影响 0 行不报错，随后回查记录，不存在时返回 `404`。
- 后台创建与更新话题不发布事件，worker 不感知话题变更；搜索索引中的话题数据仅在 worker 启动时全量重建（仅索引 `status = enabled` 的话题）。

### 3.5 帖子绑定话题与计数维护

- 发帖（`PostService.Create`）时，`topicIDs` 先经 `dedup` 去重：过滤 `<= 0` 的 ID 并按原顺序去重。
- `PostRepository.CreateWithRelations` 在单个事务内完成：插入 `posts` 记录、插入图片、插入标签关联并递增标签计数、逐条插入 `post_topics` 关联并对每条执行 `topics.post_count = post_count + 1`、递增用户/圈子发帖数。事务失败整体回滚。
- 数据库层兜底去重：`post_topics` 的 `uk_post_topic (post_id, topic_id)` 唯一索引阻止同一帖子重复绑定同一话题。
- 绑定前不校验话题存在性与状态：传入不存在的 `topicID` 仍会插入 `post_topics` 记录（表上无外键约束），对应的计数 UPDATE 影响 0 行不报错，产生无话题对应的孤儿关联；已禁用的话题同样可被绑定。
- 帖子装配时的批量回查 `FindByPostIDs` 通过 `JOIN topics` 仅取 `id`/`name` 两个字段，按 `postID` 分组；孤儿关联记录因 JOIN 无匹配被自然排除，不影响帖子返回。
- 帖子编辑接口不支持修改话题绑定（`UpdatePostRequest` 无话题字段，service 层不触达 `post_topics`）。
- 帖子隐藏与软删除均不删除 `post_topics` 记录，也不回减 `topics.post_count`；被删帖子因状态过滤不再出现在话题帖列表，但话题计数与关联记录保留。

### 3.6 参与人数计数

`topics.participant_count` 在全部代码中没有任何写入或递增逻辑，仅在 `hot` 排序、话题榜单、搜索索引文档与 DTO 展示中读取，初始值为默认值 0（种子数据亦不设置）。因此 `hot` 分类与话题榜单在当前数据状态下实际退化为按 `post_count` 排序。

### 3.7 其他模块对本模块数据的复用

| 使用方 | 行为 |
|---|---|
| 搜索模块（`SearchService`） | `type = topic` 时查询 ES 话题索引取 ID，再回库 `FindByID` 装配（不过滤状态） |
| 榜单模块（`RankService`） | 话题热门榜优先读 Redis ZSet；为空时回退 DB 查询（`status = enabled`，`participant_count DESC, post_count DESC`），回查 `FindByID` 装配名称与计数 |
| 关注模块（`FollowRepository.ParticipatedTopics`） | 以当前用户未软删帖子的 `post_topics` 子查询推导"用户发帖涉及的话题"，作为"关注的话题"数据源；过滤 `status = enabled`，按 `id DESC` 排序 |
| worker（`indexTopic`） | worker 启动时将全部 `status = enabled` 的话题全量写入 ES 话题索引 |

系统不存在独立的话题关注（订阅）实体：用户与话题的关系仅由发帖行为间接产生。

## 4. 数据模型

表 `topics`（`internal/model/topic.go`）：

| 字段 | 类型/约束 | 说明 |
|---|---|---|
| `id` | 主键 | 话题 ID |
| `name` | `size:128`，唯一索引 `uk_topic_name` | 话题名称，重名由唯一索引拒绝 |
| `description` | `size:500` | 话题描述 |
| `cover_url` | `size:512` | 封面图片 URL |
| `is_official` | 默认 `false` | 是否官方话题，`official` 分类的过滤依据 |
| `is_recommended` | 默认 `false`，索引 | 是否推荐，仅后台可更新，当前查询侧无使用 |
| `participant_count` | 默认 `0` | 参与人数，无写入逻辑，恒为 0 |
| `post_count` | 默认 `0` | 帖子数冗余计数，绑定帖子时事务内 +1，不回减 |
| `status` | `size:32`，默认 `enabled`，索引 | `enabled` / `disabled`，广场列表与关注推导均过滤 `enabled` |
| `created_at` / `updated_at` | 时间戳 | 创建/更新时间 |

表 `post_topics`（`internal/model/post.go` 中的 `PostTopic`）：

| 字段 | 类型/约束 | 说明 |
|---|---|---|
| `id` | 主键 | 关联记录 ID |
| `post_id` | 唯一索引 `uk_post_topic`（联合） | 帖子 ID |
| `topic_id` | 唯一索引 `uk_post_topic`（联合），索引 | 话题 ID |
| `created_at` | 时间戳 | 绑定时间 |

两表迁移注册于 `internal/model/all.go`；`post_topics` 无外键约束，无 `deleted_at` 软删标记，帖子删除后关联记录物理保留。

## 5. 配置项与限制

本模块无自有配置项。搜索索引行为受全局搜索组件装配影响（未启用搜索组件时 `indexTopic` 直接返回）。

当前实现的限制：

- `participant_count` 无任何维护逻辑，`hot` 分类与话题榜单实际由 `post_count` 决定排序。
- 绑定话题不校验存在性与状态，可产生指向不存在话题的孤儿 `post_topics` 记录。
- 帖子删除、隐藏不回减 `topics.post_count`，计数只增不减，与话题下实际可见帖子数产生偏差；帖子编辑不支持调整话题绑定。
- 话题详情接口不区分话题状态，`disabled` 话题仍可按 ID 访问。
- 后台创建话题把所有数据库错误统一映射为 `409 话题已存在`；更新接口对 `status` 取值不做枚举校验。
- 话题创建/更新不发布事件，搜索索引仅在 worker 启动时全量重建，后台改动不实时同步到索引。
- `is_recommended` 字段有索引但查询侧无任何使用。

## 6. 历史版本

| 版本 | 日期 | 维护者 | 说明 |
|---|---|---|---|
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
