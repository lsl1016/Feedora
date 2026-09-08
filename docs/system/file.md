---
title: 文件模块功能文档
date: 2026-09-09
version: v1.0
type: system
module: file
maintainer: Feedora 项目组
status: active
related_code:
  - internal/router/file_router.go
  - internal/api/file_api.go
  - internal/dto/file_dto.go
  - internal/service/file_service.go
  - internal/repository/file_repository.go
  - internal/model/file.go
  - pkg/ossx/storage.go
  - pkg/utils/file.go
summary: file 模块提供带鉴权的图片上传接口，经类型白名单与 5MB 大小校验后写入本地磁盘存储并落库 files 表，返回拼装后的公网访问 URL。
---

## 1. 模块概述

file 模块提供图片文件上传能力，服务于帖子图片、用户头像、圈子头像、话题封面、活动封面等业务场景。模块仅暴露一个需要登录鉴权的上传接口，上传经过大小与 MIME 类型双重校验后，通过可插拔的 `ossx.Storage` 抽象写入对象存储（当前为本地磁盘模拟实现），并在 `files` 表中记录文件元数据。接口返回可直接访问的文件 URL。

## 2. 接口清单

| 路径 | 方法 | 功能 | 控制器 |
|------|------|------|--------|
| `/files/upload` | POST | 上传图片文件（`multipart/form-data`，字段 `file` 必填、`bizType` 可选），需 Bearer 鉴权 | `FileAPI.Upload` |

完整路径为 `POST /api/v1/files/upload`，注册于 `internal/router/file_router.go` 的 `registerFile`，挂载 `authMW` 鉴权中间件。本模块无删除、查询接口。

## 3. 核心逻辑

### 3.1 上传校验

`FileService.Upload`（`internal/service/file_service.go`）执行以下校验，任一失败即拒绝上传：

| 校验项 | 规则 | 失败返回 |
|--------|------|----------|
| 大小上限 | `Size > 5MB`（常量 `maxUploadSize = 5 << 20`）拒绝 | 422 `文件大小不能超过 5MB` |
| MIME 白名单 | `image/png`、`image/jpeg`、`image/jpg`、`image/gif`、`image/webp` | MIME 不在白名单时降级为扩展名校验，仍不通过则返回 422 `仅支持上传图片文件` |
| 扩展名白名单（降级路径） | 取原始文件名扩展名并转小写，允许 `.png`、`.jpg`、`.jpeg`、`.gif`、`.webp` | 同上 |

校验细节：

- MIME 类型取自 multipart 文件头的 `Content-Type`，由客户端声明，服务端不读取文件内容做魔数校验。
- 大小取自 multipart 文件头声明的 `Size`，不按实际读取字节数复核。
- MIME 白名单命中时扩展名直接采用映射值（`image/jpeg` 与 `image/jpg` 均映射为 `.jpg`；`.jpeg` 扩展名仅在 MIME 未命中白名单、走文件名降级路径时保留）。
- `bizType` 为空时默认取 `post_image`。

### 3.2 存储路径（objectKey）生成

`objectKey` 格式为 `{顶层目录}/{yyyyMMdd}/{UUID}{ext}`，例如 `post/20260909/3f2b...c1.png`。其中日期取服务器当前时间，UUID 由 `utils.UUID()` 生成，原始文件名不参与 objectKey，同名文件不会相互覆盖。顶层目录按 `bizType` 映射：

| bizType | 顶层目录 |
|---------|----------|
| `avatar` | `avatar` |
| `circle_avatar` | `circle` |
| `topic_cover` | `topic` |
| `activity_cover` | `activity` |
| 其他（含默认 `post_image`） | `post` |

### 3.3 可插拔存储（ossx）

`pkg/ossx/storage.go` 定义存储抽象接口：

```go
type Storage interface {
    Upload(objectKey string, reader io.Reader, contentType string) (url string, err error)
    Delete(objectKey string) error
    PublicURL(objectKey string) string
}
```

当前仅有 `LocalStorage` 一个实现：

