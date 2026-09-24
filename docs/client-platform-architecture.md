# Feedora Client Platform 架构冻结

## 目标

Web 与 Desktop 作为同一 Client Platform 的两个宿主应用：

- Web 保持社区、公开内容和在线交互为中心；
- Desktop 是本地知识工作台 + 在线社区 + AI Agent；
- 共享领域契约、UseCase、Feature 和 Design Tokens；
- 不共享平台实现，不让 React Feature 直接依赖 Wails、SQLite 或 Browser API。

## 目录

```text
Feedora/
├── Feedora-backend/           # 服务端，本次不重构
├── apps/
│   ├── web/                   # 原 Feedora-frontend
│   └── desktop/
│       ├── frontend/          # React Desktop Shell
│       ├── internal/
│       │   ├── app/
│       │   ├── bridge/
│       │   ├── workspace/
│       │   ├── repository/
│       │   ├── search/
│       │   ├── indexing/
│       │   ├── runtime/
│       │   ├── securestore/
│       │   ├── native/
│       │   └── jobs/
│       └── wails.json
├── packages/
│   ├── contracts/
│   ├── app-core/
│   ├── api-client/
│   ├── platform-web/
│   ├── platform-desktop/
│   ├── features/
│   ├── ui/
│   └── editor/
└── docs/
```

## 依赖方向

```text
apps/web ─────────────┐
                     ├──> features / ui / editor
apps/desktop ─────────┘             │
                                    ▼
                                 app-core
                                    │
                                    ▼
                                   ports
                     ┌──────────────┴──────────────┐
                     ▼                             ▼
              platform-web                 platform-desktop
                     │                             │
                     ▼                             ▼
              Feedora Server                    Wails/Go
                                                   │
                                 ┌─────────────────┼───────────────┐
                                 ▼                 ▼               ▼
                              SQLite              Git        Agent Runtime
```

禁止反向依赖：

- app-core -> React / Axios / Wails
- features -> SQLite / Wails
- ui -> business service
- Web 页面直接调用 Desktop Binding
- Desktop React 直接访问本地数据库

## Port

首批 Port：

- CommunityPort
- WorkspacePort
- RepositoryPort
- SearchPort
- AgentRuntimePort
- NativePort
- SecureStoragePort

Web 使用 HTTP/Browser Adapter；Desktop 使用 Wails Adapter。

## Capability

不要在组件里散落 `if (isDesktop)`。客户端能力通过 Capability Registry 表达：

- community
- remote-search
- remote-agent-runtime
- local-workspace
- local-files
- repository
- local-search
- secure-storage
- native-notification
- local-agent-runtime

## Desktop 固定 Shell

```text
┌─────────────────────────────────────────────────────────────┐
│ Native Title / Command Search                              │
├─────┬──────────────┬────────────────────┬──────────────────┤
│ Act │ Context      │ Main Canvas        │ AI Panel         │
│ Rail│ Sidebar      │                    │                  │
├─────┴──────────────┴────────────────────┴──────────────────┤
│ Status Bar                                                  │
└─────────────────────────────────────────────────────────────┘
```

一级入口：社区 / 工作台 / 仓库 / 搜索 / 通知 / 设置。

## 迁移策略

1. 原 `Feedora-frontend` 迁到 `apps/web`，优先保证 Web 行为不变。
2. `src/types.ts` 迁到 `packages/contracts`，Web 保留兼容 re-export。
3. endpoint mapping 从 Web App 抽到 `packages/api-client`；Browser 鉴权/错误 UI 仍留在 Web Adapter。
4. 建立 app-core Port 和 Capability。
5. 建立 Desktop Wails Shell 与 Bridge 模块边界。
6. 后续逐页按照已冻结 UI 实现 Workspace / Repository / Unified Search / Notification / Settings。

## V1 数据所有权

- Community：Feedora Server 是事实源。
- Desktop Workspace：SQLite + 本地文件系统。
- Git：只读扫描和索引，V1 不做 IDE 写代码。
- Search：统一查询，不统一存储。
- Secret：Desktop 必须使用系统安全存储，不写 localStorage/SQLite 明文。
- Agent：独立 Runtime，通过 Adapter/Protocol 接入。
