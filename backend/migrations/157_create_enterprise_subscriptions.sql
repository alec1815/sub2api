-- 157_create_enterprise_subscriptions.sql
-- 企业套餐表 —— 企业功能 P1 新增
-- 对应 Ent schema: ent/schema/enterprise_subscription.go
-- 与企业管理员个人套餐完全隔离，无软删除（status 管理生命周期）

CREATE TABLE IF NOT EXISTS enterprise_subscriptions (
    id                  BIGSERIAL PRIMARY KEY,
    enterprise_id       BIGINT NOT NULL REFERENCES enterprises(id),
    group_id            BIGINT NOT NULL REFERENCES groups(id),
    plan_id             BIGINT NOT NULL REFERENCES subscription_plans(id),
    starts_at           TIMESTAMPTZ NOT NULL,
    expires_at          TIMESTAMPTZ,
    status              VARCHAR(20) NOT NULL DEFAULT 'active',
    daily_usage_usd     DECIMAL(20,10) NOT NULL DEFAULT 0,
    weekly_usage_usd    DECIMAL(20,10) NOT NULL DEFAULT 0,
    monthly_usage_usd   DECIMAL(20,10) NOT NULL DEFAULT 0,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- 查企业有效套餐
CREATE INDEX IF NOT EXISTS idx_enterprise_subs_ent_id_status ON enterprise_subscriptions(enterprise_id, status);

-- 查企业某分组套餐
CREATE INDEX IF NOT EXISTS idx_enterprise_subs_ent_id_group_id ON enterprise_subscriptions(enterprise_id, group_id);
