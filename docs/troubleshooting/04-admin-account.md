# 04 — 阶段 ④：管理员账号创建

> **目标**：通过 `AUTO_SETUP=true` 自动安装流程创建管理员账号
>
> **这是问题最多的一步**——先后遇到了 Windows 进程管理、文件路径、build tag 等多个障碍，尝试了 4 种方案才最终成功。

---

## 前置操作

确认根因后，删除之前手动创建的 `config.yaml`（可能存在的 `.installed` 也删除）：

```powershell
cd D:\Project\sub2api\backend
Remove-Item config.yaml -Force -ErrorAction SilentlyContinue
Remove-Item .installed -Force -ErrorAction SilentlyContinue
```

---

## 尝试 1：直接 go run ❌

```powershell
$env:AUTO_SETUP="true"
$env:DATABASE_HOST="127.0.0.1"
# ... 其他环境变量 ...
go run ./cmd/server/
```

**问题**：`go run ./cmd/server/` 启动 HTTP 服务器后**进程持续运行不退出**，无法捕获自动安装的输出，也不确定安装是否成功。

**思考**：Go 的 HTTP server 是阻塞式的，`ListenAndServe()` 不返回，我需要让安装完成后的日志可见。

---

## 尝试 2：PowerShell Start-Job ❌

```powershell
$job = Start-Job -ScriptBlock {
    Set-Location $using:PWD
    go run ./cmd/server/ 2>&1
}
Start-Sleep -Seconds 15
Receive-Job $job
```

**问题**：`Start-Job` 创建了新的 PowerShell 进程，但 `go run` 的输出无法正确重定向回来。`Receive-Job` 拿到的是空内容。

**思考**：`Start-Job` 的 ScriptBlock 中是独立的 runspace，环境变量、工作目录的传递都有不确定性。Go 的 stdout 可能被缓冲，不会及时 flush 到 PowerShell job 的管道中。

---

## 尝试 3：cmd /c start /b ❌

```powershell
cmd /c "cd /d D:\Project\sub2api\backend && set AUTO_SETUP=true && ... && start /b go run ./cmd/server/ > setup.log 2>&1"
ping -n 20 127.0.0.1 > nul  # 等待 20 秒
type setup.log
```

**问题**：
1. `go run` 需要编译 + 运行，启动时间不确定（~30s+）
2. `ping -n 20` 等待时间不够时，日志文件还是空的
3. `start /b` 在 Windows 上对 Go 进程的重定向不可靠
4. cmd 环境变量的 `set` 与 PowerShell 的环境隔离

**思考**：我需要绕过"启动服务器"这个干扰因素。Go 的服务器和安装逻辑在同一个二进制中，但安装逻辑应该在服务器启动之前完成。

---

## 尝试 4：写辅助脚本 setup_run.go ✅

### 思路

既然 `setup.Install()` 是一个可导出的函数，我可以写一个最小脚本直接调用它，完全绕过 HTTP 服务器。

### 第一版（有问题）

```go
// +build ignore

package main

import (
    "os"
    "github.com/Wei-Shaw/sub2api/internal/setup"
)

func main() {
    os.Setenv("AUTO_SETUP", "true")
    // ... env vars ...
    
    // 检查是否需要安装
    if !setup.NeedsSetup() {
        println("Already installed, exiting")
        return
    }
    
    err := setup.AutoSetupFromEnv()
    if err != nil {
        panic(err)
    }
    println("Auto setup completed!")
}
```

### ⚠️ 问题 11：`+build ignore` 导致无法运行

**现象**：`go run setup_run.go` 报错找不到 `main` 包。

**原因**：文件头部写了 `// +build ignore`，这是 Go 的构建标签，告诉编译器忽略此文件。

**思考**：`// +build ignore` 通常用在测试工具文件中，防止它们被纳入正式构建。但我们是直接用 `go run` 执行，不需要这个标签。

**解决**：删除 `// +build ignore`：
```go
package main  // ← 直接声明 main 包
```

### 第一版运行结果

```
Auto setup enabled, configuring from environment variables...
Data directory: ./
Database connection successful
Redis connection successful
Database initialized successfully
Admin user already exists, skipping admin bootstrap   ← 管理员已存在！
Configuration file created
/app/data/.installed: Access is denied                  ← 路径错误！
```

### 两个关键发现

1. **"Admin user already exists"**——在之前的某次尝试中，管理员已经被创建了！
2. **`/app/data/.installed: Access is denied`**——Windows 无法写入 `/app/data/` 路径

---

## ⚠️ 问题 12（致命）：DATA_DIR 默认路径错误

### 现象

```
Access is denied: /app/data/.installed
```

`config.yaml` 和 `.installed` 被写入到了错误的位置（或写入失败）。

