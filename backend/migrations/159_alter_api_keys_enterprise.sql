-- 159_alter_api_keys_enterprise.sql
-- api_keys 增加企业分配字段 —— 企业功能 P1 新增
-- 对应 Ent schema 变更: ent/schema/api_key.go

-- 保留 group_id 1:1 字段（废弃但保留），不删除
-- 新逻辑中分组关系从 api_key_groups 中间表读取

-- 1. assigned_to → enterprise_members.id（NULL=个人Key/管理员自用Key，有值=企业分配Key）
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS assigned_to BIGINT;
CREATE INDEX IF NOT EXISTS idx_api_keys_assigned_to ON api_keys(assigned_to);
CREATE INDEX IF NOT EXISTS idx_api_keys_assigned_to_status ON api_keys(assigned_to, status);

-- 2. usage_purpose —— 用途说明
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS usage_purpose VARCHAR(200) NOT NULL DEFAULT '';

-- 3. bound_tool —— 绑定工具
ALTER TABLE api_keys ADD COLUMN IF NOT EXISTS bound_tool VARCHAR(50) NOT NULL DEFAULT '';
