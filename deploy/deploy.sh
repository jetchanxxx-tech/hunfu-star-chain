#!/usr/bin/env bash
set -euo pipefail

# ============================================================
# 惠福星链 · 部署脚本
# 用法: ./deploy/deploy.sh [staging|production]
# ============================================================

ENV="${1:-production}"
APP="huifu-server"
DEPLOY_DIR="/opt/huifu"
RELEASE_DIR="${DEPLOY_DIR}/releases/$(date +%Y%m%d%H%M%S)"
CURRENT_LINK="${DEPLOY_DIR}/current"
BACKUP_DIR="${DEPLOY_DIR}/backups"
BINARY="bin/${APP}"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; NC='\033[0m'

log()  { echo -e "${GREEN}[$(date +'%H:%M:%S')]${NC} $*"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
err()  { echo -e "${RED}[ERROR]${NC} $*"; exit 1; }

# --- Pre-flight checks ---
log "开始部署 (环境: ${ENV})"

command -v go >/dev/null 2>&1 || err "go 未安装"
command -v node >/dev/null 2>&1 || warn "node 未安装，跳过前端构建"
command -v migrate >/dev/null 2>&1 || warn "golang-migrate 未安装，跳过数据库迁移"

# --- Build backend ---
log "构建 Go 后端..."
go build -ldflags="-s -w" -o "${BINARY}" ./cmd/server

# --- Build admin dashboard (base=/admin/) ---
if command -v node >/dev/null 2>&1; then
    log "构建管理后台 (base=/admin/)..."
    cd admin
    npm ci --production=false
    npm run build
    cd ..
fi

# --- Build H5 用户端 ---
if command -v node >/dev/null 2>&1 && [ -d "miniprogram" ] && [ -f "miniprogram/package.json" ]; then
    log "构建 H5 用户端..."
    cd miniprogram
    npm ci --production=false
    npm run build:h5
    cd ..
fi

# --- Prepare release directory ---
log "准备发布目录: ${RELEASE_DIR}"
mkdir -p "${RELEASE_DIR}" "${BACKUP_DIR}"
cp "${BINARY}" "${RELEASE_DIR}/${APP}"
cp config.yaml "${RELEASE_DIR}/"
cp -r migrations "${RELEASE_DIR}/"

# Static files → 固定路径 (nginx 直接引用，不经过 current 软链)
STATIC_ADMIN="${DEPLOY_DIR}/admin-dist"
STATIC_H5="${DEPLOY_DIR}/h5-dist"

if [ -d "admin/dist" ]; then
    log "部署管理后台 → ${STATIC_ADMIN}"
    # 备份旧版本
    if [ -d "${STATIC_ADMIN}" ]; then
        cp -r "${STATIC_ADMIN}" "${BACKUP_DIR}/admin-dist-$(date +%Y%m%d%H%M%S)" 2>/dev/null || true
    fi
    rm -rf "${STATIC_ADMIN}"
    cp -r admin/dist "${STATIC_ADMIN}"
fi

if [ -d "miniprogram/dist" ]; then
    log "部署 H5 用户端 → ${STATIC_H5}"
    if [ -d "${STATIC_H5}" ]; then
        cp -r "${STATIC_H5}" "${BACKUP_DIR}/h5-dist-$(date +%Y%m%d%H%M%S)" 2>/dev/null || true
    fi
    rm -rf "${STATIC_H5}"
    cp -r miniprogram/dist "${STATIC_H5}"
fi

# --- Database migration ---
if command -v migrate >/dev/null 2>&1; then
    log "执行数据库迁移..."
    source "${DEPLOY_DIR}/.env" 2>/dev/null || true
    migrate -path "${RELEASE_DIR}/migrations" -database "${MYSQL_DSN}" up || warn "数据库迁移失败"
fi

# --- Atomic swap ---
log "切换版本..."
if [ -L "${CURRENT_LINK}" ]; then
    OLD=$(readlink -f "${CURRENT_LINK}")
    cp -r "${OLD}" "${BACKUP_DIR}/$(basename ${OLD})" 2>/dev/null || true
fi

ln -sfn "${RELEASE_DIR}" "${CURRENT_LINK}"

# --- Restart service ---
if command -v systemctl >/dev/null 2>&1; then
    log "重启服务..."
    systemctl restart "${APP}" || err "服务重启失败"

    # Wait for healthy
    for i in $(seq 1 30); do
        if curl -sf http://localhost:8080/api/health >/dev/null 2>&1; then
            log "服务健康检查通过"
            break
        fi
        if [ $i -eq 30 ]; then
            err "服务健康检查超时，回滚..."
        fi
        sleep 1
    done

    systemctl status "${APP}" --no-pager
elif command -v supervisorctl >/dev/null 2>&1; then
    supervisorctl restart "${APP}"
fi

# --- Cleanup old releases (keep last 5) ---
log "清理旧版本..."
ls -dt ${DEPLOY_DIR}/releases/*/ 2>/dev/null | tail -n +6 | xargs rm -rf 2>/dev/null || true
ls -dt ${BACKUP_DIR}/*/ 2>/dev/null | tail -n +6 | xargs rm -rf 2>/dev/null || true

log "部署完成!"
