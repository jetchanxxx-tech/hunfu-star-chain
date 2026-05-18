#!/usr/bin/env bash
set -euo pipefail

# ============================================================
#  惠福星链 · 一键安装部署脚本
#  用法: curl -fsSL <url>/install.sh | bash
#       或 chmod +x install.sh && ./install.sh
# ============================================================

APP="huifu-server"
APP_USER="huifu"
INSTALL_DIR="/opt/huifu"
DATA_DIR="/data/huifu/files"
LOG_DIR="/var/log/huifu"
VENV=""
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

RED='\033[0;31m'; GREEN='\033[0;32m'; YELLOW='\033[1;33m'; CYAN='\033[0;36m'; NC='\033[0m'

log()    { echo -e "${GREEN}[✓]${NC} $*"; }
warn()   { echo -e "${YELLOW}[!]${NC} $*"; }
err()    { echo -e "${RED}[✗]${NC} $*"; exit 1; }
info()   { echo -e "${CYAN}[→]${NC} $*"; }
section(){ echo -e "\n${CYAN}══╡ $* ╞══${NC}"; }

# ==================== 系统检测 ====================
detect_os() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        OS=$ID
        VER=$VERSION_ID
    elif [ -f /etc/redhat-release ]; then
        OS="rhel"
    else
        OS="unknown"
    fi
    log "检测到系统: ${OS} ${VER:-unknown}"

    case $OS in
        ubuntu|debian) PKG_MGR="apt-get"; PKG_UPDATE="apt-get update -qq" ;;
        centos|rhel|fedora|rocky|almalinux) PKG_MGR="dnf"; PKG_UPDATE="dnf check-update -q || true" ;;
        *) err "不支持的系统: $OS" ;;
    esac
}

# ==================== 组件检测与安装 ====================
check_cmd() { command -v "$1" >/dev/null 2>&1; }

install_mysql() {
    if check_cmd mysqld || check_cmd mariadbd || systemctl is-active --quiet mysql 2>/dev/null || systemctl is-active --quiet mariadb 2>/dev/null; then
        log "MySQL/MariaDB 已运行，跳过安装"
        return 0
    fi
    warn "MySQL 未安装，正在安装..."
    case $OS in
        ubuntu|debian)
            $PKG_MGR install -y mysql-server-8.0 2>/dev/null || $PKG_MGR install -y mysql-server
            ;;
        centos|rhel|fedora|rocky|almalinux)
            $PKG_MGR install -y mysql-server
            ;;
    esac
    systemctl enable --now mysql || systemctl enable --now mysqld
    log "MySQL 安装完成"

    # 创建数据库（如果不存在）
    mysql -u root -e "CREATE DATABASE IF NOT EXISTS huifu CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;" 2>/dev/null || {
        # 尝试用临时密码
        TEMP_PWD=$(grep 'temporary password' /var/log/mysqld.log 2>/dev/null | tail -1 | awk '{print $NF}' || echo "")
        if [ -n "$TEMP_PWD" ]; then
            warn "MySQL 初始密码: $TEMP_PWD，请登录后修改"
        fi
    }
}

install_redis() {
    if systemctl is-active --quiet redis 2>/dev/null || systemctl is-active --quiet redis-server 2>/dev/null; then
        log "Redis 已运行，跳过安装"
        return 0
    fi
    warn "Redis 未安装，正在安装..."
    case $OS in
        ubuntu|debian) $PKG_MGR install -y redis-server ;;
        centos|rhel|fedora|rocky|almalinux) $PKG_MGR install -y redis ;;
    esac
    systemctl enable --now redis || systemctl enable --now redis-server
    log "Redis 安装完成"
}

install_golang() {
    GO_VER="1.22.10"
    if check_cmd go; then
        CURRENT=$(go version | grep -oP 'go\K[0-9.]+')
        if [ "$(printf '%s\n' "1.22" "$CURRENT" | sort -V | head -1)" = "1.22" ]; then
            log "Go ${CURRENT} >= 1.22，跳过安装"
            # 仍然配置国内代理
            go env -w GOPROXY=https://goproxy.cn,direct 2>/dev/null || true
            return 0
        fi
    fi
    warn "安装 Go ${GO_VER}..."
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)  GO_ARCH="amd64" ;;
        aarch64) GO_ARCH="arm64" ;;
        *) err "不支持的架构: $ARCH" ;;
    esac
    # 国内优先用镜像，国外直连
    GO_URL="https://golang.google.cn/dl/go${GO_VER}.linux-${GO_ARCH}.tar.gz"
    curl -fsSL "$GO_URL" -o /tmp/go.tar.gz 2>/dev/null || \
    curl -fsSL "https://go.dev/dl/go${GO_VER}.linux-${GO_ARCH}.tar.gz" -o /tmp/go.tar.gz
    rm -rf /usr/local/go
    tar -C /usr/local -xzf /tmp/go.tar.gz
    rm -f /tmp/go.tar.gz
    export PATH=/usr/local/go/bin:$PATH
    echo 'export PATH=/usr/local/go/bin:$PATH' > /etc/profile.d/go.sh
    # 配置国内 Go 模块代理
    go env -w GOPROXY=https://goproxy.cn,direct
    go env -w GONOSUMDB=*
    log "Go ${GO_VER} 安装完成 (GOPROXY=goproxy.cn)"
}

