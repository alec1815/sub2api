# 00 — 全流程概览

## 时序总图

```
┌─ 阶段 ①：环境准备 ────────────────────────────────────────────────────┐
│                                                                         │
│  PostgreSQL 16 ──✅ 安装 (D:\Develop\PostgreSQL\16)                    │
│  Redis         ──✅ 安装 (D:\Develop\Redis)                           │
│  Node.js       ──✅ v22.15.0 (D:\Develop\Nodejs)                      │
│  pnpm          ──⚠️ 需要 corepack enable → ✅                          │
│  pnpm install  ──⚠️ npm 残留 node_modules 冲突 → 删除重装 → ✅         │
│  pnpm-lock.yaml──⚠️ 未同步 → pnpm install 更新 → ✅                    │
│                                                                         │
├─ 阶段 ②：数据库初始化 ────────────────────────────────────────────────┤
│                                                                         │
│  psql 连接     ──⚠️ localhost → IPv6 问题 → 改用 127.0.0.1 → ✅        │
│  pg_hba.conf   ──⚠️ 忘记密码 → trust 模式重置 → ✅                     │
│  创建用户/库   ──✅ sub2api / sub2api                                  │
│  执行 SQL      ──⚠️ PowerShell $2a$ 转义 → 写入文件再执行 → ✅         │
│                                                                         │
├─ 阶段 ③：后端首次启动 ────────────────────────────────────────────────┤
│                                                                         │
│  config.yaml   ──⚠️ 手动创建 → 跳过 setup → 管理员未创建 ❌              │
│  排查原因      ──🔍 读 NeedsSetup() 逻辑 → 发现根本原因                │
│                                                                         │
├─ 阶段 ④：管理员创建 ──────────────────────────────────────────────────┤
│                                                                         │
│  删除 config.yaml ──✅                                                  │
│  AUTO_SETUP      ──⚠️ go run 不退 → Start-Job 无用 → cmd /c 仍失败     │
│  辅助脚本         ──🔧 setup_run.go → +build ignore 冲突 → 修复 → ✅    │
│  DATA_DIR 缺失   ──⚠️ C:\app\data Access denied → 设置 DATA_DIR → ✅   │
│  管理员已存在     ──✅ (首次尝试已创建)                                 │
│                                                                         │
├─ 阶段 ⑤：登录验证 ────────────────────────────────────────────────────┤
│                                                                         │
│  psql 验证       ──✅ id=1, email=admin@sub2api.local, role=admin      │
│  config.yaml     ──✅ 正确路径 D:\Project\sub2api\backend              │
│  .installed      ──✅ lock 文件已创建                                  │
│  health check    ──✅ {"status":"ok"}                                  │
│  前端登录         ──✅ admin@sub2api.local / admin123                  │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

## 问题清单——共 14 个问题

| # | 阶段 | 问题 | 根因 | 严重度 |
|---|------|------|------|--------|
| 1 | ① | pnpm 不可用 | Node.js 内置但需 `corepack enable` | 低 |
| 2 | ① | pnpm install 报 EPERM | npm 残留 `node_modules` | 低 |
| 3 | ① | pnpm-lock.yaml 未同步 | 手动改 package.json 但未 `pnpm install` | 低 |
| 4 | ② | psql 连 localhost 超时 | Windows IPv6 优先，pg_hba 只配了 IPv4 | 中 |
| 5 | ② | 忘记 PG 密码 | 首次配置后未记录 | 高 |
| 6 | ② | bcrypt hash `$` 被吃掉 | PowerShell 变量展开 | 中 |
| 7 | ③ | 管理员账号不存在 | 手动创建 config.yaml 跳过了 Install() | **致命** |
| 8 | ④ | `go run` 不退无法捕获输出 | Go HTTP server 持续运行 | 中 |
| 9 | ④ | PowerShell Start-Job 无效 | Go 进程的 stdout/stderr 重定向问题 | 中 |
| 10 | ④ | cmd /c start /b 无输出 | Go 编译+运行耗时超出等待 | 中 |
| 11 | ④ | setup_run.go `+build ignore` | 文件写了 build tag 无法直接运行 | 中 |
| 12 | ④ | `Access is denied: /app/data/` | Windows 上 DATA_DIR 默认为 `/app/data` | **致命** |
| 13 | ④ | Admin already exists | 第 7 步删 config.yaml 后第一次 go run 已创建 | 信息 |
| 14 | ⑤ | 前端 dev 未启动 | 之前一直聚焦后端问题 | 低 |

## 经验总结

### 1. "先读代码再动手" 是最高效的排查方法

**场景**：管理员不存在的根本原因是什么？

- ❌ 盲目操作：反复重启、乱改配置、查 Google
- ✅ 正确做法：**读代码** `main.go` → 找到 `NeedsSetup()` → 理解启动决策树 → 定位真因

```
读代码耗时: ~5 分钟
盲目操作: 可能浪费数小时
```

### 2. Windows 路径问题需特别警惕

`GetDataDir()` 优先级：
```
1. $DATA_DIR 环境变量           ← 我们需要的
2. /app/data (Docker 环境)      ← Linux 路径！Windows 上导致 Access denied
3. 当前目录 "."                 ← 最后一个回退
```

**教训**：任何工具在 Windows 上使用时，必须先检查是否有 Linux-only 默认路径假设。

### 3. 脚本工具要可靠，多次尝试不要羞于换方案

```
尝试 1: 直接 go run → 进程不退，无法看到输出
尝试 2: PowerShell Start-Job → 无法捕获 Go 子进程输出
尝试 3: cmd /c start /b → 同样无法可靠捕获
尝试 4: 写 setup_run.go 辅助脚本 → ✅ 直接调用 Install() 函数，输出清晰
```

### 4. 排错时的信息收集顺序

```
1. 查数据库状态（数据在不在？）              → psql SELECT
2. 读启动代码理解流程（为什么会这样？）       → main.go / setup.go
3. 找配置相关逻辑（条件判断是什么？）         → NeedsSetup() / GetDataDir()
4. 用最小脚本验证（绕过干扰因素）             → setup_run.go
5. 确认结果（数据库 + 文件系统双重验证）       → psql + Test-Path
```

### 5. 文档与排错应同步记录

在解决问题时同步写文档，而不是事后回忆。关键信息包括：
- 当时看到的错误信息（原始输出）
- 读了哪个文件的哪段代码
- 尝试了哪些方案（包括失败的）
- 最终生效的方案和原因

这直接产出了本排错目录的内容。

---

> **下一步**：按阶段阅读 [01-env-prep.md](./01-env-prep.md) → [02-database-setup.md](./02-database-setup.md) → ...
