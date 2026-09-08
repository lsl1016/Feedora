---
title: 圈子模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: circle
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/circle_router.go
  - internal/api/circle_api.go
  - internal/dto/circle_dto.go
  - internal/service/circle_service.go
  - internal/repository/circle_repository.go
  - internal/model/circle.go
summary: circle 模块提供圈子的创建与检索、加入/退出、成员角色与状态管理（设角色、禁言/解禁、移除）以及成员与圈子帖子列表查询，成员数在成员变更时同步维护，圈子详情走 15 分钟 Cache Aside 缓存，创建与加入通过 circle.events 事件驱动搜索索引、通知与积分。
---

## 1. 模块概述

circle 模块负责圈子（Circle）这一内容分组单元的管理，覆盖圈子的创建与多维度检索、加入/退出、成员角色与状态管理（设置角色、禁言/解禁、移除成员）、成员列表与圈子内帖子列表查询，由 `CircleAPI` / `CircleService` / `CircleRepository` 竖切片实现。圈主（`owner`）与管理员（`moderator`）两类角色具备圈子管理权限；成员状态在 `normal` / `muted` / `removed` 之间迁移。向圈子发帖的权限校验位于 post 模块的 `checkCirclePostPermission`。本模块不含圈子公告功能，代码中不存在公告实体与公告接口。

## 2. 接口清单

路由统一挂载于 `/api/v1` 前缀下，注册于 `internal/router/circle_router.go`。查询类接口（`List` / `Get` / `Members` / `Posts`）匿名可访问，未登录时视图用户 ID 取 0；其余接口经过 `authMW` 认证中间件。

| 路径 | 方法 | 功能 | 控制器 |
|---|---|---|---|
| `/circles` | GET | 分页查询圈子列表（scope/keyword/category/sort 筛选排序） | `CircleAPI.List` |
| `/circles` | POST | 创建圈子（创建者自动成为圈主） | `CircleAPI.Create` |
| `/circles/:circleId` | GET | 查询圈子详情（叠加当前用户成员身份） | `CircleAPI.Get` |
| `/circles/:circleId/join` | POST | 加入圈子（幂等） | `CircleAPI.Join` |
| `/circles/:circleId/leave` | POST | 退出圈子（圈主不可退出） | `CircleAPI.Leave` |
| `/circles/:circleId/members` | GET | 分页查询圈子成员列表 | `CircleAPI.Members` |
| `/circles/:circleId/posts` | GET | 分页查询圈子内帖子列表 | `CircleAPI.Posts` |
| `/circles/:circleId/members/:userId/role` | PUT | 设置成员角色（仅圈主） | `CircleAPI.SetRole` |
| `/circles/:circleId/members/:userId/mute` | PUT | 禁言成员（圈主或管理员） | `CircleAPI.Mute` |
| `/circles/:circleId/members/:userId/unmute` | PUT | 解除禁言（圈主或管理员） | `CircleAPI.Unmute` |
| `/circles/:circleId/members/:userId` | DELETE | 移除成员（圈主或管理员，软移除） | `CircleAPI.Remove` |

后台管理端另有 `GET /admin/circles`（`AdminAPI`），走 `AdminService.ListCircles` → `CircleRepository.ListAll` 独立查询路径，按 `id ASC` 返回全部圈子，不做状态过滤。

## 3. 核心逻辑

### 成员身份状态机

成员记录为 `circle_members` 表中 `(circle_id, user_id)` 唯一的一条记录，状态迁移如下：

| 起始 | 操作 | 终态与副作用 |
|---|---|---|
| 无记录 | `Join` | 写入 `role=member`、`status=normal` 记录，`member_count + 1`，失效详情缓存，发布 `CircleJoined` 事件 |
| `removed` | `Join` | 原记录改回 `status=normal`（保留原 `role`），`member_count + 1`，失效详情缓存；不发布事件 |
| `normal` / `muted` | `Join` | 幂等，直接返回成功，无任何变更 |
| `normal` / `muted` | `Leave` | 物理删除成员记录（`DeleteMember`），`member_count - 1`，失效详情缓存；圈主（`role=owner`）返回 422「圈主不能退出圈子」 |
| `normal` | `Mute` | `status=muted`，写入 `mute_reason`，`duration > 0` 时写入 `muted_until = now + duration 天` |
| `muted` | `Unmute` | `status=normal`，清空 `mute_reason` 与 `muted_until` |
| 任意非圈主 | `Remove` | `status=removed`（记录保留），`member_count - 1`，失效详情缓存；`role=owner` 返回 422「不能移除圈主」 |

