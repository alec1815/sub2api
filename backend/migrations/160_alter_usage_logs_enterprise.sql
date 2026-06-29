-- 160_alter_usage_logs_enterprise.sql
-- usage_logs 增加企业归属字段 —— 企业功能 P1 新增
-- 对应 Ent schema 变更: ent/schema/usage_log.go

-- 1. enterprise_id → enterprises.id（NULL=个人消费，有值=企业消费）
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS enterprise_id BIGINT;
CREATE INDEX IF NOT EXISTS idx_usage_logs_enterprise_id_created_at ON usage_logs(enterprise_id, created_at);

-- 2. pool_type —— 资金池类型（personal / enterprise），便于审计追溯
ALTER TABLE usage_logs ADD COLUMN IF NOT EXISTS pool_type VARCHAR(20) NOT NULL DEFAULT 'personal';
CREATE INDEX IF NOT EXISTS idx_usage_logs_enterprise_id_pool_type ON usage_logs(enterprise_id, pool_type);
