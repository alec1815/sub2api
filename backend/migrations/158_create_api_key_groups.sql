-- 158_create_api_key_groups.sql
-- Key-分组中间表（M:N）—— 企业功能 P1 新增
-- 对应 Ent schema: ent/schema/api_key_group.go
-- 替代 api_keys.group_id 的 1:1 关联，支持一个 Key 关联多个分组

CREATE TABLE IF NOT EXISTS api_key_groups (
    id              BIGSERIAL PRIMARY KEY,
    api_key_id      BIGINT NOT NULL REFERENCES api_keys(id),
    group_id        BIGINT NOT NULL REFERENCES groups(id),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 唯一约束：同一 Key 不可重复关联同一分组
CREATE UNIQUE INDEX IF NOT EXISTS api_key_groups_key_group_unique
    ON api_key_groups(api_key_id, group_id);

-- 查某分组下所有 Key
CREATE INDEX IF NOT EXISTS idx_api_key_groups_group_id ON api_key_groups(group_id);

-- 数据迁移：将 api_keys.group_id 的已有 1:1 关联迁移到中间表
INSERT INTO api_key_groups (api_key_id, group_id, created_at)
    SELECT id, group_id, NOW()
    FROM api_keys
    WHERE group_id IS NOT NULL
    AND NOT EXISTS (
        SELECT 1 FROM api_key_groups agk
        WHERE agk.api_key_id = api_keys.id AND agk.group_id = api_keys.group_id
    );
