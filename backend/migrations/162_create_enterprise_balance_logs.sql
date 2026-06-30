-- 企业余额变更日志（审计表）
CREATE TABLE IF NOT EXISTS enterprise_balance_logs (
    id BIGSERIAL PRIMARY KEY,
    enterprise_id BIGINT NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    operation VARCHAR(20) NOT NULL,
    notes TEXT DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ebl_enterprise_id ON enterprise_balance_logs(enterprise_id);
CREATE INDEX IF NOT EXISTS idx_ebl_created_at ON enterprise_balance_logs(created_at DESC);
