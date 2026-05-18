-- 001: 初始核心表结构
-- 惠福星链 · 全病程健康协同平台

-- 家庭
CREATE TABLE families (
    id                BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    family_uuid       CHAR(36) NOT NULL UNIQUE COMMENT '对外 UUID',
    name              VARCHAR(100) COMMENT '家庭名称',
    primary_member_id BIGINT UNSIGNED NOT NULL COMMENT '主账号成员 ID',
    created_at        DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at        DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_primary (primary_member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭';

-- 家庭成员
CREATE TABLE family_members (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_uuid     CHAR(36) NOT NULL UNIQUE COMMENT '对外 UUID',
    family_id       BIGINT UNSIGNED NOT NULL,
    relation        ENUM('self','spouse','child','parent','other') NOT NULL DEFAULT 'self',
    nickname        VARCHAR(100) COMMENT '昵称',
    real_name_hash  VARCHAR(64) NOT NULL COMMENT 'SHA256(姓名) 脱敏',
    phone_hash      VARCHAR(64) NOT NULL COMMENT 'SHA256(手机号) 脱敏',
    encrypted_name  TEXT COMMENT 'AES-256 加密完整姓名',
    encrypted_phone TEXT COMMENT 'AES-256 加密完整手机号',
    gender          TINYINT COMMENT '0:女 1:男',
    birth_date      DATE,
    health_card     VARCHAR(50) COMMENT '就诊卡号',
    avatar_url      VARCHAR(500),
    status          ENUM('active','inactive') DEFAULT 'active',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (family_id) REFERENCES families(id),
    INDEX idx_family (family_id),
    INDEX idx_phone_hash (phone_hash)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭成员';

-- 微信用户绑定
CREATE TABLE wechat_bindings (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id   BIGINT UNSIGNED NOT NULL,
    union_id    VARCHAR(100) COMMENT '微信开放平台 union_id',
    openid_mp   VARCHAR(100) COMMENT '小程序 openid',
    openid_wework VARCHAR(100) COMMENT '企业微信 openid',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    UNIQUE INDEX idx_openid_mp (openid_mp),
    UNIQUE INDEX idx_union_id (union_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='微信用户绑定';

-- 时间轴事件
CREATE TABLE timeline_events (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id   BIGINT UNSIGNED NOT NULL,
    event_type  VARCHAR(30) NOT NULL COMMENT 'first_prenatal/nt/ogtt/delivery/42day/vaccine',
    event_date  DATE NOT NULL,
    event_data  JSON COMMENT '灵活扩展: 报告ID/异常标记/备注',
    source      ENUM('auto','manual','ai') DEFAULT 'auto',
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    INDEX idx_member_date (member_id, event_date),
    INDEX idx_type (event_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='时间轴事件';

-- 检验/检查报告（云端脱敏存储）
CREATE TABLE health_reports (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id       BIGINT UNSIGNED NOT NULL,
    report_type     VARCHAR(20) NOT NULL COMMENT 'lab(检验)/imaging(检查)/discharge(出院小结)',
    hospital_code   VARCHAR(20) COMMENT '医院编码',
    report_no       VARCHAR(100) COMMENT '医院报告号',
    summary         JSON COMMENT '报告摘要(脱敏指标)',
    abnormal_flags  JSON COMMENT '异常指标列表',
    report_date     DATE,
    file_url        VARCHAR(500) COMMENT '原始报告文件(PDF)',
    source          ENUM('his','lis','pacs','manual') DEFAULT 'manual',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    INDEX idx_member_type (member_id, report_type, report_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='健康报告';

-- 服务包定义
CREATE TABLE service_packages (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    package_uuid    CHAR(36) NOT NULL UNIQUE,
    name            VARCHAR(200) NOT NULL,
    description     TEXT,
    level           ENUM('VIP','VVIP') NOT NULL DEFAULT 'VIP',
    price           DECIMAL(10,2) NOT NULL,
    cover_image     VARCHAR(500),
    benefits        JSON COMMENT '权益列表: [{type, name, count}]',
    status          ENUM('draft','online','offline') DEFAULT 'draft',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='服务包';

-- 会员购买记录
CREATE TABLE member_packages (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id       BIGINT UNSIGNED NOT NULL,
    package_id      BIGINT UNSIGNED NOT NULL,
    order_no        VARCHAR(50) NOT NULL UNIQUE,
    amount          DECIMAL(10,2) NOT NULL,
    start_date      DATE NOT NULL,
    end_date        DATE NOT NULL,
    status          ENUM('active','expired','refunded') DEFAULT 'active',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    FOREIGN KEY (package_id) REFERENCES service_packages(id),
    INDEX idx_member_status (member_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='会员购买记录';

-- 用户权益（原子扣减模式）
CREATE TABLE user_entitlements (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id       BIGINT UNSIGNED NOT NULL,
    package_id      BIGINT UNSIGNED NOT NULL,
    benefit_type    VARCHAR(30) NOT NULL COMMENT 'escort/consult/report_explain',
    total           INT UNSIGNED NOT NULL,
    consumed        INT UNSIGNED DEFAULT 0,
    valid_until     DATE,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    INDEX idx_member_package (member_id, package_id),
    CONSTRAINT chk_consumed CHECK (consumed <= total)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户权益';

-- AI 对话记录
CREATE TABLE ai_conversations (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id       BIGINT UNSIGNED NOT NULL,
    session_id      CHAR(36) NOT NULL,
    provider        VARCHAR(20) NOT NULL COMMENT 'deepseek/qwen',
    role            ENUM('user','assistant','system') NOT NULL,
    content         TEXT NOT NULL,
    tokens_in       INT UNSIGNED DEFAULT 0,
    tokens_out      INT UNSIGNED DEFAULT 0,
    emotion_tag     VARCHAR(20) COMMENT 'anxious/normal/complaint',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_session (session_id, created_at),
    INDEX idx_member (member_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='AI对话记录';

-- FAQ 知识库
CREATE TABLE faq_entries (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    category        VARCHAR(50) NOT NULL COMMENT 'prenatal/postpartum/pediatrics/general',
    question        VARCHAR(500) NOT NULL,
    answer          TEXT NOT NULL,
    keywords        VARCHAR(500) COMMENT '逗号分隔关键词用于匹配',
    priority        INT DEFAULT 0 COMMENT '排序权重',
    status          ENUM('published','hidden') DEFAULT 'published',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FULLTEXT INDEX ft_question (question),
    INDEX idx_category (category, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='FAQ知识库';

-- 随访任务
CREATE TABLE followup_tasks (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    member_id       BIGINT UNSIGNED NOT NULL,
    trigger_type    VARCHAR(30) NOT NULL COMMENT 'gestation_week/age_month/lab_abnormal/no_show',
    trigger_value   VARCHAR(100) COMMENT '触发条件值',
    title           VARCHAR(200) NOT NULL,
    status          ENUM('pending','in_progress','completed','cancelled') DEFAULT 'pending',
    assigned_to     BIGINT UNSIGNED COMMENT '分配的管家(医疗/服务)',
    due_date        DATE,
    completed_at    DATETIME,
    notes           TEXT,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    INDEX idx_status_due (status, due_date),
    INDEX idx_member (member_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='随访任务';
