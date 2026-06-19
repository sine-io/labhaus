# Labhaus 部署指南

## 环境要求

- **操作系统**: Linux (推荐 Ubuntu 20.04+)
- **Node.js**: >= 20.0.0
- **pnpm**: >= 9.0.0
- **PostgreSQL**: >= 14
- **Docker** (可选，推荐用于开发环境)

## 部署方式

### 方式 1: Docker Compose (推荐用于开发)

```bash
# 克隆项目
git clone https://github.com/sine-io/labhaus.git
cd labhaus

# 安装依赖
pnpm install

# 启动所有服务
docker compose up -d

# 运行数据库迁移
cd apps/api
pnpm migrate

# 导入样式库数据（可选）
pnpm import-styles

# 启动 API 服务
pnpm dev
```

访问：
- API: http://localhost:3001
- PostgreSQL: localhost:5432
- Redis: localhost:6379
- MinIO: http://localhost:9001

### 方式 2: 生产部署

#### 1. 安装依赖

```bash
# 安装 Node.js 20+ 和 pnpm
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt-get install -y nodejs
npm install -g pnpm

# 克隆项目
git clone https://github.com/sine-io/labhaus.git
cd labhaus
pnpm install
```

#### 2. 配置环境变量

```bash
cp .env.example .env
```

编辑 `.env`：

```bash
# Database
DATABASE_URL=postgresql://labhaus:***@localhost:5432/labhaus

# Redis
REDIS_URL=redis://localhost:6379

# MinIO / S3
MINIO_ENDPOINT=s3.amazonaws.com
MINIO_ACCESS_KEY=your-acces...T

# JWT
JWT_SECRET=your-strong-s...duction
JWT_REFRESH_SECRET=your-refresh-secret-change-in-production
JWT_EXPIRES_IN=1h
JWT_REFRESH_EXPIRES_IN=7d

# API
API_PORT=3001
API_HOST=0.0.0.0
NODE_ENV=production

# CORS
CORS_ORIGIN=https://your-frontend-domain.com
```

#### 3. 构建

```bash
pnpm build
```

#### 4. 数据库设置

```bash
# 创建数据库
sudo -u postgres createdb labhaus
sudo -u postgres createuser labhaus -P

# 运行迁移
cd apps/api
pnpm migrate
```

#### 5. 启动服务

使用 PM2 管理进程：

```bash
# 安装 PM2
npm install -g pm2

# 启动 API
cd apps/api
pm2 start dist/index.js --name labhaus-api

# 保存 PM2 配置
pm2 save
pm2 startup
```

或使用 systemd：

```ini
# /etc/systemd/system/labhaus-api.service
[Unit]
Description=Labhaus API
After=network.target postgresql.service

[Service]
Type=simple
User=labhaus
WorkingDirectory=/home/labhaus/labhaus/apps/api
Environment="NODE_ENV=production"
EnvironmentFile=/home/labhaus/labhaus/.env
ExecStart=/usr/bin/node dist/index.js
Restart=always

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl enable labhaus-api
sudo systemctl start labhaus-api
sudo systemctl status labhaus-api
```

#### 6. Nginx 反向代理

```nginx
server {
    listen 80;
    server_name api.labhaus.io;

    location / {
        proxy_pass http://localhost:3001;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_cache_bypass $http_upgrade;
    }
}
```

```bash
sudo systemctl reload nginx
```

#### 7. SSL 证书 (Let's Encrypt)

```bash
sudo apt-get install certbot python3-certbot-nginx
sudo certbot --nginx -d api.labhaus.io
```

## 健康检查

```bash
curl http://localhost:3001/api/health
```

预期响应：
```json
{
  "status": "ok",
  "timestamp": "2026-06-19T..."
}
```

## 监控

### 日志

```bash
# PM2 日志
pm2 logs labhaus-api

# systemd 日志
sudo journalctl -u labhaus-api -f
```

### 性能监控

```bash
# PM2 监控
pm2 monit

# 或使用外部监控工具
# - Prometheus + Grafana
# - Datadog
# - New Relic
```

## 备份

### 数据库备份

```bash
# 每日备份脚本
#!/bin/bash
BACKUP_DIR=/var/backups/labhaus
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump labhaus | gzip > $BACKUP_DIR/labhaus_$DATE.sql.gz

# 保留最近 7 天
find $BACKUP_DIR -name "labhaus_*.sql.gz" -mtime +7 -delete
```

添加到 crontab：
```bash
0 2 * * * /path/to/backup-script.sh
```

## 故障排查

### API 无法启动

```bash
# 检查端口占用
sudo lsof -i :3001

# 检查数据库连接
psql $DATABASE_URL -c "SELECT 1"

# 查看日志
pm2 logs labhaus-api --lines 100
```

### 数据库连接失败

```bash
# 检查 PostgreSQL 服务
sudo systemctl status postgresql

# 检查连接字符串
echo $DATABASE_URL
```

### 性能问题

```bash
# 检查数据库查询
# 在 psql 中
EXPLAIN ANALYZE SELECT * FROM styles LIMIT 10;

# 检查索引
\d+ styles

# 添加缺失索引
CREATE INDEX IF NOT EXISTS idx_name ON table(column);
```

## 更新部署

```bash
# 拉取最新代码
cd ~/labhaus
git pull origin master

# 安装依赖
pnpm install

# 运行新迁移（如果有）
cd apps/api
pnpm migrate

# 重新构建
cd ../..
pnpm build

# 重启服务
pm2 restart labhaus-api
```

## 安全建议

1. **更改默认密钥**: 确保 `JWT_SECRET` 使用强随机字符串
2. **限制数据库访问**: 仅允许本地连接或特定 IP
3. **启用防火墙**: 只开放必要端口 (80, 443, 22)
4. **定期更新**: 及时更新系统和依赖包
5. **使用 HTTPS**: 生产环境必须启用 SSL
6. **备份策略**: 定期备份数据库和重要文件

## 扩展部署

### 负载均衡

使用多个 API 实例 + Nginx 负载均衡：

```nginx
upstream labhaus_api {
    server localhost:3001;
    server localhost:3002;
    server localhost:3003;
}

server {
    location / {
        proxy_pass http://labhaus_api;
    }
}
```

### Docker 生产部署

```bash
# 构建生产镜像
docker build -t labhaus-api:latest .

# 运行
docker run -d \
  --name labhaus-api \
  -p 3001:3001 \
  --env-file .env \
  labhaus-api:latest
```