install_node() {
    if check_cmd node; then
        NODE_VER=$(node -v | grep -oP '\d+' | head -1)
        if [ "$NODE_VER" -ge 18 ]; then
            log "Node.js $(node -v) >= 18，跳过安装"
            return 0
        fi
    fi
    warn "安装 Node.js 20 LTS..."
    case $OS in
        ubuntu|debian)
            curl -fsSL https://deb.nodesource.com/setup_20.x | bash - 2>/dev/null || \
            curl -fsSL https://mirrors.tuna.tsinghua.edu.cn/nodesource/deb/setup_20.x | bash -
            $PKG_MGR install -y nodejs
            ;;
        centos|rhel|fedora|rocky|almalinux)
            curl -fsSL https://rpm.nodesource.com/setup_20.x | bash -
            $PKG_MGR install -y nodejs
            ;;
    esac
    # 国内镜像加速
    npm config set registry https://registry.npmmirror.com 2>/dev/null || true
    log "Node.js $(node -v) 安装完成"
}

install_migrate() {
    if check_cmd migrate; then
        log "golang-migrate 已安装，跳过"
        return 0
    fi
    warn "安装 golang-migrate..."
    MIG_VER="4.17.0"
    ARCH=$(uname -m)
    case $ARCH in
        x86_64)  MIG_ARCH="amd64" ;;
        aarch64) MIG_ARCH="arm64" ;;
        *) err "不支持的架构: $ARCH" ;;
    esac
    # 国内优先用镜像
    MIG_URL="https://github.com/golang-migrate/migrate/releases/download/v${MIG_VER}/migrate.linux-${MIG_ARCH}.tar.gz"
    curl -fsSL "https://ghproxy.com/$MIG_URL" -o /tmp/migrate.tar.gz 2>/dev/null || \
    curl -fsSL "$MIG_URL" -o /tmp/migrate.tar.gz
    tar -C /usr/local/bin -xzf /tmp/migrate.tar.gz migrate
    rm -f /tmp/migrate.tar.gz
    log "golang-migrate 安装完成"
}

install_all() {
    section "1/6 系统更新与基础依赖"
    $PKG_UPDATE
    case $OS in
        ubuntu|debian) $PKG_MGR install -y curl wget tar gzip ca-certificates systemd 2>/dev/null || true ;;
        centos|rhel|fedora|rocky|almalinux) $PKG_MGR install -y curl wget tar gzip ca-certificates systemd 2>/dev/null || true ;;
    esac

    section "2/6 安装 MySQL 8.0"
    install_mysql

    section "3/6 安装 Redis 7"
    install_redis

    section "4/6 安装 Go 1.22"
    install_golang

    section "5/6 安装 Node.js 20"
    install_node

    section "6/6 安装 golang-migrate"
    install_migrate

    export PATH=/usr/local/go/bin:$PATH
}

# ==================== 用户与目录 ====================
setup_user_and_dirs() {
    info "创建用户与目录..."
    id -u "$APP_USER" >/dev/null 2>&1 || useradd -r -s /bin/false "$APP_USER"
    mkdir -p "$INSTALL_DIR" "$DATA_DIR" "$LOG_DIR" "${INSTALL_DIR}/backups"
    chown -R "$APP_USER":"$APP_USER" "$INSTALL_DIR" "$DATA_DIR" "$LOG_DIR"
    log "目录创建完成"
}

# ==================== 环境配置 ====================
setup_env() {
    if [ -f "${INSTALL_DIR}/.env" ]; then
        log ".env 已存在，跳过"
        return 0
    fi

    if [ -f "${SCRIPT_DIR}/.env.example" ]; then
        cp "${SCRIPT_DIR}/.env.example" "${INSTALL_DIR}/.env"
    elif curl -fsSL "https://raw.githubusercontent.com/jetchanxxx-tech/hunfu-star-chain/master/.env.example" -o "${INSTALL_DIR}/.env" 2>/dev/null || \
         curl -fsSL "https://ghproxy.com/https://raw.githubusercontent.com/jetchanxxx-tech/hunfu-star-chain/master/.env.example" -o "${INSTALL_DIR}/.env" 2>/dev/null; then
        log "已下载 .env.example"
    else
        warn "无法获取 .env.example，手动创建默认配置"
        cat > "${INSTALL_DIR}/.env" << 'ENVEOF'
MYSQL_DSN=root:CHANGE_ME@tcp(127.0.0.1:3306)/huifu?charset=utf8mb4&parseTime=True&loc=Local
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
DEEPSEEK_API_KEY=
QWEN_API_KEY=
WX_MP_APPID=
WX_MP_SECRET=
DATA_DIR=/data/huifu/files
ENVEOF
    fi

    chown "$APP_USER":"$APP_USER" "${INSTALL_DIR}/.env"
    chmod 600 "${INSTALL_DIR}/.env"
    log ".env 已创建 → ${INSTALL_DIR}/.env"

    read -rp "是否现在编辑 .env？[y/N] " ans
    [ "$ans" = "y" ] || [ "$ans" = "Y" ] && ${EDITOR:-vim} "${INSTALL_DIR}/.env"
}

