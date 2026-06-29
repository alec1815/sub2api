-- 154_create_enterprises.sql
-- 企业表 —— 企业功能 P1 新增
-- 对应 Ent schema: ent/schema/enterprise.go

CREATE TABLE IF NOT EXISTS enterprises (
    id              BIGSERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    short_name      VARCHAR(100) NOT NULL DEFAULT '',
    credit_code     VARCHAR(50) NOT NULL DEFAULT '',
    address         VARCHAR(500) NOT NULL DEFAULT '',
    scale           VARCHAR(20) NOT NULL DEFAULT '',
    industry        VARCHAR(50) NOT NULL DEFAULT '',
    parent_id       BIGINT NOT NULL DEFAULT 0,
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    contact_name    VARCHAR(100) NOT NULL DEFAULT '',
    contact_phone   VARCHAR(50) NOT NULL DEFAULT '',
    contact_email   VARCHAR(255) NOT NULL DEFAULT '',
    notes           TEXT NOT NULL DEFAULT '',
    balance         DECIMAL(20,8) NOT NULL DEFAULT 0,
    total_recharged DECIMAL(20,8) NOT NULL DEFAULT 0,
    admin_user_id   BIGINT NOT NULL REFERENCES users(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- 通用索引
CREATE INDEX IF NOT EXISTS idx_enterprises_parent_id ON enterprises(parent_id);
CREATE INDEX IF NOT EXISTS idx_enterprises_status ON enterprises(status);
CREATE INDEX IF NOT EXISTS idx_enterprises_admin_user_id ON enterprises(admin_user_id);
CREATE INDEX IF NOT EXISTS idx_enterprises_deleted_at ON enterprises(deleted_at);

-- 部分唯一索引：同一未删除企业名称唯一
-- 参见 016_soft_delete_partial_unique_indexes.sql 模式
CREATE UNIQUE INDEX IF NOT EXISTS enterprises_name_unique_active
    ON enterprises(name)
    WHERE deleted_at IS NULL;
