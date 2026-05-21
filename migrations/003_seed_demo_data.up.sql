-- ============================================================
-- 003: 演示数据种子 (50用户 + 完整业务模拟数据)
-- 惠福星链 · 全病程健康协同平台
-- ============================================================

-- ============================================================
-- 1. 管理员用户表 + 演示账户
-- ============================================================
CREATE TABLE IF NOT EXISTS admin_users (
    id          BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username    VARCHAR(50) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL COMMENT 'bcrypt hash',
    role        VARCHAR(30) NOT NULL COMMENT 'super_admin/steward/doctor/operator',
    real_name   VARCHAR(50),
    phone       VARCHAR(20),
    status      ENUM('active','disabled') DEFAULT 'active',
    last_login  DATETIME,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='管理后台用户';

-- 确保 002 迁移中的表存在（兼容未执行 002 的场景）
CREATE TABLE IF NOT EXISTS member_authorizations (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    grantor_id      BIGINT UNSIGNED NOT NULL,
    grantee_id      BIGINT UNSIGNED NOT NULL,
    auth_scope      JSON NOT NULL,
    status          ENUM('pending','active','rejected','revoked','expired') DEFAULT 'pending',
    request_msg     VARCHAR(500),
    reject_reason   VARCHAR(500),
    valid_until     DATE,
    requested_at    DATETIME DEFAULT CURRENT_TIMESTAMP,
    responded_at    DATETIME,
    revoked_at      DATETIME,
    FOREIGN KEY (grantor_id) REFERENCES family_members(id),
    FOREIGN KEY (grantee_id) REFERENCES family_members(id),
    INDEX idx_grantor (grantor_id, status),
    INDEX idx_grantee (grantee_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='家庭成员授权记录';

CREATE TABLE IF NOT EXISTS verification_records (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    entitlement_id  BIGINT UNSIGNED NOT NULL,
    member_id       BIGINT UNSIGNED NOT NULL,
    steward_id      BIGINT UNSIGNED,
    qr_nonce        VARCHAR(64) NOT NULL UNIQUE,
    benefit_type    VARCHAR(30) NOT NULL,
    verify_count    INT UNSIGNED DEFAULT 1,
    status          ENUM('success','failed','confirmed') DEFAULT 'success',
    fail_reason     VARCHAR(500),
    verified_at     DATETIME DEFAULT CURRENT_TIMESTAMP,
    confirmed_at    DATETIME,
    FOREIGN KEY (entitlement_id) REFERENCES user_entitlements(id),
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    INDEX idx_member_time (member_id, verified_at),
    INDEX idx_qr_nonce (qr_nonce)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='核销记录';

CREATE TABLE IF NOT EXISTS voice_call_logs (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    task_id         BIGINT UNSIGNED,
    member_id       BIGINT UNSIGNED NOT NULL,
    phone           VARCHAR(20) NOT NULL,
    call_type       VARCHAR(20) NOT NULL,
    provider        VARCHAR(30) NOT NULL,
    template_code   VARCHAR(50),
    dialogue_text   TEXT,
    call_status     VARCHAR(20) DEFAULT 'pending',
    duration_seconds INT DEFAULT 0,
    transcript      TEXT,
    intent_tags     JSON,
    retry_count     INT DEFAULT 0,
    fail_reason     VARCHAR(500),
    called_at       DATETIME,
    answered_at     DATETIME,
    finished_at     DATETIME,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (task_id) REFERENCES followup_tasks(id),
    FOREIGN KEY (member_id) REFERENCES family_members(id),
    INDEX idx_member_time (member_id, created_at),
    INDEX idx_status (call_status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='语音通话记录';

CREATE TABLE IF NOT EXISTS data_sync_logs (
    id              BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    source          VARCHAR(30) NOT NULL,
    sync_type       VARCHAR(20) NOT NULL,
    total_count     INT DEFAULT 0,
    success_count   INT DEFAULT 0,
    fail_count      INT DEFAULT 0,
    started_at      DATETIME,
    finished_at     DATETIME,
    created_at      DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_source_time (source, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='数据同步日志';

INSERT INTO admin_users (username, password, role, real_name, phone, status) VALUES
('admin',       '$2b$10$hc1RjGxjKrMSinSs3Y8GFObAz9uJijm5p/8izw0MTmEQUyGpjyqf.',  'super_admin', '系统管理员', '13800000001', 'active'),
('steward01',   '$2b$10$ukJv2hhYQLOQb8LN7dJ1AOdeNkztho4LD7T8XyFGUKi5xSuDaiP.i',  'steward',     '张管家',     '13800000002', 'active'),
('steward02',   '$2b$10$ukJv2hhYQLOQb8LN7dJ1AOdeNkztho4LD7T8XyFGUKi5xSuDaiP.i',  'steward',     '李管家',     '13800000003', 'active'),
('doctor01',    '$2b$10$MgzSw6tvYVq7bOC4CoUsf.hRaCvzpaA9ZU.wrqH4SPG6rRl2pbTbm',  'doctor',      '王医生',     '13800000004', 'active'),
('operator01',  '$2b$10$eMWhrkxUWpyusBiwgA.68uMicBr9GdTchBE7mdCjdi6iP7pyEv2t6',  'operator',    '赵运营',     '13800000005', 'active'),
('nurse01',     '$2b$10$vO/ikkn7ME9TCXgmC1BY9OUElcy9kD1ru6s.igs9/RRnlrTN4ErTu',  'steward',     '刘护士',     '13800000006', 'active');

-- ============================================================
-- 2. 服务包
-- ============================================================
INSERT INTO service_packages (package_uuid, name, description, level, price, benefits, status) VALUES
('pkg-0000-0000-0000-0001', 'VIP 孕期护航版',
 '孕前检查、NT、四维、糖耐量等全程产检咨询，含12次在线问诊、24次报告解读',
 'VIP', 2999.00,
 '[{"type":"consult","name":"在线问诊","count":12},{"type":"report_explain","name":"报告解读","count":24},{"type":"escort","name":"陪诊服务","count":1}]',
 'online'),

('pkg-0000-0000-0000-0002', 'VIP 产后康复版',
 '产后42天复查指导、盆底肌康复、母乳喂养咨询，含6次在线问诊、12次报告解读',
 'VIP', 3999.00,
 '[{"type":"consult","name":"在线问诊","count":6},{"type":"report_explain","name":"报告解读","count":12},{"type":"escort","name":"陪诊服务","count":2}]',
 'online'),

('pkg-0000-0000-0000-0003', 'VVIP 尊享全程版',
 '孕前-孕期-产后全程1对1管家服务，含36次在线问诊、72次报告解读、6次陪诊、AI智能提醒',
 'VVIP', 9999.00,
 '[{"type":"consult","name":"在线问诊","count":36},{"type":"report_explain","name":"报告解读","count":72},{"type":"escort","name":"陪诊服务","count":6}]',
 'online'),

('pkg-0000-0000-0000-0004', 'VVIP 全家守护版',
 '全家4人共享，含48次在线问诊、96次报告解读、12次陪诊、疫苗接种提醒、儿童发育评估',
 'VVIP', 19999.00,
 '[{"type":"consult","name":"在线问诊","count":48},{"type":"report_explain","name":"报告解读","count":96},{"type":"escort","name":"陪诊服务","count":12}]',
 'online'),

('pkg-0000-0000-0000-0005', 'VIP 儿科关爱版',
 '儿童体检、疫苗接种提醒、发育评估，含8次在线问诊、16次报告解读',
 'VIP', 1999.00,
 '[{"type":"consult","name":"在线问诊","count":8},{"type":"report_explain","name":"报告解读","count":16}]',
 'online'),

('pkg-0000-0000-0000-0006', 'VIP 慢病管理版',
 '高血压/糖尿病等慢病管理，含12次在线问诊、24次报告解读、2次陪诊',
 'VIP', 4999.00,
 '[{"type":"consult","name":"在线问诊","count":12},{"type":"report_explain","name":"报告解读","count":24},{"type":"escort","name":"陪诊服务","count":2}]',
 'online');

-- ============================================================
-- 3. 家庭 + 家庭成员 (16个家庭, 50个成员)
-- 角色分布: self=16, spouse=12, child=16, parent=4, other=2
-- 状态分布: active=42, inactive=8
-- ============================================================

-- F1: 张明家庭 (4人)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000001', '张明的家', 1);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(1,  'm0000001-aaaa-bbbb-cccc-000000000001', 1, 'self',   '大明', SHA2('张明',256), SHA2('13900010001',256), 1, '1990-03-15', 'HC001001', 'active'),
(2,  'm0000001-aaaa-bbbb-cccc-000000000002', 1, 'spouse', '娜娜', SHA2('李娜',256), SHA2('13900010002',256), 0, '1992-07-22', 'HC001002', 'active'),
(3,  'm0000001-aaaa-bbbb-cccc-000000000003', 1, 'child',  '小明', SHA2('张小明',256), SHA2('13900010003',256), 1, '2025-01-10', NULL, 'active'),
(4,  'm0000001-aaaa-bbbb-cccc-000000000004', 1, 'parent', '张爸', SHA2('张德福',256), SHA2('13900010004',256), 1, '1965-11-03', 'HC001004', 'active');

-- F2: 王芳家庭 (1人·单亲妈妈)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000002', '王芳的家', 5);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(5,  'm0000001-aaaa-bbbb-cccc-000000000005', 2, 'self', '芳芳', SHA2('王芳',256), SHA2('13900020001',256), 0, '1994-05-18', 'HC002001', 'active');

-- F3: 刘伟家庭 (4人)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000003', '刘伟的家', 6);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(6,  'm0000001-aaaa-bbbb-cccc-000000000006', 3, 'self',   '大伟', SHA2('刘伟',256), SHA2('13900030001',256), 1, '1988-09-12', 'HC003001', 'active'),
(7,  'm0000001-aaaa-bbbb-cccc-000000000007', 3, 'spouse', '静静', SHA2('陈静',256), SHA2('13900030002',256), 0, '1990-01-28', 'HC003002', 'active'),
(8,  'm0000001-aaaa-bbbb-cccc-000000000008', 3, 'child',  '小伟', SHA2('刘小伟',256), SHA2('13900030003',256), 1, '2023-08-05', NULL, 'active'),
(9,  'm0000001-aaaa-bbbb-cccc-000000000009', 3, 'child',  '小雨', SHA2('刘小雨',256), SHA2('13900030004',256), 0, '2025-02-14', NULL, 'active');

-- F4: 赵丽家庭 (2人·离异带娃)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000004', '赵丽的小窝', 10);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(10, 'm0000001-aaaa-bbbb-cccc-000000000010', 4, 'self',  '丽丽', SHA2('赵丽',256), SHA2('13900040001',256), 0, '1993-12-03', 'HC004001', 'active'),
(11, 'm0000001-aaaa-bbbb-cccc-000000000011', 4, 'child', '小宝', SHA2('赵小宝',256), SHA2('13900040002',256), 1, '2024-06-20', NULL, 'active');

-- F5: 孙涛家庭 (3人·新晋父母)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000005', '孙涛的家', 12);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(12, 'm0000001-aaaa-bbbb-cccc-000000000012', 5, 'self',   '涛哥', SHA2('孙涛',256), SHA2('13900050001',256), 1, '1991-06-08', 'HC005001', 'active'),
(13, 'm0000001-aaaa-bbbb-cccc-000000000013', 5, 'spouse', '婷婷', SHA2('周婷',256), SHA2('13900050002',256), 0, '1993-11-20', 'HC005002', 'active'),
(14, 'm0000001-aaaa-bbbb-cccc-000000000014', 5, 'parent', '孙妈', SHA2('钱秀兰',256), SHA2('13900050003',256), 0, '1963-04-25', 'HC005003', 'active');

-- F6: 周杰家庭 (4人)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000006', '周杰的家', 15);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(15, 'm0000001-aaaa-bbbb-cccc-000000000015', 6, 'self',   '杰哥', SHA2('周杰',256), SHA2('13900060001',256), 1, '1987-02-14', 'HC006001', 'active'),
(16, 'm0000001-aaaa-bbbb-cccc-000000000016', 6, 'spouse', '敏敏', SHA2('吴敏',256), SHA2('13900060002',256), 0, '1989-08-30', 'HC006002', 'active'),
(17, 'm0000001-aaaa-bbbb-cccc-000000000017', 6, 'child',  '小杰', SHA2('周小杰',256), SHA2('13900060003',256), 1, '2024-02-28', NULL, 'active'),
(18, 'm0000001-aaaa-bbbb-cccc-000000000018', 6, 'parent', '周爸', SHA2('周建国',256), SHA2('13900060004',256), 1, '1960-09-10', 'HC006004', 'inactive');

-- F7: 黄蓉家庭 (2人)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000007', '蓉儿的小家', 19);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(19, 'm0000001-aaaa-bbbb-cccc-000000000019', 7, 'self',  '蓉儿', SHA2('黄蓉',256), SHA2('13900070001',256), 0, '1995-08-15', 'HC007001', 'active'),
(20, 'm0000001-aaaa-bbbb-cccc-000000000020', 7, 'child', '子涵', SHA2('黄子涵',256), SHA2('13900070002',256), 0, '2024-11-02', NULL, 'active');

-- F8: 杨勇家庭 (5人·三代同堂)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000008', '杨家大家庭', 21);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(21, 'm0000001-aaaa-bbbb-cccc-000000000021', 8, 'self',   '勇哥', SHA2('杨勇',256), SHA2('13900080001',256), 1, '1986-04-02', 'HC008001', 'active'),
(22, 'm0000001-aaaa-bbbb-cccc-000000000022', 8, 'spouse', '丽丽', SHA2('徐丽',256), SHA2('13900080002',256), 0, '1988-12-17', 'HC008002', 'active'),
(23, 'm0000001-aaaa-bbbb-cccc-000000000023', 8, 'child',  '大勇', SHA2('杨小勇',256), SHA2('13900080003',256), 1, '2021-07-05', NULL, 'active'),
(24, 'm0000001-aaaa-bbbb-cccc-000000000024', 8, 'child',  '小丽', SHA2('杨小丽',256), SHA2('13900080004',256), 0, '2023-11-20', NULL, 'active'),
(25, 'm0000001-aaaa-bbbb-cccc-000000000025', 8, 'parent', '杨妈', SHA2('陈桂英',256), SHA2('13900080005',256), 0, '1962-03-08', 'HC008005', 'inactive');

-- F9: 马超家庭 (3人)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000009', '马超的家', 26);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(26, 'm0000001-aaaa-bbbb-cccc-000000000026', 9, 'self',   '超哥', SHA2('马超',256), SHA2('13900090001',256), 1, '1992-01-30', 'HC009001', 'active'),
(27, 'm0000001-aaaa-bbbb-cccc-000000000027', 9, 'spouse', '红红', SHA2('朱红',256), SHA2('13900090002',256), 0, '1994-06-12', 'HC009002', 'active'),
(28, 'm0000001-aaaa-bbbb-cccc-000000000028', 9, 'child',  '小红', SHA2('马小红',256), SHA2('13900090003',256), 0, '2025-03-01', NULL, 'active');

-- F10: 何平家庭 (1人·未婚)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000010', '何平的家', 29);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(29, 'm0000001-aaaa-bbbb-cccc-000000000029', 10, 'self', '阿平', SHA2('何平',256), SHA2('13900100001',256), 1, '1998-09-25', 'HC010001', 'active');

-- F11: 林晓家庭 (3人·二胎中)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000011', '林家小院', 30);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(30, 'm0000001-aaaa-bbbb-cccc-000000000030', 11, 'self',   '晓晓', SHA2('林晓',256), SHA2('13900110001',256), 0, '1993-03-21', 'HC011001', 'active'),
(31, 'm0000001-aaaa-bbbb-cccc-000000000031', 11, 'spouse', '宝玉', SHA2('贾宝玉',256), SHA2('13900110002',256), 1, '1991-10-07', 'HC011002', 'active'),
(32, 'm0000001-aaaa-bbbb-cccc-000000000032', 11, 'child',  '小玉', SHA2('贾小玉',256), SHA2('13900110003',256), 0, '2024-09-15', NULL, 'active');

-- F12: 陈强家庭 (4人)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000012', '陈强的家', 33);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(33, 'm0000001-aaaa-bbbb-cccc-000000000033', 12, 'self',   '强哥', SHA2('陈强',256), SHA2('13900120001',256), 1, '1985-07-19', 'HC012001', 'active'),
(34, 'm0000001-aaaa-bbbb-cccc-000000000034', 12, 'spouse', '王丽', SHA2('王丽',256), SHA2('13900120002',256), 0, '1987-05-28', 'HC012002', 'active'),
(35, 'm0000001-aaaa-bbbb-cccc-000000000035', 12, 'child',  '晨晨', SHA2('陈晨',256), SHA2('13900120003',256), 1, '2022-12-12', NULL, 'active'),
(36, 'm0000001-aaaa-bbbb-cccc-000000000036', 12, 'other',  '陈叔', SHA2('陈建国',256), SHA2('13900120004',256), 1, '1961-08-08', 'HC012004', 'active');

-- F13: 宋江家庭 (4人)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000013', '宋府', 37);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(37, 'm0000001-aaaa-bbbb-cccc-000000000037', 13, 'self',   '江哥', SHA2('宋江',256), SHA2('13900130001',256), 1, '1984-11-11', 'HC013001', 'active'),
(38, 'm0000001-aaaa-bbbb-cccc-000000000038', 13, 'spouse', '师师', SHA2('李师师',256), SHA2('13900130002',256), 0, '1986-02-14', 'HC013002', 'active'),
(39, 'm0000001-aaaa-bbbb-cccc-000000000039', 13, 'child',  '宋清', SHA2('宋清',256), SHA2('13900130003',256), 1, '2022-04-30', NULL, 'active'),
(40, 'm0000001-aaaa-bbbb-cccc-000000000040', 13, 'child',  '宋玉', SHA2('宋玉',256), SHA2('13900130004',256), 0, '2024-10-18', NULL, 'active');

-- F14: 曹操家庭 (4人)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000014', '曹家大院', 41);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(41, 'm0000001-aaaa-bbbb-cccc-000000000041', 14, 'self',   '孟德', SHA2('曹操',256), SHA2('13900140001',256), 1, '1983-05-06', 'HC014001', 'active'),
(42, 'm0000001-aaaa-bbbb-cccc-000000000042', 14, 'spouse', '丁夫人', SHA2('丁红',256), SHA2('13900140002',256), 0, '1985-09-15', 'HC014002', 'inactive'),
(43, 'm0000001-aaaa-bbbb-cccc-000000000043', 14, 'child',  '曹丕', SHA2('曹丕',256), SHA2('13900140003',256), 1, '2020-01-20', NULL, 'active'),
(44, 'm0000001-aaaa-bbbb-cccc-000000000044', 14, 'child',  '曹植', SHA2('曹植',256), SHA2('13900140004',256), 1, '2024-01-08', NULL, 'active');

-- F15: 刘备家庭 (3人)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000015', '刘皇叔家', 45);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(45, 'm0000001-aaaa-bbbb-cccc-000000000045', 15, 'self',   '玄德', SHA2('刘备',256), SHA2('13900150001',256), 1, '1982-12-31', 'HC015001', 'active'),
(46, 'm0000001-aaaa-bbbb-cccc-000000000046', 15, 'spouse', '尚香', SHA2('孙尚香',256), SHA2('13900150002',256), 0, '1984-04-08', 'HC015002', 'active'),
(47, 'm0000001-aaaa-bbbb-cccc-000000000047', 15, 'child',  '阿斗', SHA2('刘禅',256), SHA2('13900150003',256), 1, '2023-06-15', NULL, 'active');

-- F16: 吴用家庭 (3人·补充至50)
INSERT INTO families (family_uuid, name, primary_member_id) VALUES ('f0000001-aaaa-bbbb-cccc-000000000016', '吴用的家', 48);
INSERT INTO family_members (id, member_uuid, family_id, relation, nickname, real_name_hash, phone_hash, gender, birth_date, health_card, status) VALUES
(48, 'm0000001-aaaa-bbbb-cccc-000000000048', 16, 'self',   '智多星', SHA2('吴用',256), SHA2('13900160001',256), 1, '1990-10-01', 'HC016001', 'active'),
(49, 'm0000001-aaaa-bbbb-cccc-000000000049', 16, 'spouse', '公孙嫂', SHA2('公孙兰',256), SHA2('13900160002',256), 0, '1992-08-16', 'HC016002', 'active'),
(50, 'm0000001-aaaa-bbbb-cccc-000000000050', 16, 'child',  '小宝', SHA2('吴小宝',256), SHA2('13900160003',256), 1, '2025-04-22', NULL, 'active');

-- ============================================================
-- 4. 微信绑定 (部分用户)
-- ============================================================
INSERT INTO wechat_bindings (member_id, union_id, openid_mp) VALUES
(1,  'wx-union-0000000001', 'wx-openid-mp-0000000001'),
(2,  'wx-union-0000000002', 'wx-openid-mp-0000000002'),
(5,  'wx-union-0000000005', 'wx-openid-mp-0000000005'),
(6,  'wx-union-0000000006', 'wx-openid-mp-0000000006'),
(7,  'wx-union-0000000007', 'wx-openid-mp-0000000007'),
(10, 'wx-union-0000000010', 'wx-openid-mp-0000000010'),
(12, 'wx-union-0000000012', 'wx-openid-mp-0000000012'),
(13, 'wx-union-0000000013', 'wx-openid-mp-0000000013'),
(15, 'wx-union-0000000015', 'wx-openid-mp-0000000015'),
(19, 'wx-union-0000000019', 'wx-openid-mp-0000000019'),
(21, 'wx-union-0000000021', 'wx-openid-mp-0000000021'),
(26, 'wx-union-0000000026', 'wx-openid-mp-0000000026'),
(29, 'wx-union-0000000029', 'wx-openid-mp-0000000029'),
(30, 'wx-union-0000000030', 'wx-openid-mp-0000000030'),
(33, 'wx-union-0000000033', 'wx-openid-mp-0000000033'),
(37, 'wx-union-0000000037', 'wx-openid-mp-0000000037'),
(41, 'wx-union-0000000041', 'wx-openid-mp-0000000041'),
(45, 'wx-union-0000000045', 'wx-openid-mp-0000000045'),
(48, 'wx-union-0000000048', 'wx-openid-mp-0000000048');

-- ============================================================
-- 5. 会员购买记录 (多状态: active/expired/refunded)
-- ============================================================
INSERT INTO member_packages (member_id, package_id, order_no, amount, start_date, end_date, status) VALUES
-- 活跃购买
(1,  1, 'ORD202501001', 2999.00,  '2025-01-15', '2026-01-15', 'active'),
(2,  1, 'ORD202501002', 2999.00,  '2025-01-15', '2026-01-15', 'active'),
(5,  2, 'ORD202502001', 3999.00,  '2025-02-01', '2026-02-01', 'active'),
(6,  3, 'ORD202503001', 9999.00,  '2025-03-01', '2026-03-01', 'active'),
(7,  3, 'ORD202503002', 9999.00,  '2025-03-01', '2026-03-01', 'active'),
(10, 4, 'ORD202504001', 19999.00, '2025-04-01', '2026-04-01', 'active'),
(12, 1, 'ORD202504002', 2999.00,  '2025-04-10', '2026-04-10', 'active'),
(13, 1, 'ORD202504003', 2999.00,  '2025-04-10', '2026-04-10', 'active'),
(15, 3, 'ORD202505001', 9999.00,  '2025-05-01', '2026-05-01', 'active'),
(16, 3, 'ORD202505002', 9999.00,  '2025-05-01', '2026-05-01', 'active'),
(19, 5, 'ORD202505003', 1999.00,  '2025-05-15', '2026-05-15', 'active'),
(21, 4, 'ORD202506001', 19999.00, '2025-06-01', '2026-06-01', 'active'),
(22, 4, 'ORD202506002', 19999.00, '2025-06-01', '2026-06-01', 'active'),
(26, 1, 'ORD202507001', 2999.00,  '2025-07-01', '2026-07-01', 'active'),
(27, 1, 'ORD202507002', 2999.00,  '2025-07-01', '2026-07-01', 'active'),
(30, 3, 'ORD202508001', 9999.00,  '2025-08-01', '2026-08-01', 'active'),
(31, 3, 'ORD202508002', 9999.00,  '2025-08-01', '2026-08-01', 'active'),
(33, 6, 'ORD202509001', 4999.00,  '2025-09-01', '2026-09-01', 'active'),
(34, 6, 'ORD202509002', 4999.00,  '2025-09-01', '2026-09-01', 'active'),
(37, 4, 'ORD202510001', 19999.00, '2025-10-01', '2026-10-01', 'active'),
(38, 4, 'ORD202510002', 19999.00, '2025-10-01', '2026-10-01', 'active'),
(41, 4, 'ORD202511001', 19999.00, '2025-11-01', '2026-11-01', 'active'),
(45, 1, 'ORD202512001', 2999.00,  '2025-12-01', '2026-12-01', 'active'),
(46, 1, 'ORD202512002', 2999.00,  '2025-12-01', '2026-12-01', 'active'),
(48, 2, 'ORD202601001', 3999.00,  '2026-01-01', '2027-01-01', 'active'),
(49, 2, 'ORD202601002', 3999.00,  '2026-01-01', '2027-01-01', 'active'),
-- 已过期
(3,  5, 'ORD202401001', 1999.00,  '2024-01-10', '2025-01-10', 'expired'),
(8,  5, 'ORD202403001', 1999.00,  '2024-03-01', '2025-03-01', 'expired'),
(11, 5, 'ORD202406001', 1999.00,  '2024-06-20', '2025-06-20', 'expired'),
(14, 1, 'ORD202501003', 2999.00,  '2025-01-01', '2025-02-01', 'expired'),
(17, 5, 'ORD202402001', 1999.00,  '2024-02-28', '2025-02-28', 'expired'),
(20, 5, 'ORD202411001', 1999.00,  '2024-11-02', '2025-11-02', 'expired'),
(23, 5, 'ORD202407001', 1999.00,  '2024-07-05', '2025-07-05', 'expired'),
(25, 6, 'ORD202503003', 4999.00,  '2025-03-08', '2025-09-08', 'expired'),
(32, 5, 'ORD202409001', 1999.00,  '2024-09-15', '2025-09-15', 'expired'),
(35, 5, 'ORD202412001', 1999.00,  '2024-12-12', '2025-12-12', 'expired'),
(39, 5, 'ORD202404001', 1999.00,  '2024-04-30', '2025-04-30', 'expired'),
-- 已退款
(29, 2, 'ORD202509003', 3999.00,  '2025-09-25', '2026-09-25', 'refunded'),
(36, 1, 'ORD202508003', 2999.00,  '2025-08-08', '2026-08-08', 'refunded'),
(42, 6, 'ORD202511002', 4999.00,  '2025-11-15', '2026-11-15', 'refunded');

-- ============================================================
-- 6. 用户权益 (基于活跃购买包)
-- ============================================================
INSERT INTO user_entitlements (member_id, package_id, benefit_type, total, consumed, valid_until) VALUES
-- 张明 (VIP孕期护航版)
(1, 1, 'consult',        12,  3, '2026-01-15'),
(1, 1, 'report_explain', 24,  8, '2026-01-15'),
(1, 1, 'escort',          1,  0, '2026-01-15'),
-- 李娜
(2, 1, 'consult',        12,  5, '2026-01-15'),
(2, 1, 'report_explain', 24, 12, '2026-01-15'),
(2, 1, 'escort',          1,  1, '2026-01-15'),
-- 王芳 (VIP产后康复版)
(5, 2, 'consult',         6,  2, '2026-02-01'),
(5, 2, 'report_explain', 12,  4, '2026-02-01'),
(5, 2, 'escort',          2,  1, '2026-02-01'),
-- 刘伟 (VVIP尊享全程版)
(6, 3, 'consult',        36, 10, '2026-03-01'),
(6, 3, 'report_explain', 72, 20, '2026-03-01'),
(6, 3, 'escort',          6,  1, '2026-03-01'),
-- 陈静
(7, 3, 'consult',        36,  8, '2026-03-01'),
(7, 3, 'report_explain', 72, 15, '2026-03-01'),
(7, 3, 'escort',          6,  0, '2026-03-01'),
-- 赵丽 (VVIP全家守护版)
(10, 4, 'consult',        48,  6, '2026-04-01'),
(10, 4, 'report_explain', 96, 18, '2026-04-01'),
(10, 4, 'escort',         12,  2, '2026-04-01'),
-- 孙涛 (VIP孕期护航版)
(12, 1, 'consult',        12,  4, '2026-04-10'),
(12, 1, 'report_explain', 24,  6, '2026-04-10'),
(12, 1, 'escort',          1,  0, '2026-04-10'),
-- 周杰 (VVIP尊享全程版)
(15, 3, 'consult',        36, 15, '2026-05-01'),
(15, 3, 'report_explain', 72, 30, '2026-05-01'),
(15, 3, 'escort',          6,  2, '2026-05-01'),
-- 黄蓉 (VIP儿科关爱版) - 子涵用
(19, 5, 'consult',         8,  3, '2026-05-15'),
(19, 5, 'report_explain', 16,  5, '2026-05-15'),
-- 杨勇 (VVIP全家守护版)
(21, 4, 'consult',        48, 20, '2026-06-01'),
(21, 4, 'report_explain', 96, 40, '2026-06-01'),
(21, 4, 'escort',         12,  3, '2026-06-01'),
-- 林晓 (VVIP尊享全程版)
(30, 3, 'consult',        36,  5, '2026-08-01'),
(30, 3, 'report_explain', 72, 10, '2026-08-01'),
(30, 3, 'escort',          6,  1, '2026-08-01'),
-- 陈强 (VIP慢病管理版)
(33, 6, 'consult',        12,  5, '2026-09-01'),
(33, 6, 'report_explain', 24, 10, '2026-09-01'),
(33, 6, 'escort',          2,  1, '2026-09-01'),
-- 宋江 (VVIP全家守护版)
(37, 4, 'consult',        48, 12, '2026-10-01'),
(37, 4, 'report_explain', 96, 25, '2026-10-01'),
(37, 4, 'escort',         12,  1, '2026-10-01'),
-- 曹操 (VVIP全家守护版)
(41, 4, 'consult',        48, 18, '2026-11-01'),
(41, 4, 'report_explain', 96, 35, '2026-11-01'),
(41, 4, 'escort',         12,  4, '2026-11-01'),
-- 刘备 (VIP孕期护航版)
(45, 1, 'consult',        12,  3, '2026-12-01'),
(45, 1, 'report_explain', 24,  7, '2026-12-01'),
(45, 1, 'escort',          1,  0, '2026-12-01'),
-- 吴用 (VIP产后康复版)
(48, 2, 'consult',         6,  1, '2027-01-01'),
(48, 2, 'report_explain', 12,  3, '2027-01-01'),
(48, 2, 'escort',          2,  0, '2027-01-01');

-- ============================================================
-- 7. 时间轴事件 (模拟真实孕期/产后/育儿场景)
-- ============================================================
INSERT INTO timeline_events (member_id, event_type, event_date, event_data, source) VALUES
-- 李娜(2) - 孕期完整时间轴
(2, 'first_prenatal', '2024-12-01', '{"gestation_week":8,"notes":"首次产检确认宫内妊娠"}', 'auto'),
(2, 'nt',             '2025-01-05', '{"gestation_week":12,"nt_value":"1.2mm","result":"正常"}', 'auto'),
(2, 'early_tang',     '2025-01-28', '{"gestation_week":16,"result":"低风险"}', 'auto'),
(2, 'quad_d',         '2025-03-10', '{"gestation_week":22,"result":"未见异常"}', 'auto'),
(2, 'ogtt',           '2025-03-28', '{"gestation_week":25,"glucose_0h":4.8,"glucose_1h":9.2,"glucose_2h":7.8,"result":"正常"}', 'auto'),
(2, 'delivery',       '2025-07-10', '{"delivery_type":"顺产","gestation_week":39,"baby_weight":3400,"note":"母子平安"}', 'manual'),
(2, '42day',          '2025-08-21', '{"notes":"产妇恢复良好，新生儿体重4.5kg"}', 'auto'),
-- 小明(3) - 新生儿疫苗
(3, 'vaccine_2m',     '2025-03-12', '{"vaccine":"脊灰疫苗第1剂","hospital":"市妇幼"}', 'manual'),
(3, 'vaccine_3m',     '2025-04-12', '{"vaccine":"百白破第1剂","hospital":"市妇幼"}', 'manual'),
-- 陈静(7) - 孕期
(7, 'first_prenatal', '2024-09-15', '{"gestation_week":7,"notes":"确认妊娠"}', 'auto'),
(7, 'nt',             '2024-10-20', '{"gestation_week":12,"nt_value":"1.5mm","result":"正常"}', 'auto'),
(7, 'early_tang',     '2024-11-15', '{"gestation_week":16,"result":"低风险"}', 'auto'),
(7, 'quad_d',         '2025-01-05', '{"gestation_week":22,"result":"未见异常"}', 'auto'),
(7, 'ogtt',           '2025-02-14', '{"gestation_week":25,"glucose_0h":5.1,"glucose_1h":10.3,"glucose_2h":8.5,"result":"异常-妊娠期糖尿病"}', 'auto'),
(7, 'delivery',       '2025-06-02', '{"delivery_type":"剖宫产","gestation_week":38,"baby_weight":3100,"note":"因GDM行剖宫产"}', 'manual'),
(7, '42day',          '2025-07-14', '{"notes":"血糖恢复正常，新生儿体重3.8kg"}', 'auto'),
-- 周婷(13) - 孕期
(13, 'first_prenatal', '2025-01-08', '{"gestation_week":6,"notes":"首次产检"}', 'auto'),
(13, 'nt',             '2025-02-15', '{"gestation_week":12,"nt_value":"1.0mm","result":"正常"}', 'auto'),
(13, 'early_tang',     '2025-03-20', '{"gestation_week":16,"result":"低风险"}', 'auto'),
(13, 'quad_d',         '2025-05-01', '{"gestation_week":22,"result":"未见异常"}', 'auto'),
-- 吴敏(16) - 已完成
(16, 'first_prenatal', '2023-10-10', '{"gestation_week":8}', 'auto'),
(16, 'nt',             '2023-11-15', '{"gestation_week":12,"nt_value":"1.3mm"}', 'auto'),
(16, 'delivery',       '2024-05-28', '{"delivery_type":"顺产","gestation_week":39}', 'manual'),
(16, '42day',          '2024-07-09', '{"notes":"产后恢复良好"}', 'auto'),
-- 徐丽(22) - 二胎
(22, 'first_prenatal', '2022-06-01', '{"gestation_week":8}', 'auto'),
(22, 'nt',             '2022-07-05', '{"gestation_week":12}', 'auto'),
(22, 'delivery',       '2023-02-14', '{"delivery_type":"顺产","gestation_week":40}', 'manual'),
(22, '42day',          '2023-03-28', '{"notes":"恢复良好"}', 'auto'),
-- 再孕
(22, 'first_prenatal', '2024-11-01', '{"gestation_week":7}', 'auto'),
(22, 'nt',             '2024-12-10', '{"gestation_week":12,"nt_value":"0.9mm"}', 'auto'),
(22, 'quad_d',         '2025-02-20', '{"gestation_week":22,"result":"未见异常"}', 'auto'),
-- 朱红(27) - 孕晚期
(27, 'first_prenatal', '2024-10-01', '{"gestation_week":7,"notes":"确认妊娠"}', 'auto'),
(27, 'nt',             '2024-11-08', '{"gestation_week":12,"nt_value":"1.4mm"}', 'auto'),
(27, 'early_tang',     '2024-12-10', '{"gestation_week":16,"result":"临界风险"}', 'auto'),
(27, 'quad_d',         '2025-01-15', '{"gestation_week":22,"result":"未见异常"}', 'auto'),
(27, 'ogtt',           '2025-03-01', '{"gestation_week":26,"result":"正常"}', 'auto'),
-- 孙尚香(46) - 备孕中
(46, 'first_prenatal', '2025-05-01', '{"gestation_week":6,"notes":"确认妊娠"}', 'auto'),
(46, 'nt',             '2025-06-08', '{"gestation_week":12,"nt_value":"1.1mm","result":"正常"}', 'auto'),
-- 儿童疫苗
(8,  'vaccine_2m',     '2023-10-05', '{"vaccine":"脊灰第1剂"}', 'manual'),
(11, 'vaccine_2m',     '2024-08-22', '{"vaccine":"脊灰第1剂"}', 'manual'),
(11, 'vaccine_3m',     '2024-09-22', '{"vaccine":"百白破第1剂"}', 'manual'),
(17, 'vaccine_2m',     '2024-05-01', '{"vaccine":"脊灰第1剂"}', 'manual'),
(17, 'vaccine_3m',     '2024-06-01', '{"vaccine":"百白破第1剂"}', 'manual'),
(20, 'vaccine_2m',     '2025-01-05', '{"vaccine":"脊灰第1剂"}', 'manual'),
(23, 'vaccine_2m',     '2021-09-05', '{"vaccine":"脊灰第1剂"}', 'manual'),
(24, 'vaccine_2m',     '2024-01-22', '{"vaccine":"脊灰第1剂"}', 'manual'),
(24, 'vaccine_3m',     '2024-02-22', '{"vaccine":"百白破第1剂"}', 'manual'),
(28, 'vaccine_2m',     '2025-05-01', '{"vaccine":"脊灰第1剂"}', 'manual'),
(32, 'vaccine_2m',     '2024-11-16', '{"vaccine":"脊灰第1剂"}', 'manual'),
(35, 'vaccine_2m',     '2023-02-14', '{"vaccine":"脊灰第1剂"}', 'manual'),
(39, 'vaccine_2m',     '2022-07-01', '{"vaccine":"脊灰第1剂"}', 'manual'),
(43, 'vaccine_2m',     '2020-03-22', '{"vaccine":"脊灰第1剂"}', 'manual'),
(44, 'vaccine_2m',     '2024-03-10', '{"vaccine":"脊灰第1剂"}', 'manual'),
(47, 'vaccine_2m',     '2023-08-17', '{"vaccine":"脊灰第1剂"}', 'manual'),
(47, 'vaccine_3m',     '2023-09-17', '{"vaccine":"百白破第1剂"}', 'manual'),
(50, 'vaccine_2m',     '2025-06-23', '{"vaccine":"脊灰第1剂"}', 'manual');

-- ============================================================
-- 8. 健康报告 (检验/检查/出院小结)
-- ============================================================
INSERT INTO health_reports (member_id, report_type, hospital_code, report_no, summary, abnormal_flags, report_date, source) VALUES
-- 李娜 - 孕期检查
(2, 'lab',     'HOSP001', 'RPT202501001', '{"hgb":125,"glucose":5.1,"alt":18,"creatinine":55}', NULL, '2025-01-05', 'lis'),
(2, 'imaging', 'HOSP001', 'RPT202503001', '{"procedure":"NT超声","result":"NT 1.2mm，正常范围","crl":55}', NULL, '2025-01-05', 'pacs'),
(2, 'imaging', 'HOSP001', 'RPT202503002', '{"procedure":"四维彩超","result":"胎儿结构未见异常","bpd":56,"fl":45}', NULL, '2025-03-10', 'pacs'),
(2, 'lab',     'HOSP001', 'RPT202503003', '{"glucose_0h":4.8,"glucose_1h":9.2,"glucose_2h":7.8,"result":"OGTT正常"}', NULL, '2025-03-28', 'lis'),
(2, 'discharge','HOSP001', 'RPT202507001', '{"diagnosis":"足月顺产","procedure":"阴道分娩","length_of_stay":3}', NULL, '2025-07-13', 'his'),
-- 陈静 - 妊娠期糖尿病
(7, 'lab',     'HOSP001', 'RPT202410001', '{"hgb":118,"glucose":5.3,"alt":22}', '["glucose偏高"]', '2024-10-20', 'lis'),
(7, 'lab',     'HOSP001', 'RPT202502001', '{"glucose_0h":5.1,"glucose_1h":10.3,"glucose_2h":8.5}', '["糖耐量异常","妊娠期糖尿病"]', '2025-02-14', 'lis'),
(7, 'discharge','HOSP001', 'RPT202506001', '{"diagnosis":"GDM+剖宫产","procedure":"子宫下段剖宫产术","length_of_stay":5}', '["妊娠期糖尿病"]', '2025-06-07', 'his'),
-- 张德福(4) - 慢病
(4, 'lab',     'HOSP002', 'RPT202503001', '{"hba1c":7.2,"glucose_fasting":6.8,"tc":5.6,"ldl":3.8}', '["糖化血红蛋白偏高","空腹血糖偏高","LDL偏高"]', '2025-03-15', 'lis'),
(4, 'lab',     'HOSP002', 'RPT202505001', '{"hba1c":6.9,"glucose_fasting":6.2,"tc":5.3}', '["糖化血红蛋白偏高"]', '2025-05-20', 'lis'),
-- 周婷 - 孕期常规
(13, 'lab',    'HOSP003', 'RPT202502001', '{"hgb":130,"glucose":4.6,"alt":15}', NULL, '2025-02-15', 'lis'),
(13, 'imaging','HOSP003', 'RPT202505001', '{"procedure":"四维彩超","result":"胎儿结构未见异常"}', NULL, '2025-05-01', 'pacs'),
-- 周建国(18) - 老年体检
(18, 'lab',    'HOSP002', 'RPT202503004', '{"hba1c":6.5,"creatinine":98,"alt":35}', '["糖化血红蛋白临界","ALT偏高"]', '2025-03-10', 'lis'),
(18, 'imaging','HOSP002', 'RPT202503005', '{"procedure":"腹部B超","result":"轻度脂肪肝"}', '["轻度脂肪肝"]', '2025-03-10', 'pacs'),
-- 杨妈(25) - 老年慢病
(25, 'lab',    'HOSP002', 'RPT202506003', '{"hba1c":8.1,"glucose_fasting":8.5,"creatinine":110}', '["糖化血红蛋白显著偏高","空腹血糖偏高"]', '2025-06-15', 'lis'),
-- 其他成员
(1,  'lab',    'HOSP001', 'RPT202501003', '{"hgb":155,"alt":28,"creatinine":72}', NULL, '2025-01-20', 'lis'),
(5,  'lab',    'HOSP003', 'RPT202502002', '{"hgb":120,"tsh":2.5,"ft4":15}', NULL, '2025-02-05', 'lis'),
(16, 'lab',    'HOSP001', 'RPT202505003', '{"hgb":125,"hct":0.38}', NULL, '2025-05-05', 'lis'),
(22, 'imaging','HOSP003', 'RPT202502002', '{"procedure":"B超","result":"宫内早孕"}', NULL, '2025-02-20', 'pacs'),
(33, 'lab',    'HOSP002', 'RPT202509003', '{"tc":6.2,"ldl":4.1,"glucose":5.8}', '["TC偏高","LDL偏高"]', '2025-09-05', 'lis'),
(34, 'lab',    'HOSP002', 'RPT202509004', '{"hgb":130,"alt":20}', NULL, '2025-09-05', 'lis'),
(42, 'lab',    'HOSP002', 'RPT202511003', '{"hba1c":5.8,"glucose":5.2}', NULL, '2025-11-20', 'lis');

-- ============================================================
-- 9. AI 对话记录
-- ============================================================
INSERT INTO ai_conversations (member_id, session_id, provider, role, content, tokens_in, tokens_out, emotion_tag, created_at) VALUES
-- 李娜 - 孕期焦虑咨询
(2, 'sess-0001-aaaa-000000000001', 'deepseek', 'user',      '我最近感觉胎动减少了，是不是有问题？现在孕28周', 0, 0, 'anxious', '2025-03-25 09:15:00'),
(2, 'sess-0001-aaaa-000000000001', 'deepseek', 'assistant', '胎动减少需要关注。孕28周胎儿每天应有10次以上胎动，建议您选择一个固定时间监测。如果2小时内胎动不足10次，请及时到院检查。', 45, 120, NULL, '2025-03-25 09:15:05'),
-- 陈静 - GDM 咨询
(7, 'sess-0002-bbbb-000000000002', 'deepseek', 'user',      '血糖高了，GDM饮食应该怎么安排？', 0, 0, 'anxious', '2025-02-15 14:30:00'),
(7, 'sess-0002-bbbb-000000000002', 'deepseek', 'assistant', 'GDM饮食建议：1) 少量多餐，每天5-6餐；2) 主食粗细搭配，控制每餐碳水；3) 先吃菜后吃主食；4) 每日监测血糖4次。具体食谱可以咨询营养科。', 38, 150, NULL, '2025-02-15 14:30:08'),
(7, 'sess-0003-cccc-000000000003', 'qwen',     'user',      '宝宝出生后需要做哪些筛查？', 0, 0, 'normal', '2025-06-03 10:00:00'),
(7, 'sess-0003-cccc-000000000003', 'qwen',     'assistant', '新生儿筛查包括：1) 先天性甲减；2) 苯丙酮尿症；3) G6PD缺乏症；4) 听力筛查；5) 视力筛查。一般在出生72小时后采足跟血。', 40, 130, NULL, '2025-06-03 10:00:06'),
-- 周婷 - 孕期不适
(13, 'sess-0004-dddd-000000000004', 'deepseek', 'user',      '孕晚期腰部酸痛，有什么缓解办法吗？', 0, 0, 'normal', '2025-04-15 20:45:00'),
(13, 'sess-0004-dddd-000000000004', 'deepseek', 'assistant', '可以尝试：1) 左侧卧位休息；2) 孕妇瑜伽/游泳；3) 热敷（避开腹部）；4) 穿低跟鞋。如果疼痛剧烈或伴随其他症状请及时就医。', 30, 110, NULL, '2025-04-15 20:45:05'),
-- 王芳 - 产后咨询
(5, 'sess-0005-eeee-000000000005', 'qwen',     'user',      '产后多久可以开始运动？我是顺产', 0, 0, 'normal', '2025-08-25 16:20:00'),
(5, 'sess-0005-eeee-000000000005', 'qwen',     'assistant', '顺产后一般6周可以开始轻度运动，如散步、盆底肌训练。8-12周可逐步增加强度。建议从凯格尔运动开始。剖宫产需要更长时间，建议咨询医生。', 35, 125, NULL, '2025-08-25 16:20:07'),
-- 刘备(45) - 常规咨询
(45, 'sess-0006-ffff-000000000006', 'deepseek', 'user',      '备孕期需要补充什么营养？', 0, 0, 'normal', '2025-04-01 11:00:00'),
(45, 'sess-0006-ffff-000000000006', 'deepseek', 'assistant', '备孕期建议补充：1) 叶酸0.4mg/天（孕前3个月起）；2) 铁剂（如有贫血）；3) 碘；4) 维生素D；5) 均衡饮食。双方戒烟戒酒，适当运动。', 40, 140, NULL, '2025-04-01 11:00:08'),
-- 黄蓉 - 宝宝护理
(19, 'sess-0007-1111-000000000007', 'deepseek', 'user',      '宝宝红屁股反复发作怎么办？', 0, 0, 'anxious', '2025-06-10 08:35:00'),
(19, 'sess-0007-1111-000000000007', 'deepseek', 'assistant', '红屁屁预防和护理：1) 勤换尿布，2-3小时一次；2) 温水清洗，轻轻蘸干；3) 涂护臀膏形成保护层；4) 适当光屁屁通风；5) 如糜烂渗液需就医，可能真菌感染。', 35, 150, NULL, '2025-06-10 08:35:06'),
-- 曹操 - 体检咨询
(41, 'sess-0008-2222-000000000008', 'qwen',     'user',      '体检报告显示轻度脂肪肝，严重吗？怎么改善？', 0, 0, 'normal', '2025-06-20 13:00:00'),
(41, 'sess-0008-2222-000000000008', 'qwen',     'assistant', '轻度脂肪肝是可逆的。建议：1) 控制饮食，减少高脂高糖食物；2) 每周至少150分钟有氧运动；3) 减重5-10%；4) 戒酒；5) 3-6个月后复查肝功能+B超。', 42, 150, NULL, '2025-06-20 13:00:07'),
-- 急诊场景 - 出血咨询（触发紧急关键词）
(2, 'sess-0009-3333-000000000009', 'deepseek', 'user',      '孕35周突然发现有一点出血，肚子不疼', 0, 0, 'anxious', '2025-06-01 03:15:00'),
(2, 'sess-0009-3333-000000000009', 'deepseek', 'assistant', '⚠️ 您提到了"出血"症状。孕晚期出血需要紧急医学评估，可能的原因包括前置胎盘、胎盘早剥等。请立即前往最近医院急诊科就诊，不要等待！', 42, 100, NULL, '2025-06-01 03:15:05'),
-- 投诉场景
(21, 'sess-0010-4444-000000000010', 'deepseek', 'user',      '你们的预约系统总是约不上，太差了', 0, 0, 'complaint', '2025-06-05 10:00:00'),
(21, 'sess-0010-4444-000000000010', 'deepseek', 'assistant', '非常抱歉给您带来不便。我已经记录您的反馈。目前高峰期预约较紧张，建议您可以尝试早晨7点放号时段预约，或联系您的专属管家协助预约。', 20, 80, NULL, '2025-06-05 10:00:05');

-- ============================================================
-- 10. FAQ 知识库 (增强已有数据)
-- ============================================================
INSERT INTO faq_entries (category, question, answer, keywords, priority, status) VALUES
('prenatal',    'NT检查是什么？什么时候做？', 'NT（颈项透明层）检查是孕早期超声筛查，测量胎儿颈后透明层厚度。最佳检查时间为孕11-13周+6天。NT值<2.5mm为正常范围。', 'NT,颈项透明层,早筛', 10, 'published'),
('prenatal',    '糖耐量检查要空腹吗？', '需要空腹8-12小时。检查当天先抽空腹血糖，然后喝75g葡萄糖水，1小时和2小时后分别再抽血。', '糖耐量,OGTT,空腹', 8, 'published'),
('prenatal',    '孕期可以同房吗？', '一般情况下孕期可以同房，但有以下情况需避免：1) 有流产史；2) 前置胎盘；3) 出血/破水；4) 孕晚期最后4周。建议咨询医生。', '同房,性生活,注意事项', 5, 'published'),
('postpartum',  '产后多久可以洗澡？', '顺产产后1-2天即可淋浴，剖宫产产后3-5天拆线后可淋浴（保持伤口干燥）。水温不宜过高，时间5-10分钟为宜。', '洗澡,淋浴,产褥期', 7, 'published'),
('postpartum',  '母乳不足怎么办？', '增加乳汁分泌的方法：1) 多吸吮，按需哺乳；2) 保证充足睡眠和水分；3) 喝催乳汤（鲫鱼汤、猪蹄汤）；4) 保持心情愉快；5) 可咨询催乳师。', '母乳,奶水不足,催乳', 9, 'published'),
('postpartum',  '产后42天复查查什么？', '检查项目包括：1) 妇科检查（子宫恢复、伤口愈合）；2) 盆底肌评估；3) 乳腺检查；4) 血常规；5) 心理评估。新生儿同时需做生长发育评估。', '42天,复查,产后检查', 8, 'published'),
('pediatrics',  '宝宝发烧怎么办？', '38.5℃以下可物理降温（温水擦浴、减少衣物）；38.5℃以上可使用退烧药（对乙酰氨基酚/布洛芬）。3个月以下婴儿发烧需立即就医。', '发烧,发热,退烧', 10, 'published'),
('pediatrics',  '疫苗接种后发热正常吗？', '部分疫苗（如百白破、麻腮风）接种后24-48小时内可能出现低热，一般38.5℃以下持续1-2天。多喝水、物理降温即可。如果高热不退或超过3天需就医。', '疫苗,发热,反应', 8, 'published'),
('pediatrics',  '宝宝多大可以添加辅食？', '一般6个月开始添加。信号：能独坐、对食物感兴趣、挺舌反射消失。第一口推荐高铁米粉，由稀到稠，单一品种逐步添加。', '辅食,添加,月龄', 9, 'published'),
('general',     '如何绑定家庭成员？', '在家庭页面点击"添加成员"，填写昵称、关系、手机号。被添加人会收到授权请求，授权通过后即可查看相关健康数据。', '绑定,添加,家庭成员', 7, 'published'),
('general',     '如何预约陪诊服务？', '在服务包页面选择已购买的陪诊权益，选择就诊医院、日期、科室。提交后管家将在30分钟内与您联系确认。', '陪诊,预约,服务', 7, 'published'),
('general',     '报告解读需要多长时间？', '常规报告提交后24小时内回复，标注为"紧急"的报告2小时内回复。疑难报告会转接专科医生，48小时内回复。', '报告解读,时间,回复', 6, 'published'),
('vaccine',     '什么是疫苗接种程序？', '国家规定的一类疫苗免费接种。常见：出生-B肝第1剂，2月龄-脊灰第1剂，3月龄-百白破第1剂。二类疫苗自愿自费接种。', '疫苗,接种,程序,计划', 8, 'published'),
('vaccine',     '自费疫苗值得打吗？', '建议接种的自费疫苗：13价肺炎球菌疫苗、轮状病毒疫苗、水痘疫苗、手足口疫苗、流感疫苗。可根据宝宝体质和经济条件选择。', '自费,疫苗,二类,推荐', 7, 'published');

-- ============================================================
-- 11. 随访任务 (模拟运行中的状态)
-- ============================================================
INSERT INTO followup_tasks (member_id, trigger_type, trigger_value, title, status, assigned_to, due_date, completed_at, notes) VALUES
-- 待处理
(2,  'gestation_week', '28', '孕28周常规随访', 'pending', 2, '2025-03-25', NULL, '常规产检提醒'),
(7,  'lab_abnormal',   'ogtt', 'GDM血糖随访-陈静', 'in_progress', 2, '2025-02-20', NULL, 'OGTT异常，需饮食指导'),
(13, 'gestation_week', '24', '孕24周糖耐量检查提醒', 'pending', 3, '2025-05-10', NULL, NULL),
(4,  'lab_abnormal',   'hba1c', '糖尿病管理随访-张德福', 'in_progress', 3, '2025-03-20', NULL, '血糖控制欠佳，需调整方案'),
(18, 'no_show',        '体检', '未按时复查提醒-周建国', 'pending', 2, '2025-04-01', NULL, '3月体检未到，电话联系'),
(25, 'lab_abnormal',   'hba1c', '血糖异常随访-陈桂英', 'in_progress', 3, '2025-06-20', NULL, 'HbA1c 8.1% 需重点管理'),
(33, 'lab_abnormal',   'ldl', '血脂异常随访-陈强', 'pending', 3, '2025-09-10', NULL, 'LDL 4.1 偏高'),
(46, 'gestation_week', '13', '孕13周NT检查提醒-孙尚香', 'pending', 2, '2025-06-15', NULL, NULL),
(48, 'gestation_week', '8', '早孕确认随访-公孙兰', 'pending', 2, '2025-08-20', NULL, NULL),
-- 已完成
(2,  'gestation_week', '12', '孕12周NT检查提醒-李娜', 'completed', 2, '2025-01-01', '2025-01-06 15:00:00', '已完成NT检查，结果正常'),
(6,  'gestation_week', '8', '首次产检建档提醒-陈静', 'completed', 2, '2024-09-10', '2024-09-16 10:30:00', '已建档'),
(7,  'gestation_week', '7', '产后42天复查提醒-陈静', 'completed', 2, '2025-07-10', '2025-07-15 11:00:00', '复查完成，血糖正常'),
(13, 'gestation_week', '12', '孕12周NT检查提醒-周婷', 'completed', 3, '2025-02-10', '2025-02-16 14:00:00', 'NT 1.0mm，正常'),
(16, 'gestation_week', '12', '孕12周NT检查提醒-吴敏', 'completed', 2, '2023-11-10', '2023-11-16 11:00:00', 'NT结果正常'),
(16, 'age_month',      '2', '2月龄疫苗提醒-周小杰', 'completed', 3, '2024-04-28', '2024-04-29 09:00:00', '已完成脊灰第1剂'),
-- 已取消
(11, 'no_show',        '疫苗', '未按时接种提醒-赵小宝', 'cancelled', 3, '2024-08-20', NULL, '联系后已自行去社区接种'),
(20, 'no_show',        '体检', '未按时体检提醒-黄子涵', 'cancelled', 2, '2025-02-15', NULL, '已在外院完成体检');

-- ============================================================
-- 12. 核销记录
-- ============================================================
INSERT INTO verification_records (entitlement_id, member_id, steward_id, qr_nonce, benefit_type, verify_count, status, verified_at, confirmed_at) VALUES
-- 已确认核销
(1,  1,  2, 'qr-nonce-20250120001', 'consult',        1, 'confirmed', '2025-01-20 10:00:00', '2025-01-20 10:05:00'),
(1,  1,  2, 'qr-nonce-20250215001', 'consult',        1, 'confirmed', '2025-02-15 14:00:00', '2025-02-15 14:03:00'),
(2,  1,  3, 'qr-nonce-20250125001', 'report_explain', 1, 'confirmed', '2025-01-25 09:00:00', '2025-01-25 09:02:00'),
(5,  2,  2, 'qr-nonce-20250120002', 'report_explain', 1, 'confirmed', '2025-01-20 11:00:00', '2025-01-20 11:03:00'),
(5,  2,  2, 'qr-nonce-20250310001', 'report_explain', 1, 'confirmed', '2025-03-10 16:00:00', '2025-03-10 16:02:00'),
(8,  5,  3, 'qr-nonce-20250205001', 'escort',         1, 'confirmed', '2025-02-05 08:00:00', '2025-02-05 17:00:00'),
(10, 6,  2, 'qr-nonce-20250315001', 'consult',        1, 'confirmed', '2025-03-15 10:00:00', '2025-03-15 10:05:00'),
(37, 21, 2, 'qr-nonce-20250610001', 'consult',        1, 'confirmed', '2025-06-10 14:00:00', '2025-06-10 14:03:00'),
(38, 21, 2, 'qr-nonce-20250612001', 'report_explain', 1, 'confirmed', '2025-06-12 15:00:00', '2025-06-12 15:02:00'),
(44, 33, 3, 'qr-nonce-20250905001', 'escort',         1, 'confirmed', '2025-09-05 08:30:00', '2025-09-05 12:00:00'),
-- 仅成功未确认
(16, 12, 3, 'qr-nonce-20250420001', 'consult',        1, 'success', '2025-04-20 10:00:00', NULL),
(2,  1,  2, 'qr-nonce-20250501001', 'report_explain', 1, 'success', '2025-05-01 09:00:00', NULL),
-- 失败核销
(3,  1,  2, 'qr-nonce-20250401001', 'escort',         1, 'failed', '2025-04-01 10:00:00', NULL);

-- ============================================================
-- 13. 语音通话记录
-- ============================================================
INSERT INTO voice_call_logs (task_id, member_id, phone, call_type, provider, template_code, dialogue_text, call_status, duration_seconds, transcript, intent_tags, called_at, answered_at, finished_at) VALUES
(5,  18, '13900060004', 'auto', 'aliyun', 'tpl-no-show-01',
 '您好，这里是惠福星链健康管理中心。您有随访任务未完成，请尽快完成体检复查。',
 'completed', 45, '用户表示已收到提醒，将于下周前往医院复查。',
 '["接受提醒","计划复查"]', '2025-04-01 10:00:00', '2025-04-01 10:00:05', '2025-04-01 10:00:50'),

(3,  7,  '13900030002', 'manual', 'aliyun', 'tpl-gdm-01',
 '您好陈女士，我是您的专属健康管家。关于您的血糖管理，我来为您做一下饮食指导。',
 'completed', 180, '管家详细介绍了GDM饮食方案，陈女士表示理解并会记录血糖。约定3天后回访。',
 '["接受指导","配合管理","约定回访"]', '2025-02-20 15:00:00', '2025-02-20 15:00:03', '2025-02-20 15:03:03'),

(6,  25, '13900080005', 'auto', 'aliyun', 'tpl-dm-01',
 '您好，您的血糖检测结果偏高，请及时回院复查。需要帮助请按1，转人工请按2。',
 'no_answer', 0, NULL, NULL,
 '["无人接听"]', '2025-06-20 10:00:00', NULL, '2025-06-20 10:00:30'),

(4,  4,  '13900010004', 'auto', 'aliyun', 'tpl-dm-recheck',
 '您好张先生，您的血糖仍需进一步控制。请按时服药并来院复查。',
 'completed', 30, '用户确认已收到提醒，表示会按时复查。',
 '["接受提醒","计划复查"]', '2025-03-22 10:00:00', '2025-03-22 10:00:05', '2025-03-22 10:00:35'),

(NULL, 7, '13900030002', 'manual', 'aliyun', 'tpl-postpartum',
 '陈女士您好，产后恢复得怎么样？宝宝吃奶正常吗？我来做一下42天复查的预约。',
 'completed', 120, '确认产后恢复良好，宝宝纯母乳喂养。已预约7月15日复查。',
 '["产后随访","预约复查","满意度高"]', '2025-07-10 15:00:00', '2025-07-10 15:00:03', '2025-07-10 15:02:03'),

(6,  25, '13900080005', 'auto', 'aliyun', 'tpl-dm-01',
 '您好，血糖复查提醒第二次。请尽快安排时间回院复查。',
 'rejected', 3, '用户直接挂断。',
 '["拒接"]', '2025-06-21 10:00:00', '2025-06-21 10:00:05', '2025-06-21 10:00:08'),

(NULL, 13, '13900050002', 'manual', 'aliyun', 'tpl-ogtt-remind',
 '周女士您好，即将进入糖耐量检查周期，我先为您介绍一下注意事项。',
 'completed', 150, '详细讲解了OGTT检查流程和注意事项，用户表示明白。',
 '["接受指导","满意度高"]', '2025-05-08 11:00:00', '2025-05-08 11:00:02', '2025-05-08 11:02:32'),

(1, 2, '13900010002', 'auto', 'aliyun', 'tpl-gestation-28',
 '李女士您好，孕28周常规随访提醒。请按时产检并关注胎动。',
 'completed', 20, '已收到提醒。',
 '["接受提醒"]', '2025-03-25 09:00:00', '2025-03-25 09:00:03', '2025-03-25 09:00:23');

-- ============================================================
-- 14. 数据同步日志 (模拟历史同步)
-- ============================================================
INSERT INTO data_sync_logs (source, sync_type, total_count, success_count, fail_count, started_at, finished_at) VALUES
('his',  'full',         500, 498, 2, '2025-01-01 02:00:00', '2025-01-01 02:15:00'),
('his',  'incremental',   15,  15, 0, '2025-01-02 02:00:00', '2025-01-02 02:02:00'),
('lis',  'full',         300, 300, 0, '2025-01-01 02:30:00', '2025-01-01 02:40:00'),
('lis',  'incremental',   20,  18, 2, '2025-01-02 02:30:00', '2025-01-02 02:32:00'),
('pacs', 'full',         100, 100, 0, '2025-01-01 03:00:00', '2025-01-01 03:05:00');

-- ============================================================
-- 15. 授权记录 (示例)
-- ============================================================
INSERT INTO member_authorizations (grantor_id, grantee_id, auth_scope, status, request_msg, valid_until, requested_at, responded_at) VALUES
(1,  2,  '["report","timeline","reminder"]', 'active',   '家庭共享授权', '2026-01-15', '2025-01-15 10:00:00', '2025-01-15 10:30:00'),
(6,  7,  '["report","timeline","reminder"]', 'active',   '孕期数据共享', '2026-03-01', '2025-03-01 09:00:00', '2025-03-01 09:10:00'),
(12, 13, '["report","timeline","reminder"]', 'active',   '孕期数据共享', '2026-04-10', '2025-04-10 08:00:00', '2025-04-10 08:15:00'),
(21, 22, '["report","timeline","reminder"]', 'active',   '孩子管理共享', '2026-06-01', '2025-06-01 10:00:00', '2025-06-01 10:05:00'),
(33, 36, '["report"]',                        'pending',  '请帮忙查看体检报告', NULL, '2025-09-05 14:00:00', NULL),
(41, 42, '["report","timeline"]',             'revoked',  '临时共享', NULL, '2025-11-10 11:00:00', '2025-12-01 08:00:00');

-- ============================================================
-- END OF 003 SEED DATA
-- ============================================================