# 后端接入说明

本文档说明社区 V2.1 前端项目如何接入真实后端服务。前端已经支持 Mock / Real API 切换，不需要改页面代码。

## 1. 前端运行模式

通过环境变量控制：

```env
VITE_API_MODE=mock
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

| 模式 | 说明 |
|---|---|
| `mock` | 使用当前前端内存 Mock 数据，适合纯前端开发 |
| `real` | 通过 Axios 请求真实后端，适合联调和生产部署 |

## 2. 后端统一响应格式

```ts
interface ApiResponse<T> {
  code: number;
  message: string;
  data: T;
  traceId?: string;
  timestamp?: string;
}
```

成功示例：

```json
{
  "code": 0,
  "message": "success",
  "data": {},
  "traceId": "trace-abc-123",
  "timestamp": "2026-07-05T10:00:00+08:00"
}
```

分页结构：

```ts
interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  pageSize: number;
}
```

## 3. Token 约定

登录成功后，前端保存 token 到：

```text
community_v21_token
```

请求时自动携带：

```http
Authorization: Bearer <token>
```

401 时前端会清除登录态，并跳转登录页。

## 4. CORS 配置

开发环境需要允许：

```text
Access-Control-Allow-Origin: http://localhost:5173
Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
Access-Control-Allow-Headers: Content-Type, Authorization
Access-Control-Allow-Credentials: true
```

Gin 示例：

```go
r.Use(cors.New(cors.Config{
    AllowOrigins:     []string{"http://localhost:5173"},
    AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
    ExposeHeaders:    []string{"Content-Length"},
    AllowCredentials: true,
}))
```

## 5. 后端需要实现的模块

| 模块 | 说明 |
|---|---|
| Auth | 登录、注册、当前用户 |
| User | 用户资料、我的帖子、我的评论、收藏、点赞 |
| Post | 内容流、帖子详情、发帖、编辑、隐藏、分享、转发 |
| Comment | 评论、回复、删除评论 |
| Tag | 标签列表、标签聚合 |
| Circle | 圈子、成员管理、聊天、公告、智能摘要 |
| Topic | 话题列表、话题详情、话题内容 |
| Search | 搜索、搜索联想、热门搜索 |
| Growth | 签到、任务、排行榜、勋章 |
| Activity | 活动列表、活动详情、投票、打卡 |
| Workspace | 笔记、知识库、协作成员 |
| AI | 用户模型配置、AI 会话、话题 Agent |
| Admin | 用户、帖子、评论、审核、举报、标签、圈子、话题 |
| Admin AI | 基础模型、Token 配置、Token 使用情况 |

## 6. 错误码约定

| code | 说明 |
|---:|---|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未登录或 token 过期 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源冲突 |
| 422 | 业务校验失败 |
| 429 | 请求过于频繁 |
| 500 | 服务内部错误 |
| 10001 | 账号或密码错误 |
| 10002 | 账号已被封禁 |
| 20001 | 帖子不存在 |
| 30001 | 圈子不存在 |
| 30002 | 你已被该圈子禁言 |
| 40001 | AI 模型未配置 |
| 40002 | Token 额度不足 |

## 7. 本地联调步骤

1. 启动后端服务，监听 `http://localhost:8080`。
2. 后端接口统一前缀为 `/api/v1`。
3. 前端 `.env.development` 设置：

```env
VITE_API_MODE=real
VITE_API_BASE_URL=http://localhost:8080/api/v1
```

4. 前端运行：

```bash
npm install
npm run dev
```

5. 浏览器访问 `http://localhost:5173`。

## 8. 文件上传

接口：

```http
POST /api/v1/files/upload
Content-Type: multipart/form-data
```

字段：

| 字段 | 说明 |
|---|---|
| file | 文件 |
| bizType | `avatar/post_image/circle_avatar/topic_cover/activity_cover` |

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "url": "https://cdn.example.com/images/xxx.png",
    "filename": "xxx.png",
    "size": 123456,
    "mimeType": "image/png"
  }
}
```
