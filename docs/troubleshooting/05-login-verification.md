# 05 — 阶段 ⑤：登录验证

> **目标**：从数据库、文件系统和 HTTP 三个层面验证安装结果，最终通过前端登录

---

## 1. 数据库验证

```powershell
& "D:\Develop\PostgreSQL\16\bin\psql.exe" -U sub2api -h 127.0.0.1 -d sub2api -c "SELECT id, email, role, status FROM users;"
```

预期结果：
```
 id |        email        | role  | status
----+---------------------+-------+--------
  1 | admin@sub2api.local | admin | active
```

✅ 管理员账号存在且状态为 active。

---

## 2. 文件系统验证

```powershell
# config.yaml 存在
Test-Path "D:\Project\sub2api\backend\config.yaml"
# 应返回: True

# .installed 存在
Test-Path "D:\Project\sub2api\backend\.installed"
# 应返回: True
```

### config.yaml 内容检查

```yaml
server:
  host: 0.0.0.0
  port: 8080
  mode: debug

database:
  host: 127.0.0.1
  port: 5432
  user: sub2api
  password: sub2api
  dbname: sub2api
  sslmode: disable

redis:
  host: 127.0.0.1
  port: 6379
  password: ""
  db: 0

timezone: Asia/Shanghai
```

✅ 配置文件完整且配置正确。

---

## 3. 后端健康检查

```powershell
Invoke-WebRequest -Uri "http://localhost:8080/health" -UseBasicParsing
```

预期结果：
```
StatusCode: 200
Content: {"status":"ok"}
```

✅ 后端正常运行。

---

## 4. 启动后端（长期运行）

```powershell
cd D:\Project\sub2api\backend
$env:DATA_DIR="D:\Project\sub2api\backend"  # 这个必须设置！
go run ./cmd/server/
```

> **注意**：每次启动都需要设置 `DATA_DIR`。可以考虑写入 PowerShell Profile 或创建启动脚本。

**启动日志**：
```
Server started on 0.0.0.0:8080
```

---

## 5. 启动前端

```powershell
# 新开一个终端
cd D:\Project\sub2api\frontend
pnpm dev
```

---

## 6. 登录

打开浏览器访问 `http://localhost:8080`（或前端 dev server 地址），用以下凭据登录：

| 字段 | 值 |
|------|-----|
| 邮箱 | `admin@sub2api.local` |
| 密码 | `admin123` |

---

## ⚠️ 问题 14：前端 dev 需要额外启动

**现象**：只启动了后端，前端页面无法访问。

**原因**：后端是 API 服务（`:8080`），前端是 Vue dev server（独立的端口），需要分别启动。

**解决**：两个终端分别启动前后端。

---

## 最终状态确认

| 检查项 | 状态 | 命令/验证方式 |
|--------|------|---------------|
| PostgreSQL 运行 | ✅ | `Get-Service postgresql-x64-16` |
| Redis 运行 | ✅ | `redis-cli PING` → `PONG` |
| 数据库有管理员 | ✅ | `SELECT * FROM users` |
| config.yaml 存在 | ✅ | `Test-Path config.yaml` |
| .installed 存在 | ✅ | `Test-Path .installed` |
| 后端 health | ✅ | `curl localhost:8080/health` |
| 前端可访问 | ✅ | 浏览器打开 `http://localhost:8080` |
| 可以登录 | ✅ | `admin@sub2api.local` / `admin123` |

---

## 完整的启动命令（备忘）

```powershell
# ====== 终端 1：后端 ======
cd D:\Project\sub2api\backend
$env:DATA_DIR="D:\Project\sub2api\backend"
go run ./cmd/server/

# ====== 终端 2：前端 ======
cd D:\Project\sub2api\frontend
pnpm dev
```

---

## 阶段小结

🎉 **全部完成！** 从零搭建开发环境到成功登录中转站系统，共解决 14 个问题，覆盖环境、数据库、后端启动、管理员创建、登录验证五大阶段。

> 附录：[appendix.md](./appendix.md) — 关键代码解读与参考资料
