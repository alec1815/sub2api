-- 企业表新增并发数列
ALTER TABLE enterprises ADD COLUMN IF NOT EXISTS concurrency INT NOT NULL DEFAULT 0;
COMMENT ON COLUMN enterprises.concurrency IS '企业级并发上限，0=不限制';
