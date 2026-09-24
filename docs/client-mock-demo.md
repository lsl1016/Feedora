# Client Mock Demo

当前 `feat/client-platform-refactor` 在不启动 `Feedora-backend` 的情况下即可完整演示 Web 与 Desktop UI。

## 安装

```bash
pnpm install
```

## Web

```bash
pnpm dev:web
```

默认地址：`http://localhost:5173`。

Web 默认 `VITE_API_MODE=mock`，Mock 数据来自 `packages/mock-data`。如需连接真实服务端，在环境变量中切换：

```text
VITE_API_MODE=real
VITE_API_BASE_URL=http://localhost:8090/api/v1
```

## Desktop UI 浏览器调试

```bash
pnpm dev:desktop-ui
```

默认地址：`http://localhost:5174`。

Desktop UI 当前默认使用 Zustand + `@feedora/mock-data`，覆盖：

- 社区 Feed
- 工作台 Dashboard
- Markdown 笔记编辑器
- 知识库
- Repository Reader
- Unified Search
- Notification Center
- Settings
- AI Panel / Context

## Wails

```bash
pnpm build:desktop-ui
cd apps/desktop
wails dev
```

当前 Wails Go Bridge 保留真实平台边界，业务数据仍由 Mock UI 驱动。后续替换 WorkspacePort / RepositoryPort / SearchPort / AgentRuntimePort 时不需要重写页面。

## E2E

```bash
pnpm exec playwright install chromium
pnpm test:e2e
```

E2E 会同时启动 Web 和 Desktop UI，验证跨页面关键交互。
