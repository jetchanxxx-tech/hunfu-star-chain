-- 003: 演示数据回滚
-- 清理所有种子数据（保留表结构）
DELETE FROM authorization_audit_log WHERE 1=1;
DELETE FROM member_authorizations WHERE 1=1;
DELETE FROM data_sync_logs WHERE 1=1;
DELETE FROM voice_call_logs WHERE 1=1;
DELETE FROM verification_records WHERE 1=1;
DELETE FROM followup_tasks WHERE 1=1;
DELETE FROM faq_entries WHERE category IN ('prenatal','postpartum','pediatrics','general','vaccine');
DELETE FROM ai_conversations WHERE 1=1;
DELETE FROM health_reports WHERE 1=1;
DELETE FROM timeline_events WHERE 1=1;
DELETE FROM user_entitlements WHERE 1=1;
DELETE FROM member_packages WHERE 1=1;
DELETE FROM wechat_bindings WHERE 1=1;
DELETE FROM family_members WHERE 1=1;
DELETE FROM families WHERE 1=1;
DELETE FROM service_packages WHERE package_uuid LIKE 'pkg-%';
DELETE FROM admin_users WHERE 1=1;
DROP TABLE IF EXISTS admin_users;