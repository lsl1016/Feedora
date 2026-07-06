# 环境变量说明

项目支持 Mock 模式和真实后端 API 模式。

## 变量列表

| 变量 | 示例 | 说明 |
|---|---|---|
| `VITE_API_MODE` | `mock` / `real` | API 模式，`mock` 使用本地内存 Mock，`real` 请求真实后端 |
| `VITE_API_BASE_URL` | `http://localhost:8080/api/v1` | 后端接口基础地址 |
| `VITE_API_TIMEOUT` | `15000` | Axios 超时时间，单位毫秒 |
| `VITE_ENABLE_API_LOG` | `true` / `false` | 是否在控制台输出接口日志 |

## 开发环境

```env
VITE_API_MODE=mock
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_API_TIMEOUT=15000
VITE_ENABLE_API_LOG=true
```

## 联调真实后端

```env
VITE_API_MODE=real
VITE_API_BASE_URL=http://localhost:8080/api/v1
VITE_API_TIMEOUT=15000
VITE_ENABLE_API_LOG=true
```

## 生产环境

```env
VITE_API_MODE=real
VITE_API_BASE_URL=https://api.example.com/api/v1
VITE_API_TIMEOUT=15000
VITE_ENABLE_API_LOG=false
```

## 切换说明

页面组件不需要感知当前是 Mock 还是 Real API。组件只调用统一的 `api` 对象。

```ts
import { api } from '@/api';

const posts = await api.getPosts({ page: 1, pageSize: 10 });
```

`VITE_API_MODE=mock` 时走本地 Mock 数据；`VITE_API_MODE=real` 时走 Axios 请求。
