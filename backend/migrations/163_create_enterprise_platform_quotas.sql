-- 企业平台限额表
CREATE TABLE IF NOT EXISTS enterprise_platform_quotas (
    id BIGSERIAL PRIMARY KEY,
    enterprise_id BIGINT NOT NULL,
    platform VARCHAR(50) NOT NULL,
    daily_limit_usd DOUBLE PRECISION,
    weekly_limit_usd DOUBLE PRECISION,
    monthly_limit_usd DOUBLE PRECISION,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_epq_enterprise_platform ON enterprise_platform_quotas(enterprise_id, platform);
CREATE INDEX IF NOT EXISTS idx_epq_enterprise_id ON enterprise_platform_quotas(enterprise_id);
