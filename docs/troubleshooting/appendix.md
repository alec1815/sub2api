# 附录 — 关键代码解读 & 参考资料

## 1. 启动流程核心代码

### main.go — 入口决策树

```go
// backend/cmd/server/main.go

func main() {
    // ...

    // --setup 标志 → CLI 安装模式
    if *setupMode {
        if err := setup.RunCLI(); err != nil {
            log.Fatalf("Setup failed: %v", err)
        }
        return
    }

    // config.yaml 不存在时 → 需要安装
    if setup.NeedsSetup() {
        // 环境变量 AUTO_SETUP=true → 自动安装
        if setup.AutoSetupEnabled() {
            if err := setup.AutoSetupFromEnv(); err != nil {
                log.Fatalf("Auto setup failed: %v", err)
            }
        } else {
            // Web 安装向导
            runSetupServer()
            return
        }
    }

    // 正常模式
    runMainServer()
}
```

### GetDataDir — 路径优先级

```go
// backend/internal/setup/setup.go

func GetDataDir() string {
    // 1️⃣ DATA_DIR 环境变量 (最优先)
    if dir := os.Getenv("DATA_DIR"); dir != "" {
        return dir
    }

    // 2️⃣ /app/data (Docker 环境，写成静态路径是硬编码了 Linux 路径)
    dockerDataDir := "/app/data"
    if info, err := os.Stat(dockerDataDir); err == nil && info.IsDir() {
        testFile := dockerDataDir + "/.write_test"
        if f, err := os.Create(testFile); err == nil {
            f.Close()
            os.Remove(testFile)
            return dockerDataDir
        }
    }

    // 3️⃣ 当前目录
    return "."
}
```

> **Windows 注意事项**：如果不设置 `DATA_DIR`，且在 Windows 上没有 `/app/data` 目录，则回退到 `.`（当前工作目录）。如果 `C:\app\data` 存在但不可写，报 `Access is denied`。

### NeedsSetup — 双重锁

```go
func NeedsSetup() bool {
    // 锁 1: config.yaml
    if _, err := os.Stat(GetConfigFilePath()); !os.IsNotExist(err) {
        return false
    }

    // 锁 2: .installed (防止攻击者删除 config.yaml 触发重装)
    if _, err := os.Stat(GetInstallLockPath()); !os.IsNotExist(err) {
        return false
    }

    return true
}
```

### Install — 执行顺序

```go
func Install(cfg *SetupConfig) error {
    // 1. 安全性检查
    if !NeedsSetup() {
        return fmt.Errorf("system is already installed")
    }

    // 2. 生成 JWT key
    if cfg.JWT.Secret == "" {
        cfg.JWT.Secret, _ = generateSecret(32)
    }

    // 3. 数据库连接测试 (自动建库)
    TestDatabaseConnection(&cfg.Database)

    // 4. Redis 连接测试
    TestRedisConnection(&cfg.Redis)

    // 5. 执行 migration
    initializeDatabase(cfg)

    // 6. ★ 创建管理员 (唯一入口)
    createAdminUser(cfg)

    // 7. 写 config.yaml
    writeConfigFile(cfg)

    // 8. 写 .installed
    createInstallLock()

    return nil
}
```

**管理员在 Install() 的第 6 步被创建**——这是全代码库中创建管理员的唯一入口。

---

## 2. 启动模式的完整枚举

| 命令/条件 | 执行的函数 | 是否创建管理员 | 使用场景 |
|-----------|-----------|:---:|----------|
| `go run ./cmd/server/` (首次，无 config.yaml) | `runSetupServer()` | ✅ (Web 向导) | 首次部署 |
| `go run ./cmd/server/ --setup` | `setup.RunCLI()` | ✅ (CLI 交互) | 命令行安装 |
| `AUTO_SETUP=true go run ./cmd/server/` | `AutoSetupFromEnv()` | ✅ | Docker/自动化 |
| `go run ./cmd/server/` (config.yaml 已存在) | `runMainServer()` | ❌ | 日常运行 |
| 直接调用 `setup.Install(cfg)` | `Install()` | ✅ | 编程触发 |

