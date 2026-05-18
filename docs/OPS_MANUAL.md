# 惠福星链 · 运维手册

## 一、系统架构概览

```
┌─────────────────────────────────────────────────────────────┐
│                        互联网用户                            │
├───────────────┬───────────────────────────────┬─────────────┤
│  微信小程序     │       管理后台 (Vue3)          │   API 消费者 │
└───────┬───────┴───────────────┬───────────────┴──────┬──────┘
        │                       │                      │
   ┌────▼───────────────────────▼──────────────────────▼──────┐
   │                    Nginx (TLS 终止 + 反代)                 │
   │              :443 (HTTPS) / :80 → 301 HTTPS              │
   └───────────────────────────┬──────────────────────────────┘
                               │
   ┌───────────────────────────▼──────────────────────────────┐
   │              惠福星链 Go 服务 (:8080)                      │
   │  ┌──────────┬──────────┬──────────┬──────────┐          │
   │  │ 时光轴   │ 同心圆    │ 灵犀AI    │ 星盘     │          │
   │  └──────────┴──────────┴──────────┴──────────┘          │
   └───┬───────────┬───────────────┬──────────────────────────┘
       │           │               │
  ┌────▼────┐ ┌───▼────┐  ┌───────▼──────┐
  │ MySQL   │ │ Redis  │  │  RabbitMQ     │
  │ 8.0     │ │ 7      │  │  (随访/通知)   │
  └─────────┘ └────────┘  └──────────────┘
```

| 组件 | 端口 | 用途 |
|------|------|------|
| Nginx | 80/443 | 反向代理 + TLS |
| Go Server | 8080 | 核心业务 API |
| MySQL | 3306 | 主数据库 |
| Redis | 6379 | 缓存/会话 |
| RabbitMQ | 5672/15672 | 消息队列/管理面板 |
| Admin | 3000 (dev) | 管理后台开发模式 |

---

## 二、快速部署

### 2.1 环境要求

| 软件 | 最低版本 | 说明 |
|------|----------|------|
| Go | 1.22 | 编译后端 |
| Node.js | 18+ | 编译管理后台 |
| MySQL | 8.0 | 数据库 |
| Redis | 7 | 缓存（可选） |
| RabbitMQ | 3.x | 消息队列（可选） |
| Nginx | 1.24+ | 反向代理 |
| golang-migrate | 4.x | 数据库迁移 |

### 2.2 Docker Compose 一键部署

```bash
# 1. 配置环境变量
cp .env.example .env
vim .env   # 填入实际密钥

# 2. 启动全部服务
docker-compose -f deploy/docker-compose.yml up -d

# 3. 查看状态
docker-compose -f deploy/docker-compose.yml ps

# 4. 查看日志
docker-compose -f deploy/docker-compose.yml logs -f server
```

### 2.3 手动部署（自有服务器）

```bash
# 1. 编译
go build -o bin/huifu-server ./cmd/server
cd admin && npm ci && npm run build && cd ..

# 2. 创建部署目录
sudo mkdir -p /opt/huifu /data/huifu/files /var/log/huifu
sudo useradd -r -s /bin/false huifu 2>/dev/null || true
sudo chown -R huifu:huifu /opt/huifu /data/huifu /var/log/huifu

# 3. 复制文件
sudo cp bin/huifu-server /opt/huifu/
sudo cp config.yaml /opt/huifu/
sudo cp -r migrations /opt/huifu/
sudo cp -r admin/dist /opt/huifu/  # 如果有

# 4. 配置环境变量
sudo cp .env.example /opt/huifu/.env
sudo vim /opt/huifu/.env

# 5. 执行数据库迁移
migrate -path /opt/huifu/migrations -database "${MYSQL_DSN}" up

# 6. 安装 systemd 服务
sudo cp deploy/huifu-server.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable huifu-server
sudo systemctl start huifu-server

# 7. 配置 Nginx
sudo cp deploy/nginx.conf /etc/nginx/sites-available/huifu
sudo ln -sf /etc/nginx/sites-available/huifu /etc/nginx/sites-enabled/
sudo nginx -t && sudo systemctl reload nginx
```

### 2.4 一键部署脚本

```bash
# 用法: deploy.sh [staging|production]
./deploy/deploy.sh production
```

脚本会依次执行：**编译 → 前端构建 → 数据库迁移 → 原子切换版本 → 重启服务 → 健康检查 → 清理旧版本**

### 2.5 回滚

```bash
# 回滚到上一个版本
OLD=$(ls -dt /opt/huifu/backups/*/ 2>/dev/null | head -1)
if [ -n "$OLD" ]; then
    ln -sfn "$OLD" /opt/huifu/current
    systemctl restart huifu-server
fi
```

---

## 三、健康检查

### 3.1 健康端点

| 端点 | 用途 | 频率 |
|------|------|------|
| `GET /api/health` | 综合健康（含 MySQL ping） | Prometheus 每 15s |
| `curl http://localhost:8080/api/health` | 手动检查 | 部署后 |

### 3.2 检查脚本

