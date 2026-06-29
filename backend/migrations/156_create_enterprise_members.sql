-- 156_create_enterprise_members.sql
-- 企业成员表 —— 企业功能 P1 新增
-- 对应 Ent schema: ent/schema/enterprise_member.go
-- 本期约束：一人一企业（1:1）

CREATE TABLE IF NOT EXISTS enterprise_members (
    id              BIGSERIAL PRIMARY KEY,
    enterprise_id   BIGINT NOT NULL REFERENCES enterprises(id),
    user_id         BIGINT NOT NULL REFERENCES users(id),
    role            VARCHAR(20) NOT NULL DEFAULT 'enterprise_member',
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    department_id   BIGINT REFERENCES departments(id),
    concurrency     INT NOT NULL DEFAULT 0,
    rpm_limit       INT NOT NULL DEFAULT 0,
    notes           TEXT NOT NULL DEFAULT '',
    joined_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    unbound_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- 通用索引
CREATE INDEX IF NOT EXISTS idx_enterprise_members_user_id ON enterprise_members(user_id);
CREATE INDEX IF NOT EXISTS idx_enterprise_members_enterprise_id ON enterprise_members(enterprise_id);
CREATE INDEX IF NOT EXISTS idx_enterprise_members_role ON enterprise_members(role);
CREATE INDEX IF NOT EXISTS idx_enterprise_members_enterprise_id_role ON enterprise_members(enterprise_id, role);
CREATE INDEX IF NOT EXISTS idx_enterprise_members_department_id ON enterprise_members(department_id);
CREATE INDEX IF NOT EXISTS idx_enterprise_members_deleted_at ON enterprise_members(deleted_at);

-- 部分唯一索引：同一企业内同一用户只能有一条未删除记录
CREATE UNIQUE INDEX IF NOT EXISTS enterprise_members_ent_user_unique_active
    ON enterprise_members(enterprise_id, user_id)
    WHERE deleted_at IS NULL;

-- 部分唯一索引：本期 1:1 约束 —— 一个活跃用户只能在一个企业
CREATE UNIQUE INDEX IF NOT EXISTS enterprise_members_user_unique_active
    ON enterprise_members(user_id)
    WHERE status = 'active' AND deleted_at IS NULL;
