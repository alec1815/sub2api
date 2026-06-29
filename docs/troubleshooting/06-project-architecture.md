# 06 — 项目架构分析（基于官方技术博客）

> **来源**: CSDN 博客《SUB2API 源码技术分析与搭建教程：把 AI 订阅变成可管理的 API 网关》
>
> 作者: 陌陌龙（项目作者 Wei-Shaw）| 发布时间: 2026-05-20
>
> 本文结合博客内容与我们已搭建的开发环境，梳理项目的核心架构和关键设计。

---

## 1. 项目定位

**Sub2API = 面向 AI 服务的 SaaS API 网关平台**

核心逻辑：
```
管理员接入上游 AI 账号/API Key
        ↓
生成平台 API Key 发给用户
        ↓
用户调用统一入口 → 系统自动完成:
    鉴权 → 选账号 → 转发 → 用量记录 → 扣费 → 限流 → 监控
```

与普通反向代理（Nginx）的关键区别：

| 能力 | Nginx 反代 | Sub2API |
|------|:---:|:---:|
| 基础转发 | ✅ | ✅ |
| API Key 鉴权 | ❌ | ✅ (平台级) |
| 多账号智能调度 | ❌ | ✅ (负载均衡+粘性) |
| Token 级计费扣费 | ❌ | ✅ |
| 用户订阅管理 | ❌ | ✅ |
| 限流 (多时间窗口) | ❌ | ✅ (Redis Lua 原子) |
| 管理后台 | ❌ | ✅ (Vue3 完整后台) |
| 支付/兑换码 | ❌ | ✅ |

---

## 2. 技术栈总览

```
┌────────────────────────────────────┐
│           前端 (Vue 3)              │
│  Vite + Pinia + Vue Router         │
│  TailwindCSS + TypeScript          │
├────────────────────────────────────┤
│           后端 (Go)                 │
│  Gin (HTTP)                        │
│  Ent (ORM, 代码生成)               │
│  Wire (依赖注入)                   │
├────────────────────────────────────┤
│         数据层                      │
│  PostgreSQL (用户/账号/计费/日志)   │
│  Redis (限流/缓存/会话粘性)        │
├────────────────────────────────────┤
│         部署 (Docker)               │
│  前端嵌入后端二进制 → 单镜像       │
│  systemd / Docker Compose          │
└────────────────────────────────────┘
```

### 与我们本地环境的对应

| 层级 | 博客推荐 | 我们的实际 |
|------|----------|------------|
| 后端 | `go build -o bin/server` | `go run ./cmd/server/` |
| 前端 | `pnpm run dev` | `pnpm dev` |
| 数据库 | Docker 内 PG | Windows 本地 PG 16 |
| Redis | Docker 内 Redis | Windows 本地 Redis |
| 部署 | Docker Compose | 本地开发模式 |

---

## 3. 运行模式：Standard vs Simple

与我们的代码对应：
```go
// cfg.RunMode == config.RunModeSimple
// → 隐藏计费、订阅等 SaaS 功能，适合内部使用
// cfg.RunMode == "standard" (默认)
// → 完整 SaaS 平台，含计费、余额检查、配额管理
```

我们当前使用的是 **debug 模式**（等同于 standard 下的开发模式）。

---

## 4. 一次 AI 调用的完整链路（7 步）

这是博客中最有价值的部分——**理解整个系统的数据流**：

```
用户请求
  │
  ▼
[1] Gin 路由匹配
  ├─ /v1/chat/completions       (OpenAI Chat)
  ├─ /v1/messages                (Claude Messages)
  ├─ /v1beta/models/...          (Gemini 原生)
  ├─ /v1/images/generations      (图片生成)
  └─ /backend-api/codex/...     (Codex 直连)
  │
  ▼
[2] API Key 鉴权 ← api_key_auth.go
  ├─ 查数据库 (PostgreSQL + Redis 缓存)
  ├─ 错误类型: NOT_FOUND / EXPIRED / QUOTA_EXHAUSTED / RATE_EXCEEDED
  └─ 通过 → 继续
  │
  ▼
[3] 检查分组、订阅、余额 ← billing_service.go
  ├─ 余额检查
  ├─ 分组订阅检查
  └─ 多时间窗口用量 (5小时/1天/7天)
  │
  ▼
[4] 选择上游账号（调度器核心） ← openai_account_scheduler.go
  ├─ session_hash → 粘性会话（同会话用同一上游）
  ├─ previous_response_id → 上下文响应链
  ├─ EWMA 错误率指标 → 故障账号自动降权
  ├─ TTFT (首Token延迟) EWMA → 性能加权
  ├─ 模块/传输方式/图片能力匹配
  └─ 排除黑名单账号
  │
  ▼
[5] 构造并转发上游请求 ← openai_gateway_service.go
  ├─ 不同协议兼容 (OpenAI/Claude/Gemini/Codex)
  ├─ OAuth vs API Key 不同上游地址
  ├─ 安全头白名单（只透传允许的头）
  └─ WebSocket 支持 (Codex CLI)
  │
  ▼
[6] 记录用量、扣费 ← billing_service.go
  ├─ ModelPricing: 输入/输出/cache/long_context token
  ├─ Token 级别精确计费
  └─ 写入 PostgreSQL
  │
  ▼
[7] 返回响应给用户
```

