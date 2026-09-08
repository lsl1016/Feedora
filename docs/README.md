# Feedora 文档索引

Feedora 是开发者知识社区 V2.1 的 monorepo，包含 Go 后端（`Feedora-backend/`）与 React 前端（`Feedora-frontend/`）。

写作规范见 [CONVENTION.md](CONVENTION.md)，提交前用 `bash scripts/check-docs.sh` 校验。

## 系统功能文档（docs/system/）

描述各模块**当前**的接口、逻辑与数据模型，随代码演进持续更新。

| 文档 | 模块 | 版本 | 最近更新 | 说明 |
|------|------|------|----------|------|
| [auth.md](system/auth.md) | auth | v1.0 | 2026-09-09 | 注册、登录、JWT 鉴权 |
| [user.md](system/user.md) | user | v1.0 | 2026-09-09 | 用户资料与个人主页 |
| [post.md](system/post.md) | post | v1.1 | 2026-09-09 | 帖子发布、浏览、详情 |
| [comment.md](system/comment.md) | comment | v1.1 | 2026-09-09 | 评论与回复 |
| [circle.md](system/circle.md) | circle | v1.0 | 2026-09-09 | 圈子、成员、公告 |
| [topic.md](system/topic.md) | topic | v1.1 | 2026-09-09 | 话题 |
| [tag.md](system/tag.md) | tag | v1.0 | 2026-09-09 | 标签 |
| [follow.md](system/follow.md) | follow | v1.0 | 2026-09-09 | 用户关注与关注流聚合 |
| [interaction.md](system/interaction.md) | interaction | v1.1 | 2026-09-09 | 点赞与收藏 |
| [notification.md](system/notification.md) | notification | v1.1 | 2026-09-09 | 站内通知 |
| [file.md](system/file.md) | file | v1.1 | 2026-09-09 | 文件上传与存储 |
| [search.md](system/search.md) | search | v1.1 | 2026-09-09 | Elasticsearch 全文搜索 |
| [rank.md](system/rank.md) | rank | v1.1 | 2026-09-09 | 热门榜单 |
| [growth.md](system/growth.md) | growth | v1.1 | 2026-09-09 | 签到、积分、任务、等级 |
| [admin.md](system/admin.md) | admin | v1.1 | 2026-09-09 | 后台管理 |
| [infrastructure.md](system/infrastructure.md) | infrastructure | v1.1 | 2026-09-09 | worker/cache/event/pkg 横切基础设施 |

## 历史设计文档（只读快照）

各开发阶段的设计与规范记录，反映**编写时点**的设计意图，不追随代码更新。

`Feedora-backend/docs/`：

| 文档 | 说明 |
|------|------|
| `阶段一后端MVP设计文档.md` | 阶段一（Gin + GORM + JWT 核心闭环）技术设计 |
| `阶段二后端MVP设计文档.md` | 阶段二（Redis + Kafka Worker + Elasticsearch）技术设计 |
| `项目整体原型文档.md` | V2.1 后端整体技术方案（八大能力全景） |
| `功能边界说明.md` | 各阶段功能边界划分 |
| `observability-design.md` 等 6 篇 | 日志/指标/追踪/告警可观测性设计与实施指南 |
| `日志记录规范.md` | 中文日志记录规范（与 `logging-spec.md` 主题重叠，待合并） |
| `docs.go` | swaggo 自动生成的 Swagger 产物，由 `internal/router/router.go` 导入注册 Swagger UI。属生成代码而非文档，重新执行 `swag init` 会覆盖 |

`Feedora-frontend/docs/`：

| 文档 | 说明 |
|------|------|
| `api-contract.md` | 前后端接口契约（**部分过期**：缺 follow 模块 10 条路由，Workspace/AI 章节后端未实现） |
| `backend-integration.md` | 前端接入真实后端的适配说明 |
| `deployment.md` / `env-config.md` | 部署与环境变量说明 |
| `项目原型文档.md` | V2.1 最终产品功能文档 |
| `归档/` | 一期 MVP、V2 增强版历史文档 |

## 变更记录（docs/changelog/）

**按需生成**：仅当用户明确要求时才写变更记录（信号词见 [CONVENTION.md](CONVENTION.md) §7），日常代码变更不自动生成。同模块积累 5 篇 active 后按 CONVENTION §6 整合进系统文档；`merged`/`archived` 状态的文件允许删除。现有 11 篇（均 active）：

| 文档 | 模块 | 日期 |
|------|------|------|
| [20260909_v1.0_修复帖子计数双重递增与删帖不对称回减.md](changelog/20260909_v1.0_修复帖子计数双重递增与删帖不对称回减.md) | post | 2026-09-09 |
| [20260909_v1.0_评论删除事务化并补齐用户评论数回减.md](changelog/20260909_v1.0_评论删除事务化并补齐用户评论数回减.md) | comment | 2026-09-09 |
| [20260909_v1.0_取消收藏补发PostUnfavorited事件.md](changelog/20260909_v1.0_取消收藏补发PostUnfavorited事件.md) | interaction | 2026-09-09 |
| [20260909_v1.0_任务领取真实发放积分与等级重算.md](changelog/20260909_v1.0_任务领取真实发放积分与等级重算.md) | growth | 2026-09-09 |
| [20260909_v1.0_新增用户封禁与帖子上下架接口.md](changelog/20260909_v1.0_新增用户封禁与帖子上下架接口.md) | admin | 2026-09-09 |
| [20260909_v1.0_worker开关生效与事件链路补齐.md](changelog/20260909_v1.0_worker开关生效与事件链路补齐.md) | infrastructure | 2026-09-09 |
| [20260909_v1.0_oss_type配置生效.md](changelog/20260909_v1.0_oss_type配置生效.md) | file | 2026-09-09 |
| [20260909_v1.0_未读数改为Redis缓存旁路.md](changelog/20260909_v1.0_未读数改为Redis缓存旁路.md) | notification | 2026-09-09 |
| [20260909_v1.0_取消互动热度减分与hot_score定时回写.md](changelog/20260909_v1.0_取消互动热度减分与hot_score定时回写.md) | rank | 2026-09-09 |
| [20260909_v1.0_话题参与人数落库重算.md](changelog/20260909_v1.0_话题参与人数落库重算.md) | topic | 2026-09-09 |
| [20260909_v1.0_话题索引增量同步.md](changelog/20260909_v1.0_话题索引增量同步.md) | search | 2026-09-09 |

## 资产

- `docs/assets/` —— 图片等素材（原散落在 backend docs 的 `已有的镜像.jpg` 已移入，未被任何文档引用）