# ==================== 构建 ====================
build_backend() {
    info "构建 Go 后端..."
    cd "$SCRIPT_DIR"
    go build -ldflags="-s -w" -o bin/${APP} ./cmd/server
    log "后端编译完成 ($(du -sh bin/${APP} | cut -f1))"
}

build_admin() {
    info "构建管理后台..."
    cd "${SCRIPT_DIR}/admin"
    npm ci --silent 2>&1 | tail -3
    npm run build 2>&1 | tail -5
    cd "$SCRIPT_DIR"
    log "管理后台构建完成"
}

# ==================== 部署 ====================
deploy_files() {
    info "部署文件到 ${INSTALL_DIR}..."
    cp "${SCRIPT_DIR}/bin/${APP}" "${INSTALL_DIR}/${APP}"
    cp "${SCRIPT_DIR}/config.yaml" "${INSTALL_DIR}/"
    cp -r "${SCRIPT_DIR}/migrations" "${INSTALL_DIR}/"
    cp -r "${SCRIPT_DIR}/admin/dist" "${INSTALL_DIR}/admin-dist" 2>/dev/null || true
    chown -R "$APP_USER":"$APP_USER" "$INSTALL_DIR"
    chmod +x "${INSTALL_DIR}/${APP}"
    log "文件部署完成"
}

run_migration() {
    info "执行数据库迁移..."
    set -a; source "${INSTALL_DIR}/.env" 2>/dev/null || true; set +a
    if [ -z "${MYSQL_DSN:-}" ]; then
        warn "MYSQL_DSN 未配置，跳过迁移"
        return 0
    fi
    if migrate -path "${INSTALL_DIR}/migrations" -database "${MYSQL_DSN}" up 2>&1; then
        log "数据库迁移完成"
    else
        warn "数据库迁移失败，请检查 MYSQL_DSN 配置"
    fi
}

setup_systemd() {
    info "安装 systemd 服务..."
    cat > /etc/systemd/system/${APP}.service << SYSTEMD
[Unit]
Description=惠福星链 · 全病程健康协同平台
After=network.target mysql.service mariadb.service redis.service

[Service]
Type=simple
User=${APP_USER}
Group=${APP_USER}
WorkingDirectory=${INSTALL_DIR}
EnvironmentFile=${INSTALL_DIR}/.env
ExecStart=${INSTALL_DIR}/${APP}
Restart=on-failure
RestartSec=5s
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ProtectHome=yes
ReadWritePaths=${INSTALL_DIR} ${DATA_DIR} ${LOG_DIR}
MemoryLimit=2G
StandardOutput=journal
StandardError=journal
SyslogIdentifier=${APP}

[Install]
WantedBy=multi-user.target
SYSTEMD

    systemctl daemon-reload
    systemctl enable "${APP}"
    systemctl restart "${APP}"
    sleep 2
    if systemctl is-active --quiet "${APP}"; then
        log "服务启动成功"
    else
        warn "服务启动失败，查看日志: journalctl -u ${APP} -n 30"
    fi
}

# ==================== 健康检查 ====================
health_check() {
    info "健康检查..."
    for i in $(seq 1 15); do
        if curl -sf http://localhost:8080/api/health >/dev/null 2>&1; then
            log "✓ 服务运行正常"
            curl -s http://localhost:8080/api/health | python3 -m json.tool 2>/dev/null || curl -s http://localhost:8080/api/health
            return 0
        fi
        sleep 1
    done
    warn "健康检查超时，检查日志: journalctl -u ${APP} -n 50"
}

# ==================== 入口 ====================
main() {
    echo -e "${CYAN}"
    echo "  ╔══════════════════════════════════════╗"
    echo "  ║   惠福星链 · 全病程健康协同平台     ║"
    echo "  ║  一键安装部署 v1.0                  ║"
    echo "  ╚══════════════════════════════════════╝"
    echo -e "${NC}"

    if [ "$(id -u)" != "0" ]; then
        err "请用 root 运行: sudo ./install.sh"
    fi

    detect_os
    install_all
    setup_user_and_dirs
    setup_env
    build_backend
    build_admin
    deploy_files
    run_migration
    setup_systemd
    health_check

    section "安装完成!"
    echo ""
    echo "  后端 API:    http://$(hostname -I | awk '{print $1}'):8080/api/health"
    echo "  管理后台:    已编译到 ${INSTALL_DIR}/admin-dist/（需自行配置静态文件服务）"
    echo "  服务管理:    systemctl start|stop|restart|status ${APP}"
    echo "  查看日志:    journalctl -u ${APP} -f"
    echo "  运维手册:    ${INSTALL_DIR}/../docs/OPS_MANUAL.md"
    echo ""
    echo "  下一步: 编辑 ${INSTALL_DIR}/.env 填入微信/DeepSeek 密钥后重启服务"
    echo "          systemctl restart ${APP}"
    echo ""
}

main "$@"