### 关键源码文件（按链路顺序）

| 步骤 | 文件 | 职责 |
|:---:|------|------|
| 1 | `internal/server/router.go` | 路由注册 |
| 2 | `internal/server/gateway.go` | 网关入口 |
| 3 | `internal/handler/api_key_handler.go` | Key 鉴权 |
| 4 | `internal/handler/openai_gateway_handler.go` | 请求处理 |
| 5 | `internal/service/openai_gateway_service.go` | 转发逻辑 |
| 5 | `internal/service/openai_account_scheduler.go` | 调度器 |
| 6 | `internal/service/billing_service.go` | 计费扣费 |
| 3/6 | `internal/service/api_key_service.go` | Key 管理 |

---

## 5. 值得深入学习的技术点

### 5.1 Redis Lua 原子限流

高并发下通过 Lua 脚本原子执行"计数+设置过期时间"，避免竞态条件：

```
EVAL "计数; if 超限 then reject; else incr + expire" 
```

这是 Redis 在高并发限流场景下的经典方案。

### 5.2 粘性会话

对话场景下优先根据 `session_id`/`session_hash` 分配到同一上游账号：
- 避免 Claude 上下文链断裂
- previous_response_id 链式路由

### 5.3 EWMA (指数加权移动平均) 负载均衡

调度器不只做轮询，而是：
- EWMA 错误率 → 故障账号自动降低权重
- TTFT EWMA → 首Token延迟作为性能指标
- 综合打分 → 智能选最优账号

### 5.4 安全头白名单

只透传允许的 HTTP 头到上游，防止：
- 用户传入的杂乱头触发上游风控
- 敏感信息泄漏

### 5.5 前端嵌入后端二进制

Docker 三阶段构建：
1. 编译前端 (Vite build)
2. 编译后端 (Go build)
3. Go embed 嵌入 `dist/` → 单二进制，无需 Nginx

对应我们的代码：
```go
//go:embed all:dist
// backend/internal/web/
```

---

## 6. 新手源码阅读路线

推荐 5 步阅读顺序：

```
❶ 启动与路由
   main.go → router.go → gateway.go
           ↓
❷ 用户与 API Key
   api_key_handler.go → api_key_service.go → api_key_auth_cache.go
           ↓
❸ 网关主链路 ★核心
   openai_gateway_handler.go → openai_gateway_service.go
   → openai_account_scheduler.go
           ↓
❹ 计费与用量
   billing_service.go → account_usage_service.go
   → usage_handler.go → ops_repo*.go
           ↓
❺ 前端后台
   router/index.ts → api/ → stores/ → views/admin/
```

---

## 7. 已知信息对照表

| 知识点 | 博客提到 | 我们已遇到 |
|--------|:---:|:---:|
| AUTO_SETUP 环境变量 | ✅ | ✅ (坑12) |
| 必须用 pnpm | ✅ | ✅ (坑1, 坑2) |
| Setup Wizard vs Auto Setup | ✅ | ✅ (阶段③) |
| config.yaml / .installed 双重锁 | - | ✅ (自己分析出) |
| DATA_DIR 路径优先级 | - | ✅ (坑12) |
| Go 自动工具链 | - | ✅ (刚发现) |
| Standard vs Simple 模式 | ✅ | 🆕 |
| 7步调用链路 | ✅ | 🆕 |
| Redis Lua 限流 | ✅ | 🆕 |
| EWMA 调度算法 | ✅ | 🆕 |
| 安全头白名单 | ✅ | 🆕 |

---

## 8. 下一步学习计划建议

基于博客指导的源码阅读路线，结合我们已有的环境：

1. **阅读网关主链路**（步骤❸）：`openai_gateway_handler.go` → `openai_gateway_service.go`
2. **理解调度器**：`openai_account_scheduler.go`（EWMA、粘性会话）
3. **理解计费模型**：`billing_service.go`（ModelPricing 结构体）
4. **前端管理后台**：`views/admin/` 下的页面组件

---

> 下一篇排错文档将在后续实际操作中补充。
