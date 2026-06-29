# 02 — 阶段 ②：数据库初始化

> **目标**：PostgreSQL 中创建 sub2api 数据库，执行 migration，验证表结构

---

## ⚠️ 问题 4：psql 连接 localhost 超时/很慢

**现象**：
```powershell
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U postgres -h localhost
# 等待很久，然后才连上（或直接超时）
```

**原因**：Windows TCP/IP 栈优先尝试 IPv6 (`::1`)，而 `pg_hba.conf` 可能只配置了 IPv4 (`127.0.0.1/32`)。psql 需要等 IPv6 超时后回退到 IPv4。

**思考过程**：
- `localhost` 在 hosts 文件中同时映射到 `::1` 和 `127.0.0.1`
- psql 默认优先 IPv6
- `pg_hba.conf` 的 `host` 条目只匹配 `127.0.0.1/32`（IPv4）
- 连接失败 → 等超时 → 回退 IPv4 → 第二次成功
- 整个过程可以长达 30 秒

**解决**：始终用 `127.0.0.1` 代替 `localhost`：
```powershell
# ✅ 快
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U postgres -h 127.0.0.1

# ❌ 可能慢
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U postgres -h localhost
```

---

## ⚠️ 问题 5：忘记 PostgreSQL 密码

**现象**：不记得 `sub2api` 或 `postgres` 用户的密码。

**原因**：首次安装 PG 时设了密码但未记录。

**解决思路**：
1. PostgreSQL 通过 `pg_hba.conf` 控制认证方式
2. 可以将认证方式改为 `trust`（信任，无需密码）
3. 无密码登录后重置密码
4. 改回安全认证方式

**详细步骤**：

### Step 1：修改 pg_hba.conf
```
# 文件: D:\Develop\PostgreSQL\16\data\pg_hba.conf

# 找到并修改 IPv4 条目：
# 原始: host    all    all    127.0.0.1/32    scram-sha-256
# 改为:
host    all    all    127.0.0.1/32    trust
```

### Step 2：重启服务
```powershell
Restart-Service postgresql-x64-16
```

### Step 3：无密码登录并重置
```powershell
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U postgres -h 127.0.0.1
```
```sql
-- 在 psql 中执行：
ALTER USER sub2api WITH PASSWORD 'sub2api';
ALTER USER postgres WITH PASSWORD 'postgres';
\q
```

### Step 4：恢复认证方式
```
# pg_hba.conf 改回:
host    all    all    127.0.0.1/32    scram-sha-256

# 重启服务:
Restart-Service postgresql-x64-16
```

**关键技术点**：
- `pg_hba.conf` 是 PostgreSQL 的客户端认证配置文件
- `trust` 模式 = 不验证密码（仅本地开发用）
- `scram-sha-256` = 安全的密码认证（生产必须）
- 改配置后必须重启服务才生效

---

## 创建应用数据库（如果需要）

```powershell
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U postgres -h 127.0.0.1
```
```sql
CREATE USER sub2api WITH PASSWORD 'sub2api';
CREATE DATABASE sub2api OWNER sub2api;
GRANT ALL PRIVILEGES ON DATABASE sub2api TO sub2api;
\q
```

> 注意：在后续的 setup 安装流程中，`TestDatabaseConnection()` 也会自动检测并创建数据库（如果不存在）。

---

## ⚠️ 问题 6：bcrypt hash 在 PowerShell 的 `$` 转义问题

**现象**：用 psql 执行包含 bcrypt hash 的 INSERT 语句时，hash 值被截断或变形。

```
原始: $2a$10$XyZabc123...
结果: a$10$...  （$2 被 PowerShell 当作变量 $2 展开成了空值）
```

**原因**：PowerShell 中 `$` 是变量前缀符号。`$2a` 被解释为变量 `$2a`，不存在 → 展开为 `""`。

**思考过程**：
- 直接在命令行写 SQL 是最容易出错的方式
- 用单引号 `'...'` 可以阻止变量展开（PowerShell 特性）
- 但 bcrypt hash 里也可能有单引号，会破坏 SQL 语法
- 更安全的方式：写入 `.sql` 文件，用 `psql -f` 执行

**解决——写入 SQL 文件再执行**：
```powershell
# 创建 SQL 文件
@"
INSERT INTO users (email, password_hash, role, status)
VALUES ('admin@example.com', '`$2a`$10`$XyZabc123...', 'admin', 'active');
"@ | Out-File -FilePath temp.sql -Encoding utf8

# 用 psql 执行文件
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U sub2api -h 127.0.0.1 -d sub2api -f temp.sql

# 清理
Remove-Item temp.sql
```

> 但实际开发中，管理员应由 setup 流程自动创建，无需手动 INSERT。

---

## 数据库 Migration

后端的 setup 流程会自动执行 migration，不需要手动操作。Migration 文件位于 `backend/migrations/`。

手动执行（仅调试时）：
```powershell
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U sub2api -h 127.0.0.1 -d sub2api -f backend/migrations/xxx.sql
```

---

## 阶段小结

✅ psql 能用 `127.0.0.1` 正常连接（避免 IPv6）
✅ 密码管理流程清晰（`pg_hba.conf` trust 模式兜底）
✅ 安全执行含特殊字符的 SQL（写文件再 `-f` 执行）
✅ Migration 由 setup 流程自动处理

> 下一步：[03-backend-startup.md](./03-backend-startup.md) — 后端首次启动
