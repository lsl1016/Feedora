# Feedora Desktop

Feedora Desktop 是 Feedora Client Platform 的桌面宿主，使用 Wails v2 + Go + React。

## 边界

- React 只依赖 `@feedora/app-core` Port，不直接写 SQLite/Git/Wails 业务逻辑。
- Wails Binding 只存在于 `packages/platform-desktop` 和 Go `internal/bridge` 边界。
- 社区事实继续来自 Feedora Server。
- 私人 Workspace、Repo 索引和本地搜索属于 Desktop Local-first 能力。
- Agent Runtime 独立部署，通过 Adapter 接入，不在 Feedora Desktop 内重新实现。

## 本地启动

1. 在仓库根目录执行 `pnpm install`。
2. 安装 Wails v2 CLI。
3. 先执行 `pnpm build:desktop-ui`，然后进入 `apps/desktop` 执行 `wails dev`。

当前提交完成 Desktop Shell、Bridge 边界和模块骨架。Workspace/Repository/Search 的持久化实现将在后续工程任务中逐步落地。