`member_count` 由 `IncMemberCount` 以 SQL 表达式增量维护，减量使用 `GREATEST(member_count - ?, 0)` 兜底不为负。

### 角色与权限判断

- 角色取值：`owner` / `moderator` / `reviewer` / `member`（定义于 `internal/model/constants.go`）。`reviewer` 仅出现在成员列表排序的 `FIELD(role,'owner','moderator','reviewer','member')` 中，本模块无任何代码路径赋予该角色。
- `SetRole` 仅允许操作者为圈主（`role=owner`），否则返回 `ErrForbidden`。
- `Mute` / `Unmute` / `Remove` 通过 `requireManage` 校验操作者角色为 `owner` 或 `moderator`。
- 向圈子发帖的校验在 post 模块 `checkCirclePostPermission`（`internal/service/post_service.go`）完成：圈子存在、成员记录存在且 `status != removed`（否则 `ErrCircleNotJoined`）、`status != muted`（否则 `ErrCircleMuted`）、`post_permission == "admin_only"` 时仅 `owner` / `moderator` 可发（否则 `ErrPostNoPermission`）。`post_permission` 为其他取值时不做额外限制。

### 列表检索与圈子帖子关联

- `List` 支持 `scope`：`all`（默认，API 层兜底）、`recommended`（`is_recommended = true`）、`joined`（当前用户 `status <> removed` 的成员圈子，未登录返回空）、`created`（`owner_id = viewerID`）；`keyword` 对 `name` / `description` 做前后通配 `LIKE`；`sort=hot` 按 `member_count DESC, post_count DESC`，默认按 `created_at DESC`。所有查询固定过滤 `status = 'normal'` 的圈子。
- `Members` 排除 `status=removed` 的记录，按角色优先级（owner → moderator → reviewer → member）与 `joined_at ASC` 排序，仅支持分页，无角色/状态筛选参数。
- `Posts` 委托 `PostService.List`，以 `circle_id = :circleId` 过滤帖子并复用帖子模块的可见性与装配逻辑。
- `post_count` 在帖子创建事务内 `+1`（`internal/repository/post_repository.go`）；帖子删除无对应回减路径，`post_count` 只增不减。

### 详情缓存

`Get` 采用 Cache Aside：缓存键 `circle:detail:{circleId}`（`internal/cache/keys.go`），TTL 为代码内常量 15 分钟。缓存主体以 `viewerID=0` 装配（不含任何用户身份字段），命中后再通过 `overlayMembership` 从数据库实时叠加当前用户的 `IsJoined` / `MyRole` / `MyStatus`。`Create`（经创建后的 `Get` 回填）、`Join`、`Leave`、`Remove` 主动删除该缓存键；`SetRole` / `Mute` / `Unmute` 不删缓存，因角色与禁言状态每次请求都从成员表实时读取，缓存仅含圈主信息与计数。

### 事件

创建与加入向 `circle.events` 主题（`internal/event/topic.go`）发布事件，worker 消费：`CircleCreated` 触发搜索索引写入（`indexCircle`）与创建者 +20 积分（`create_circle`）；`CircleJoined` 触发圈主通知（圈主本人加入除外，幂等去重）、加入者 +1 积分（`join_circle`）与热度榜单统计。`Leave` / `Remove` / `SetRole` / `Mute` / `Unmute` 不发布事件。

### 创建约束

