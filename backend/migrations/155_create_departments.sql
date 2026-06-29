-- 155_create_departments.sql
-- 部门表 —— 企业功能 P1 新增
-- 对应 Ent schema: ent/schema/department.go

CREATE TABLE IF NOT EXISTS departments (
    id              BIGSERIAL PRIMARY KEY,
    enterprise_id   BIGINT NOT NULL REFERENCES enterprises(id),
    parent_id       BIGINT NOT NULL DEFAULT 0,
    name            VARCHAR(100) NOT NULL,
    order_num       INT NOT NULL DEFAULT 0,
    leader          VARCHAR(100) NOT NULL DEFAULT '',
    phone           VARCHAR(50) NOT NULL DEFAULT '',
    email           VARCHAR(255) NOT NULL DEFAULT '',
    status          VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at      TIMESTAMPTZ
);

-- 查子部门
CREATE INDEX IF NOT EXISTS idx_departments_ent_id_parent_id ON departments(enterprise_id, parent_id);
CREATE INDEX IF NOT EXISTS idx_departments_deleted_at ON departments(deleted_at);

-- 部分唯一索引：同一企业下名称唯一（已删除除外）
CREATE UNIQUE INDEX IF NOT EXISTS departments_ent_name_unique_active
    ON departments(enterprise_id, name)
    WHERE deleted_at IS NULL;
