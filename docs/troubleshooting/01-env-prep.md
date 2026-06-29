# 01 — 阶段 ①：环境准备

> **目标**：PostgreSQL + Redis + Node.js/pnpm 全部可用

---

## 1. PostgreSQL 16

### 安装与配置

```
安装路径: D:\Develop\PostgreSQL\16
服务名:   postgresql-x64-16
端口:     5432
psql:     D:\Develop\PostgreSQL\16\bin\psql.exe
数据目录: D:\Develop\PostgreSQL\16\data
```

### 关键文件

| 文件 | 路径 | 用途 |
|------|------|------|
| `pg_hba.conf` | `D:\Develop\PostgreSQL\16\data\pg_hba.conf` | 客户端认证配置 |
| `postgresql.conf` | `D:\Develop\PostgreSQL\16\data\postgresql.conf` | 服务器配置 |

### 数据库凭据

| 角色 | 用户名 | 密码 | 用途 |
|------|--------|------|------|
| 超级用户 | `postgres` | `postgres` | 管理用 |
| 应用用户 | `sub2api` | `sub2api` | 应用连接 |
| 应用数据库 | `sub2api` | - | 数据存储 |

### 使用方式

```powershell
# 用超级用户连接
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U postgres -h 127.0.0.1

# 用应用用户连接
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U sub2api -h 127.0.0.1 -d sub2api

# 服务管理
Restart-Service postgresql-x64-16
Get-Service postgresql-x64-16
```

---

## 2. Redis

```
安装路径: D:\Develop\Redis
端口:     6379
密码:     无
cli:      D:\Develop\Redis\redis-cli.exe
```

### 验证

```powershell
& "D:\Develop\Redis\redis-cli.exe" PING
# 应返回: PONG
```

---

## 3. Node.js & pnpm

```
Node.js:  v22.15.0
安装路径: D:\Develop\Nodejs
pnpm:     11.9.0（通过 corepack 启用）
```

### ⚠️ 问题 1：pnpm 命令不可用

**现象**：安装 Node.js v22 后，运行 `pnpm` 提示命令不存在。

**原因**：Node.js v22 已内置 pnpm，但需要通过 `corepack enable` 启用。

**思考过程**：
- 最初认为需要 `npm install -g pnpm`
- 查了 Node.js 22 的 release notes → 发现 corepack 内置
- `corepack enable` 比全局安装更好（版本与 Node 绑定管理）

**解决**：
```powershell
corepack enable
```

---

### ⚠️ 问题 2：pnpm install 报 EPERM

**现象**：
```
pnpm install
# Error: EPERM: operation not permitted, unlink '...\node_modules\...'
```

**原因**：之前用 `npm install` 安装了依赖，`node_modules` 目录结构和权限与 pnpm 不兼容。pnpm 用 symlink 硬链接方式，和 npm 的扁平结构冲突。

**思考过程**：
- EPERM 在 Windows 上常见于文件残留
- pnpm 的 `node_modules` 结构特殊（`.pnpm` 目录 + symlink）
- 与 npm 的扁平展开完全不同的实现原理
- 所以必须先清空再用 pnpm 安装

**解决**：
```powershell
cd D:\Project\sub2api\frontend
Remove-Item -Recurse -Force node_modules
pnpm install
```

**涉及知识点**：
- npm: 扁平展开，所有依赖直接放到 `node_modules/` 下
- pnpm: 内容寻址存储 + 硬链接，通过 `.pnpm/` 目录和 symlink 实现
- 两者 `node_modules` 结构不兼容，不能混用

---

### ⚠️ 问题 3：pnpm-lock.yaml 未同步

**现象**：`package.json` 改了依赖但 `pnpm-lock.yaml` 未更新，导致：
- CI 中 `pnpm install --frozen-lockfile` 失败
- 其他开发者拉代码后依赖不一致

**原因**：手动编辑 `package.json` 后没有运行 `pnpm install` 来更新 lock 文件。

**思考过程**：
- CI 工作流用了 `--frozen-lockfile`（严格校验 lock 文件一致性）
- 只要 lock 文件和 package.json 不同步就会报错
- 这是 pnpm 的保护机制，防止依赖漂移

**解决**：
```powershell
cd D:\Project\sub2api\frontend
pnpm install  # 自动更新 pnpm-lock.yaml
git add pnpm-lock.yaml
```

**教训**：任何手动修改 `package.json` 后，必须运行一次 `pnpm install`。

---

## 4. Go 环境

```
Go 版本: ≥ 1.25.7
项目模块: backend/go.mod
```

后端代码在 `backend/` 目录下，所有 Go 命令需在 `backend/` 目录执行。

---

## 阶段小结

✅ PostgreSQL 16 正常服务
✅ Redis 正常服务
✅ pnpm 可用，前端依赖安装成功
✅ Go 环境就绪

> 下一步：[02-database-setup.md](./02-database-setup.md) — 数据库初始化