`Create` 对名称做 `TrimSpace` 后非空校验（`ErrParams`），同名圈子返回 409「圈子名称已存在」；`joinType` 缺省补 `direct`，`postPermission` 缺省补 `all`；圈子与圈主成员记录在同一事务内写入（`CreateWithOwner`），初始 `member_count = 1`。

### AI 能力

圈子模块无 AI 能力。帖子摘要由 post 模块的 `summarize` 生成，为正文前 120 个字符的截断，不是模型生成的智能摘要。

## 4. 数据模型

表结构定义于 `internal/model/circle.go`。

### `circles`（圈子表）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | int64 | 主键 |
| `owner_id` | int64 | 圈主用户 ID，索引 |
| `name` | string(128) | 圈子名称，唯一索引 `uk_circle_name` |
| `avatar` | string(512) | 圈子头像 |
| `description` | string(500) | 圈子描述 |
| `category` | string(64) | 圈子分类 |
| `join_type` | string(32) | 加入方式，默认 `direct` |
| `post_permission` | string(32) | 发帖权限，默认 `all` |
| `rules` | text | 圈规说明 |
| `status` | string(32) | 圈子状态，默认 `normal`，索引 |
| `member_count` | int64 | 成员数，默认 0 |
| `post_count` | int64 | 帖子数，默认 0 |
| `featured_count` | int64 | 精选帖子数，默认 0 |
| `is_recommended` | bool | 是否推荐，默认 false，索引 |
| `created_at` / `updated_at` / `deleted_at` | time | 时间字段，`deleted_at` 软删除索引 |

### `circle_members`（圈子成员表）

| 字段 | 类型 | 说明 |
|---|---|---|
| `id` | int64 | 主键 |
| `circle_id` + `user_id` | int64 | 联合唯一索引 `uk_circle_user` |
| `role` | string(32) | 成员角色，默认 `member`，索引 `idx_circle_role` |
| `status` | string(32) | 成员状态，默认 `normal`，索引 `idx_circle_status` |
| `mute_reason` | string(255) | 禁言原因 |
| `muted_until` | *time.Time | 禁言截止时间，可空 |
| `joined_at` / `updated_at` | time | 时间字段 |

`circle_members` 无 `deleted_at` 字段，退出圈子为物理删除。

### 枚举取值

| 枚举 | 取值 |
|---|---|
| 成员角色 `role` | `owner`、`moderator`、`reviewer`、`member` |
| 成员状态 `status` | `normal`、`muted`、`removed` |
| 圈子状态 `circles.status` | 本模块仅写入 `normal`（复用 `CircleMemberNormal` 常量），无其他取值的写入路径 |

## 5. 配置项与限制

无新增配置。缓存 TTL（15 分钟）为 `circle_service.go` 内的代码常量，不随配置文件调整。

当前实现的限制：

- `join_type` 仅存储，`Join` 不校验该字段，不存在审批/邀请流程，任何圈子均直接加入。
- 限时禁言不自动解除：发帖校验只比较 `status == muted`，不比较 `muted_until`，到期后仍视为禁言，须调用 `unmute` 手动解除。
- `mute` 请求体的 `duration` 绑定 `min=1`，服务层 `days <= 0` 表示永久禁言的分支经 API 不可达。
- `SetRole` 不校验 `role` 取值，任意字符串会直接写入 `role` 字段；不阻止将成员设为 `owner`（可产生多圈主）；目标成员不存在时更新 0 行且不报错。
- `Mute` 不校验目标角色，可对圈主或其他管理员执行禁言。
- `post_count` 只增不减；`featured_count` 无任何写入路径，恒为默认值 0；`is_recommended` 无 API 写入路径，`scope=recommended` 依赖数据库中的既有数据。
- 被 `Remove` 的成员重新 `Join` 时保留原 `role`（曾被设为 `moderator` 的成员复活后仍为 `moderator`）。
- 圈子无编辑、删除接口；模型中的 `deleted_at` 软删除字段无本模块写入路径。
- 成员列表无角色/状态筛选参数，仅分页。

## 6. 历史版本

| 版本 | 日期 | 维护者 | 变更说明 |
|---|---|---|---|
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
