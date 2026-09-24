# 部署说明

## 1. 构建前端

```bash
npm install
npm run build
```

构建产物位于：

```text
dist/
```

## 2. 前后端分离部署

前端通过 Nginx 托管静态资源，后端独立部署为 API 服务。

生产环境 `.env.production` 示例：

```env
VITE_API_MODE=real
VITE_API_BASE_URL=https://api.example.com/api/v1
VITE_API_TIMEOUT=15000
VITE_ENABLE_API_LOG=false
```

## 3. Nginx 配置

```nginx
server {
    listen 80;
    server_name your-domain.com;

    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://backend:8080/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header Authorization $http_authorization;
    }
}
```

## 4. Dockerfile

```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:1.27-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
```

## 5. docker-compose 示例

```yaml
services:
  frontend:
    build: .
    ports:
      - "80:80"
    depends_on:
      - backend

  backend:
    image: your-backend-image
    ports:
      - "8080:8080"
```

## 6. 常见问题

### 页面刷新 404

需要 Nginx 配置：

```nginx
try_files $uri $uri/ /index.html;
```

### 真实后端请求失败

检查：

1. `VITE_API_MODE=real`
2. `VITE_API_BASE_URL` 是否正确
3. 后端是否允许 CORS
4. token 是否正确携带

### 登录后仍然 401

检查后端是否识别：

```http
Authorization: Bearer <token>
```