### 根因

`GetDataDir()` 的逻辑：

```go
func GetDataDir() string {
    // Priority 1: DATA_DIR 环境变量
    if dir := os.Getenv("DATA_DIR"); dir != "" {
        return dir
    }

    // Priority 2: /app/data (Docker 环境!)
    dockerDataDir := "/app/data"
    if info, err := os.Stat(dockerDataDir); err == nil && info.IsDir() {
        // 检查可写性...
        return dockerDataDir
    }

    // Priority 3: 当前目录
    return "."
}
```

在 Windows 上：
- `DATA_DIR` 未设置 → Priority 1 跳过
- `/app/data` 可能不存在或不可写 → 在 Windows 上 `os.Stat("/app/data")` 返回 error（目录不存在）
- 最后回退到 `.`（当前目录）

但第一次运行时（尝试 1 的 `go run`），当前目录可能是 `backend/`，文件其实已写入但被后续的删除命令清掉了。

> 实际上，`/app/data` 在 Windows 上是对应 `C:\app\data`，如果该目录存在但无法写入，就会报 `Access is denied`。但这里显示是在尝试创建 `.installed` 时失败了。

### 第二版（修复）

```go
os.Setenv("DATA_DIR", "D:/Project/sub2api/backend")  // ← 关键新增
os.Setenv("AUTO_SETUP", "true")
// ... 其他环境变量 ...
```

### 第二版运行结果

```
Auto setup enabled, configuring from environment variables...
Data directory: D:/Project/sub2api/backend
Database connection successful            ✅
Redis connection successful               ✅
Database initialized successfully         ✅
Admin user already exists                 ✅
Configuration file created                ✅
Installation lock created                 ✅
Auto setup completed successfully!        ✅
```

**全部通过！**

---

## ⚠️ 问题 13：Admin already exists（信息性）

这说明在尝试 1 的某次 `go run` 中，虽然没看到输出，但安装流程确实执行了——管理员已被创建，只是 `config.yaml` 和 `.installed` 写到了错误路径。

**验证**：
```powershell
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U sub2api -h 127.0.0.1 -d sub2api -c "SELECT id, email, role, status FROM users;"
```
```
 id |        email        | role  | status
----+---------------------+-------+--------
  1 | admin@sub2api.local | admin | active
```

管理员确实存在。

---

## 最终版本 setup_run.go（完整）

```go
package main

import (
    "os"
    "github.com/Wei-Shaw/sub2api/internal/setup"
)

func main() {
    os.Setenv("DATA_DIR", "D:/Project/sub2api/backend")
    os.Setenv("AUTO_SETUP", "true")
    os.Setenv("DATABASE_HOST", "127.0.0.1")
    os.Setenv("DATABASE_PORT", "5432")
    os.Setenv("DATABASE_USER", "sub2api")
    os.Setenv("DATABASE_PASSWORD", "sub2api")
    os.Setenv("DATABASE_DBNAME", "sub2api")
    os.Setenv("DATABASE_SSLMODE", "disable")
    os.Setenv("REDIS_HOST", "127.0.0.1")
    os.Setenv("REDIS_PORT", "6379")
    os.Setenv("ADMIN_EMAIL", "admin@sub2api.local")
    os.Setenv("ADMIN_PASSWORD", "admin123")
    os.Setenv("TIMEZONE", "Asia/Shanghai")

    if !setup.NeedsSetup() {
        println("Already installed, exiting")
        return
    }

    err := setup.AutoSetupFromEnv()
    if err != nil {
        panic(err)
    }
    println("Auto setup completed!")
}
```

> **注意**：完成安装后，这个文件被删除了（仅作为一次性工具）。

---

## 方案对比

| 方案 | 结果 | 失败原因 | 
|------|------|----------|
| `go run ./cmd/server/` | ❌ | 服务器不退，无法看到安装日志 |
| `Start-Job` | ❌ | stdout 重定向失败，输出为空 |
| `cmd /c start /b` | ❌ | 编译耗时超出等待，日志空 |
| `setup_run.go` | ✅ | 直接调用 Install()，绕过 HTTP server |

**为什么方案 4 成功？**
- `setup.Install()` 是纯函数，输出到 stdout，不启动 HTTP server
- 直接读取环境变量，不依赖配置文件
- 同步执行，所有日志实时可见
- `go run` 执行完就退出，不需要管理进程生命周期

---

## 阶段小结

✅ 通过辅助脚本成功执行了安装流程
✅ config.yaml 写入正确路径
✅ .installed lock 文件创建成功
✅ 管理员账号 `admin@sub2api.local` 已存在于数据库
✅ 所有 migration 执行完毕

> 下一步：[05-login-verification.md](./05-login-verification.md) — 登录验证