```bash
#!/bin/bash
# scripts/health_check.sh
RESP=$(curl -sf http://localhost:8080/api/health)
if [ $? -eq 0 ]; then
    echo "OK: $RESP"
else
    echo "CRITICAL: huifu-server is down"
    exit 2
fi
```

---

## 四、监控告警

### 4.1 关键指标

| 指标 | 告警阈值 | 说明 |
|------|----------|------|
| CPU 使用率 | > 80% 持续 5min | 可能需要扩容 |
| 内存使用率 | > 2GB | systemd MemoryLimit |
| API 响应时间 (P95) | > 1000ms | 检查 DB 慢查询 |
| 错误率 (5xx) | > 5% | 检查日志 |
| MySQL 连接数 | > 80% 最大连接 | 检查连接池泄漏 |
| 磁盘使用率 | > 80% | 清理旧版本/日志 |
| 服务存活性 | DOWN | 自动重启 |

### 4.2 Prometheus 指标（P1 规划）

```
# HELP huifu_http_requests_total Total HTTP requests
# HELP huifu_http_request_duration_seconds HTTP request duration
# HELP huifu_ai_tokens_used_total Total AI tokens used
# HELP huifu_member_count Current member count
```

占位文件：预留 OpenTelemetry 埋点，在 handler 层添加 `otelhttp` 中间件。

### 4.3 日志监控

```bash
# 实时查看 API 请求
journalctl -u huifu-server -f

# 按时间过滤
journalctl -u huifu-server --since "10 minutes ago"

# 查看错误日志
journalctl -u huifu-server -p err

# Nginx 访问日志
tail -f /var/log/nginx/huifu-access.log
```

---

## 五、备份与恢复

### 5.1 数据库备份

```bash
#!/bin/bash
# 每日备份脚本（建议放在 cron.daily）
BACKUP_DIR="/opt/huifu/backups/db"
mkdir -p "$BACKUP_DIR"
DATE=$(date +%Y%m%d)
mysqldump -u root -p"${MYSQL_ROOT_PASSWORD}" --single-transaction \
    --routines --triggers --events huifu | gzip > "${BACKUP_DIR}/huifu_${DATE}.sql.gz"
# 保留最近 30 天
find "$BACKUP_DIR" -name "*.sql.gz" -mtime +30 -delete
```

### 5.2 文件备份

```bash
# 用户上传文件（报告 PDF、影像）
rsync -av /data/huifu/files/ /opt/huifu/backups/files/
```

### 5.3 恢复

```bash
# 1. 停止服务
systemctl stop huifu-server

# 2. 恢复数据库
gunzip -c /opt/huifu/backups/db/huifu_20260501.sql.gz | mysql -u root -p huifu

# 3. 启动服务
systemctl start huifu-server
```

### 5.4 灾备建议

| 级别 | 方案 | RPO | RTO |
|------|------|-----|-----|
| 单机 | mysqldump + 本地文件备份 | 24h | 1h |
| 同城 | MySQL 主从复制 + rsync | 秒级 | 5min |
| 异地 | 跨云备份 + 冷备恢复 | 小时级 | 2h |

---

## 六、故障排查 SOP

### 6.1 服务无法启动

```bash
# 检查 systemd 日志
journalctl -u huifu-server -n 50 --no-pager

# 常见原因:
#   MySQL DSN 错误 → 检查 /opt/huifu/.env
#   端口占用 → lsof -i :8080
#   配置文件格式错误 → go run ./cmd/server 本地调试

# 手动启动调试
cd /opt/huifu && sudo -u huifu ./huifu-server
```

### 6.2 API 返回 500 错误

```bash
# 1. 检查 app 日志
journalctl -u huifu-server -p err --since "5 minutes ago"

# 2. 检查 MySQL 连接
mysql -u root -p -e "SHOW PROCESSLIST;"

# 3. 检查 Redis
redis-cli -a "$REDIS_PASSWORD" PING

# 4. 检查磁盘空间
df -h
```

### 6.3 AI 对话无响应

```bash
# 1. 检查 AI 配置
cat /opt/huifu/.env | grep API_KEY

# 2. 测试 DeepSeek 连通性
curl -s -H "Authorization: Bearer $DEEPSEEK_API_KEY" \
     -H "Content-Type: application/json" \
     -d '{"model":"deepseek-chat","messages":[{"role":"user","content":"hi"}]}' \
     https://api.deepseek.com/v1/chat/completions

# 3. 检查 FAQ 表
mysql -u root -p huifu -e "SELECT COUNT(*) FROM faq_entries WHERE status='published';"
```

### 6.4 微信登录失败

```bash
# 1. 检查小程序配置
cat /opt/huifu/.env | grep WX

# 2. 测试 code2session
curl "https://api.weixin.qq.com/sns/jscode2session?appid=$WX_MP_APPID&secret=$WX_MP_SECRET&js_code=TEST&grant_type=authorization_code"

# 3. 检查 wechat_bindings 表
mysql -u root -p huifu -e "SELECT COUNT(*) FROM wechat_bindings;"
```

### 6.5 紧急联系人

