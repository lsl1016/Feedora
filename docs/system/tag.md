---
title: 标签模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: tag
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/tag_router.go
  - internal/api/tag_api.go
  - internal/dto/tag_dto.go
  - internal/service/tag_service.go
  - internal/repository/tag_repository.go
  - internal/model/tag.go
  - internal/model/post.go
summary: 标签模块提供只读的标签列表（按使用次数倒序）、标签详情与标签下帖子分页查询，标签与帖子的关联关系及 use_count 计数在发帖事务中维护，标签本身的创建与修改由后台管理模块承担。
---

## 1. 模块概述

标签模块面向客户端提供三个只读接口：标签列表、标签详情、标签下帖子分页列表，路由注册于 `internal/router/tag_router.go`。标签数据本身不支持通过本模块路由创建或修改，创建、重命名、停用均由后台管理模块的 `AdminService` 完成。发帖时帖子与标签的关联关系（`post_tags` 表）及标签使用计数（`tags.use_count`）在 `PostRepository.CreateWithRelations` 的同一事务内写入和递增；标签数据另被帖子聚合装配（`PostService.Assemble`）、关注模块（`FollowRepository.UsedTags`）和搜索索引 worker 复用。

## 2. 接口清单

路由统一挂载于 `/api/v1` 前缀下（下表路径省略该前缀），`v1` 分组启用 `middleware.OptionalAuth`：尝试识别登录用户，未登录放行。三个接口均为匿名可访问，登录用户在帖子列表中额外获得点赞/收藏状态。标签的创建与修改不在本模块路由内，由后台管理模块提供 `GET /admin/tags`、`POST /admin/tags`、`PUT /admin/tags/:tagId`（注册于 `internal/router/admin_router.go`）。

| 路径 | 方法 | 功能 | 控制器 |
|---|---|---|---|
| `/tags` | GET | 标签列表（仅启用中，按使用次数倒序，不分页） | `TagAPI.List` |
| `/tags/:tagId` | GET | 标签详情（不过滤状态，已停用标签仍可查） | `TagAPI.Get` |
| `/tags/:tagId/posts` | GET | 标签下帖子分页列表（`sort` 支持 latest/hot/comment/favorite/view） | `TagAPI.Posts` |

`tagId` 路径参数经 `dto.TagIDURI` 绑定校验（`required,min=1`）；`Posts` 的分页参数缺省为 page=1、pageSize=10，pageSize 上限 100（`normalizePageRequest`）。

## 3. 核心逻辑

### 3.1 标签的创建与引用方式

标签只能由后台管理接口创建：`AdminService.CreateTag` 校验 `name` 非空后直接 `Create`，新标签状态固定为 `enabled`；创建返回任何错误（包括但不限于 `uk_tag_name` 唯一键冲突）统一映射为 409"标签已存在"。普通用户发帖请求 `CreatePostRequest.TagIDs` 只传标签 ID 列表，不存在"按标签名称创建-or-复用"的逻辑，后端也不校验所传 ID 对应的标签是否存在或处于启用状态。

### 3.2 发帖时的关联与计数维护

`PostService.Create` 先对 `TagIDs` 执行 `dedup`（去重并过滤非正整数，保持顺序），再调用 `PostRepository.CreateWithRelations` 在同一事务内：逐条写入 `post_tags`，并对每个标签执行 `use_count = use_count + 1`（`UpdateColumn` + `gorm.Expr`）。因此同一帖子对同一标签只计一次，且关联写入与计数递增保持事务一致。所传标签 ID 不存在时，插入的 `post_tags` 行不影响任何计数（UPDATE 命中 0 行），错误也不被单独检查。

### 3.3 计数只增不减

帖子软删除（`PostRepository.SoftDelete`）只更新帖子状态并软删 `posts` 记录，不回减 `use_count`，也不删除 `post_tags` 记录。帖子编辑（`PostService.Update`）仅支持修改标题、正文、可见性与图片，不支持变更已绑定的标签。因此 `use_count` 是历史绑定次数的累计值，不等于当前有效关联数。

### 3.4 热门标签排序

`TagRepository.ListEnabled` 以 `use_count DESC, id ASC` 排序返回全部 `status = enabled` 的标签，`/tags` 接口一次性返回全量，无分页、无数量上限。已停用（`disabled`）标签不出现在列表中，但仍可通过 `/tags/:tagId` 查询详情。

### 3.5 标签下帖子查询

`TagService.Posts` 不校验标签是否存在，直接委托 `PostService.List`（`ListFilter{TagID, ViewerID, Sort, Page, PageSize}`）。数据层在 `TagID > 0` 时以 `post_tags` 子查询（`id IN (SELECT post_id FROM post_tags WHERE tag_id = ?)`）过滤帖子；基础过滤条件为 `status = published` 且（`visibility = public` 或 `circle_id IS NULL`），即圈内非公开帖子不出现在标签帖子列表中。不存在的标签 ID 返回空列表而非 404。

### 3.6 批量装配与跨模块读取

`TagRepository.FindByPostIDs` 通过 `post_tags JOIN tags` 一次查询返回 postID 到标签列表的映射（仅取 id 与 name），被帖子聚合装配（`PostService.Assemble`，输出 `TagSummary`）与搜索索引 worker（提取 `tagNames` 写入索引 payload）调用。关注模块的 `FollowRepository.UsedTags` 另行基于 `post_tags` 统计用户发帖使用过的启用标签，作为"关注的标签"数据源。

## 4. 数据模型

### 4.1 tags 表（`model.Tag`）

| 字段 | 类型 | 约束/默认 | 说明 |
|---|---|---|---|
| `id` | int64 | 主键 | 标签 ID |
| `name` | string(64) | 唯一索引 `uk_tag_name` | 标签名称，全表唯一 |
| `description` | string(255) | — | 标签描述 |
| `status` | string(32) | 默认 `enabled`，普通索引 | 状态，`enabled`/`disabled` |
| `use_count` | int64 | 默认 0 | 使用次数，发帖绑定时递增 |
| `created_at` / `updated_at` | time | — | 创建/更新时间 |

### 4.2 post_tags 表（`model.PostTag`，定义于 `internal/model/post.go`）

| 字段 | 类型 | 约束/默认 | 说明 |
|---|---|---|---|
| `id` | int64 | 主键 | 关联记录 ID |
| `post_id` | int64 | 联合唯一索引 `uk_post_tag` | 帖子 ID |
| `tag_id` | int64 | 联合唯一索引 `uk_post_tag`，另有单列索引 | 标签 ID |
| `created_at` | time | — | 绑定时间 |

`uk_post_tag` 保证同一帖子对同一标签至多一条关联。该表无软删除字段，帖子软删除后关联记录与 `use_count` 均保留。

## 5. 配置项与限制

无新增配置。标签模块不引入独立配置项，不使用缓存，不发布事件。

当前实现的限制：

- `/tags` 全量返回启用标签，无分页与数量上限。
- `/tags/:tagId/posts` 不校验标签存在性，不存在或已停用的标签 ID 返回空列表而非 404。
- `use_count` 只增不减，帖子删除不回减，数值与实际有效关联数可能不一致。
- 发帖可绑定不存在或已停用的标签 ID，后端不做校验。
- 帖子编辑不支持变更标签；后台 `AdminService.UpdateTag` 重命名为已存在名称时返回 500 而非 409。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
|---|---|---|---|
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
