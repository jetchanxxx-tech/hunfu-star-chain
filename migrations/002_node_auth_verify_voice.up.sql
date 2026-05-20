-- 002: 节点时间配置 + 家庭成员授权 + 核销记录 + 死信队列 + 语音通话记录
-- 惠福星链 · 全病程健康协同平台

-- 时间轴节点模板(默认配置)
CREATE TABLE timeline_node_templates (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    node_code       VARCHAR(30) NOT NULL UNIQUE COMMENT '节点编码: nt/ogtt/delivery/42day/vaccine_2m/vaccine_3m',
    node_name       VARCHAR(100) NOT NULL COMMENT '节点名称',
    category        VARCHAR(30) NOT NULL COMMENT 'prenatal/postpartum/pediatrics/vaccine',
    default_start   VARCHAR(50) COMMENT '默认开始: 如 孕11周',
    default_end     VARCHAR(50) COMMENT '默认结束: 如 孕13周+6天',
    reminder_days   INT DEFAULT 7 COMMENT '提前提醒天数',
    sort_order      INT DEFAULT 0,
    description     VARCHAR(500) COMMENT '备注说明',
    status          ENUM('enabled','disabled') DEFAULT 'enabled',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_category (category, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='时间轴节点模板';

-- 医院级节点时间覆盖
CREATE TABLE hospital_node_overrides (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    hospital_code   VARCHAR(20) NOT NULL,
    node_code       VARCHAR(30) NOT NULL,
    start_offset    VARCHAR(50) COMMENT '覆盖开始时间',
    end_offset      VARCHAR(50) COMMENT '覆盖结束时间',
    reminder_days   INT COMMENT '覆盖提醒天数',
    is_enabled      TINYINT(1) DEFAULT 1,
    description     VARCHAR(500),
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_hospital_node (hospital_code, node_code),
    INDEX idx_hospital (hospital_code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='医院级节点时间覆盖';

-- 家庭成员授权记录
CREATE TABLE member_authorizations (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    grantor_id      BIGINT UNSIGNED NOT NULL COMMENT '授权人(主账号)',
    grantee_id      BIGINT UNSIGNED NOT NULL COMMENT '被授权人',
    auth_scope      JSON NOT NULL COMMENT '授权范围: ["report","timeline","reminder"]',
    status          ENUM('pending','active','rejected','revoked','expired') DEFAULT 'pending',
    request_msg     VARCHAR(500) COMMENT '请求说明',
    reject_reason   VARCHAR(500),
    valid_until     DATE COMMENT '授权有效期',
    requested_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    responded_at    DATETIME,
    revoked_at      DATETIME,
    FOREIGN KEY (grantor_id) REFERENCES family_members(id),
    FOREIGN KEY (grantee_id) REFERENCES family_members(id),
    INDEX idx_grantor (grantor_id, status),
    INDEX idx_grantee (grantee_id, status),
    CONSTRAINT chk_not_self CHECK (grantor_id <> grantee_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭成员授权记录';

-- 授权操作审计日志
CREATE TABLE authorization_audit_log (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    authorization_id BIGINT UNSIGNED,
    action          VARCHAR(20) NOT NULL COMMENT 'request/approve/reject/revoke/expire',
    actor_id        BIGINT UNSIGNED NOT NULL COMMENT '操作人',
    target_id       BIGINT UNSIGNED NOT NULL COMMENT '目标成员',
    detail          JSON COMMENT '操作详情',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_auth (authorization_id),
    INDEX idx_actor (actor_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='授权操作审计日志';

-- 核销记录
CREATE TABLE verification_records (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    entitlement_id  BIGINT UNSIGNED NOT NULL,
    member_id       BIGINT UNSIGNED NOT NULL,
    steward_id      BIGINT UNSIGNED COMMENT '核销管家ID',
    qr_nonce        VARCHAR(64) NOT NULL UNIQUE COMMENT '二维码一次性标识',
    benefit_type    VARCHAR(30) NOT NULL,
    verify_count    INT UNSIGNED DEFAULT 1 COMMENT '本次核销次数',
    status          ENUM('success','failed','confirmed') DEFAULT 'success',
    fail_reason     VARCHAR(500),
    verified_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    confirmed_at    DATETIME,
    FOREIGN KEY (entitlement_id) REFERENCES user_entitlements(id),
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    INDEX idx_member_time (member_id, verified_at),
    INDEX idx_steward_time (steward_id, verified_at),
    INDEX idx_qr_nonce (qr_nonce)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='核销记录';

-- 死信队列(数据对接失败记录)
CREATE TABLE dead_letter_queue (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    source          VARCHAR(30) NOT NULL COMMENT '数据源标识',
    raw_data        LONGTEXT NOT NULL COMMENT '原始数据JSON',
    error_reason    VARCHAR(500) NOT NULL,
    retry_count     INT DEFAULT 0,
    status          ENUM('pending','retrying','resolved','discarded') DEFAULT 'pending',
    resolved_at     DATETIME,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_source_status (source, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='死信队列';

-- 数据同步日志
CREATE TABLE data_sync_logs (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    source          VARCHAR(30) NOT NULL COMMENT '数据源标识',
    sync_type       VARCHAR(20) NOT NULL COMMENT 'full/incremental',
    total_count     INT DEFAULT 0,
    success_count   INT DEFAULT 0,
    fail_count      INT DEFAULT 0,
    started_at      DATETIME,
    finished_at     DATETIME,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_source_time (source, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据同步日志';

-- 语音通话记录
CREATE TABLE voice_call_logs (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    task_id         BIGINT UNSIGNED COMMENT '关联随访任务',
    member_id       BIGINT UNSIGNED NOT NULL,
    phone           VARCHAR(20) NOT NULL,
    call_type       VARCHAR(20) NOT NULL COMMENT 'auto/manual',
    provider        VARCHAR(30) NOT NULL COMMENT '语音供应商',
    template_code   VARCHAR(50) COMMENT '语音模板编码',
    dialogue_text   TEXT COMMENT 'AI生成的话术文本',
    call_status     VARCHAR(20) DEFAULT 'pending' COMMENT 'pending/ringing/in_progress/completed/failed/rejected/no_answer',
    duration_seconds INT DEFAULT 0,
    transcript      TEXT COMMENT '语音转文字记录',
    intent_tags     JSON COMMENT 'AI分析的意向标签',
    retry_count     INT DEFAULT 0,
    fail_reason     VARCHAR(500),
    called_at       DATETIME,
    answered_at     DATETIME,
    finished_at     DATETIME,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES followup_tasks(id),
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    INDEX idx_member_time (member_id, created_at),
    INDEX idx_task (task_id),
    INDEX idx_status (call_status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='语音通话记录';

-- 语音模板配置
CREATE TABLE voice_templates (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    template_code   VARCHAR(50) NOT NULL UNIQUE,
    template_name   VARCHAR(200) NOT NULL,
    category        VARCHAR(30) NOT NULL COMMENT 'auto/manual',
    provider        VARCHAR(30) NOT NULL,
    provider_tpl_id VARCHAR(100) COMMENT '供应商侧模板ID',
    llm_prompt      TEXT COMMENT 'AI话术生成prompt',
    max_retries     INT DEFAULT 3,
    retry_interval  INT DEFAULT 30 COMMENT '重试间隔(分钟)',
    daily_limit     INT DEFAULT 2 COMMENT '单用户日呼叫上限',
    status          ENUM('draft','reviewing','approved','rejected') DEFAULT 'draft',
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='语音模板配置';

-- 插入默认节点模板
INSERT INTO timeline_node_templates (node_code, node_name, category, default_start, default_end, reminder_days, sort_order, description) VALUES
('nt',          'NT检查',    'prenatal',    '孕11周',   '孕13周+6天', 7, 1, 'NT值正常范围 <2.5mm'),
('early_tang',  '早唐筛查',   'prenatal',    '孕15周',   '孕16周',     7, 2, '唐氏综合征血清学筛查'),
('ogtt',        '糖耐量',    'prenatal',    '孕24周',   '孕28周',     7, 3, '口服葡萄糖耐量试验'),
('quad_d',      '四维彩超',   'prenatal',    '孕22周',   '孕26周',     3, 4, '大排畸检查'),
('delivery',    '分娩',      'postpartum',  '预产期前几天', '预产期',   7, 5, '预产期基于末次月经/超声推算'),
('42day',       '产后42天复查','postpartum', '产后42天',  '产后42天',   7, 6, '产妇及新生儿复查'),
('vaccine_2m',  '2月龄疫苗',  'vaccine',     '出生后2月龄','出生后2月龄',7, 7, '脊灰疫苗第1剂'),
('vaccine_3m',  '3月龄疫苗',  'vaccine',     '出生后3月龄','出生后3月龄',7, 8, '百白破疫苗第1剂');