---

## 3. 环境变量清单

### 自动安装所需 (`AUTO_SETUP=true`)

| 变量 | 必填 | 默认值 | 说明 |
|------|:---:|--------|------|
| `AUTO_SETUP` | ✅ | - | 设为 `true` 启用自动安装 |
| `DATABASE_HOST` | ✅ | - | 数据库地址 |
| `DATABASE_PORT` | ❌ | `5432` | 数据库端口 |
| `DATABASE_USER` | ✅ | - | 数据库用户 |
| `DATABASE_PASSWORD` | ✅ | - | 数据库密码 |
| `DATABASE_DBNAME` | ✅ | - | 数据库名 |
| `DATABASE_SSLMODE` | ❌ | `disable` | SSL 模式 |
| `REDIS_HOST` | ✅ | - | Redis 地址 |
| `REDIS_PORT` | ❌ | `6379` | Redis 端口 |
| `REDIS_PASSWORD` | ❌ | `""` | Redis 密码 |
| `ADMIN_EMAIL` | ✅ | - | 管理员邮箱 |
| `ADMIN_PASSWORD` | ✅ | - | 管理员密码 |
| `TIMEZONE` | ❌ | `UTC` | 时区 |
| `DATA_DIR` | ❌ | 见 GetDataDir | **Windows 必须设置** |
| `SERVER_MODE` | ❌ | `release` | 服务器模式 |

### 日常启动所需

| 变量 | Windows 必填 | 说明 |
|------|:---:|------|
| `DATA_DIR` | ✅ | 指向 `backend/` 目录 |

---

## 4. 关键文件位置

```
backend/
├── cmd/server/main.go          ← 启动入口，包含决策树
├── internal/setup/
│   ├── setup.go                ← NeedsSetup / Install / AutoSetupFromEnv / GetDataDir
│   └── handler.go              ← Web 安装向导的 API handler
├── config.yaml                 ← 安装流程生成（不要手动创建！）
├── .installed                  ← 安装锁文件（不要手动删除！）
└── migrations/                 ← 数据库迁移脚本
```

---

## 5. 排查方法论总结

### 本项目的排查模式

```
发现异常现象
    ↓
查数据库/文件系统 确认事实
    ↓                ↓
读 main.go       读 setup.go
    ↓                ↓
理解启动流程    理解核心函数
    ↓                ↓
      → 定位根因 ←
            ↓
   编写最小复现/验证脚本
            ↓
     确认修复方案
            ↓
   数据库 + 文件系统双重验证
```

### 最常用的文件

1. `backend/cmd/server/main.go` — 了解"什么时候会发生什么"
2. `backend/internal/setup/setup.go` — 了解配置/安装逻辑
3. `pg_hba.conf` — 解决数据库连接认证问题

### 最常用的命令

```powershell
# 查数据库
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U sub2api -h 127.0.0.1 -d sub2api -c "SELECT ..."

# 查文件
Test-Path "D:\Project\sub2api\backend\config.yaml"

# 查进程
netstat -ano | findstr ":8080"

# 健康检查
Invoke-WebRequest -Uri "http://localhost:8080/health" -UseBasicParsing
```

---

## 6. 参考资源

| 资源 | 链接 |
|------|------|
| 上游仓库 | https://github.com/Wei-Shaw/sub2api |
| PostgreSQL 文档 | https://www.postgresql.org/docs/16/ |
| pg_hba.conf 文档 | https://www.postgresql.org/docs/16/auth-pg-hba-conf.html |
| Ent ORM 文档 | https://entgo.io/docs/getting-started |
| pnpm 文档 | https://pnpm.io/ |
| Go 官方文档 | https://go.dev/doc/ |

---

> 本附录随项目排错经验持续更新。