- 构造时对 `basePath` 执行 `MkdirAll(0755)`，`basePath` 为空时回退 `./uploads`；`publicBaseURL` 去尾部 `/`。
- `Upload` 将 objectKey 去前导 `/` 后 `filepath.Join` 到 `basePath`，逐级建目录、创建文件并 `io.Copy` 写入，成功后返回 `PublicURL(objectKey)`；`contentType` 参数被忽略。
- `Delete` 直接 `os.Remove` 对应文件，当前无任何调用方。
- `PublicURL` 拼装规则为 `{publicBaseURL}/{objectKey}`。

`pkg/ossx/minio.go` 与 `pkg/ossx/aliyun.go` 为阶段二占位文件，仅含注释、无任何实现代码。`internal/app/app.go` 在依赖装配处无条件调用 `ossx.NewLocalStorage(cfg.OSS.BasePath, cfg.OSS.PublicBaseUrl)` 构造存储实例：配置键 `oss.type` 虽已定义并被解析，但当前不参与实现选择，配置为其他值不会切换实现。

### 3.4 静态访问与 URL 返回

本地存储目录通过 `r.Static("/static", StaticDir)` 对外提供只读访问（`internal/router/router.go`），`StaticDir` 即 `storage.BasePath()`。默认配置下文件 URL 形如 `http://localhost:8090/static/post/20260909/{uuid}.png`，即 `{oss.publicBaseUrl}/{objectKey}`。上传成功后接口返回该 URL、原始文件名、大小与 MIME 类型。

### 3.5 落库行为

存储写入成功后调用 `FileRepository.Create` 写入 `files` 表。该调用的返回错误被忽略：数据库写入失败时接口仍返回成功，仅丢失文件记录，不产生孤立文件回滚。

## 4. 数据模型

`files` 表（`internal/model/file.go`，表名 `files`）：

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | int64 | 主键 |
| `user_id` | int64 | 上传用户 ID，带索引 |
| `biz_type` | string(64) | 业务类型，带索引 |
| `filename` | string(255) | 客户端原始文件名 |
| `object_key` | string(512) | 存储对象键（顶层目录/日期/UUID.扩展名） |
| `url` | string(512) | 文件公网访问地址 |
| `size` | int64 | 文件大小（字节），默认 0 |
| `mime_type` | string(128) | 客户端声明的 MIME 类型 |
| `created_at` | time | 上传时间 |

## 5. 配置项与限制

### 5.1 配置项

| 配置键 | 默认值（configs/config.yaml） | 说明 |
|--------|------------------------------|------|
| `oss.type` | `local` | 存储类型标识。当前实现不读取该值做选择，始终使用 LocalStorage |
| `oss.basePath` | `./uploads` | 本地存储根目录，空值时代码回退 `./uploads` |
| `oss.publicBaseUrl` | `http://localhost:8090/static` | 文件 URL 前缀，需与 `/static` 静态路由匹配 |

上传大小上限 5MB 与图片类型白名单为服务端硬编码常量（`internal/service/file_service.go` 的 `maxUploadSize`、`pkg/utils/file.go` 的 `imageMimeExt`），不通过配置修改。

### 5.2 当前实现的限制

- 仅支持图片上传，MIME 与扩展名均信任客户端声明，无文件内容（魔数）校验，伪造扩展名可通过校验。
- 大小校验依赖 multipart 头声明值，不按实际写入字节复核。
- 文件记录写入失败被静默忽略，接口仍返回成功。
- 无文件删除接口与文件记录清理机制，`Storage.Delete` 无调用方，存储空间只增不减。
- MinIO 与阿里云 OSS 实现为空占位，`oss.type` 配置项无实际分发作用。
- 本地磁盘存储无访问控制，`/static` 下所有文件可被匿名公开访问。
- 单机本地存储，多实例部署时文件不共享。

## 6. 历史版本

| 版本 | 日期 | 作者 | 说明 |
|------|------|------|------|
| v1.0 | 2026-09-09 | Feedora 项目组 | 初始版本 |