| 角色 | 联系方式 | 职责 |
|------|----------|------|
| 后端负责人 | -- | API / DB 问题 |
| 前端负责人 | -- | 小程序 / 管理后台问题 |
| 运维负责人 | -- | 服务器 / 网络问题 |
| 产品负责人 | -- | 业务逻辑确认 |
| 医院信息科 | -- | HIS/LIS 接口问题 |

---

## 七、日志管理

### 7.1 日志位置

| 日志 | 路径 |
|------|------|
| 应用日志 | `journalctl -u huifu-server` |
| Nginx 访问 | `/var/log/nginx/huifu-access.log` |
| Nginx 错误 | `/var/log/nginx/huifu-error.log` |

### 7.2 日志轮转（/etc/logrotate.d/huifu）

```
/var/log/nginx/huifu-*.log {
    daily
    rotate 30
    missingok
    notifempty
    compress
    delaycompress
    sharedscripts
    postrotate
        [ -f /var/run/nginx.pid ] && kill -USR1 $(cat /var/run/nginx.pid)
    endscript
}
```

### 7.3 敏感信息脱敏

- 手机号/姓名在应用层已 SHA256 + AES 双重存储
- 日志中不输出 phone/encrypted_name 字段
- AI 对话日志不追加到 Nginx access log（请求体不记录）

---

## 八、安全运维

### 8.1 证书管理

```bash
# 使用 certbot 自动续签 Let's Encrypt
certbot certonly --nginx -d api.huifu.example.com -d admin.huifu.example.com

# 自动续签 cron
0 3 * * * certbot renew --quiet --post-hook "systemctl reload nginx"
```

### 8.2 防火墙规则

```bash
# 仅开放必要端口
ufw allow 22/tcp      # SSH
ufw allow 80/tcp      # HTTP (redirect to HTTPS)
ufw allow 443/tcp     # HTTPS
ufw deny 3306         # MySQL 仅本地
ufw deny 6379         # Redis 仅本地
ufw enable
```

### 8.3 定期安全检查清单

- [ ] 系统更新：`apt update && apt upgrade`
- [ ] 证书到期检查：`certbot certificates`
- [ ] 数据库备份完整性：抽查最近 `.sql.gz` 可正常解压
- [ ] 磁盘使用率：`df -h`
- [ ] 异常登录审计：`last -20` / `faillog`
- [ ] Go 依赖漏洞：`go vet ./...`
- [ ] Nginx 配置校验：`nginx -t`

---

## 九、性能调优

### 9.1 MySQL

```sql
-- 慢查询日志
SET GLOBAL slow_query_log = ON;
SET GLOBAL long_query_time = 0.5;

-- 关键查询索引检查
EXPLAIN SELECT * FROM timeline_events WHERE member_id = ? ORDER BY event_date DESC;
EXPLAIN SELECT * FROM ai_conversations WHERE session_id = ? ORDER BY created_at;
```

### 9.2 Go 服务

```bash
# 内存分析
go tool pprof http://localhost:8080/debug/pprof/heap

# goroutine 泄漏检查
go tool pprof http://localhost:8080/debug/pprof/goroutine

# 设置 GOMAXPROCS
# /opt/huifu/.env: GOMAXPROCS=4
```

### 9.3 Nginx 调优

```nginx
# worker 进程数 = CPU 核数
worker_processes auto;
worker_connections 2048;

# 静态资源缓存
location ~* \.(js|css|png|jpg|svg)$ {
    expires 30d;
    add_header Cache-Control "public, immutable";
}
```

---

## 十、附录

### A. 常用命令速查

```bash
# 服务管理
systemctl start|stop|restart|status huifu-server

# 查看实时日志
journalctl -u huifu-server -f

# 数据库迁移
migrate -path migrations -database "$MYSQL_DSN" up
migrate -path migrations -database "$MYSQL_DSN" down 1

# 健康检查
curl http://localhost:8080/api/health | jq .

# 部署
./deploy/deploy.sh production
```

### B. 目录结构

```
/opt/huifu/
├── current → /opt/huifu/releases/20260518120000/   # 当前版本
├── releases/                                         # 各版本目录
│   └── 20260518120000/
│       ├── huifu-server
│       ├── config.yaml
│       └── migrations/
├── backups/
│   ├── db/            # 数据库备份
│   ├── files/         # 用户文件备份
│   └── 20260518120000/ # 旧版本回滚备份
└── .env               # 环境变量
/data/huifu/files/     # 用户上传文件存储
/var/log/huifu/        # 应用日志
```

### C. 基础设施配置速查

| 配置项 | 文件位置 | 说明 |
|--------|----------|------|
| 服务器环境变量 | `/opt/huifu/.env` | MYSQL_DSN, API keys |
| Nginx 配置 | `/etc/nginx/sites-enabled/huifu` | 反代规则 |
| Systemd 服务 | `/etc/systemd/system/huifu-server.service` | 守护进程 |
| Logrotate | `/etc/logrotate.d/huifu` | 日志轮转 |
| Cron 备份 | `/etc/cron.daily/huifu-backup` | 数据库备份 |
