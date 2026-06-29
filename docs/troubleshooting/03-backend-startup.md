# 03 — 阶段 ③：后端首次启动

> **目标**：成功启动后端 `go run ./cmd/server/`
>
> **这是最关键的一步**——错误的操作方式会导致管理员账号不会自动创建。

---

## 启动命令

```powershell
cd D:\Project\sub2api\backend
$env:DATA_DIR="D:\Project\sub2api\backend"  # Windows 必须设置！
go run ./cmd/server/
```

---

## ⚠️ 问题 7（致命）：管理员账号不存在

### 现象

后端启动后，可以访问 `http://localhost:8080`，但尝试登录时发现没有任何用户，无法进入系统。

### 排查过程（思考链）

#### Step 1：验证事实——查数据库

```powershell
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U sub2api -h 127.0.0.1 -d sub2api -c "SELECT id, email, role, status FROM users;"
```

结果：
```
 id | email | role | status
----+-------+------+--------
(0 rows)           ← 没有任何用户！
```

**事实确认**：管理员确实不存在。

#### Step 2：读代码——启动流程是什么？

打开 `backend/cmd/server/main.go`，阅读 `main()` 函数：

```go
func main() {
    // ... flag 解析 ...

    if setup.NeedsSetup() {             // ← ① 检查是否需要安装
        if setup.AutoSetupEnabled() {   // ← ② 检查是否自动安装
            setup.AutoSetupFromEnv()    // ← 从环境变量自动安装
        } else {
            runSetupServer()            // ← 启动安装向导
            return
        }
    }

    runMainServer()                     // ← ③ 正常启动
}
```

**关键发现**：`NeedsSetup()` 返回 `false` 时，① 和 ② 直接跳过，进入 ③——**不会创建管理员**。

#### Step 3：深入关键函数——NeedsSetup 的判断条件

打开 `backend/internal/setup/setup.go`：

```go
func NeedsSetup() bool {
    // Check 1: Config file must not exist
    if _, err := os.Stat(GetConfigFilePath()); !os.IsNotExist(err) {
        return false  // config.yaml 存在 → 不需要安装！
    }

    // Check 2: Installation lock file
    if _, err := os.Stat(GetInstallLockPath()); !os.IsNotExist(err) {
        return false  // .installed 存在 → 已经安装过！
    }

    return true
}
```

#### Step 4：回溯——为什么 config.yaml 存在？

在之前的操作中，为了让后端能启动（当时以为必须有 config.yaml），**手动创建了 `config.yaml`**。

```powershell
# 这是当时的错误操作！
New-Item config.yaml
# 手动写入了数据库配置...
```

### 根因分析

```
手动创建 config.yaml
        ↓
NeedsSetup() → "config.yaml 存在" → return false
        ↓
跳过 Install() → 跳过 createAdminUser()
        ↓
管理员账号永远不创建
```

**这是一个典型的"好心办坏事"**：为了让系统跑起来，我绕过了它的设计流程，但绕过的那个流程恰好是创建关键数据（管理员）的唯一路径。

### 学到的设计模式

这是一个常见的安全设计模式——**安装锁**：

| 检查项 | 文件 | 作用 |
|--------|------|------|
| Check 1 | `config.yaml` | 配置文件存在 = 不需要安装 |
| Check 2 | `.installed` | 显式安装锁 = 防止攻击者删 config 重装 |

双重保险：即使有人删了 `config.yaml`，`.installed` 还在，仍然不会触发重装。

---

## 启动流程决策树（完整）

```
go run ./cmd/server/
│
├─ --version 标志? → 打印版本，退出
│
├─ --setup 标志? → setup.RunCLI() → CLI 交互安装 → 退出
│
└─ 正常模式:
   │
   ├─ NeedsSetup()?
   │
   ├─ YES:
   │  ├─ AutoSetupEnabled()?
   │  │  ├─ YES → AutoSetupFromEnv() → 继续到 runMainServer()
   │  │  └─ NO  → runSetupServer() (Web 安装向导) → 退出
   │
   └─ NO:
      └─ runMainServer() (跳过一切，直接运行)
```

**管理员只在以下路径中被创建**：
- `--setup` CLI 安装
- `AutoSetupFromEnv()` 自动安装
- Web 安装向导 (`runSetupServer`)
- 直接调用 `setup.Install()`

**管理员不会在以下路径被创建**：
- `runMainServer()` ——正常启动时

---

## Install() 的执行顺序

```
Install(cfg)
│
├─ 1. NeedsSetup() 二次检查
├─ 2. TestDatabaseConnection()  → 连接测试 + 自动建库
├─ 3. TestRedisConnection()     → 连接测试
├─ 4. initializeDatabase()      → 执行 ~190 个 migration 文件
├─ 5. createAdminUser()         → ★ 创建管理员 (仅此处)
├─ 6. writeConfigFile()         → 输出 config.yaml
└─ 7. createInstallLock()       → 输出 .installed
```

**关键洞察**：`createAdminUser()` 是在 migration 之后、config 写入之前执行的。管理员用 setup 过程中提供的邮箱和密码，数据库用 migration 创建的表结构。

---

## 阶段小结

✅ 理解了 `main.go` 的启动流程
✅ 理解了 `NeedsSetup()` 的双重锁机制
✅ 发现了手动创建 `config.yaml` 是管理员不创建的根因
❌ 管理员仍未创建 → 需要在下一阶段走正确的安装流程

> 下一步：[04-admin-account.md](./04-admin-account.md) — 删除 config.yaml，走正确安装流程
