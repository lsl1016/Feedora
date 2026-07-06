# 社区 V2.1 前端项目

本项目是社区 V2.1 前端原型，技术栈为 React + TypeScript + React Router + Zustand + Axios + Ant Design + Vite。

本版本包含两部分：

1. 局部组件优化：首页筛选栏、顶部导航图标、热门榜单、首页右侧栏、排行榜、话题广场、圈子智能摘要、我的主页上半部分。
2. 后端接入适配：支持 Mock / Real API 切换，补充 Axios 封装、环境变量、接口契约和部署说明文档。

## 运行

```bash
npm install
npm run dev
```

## 构建

```bash
npm run build
```

## API 模式

复制 `.env.example`，根据需要设置：

```env
VITE_API_MODE=mock
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_API_TIMEOUT=15000
VITE_ENABLE_API_LOG=true
```

- `mock`：使用前端内存 Mock 数据。
- `real`：通过 Axios 请求真实后端。

## 默认账号

普通用户：

```text
账号：zhangsan
密码：123456
```

管理员：

```text
账号：admin
密码：123456
```

## 文档

- `docs/backend-integration.md`：后端接入说明
- `docs/api-contract.md`：接口契约
- `docs/env-config.md`：环境变量说明
- `docs/deployment.md`：部署说明
