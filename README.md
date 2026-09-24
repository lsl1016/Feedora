# Feedora

Feedora 是开发者知识社区与本地知识工作台项目。

当前客户端采用 Monorepo 组织：

```text
apps/
├── web/               # React Web 社区
└── desktop/           # Wails + Go + React Desktop
packages/
├── contracts/         # 共享数据契约
├── app-core/          # Domain / Ports / UseCases / Capabilities
├── api-client/        # 平台无关 Feedora Server endpoint mapping
├── platform-web/      # Browser/Web Adapter
├── platform-desktop/  # Wails Adapter
├── features/          # 共享 Feature 元数据与后续业务组件
├── ui/                # Design Tokens / 共享 UI
└── editor/            # Markdown/Code Editor 边界
```

服务端当前继续保留在 `Feedora-backend/`，本轮重构不改变服务端架构。

## Web

```bash
pnpm install
pnpm dev:web
pnpm build:web
```

## Desktop UI

```bash
pnpm dev:desktop-ui
pnpm build:desktop-ui
```

## Wails Desktop

先构建 Desktop UI，然后：

```bash
cd apps/desktop
wails dev
```

完整客户端架构与依赖边界见 `docs/client-platform-architecture.md`。
