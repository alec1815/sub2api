# 企业功能 — API 接口文档

> **版本**: V1.0 | **日期**: 2026-06-29
> **基础**: 基于 02-核心设计决策 / 03-数据库设计概要 / 04-开发功能概要设计
> **约定**: 所有响应格式 `{ code: 0, message: "ok", data: ... }`

---

## 目录

- [一、通用约定](#一通用约定)
- [二、平台运营方 API (admin)](#二平台运营方-api-admin)
  - [企业管理](#21-企业管理)
  - [企业成员管理](#22-企业成员管理)
  - [部门管理](#23-部门管理)
- [三、企业管理员 API (enterprise)](#三企业管理员-api-enterprise)
  - [成员管理](#31-成员管理)
  - [密钥管理](#32-密钥管理)
  - [财务管理（企业资金池）](#33-财务管理企业资金池)
  - [部门管理](#34-部门管理)
  - [企业设置](#35-企业设置)
- [四、企业成员 API (member)](#四企业成员-api-member)
- [五、个人用户 API 改造 (user)](#五个人用户-api-改造-user)
- [六、通用错误码](#六通用错误码)

---

## 一、通用约定

### 1.1 响应格式

```json
{
  "code": 0,
  "message": "ok",
  "data": { ... }
}
```

| code | 含义 |
|:--:|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未登录 / Token 过期 |
| 403 | 权限不足 |
| 404 | 资源不存在 |
| 409 | 冲突（名称重复等） |
| 422 | 业务规则校验失败 |
| 500 | 服务器内部错误 |

### 1.2 鉴权

- 所有 API 需 `Authorization: Bearer <token>` 请求头
- `admin` 路由需 `users.role = "admin"`
- `enterprise` 路由需 `enterprise_members.role = "enterprise_admin"` 且 `enterprise.status = "active"`
- `member` 路由需 `enterprise_members.status = "active"`

### 1.3 分页格式

**请求:**
```
GET /api/xxx?page=1&page_size=20
```

**响应:**
```json
{
  "items": [...],
  "total": 100,
  "page": 1,
  "page_size": 20,
  "pages": 5
}
```

---

## 二、平台运营方 API (admin)

> 前缀: `/api/admin` | 权限: `role=admin`

### 2.1 企业管理

#### 2.1.1 获取企业列表

```
GET /api/admin/enterprises
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|:--:|------|
| page | int | 否 | 页码，默认 1 |
| page_size | int | 否 | 每页条数，默认 20，最大 100 |
| search | string | 否 | 企业名称模糊搜索 |
| status | string | 否 | active / disabled |

**响应 `data`:**
```json
{
  "items": [
    {
      "id": 1,
      "name": "XX科技有限公司",
      "short_name": "XX科技",
      "credit_code": "91110108XXXXXXXXXX",
      "address": "北京市海淀区...",
      "scale": "medium",
      "industry": "internet",
      "parent_id": 0,
      "status": "active",
      "contact_name": "张三",
      "contact_phone": "13800138000",
      "contact_email": "admin@example.com",
      "notes": "",
      "balance": "5000.00000000",
      "total_recharged": "10000.00000000",
      "admin_user_id": 42,
      "admin_email": "admin@example.com",
      "member_count": 15,
      "created_at": "2026-06-01T10:00:00Z",
      "updated_at": "2026-06-15T08:30:00Z"
    }
  ],
  "total": 50
}
```

| 字段 | 说明 |
|------|------|
| scale | 枚举: micro / small / medium / large |
| industry | 枚举: internet / finance / education / healthcare / manufacturing / other |
| balance | 企业独立余额（与企业管理员个人余额分离） |
| total_recharged | 企业累计充值总额 |
| member_count | 当前活跃成员数（聚合查询） |

---

#### 2.1.2 创建企业

```
POST /api/admin/enterprises
```

**请求体:**
```json
{
  "name": "XX科技有限公司",
  "short_name": "XX科技",
  "credit_code": "91110108XXXXXXXXXX",
  "address": "北京市海淀区中关村...",
  "scale": "medium",
  "industry": "internet",
  "contact_name": "张三",
  "contact_phone": "13800138000",
  "contact_email": "admin@example.com",
  "admin_email": "admin@example.com",
  "admin_name": "张三",
  "notes": "",
  "parent_id": 0
}
```

| 字段 | 类型 | 必填 | 约束 |
|------|------|:--:|------|
| name | string(255) | ✅ | 全局唯一（含软删除后） |
| short_name | string(100) | 否 | |
| credit_code | string(50) | 否 | 统一社会信用代码 |
| address | string(500) | 否 | |
| scale | string(20) | 否 | micro/small/medium/large |
| industry | string(50) | 否 | internet/finance/education/... |
| contact_name | string(100) | 否 | |
| contact_phone | string(50) | 否 | |
| contact_email | string(255) | 否 | |
| admin_email | string(255) | ✅ | 管理员邮箱，用于创建或关联 users 记录 |
| admin_name | string(100) | 否 | |
| notes | text | 否 | |
| parent_id | int64 | 否 | 父企业 ID，0=顶级，默认 0 |

**业务规则:**
1. 企业名称唯一性校验（含已软删除企业）
2. 若 `admin_email` 已有 `users` 记录 → 关联为该企业管理员的 users 记录
3. 若 `admin_email` 不存在 → 自动注册 `users` + 发送激活邮件
4. 创建 `enterprises` 记录（`balance=0`, `total_recharged=0`, `status=active`）
5. 创建 `enterprise_members` 记录（`role=enterprise_admin`, `status=active`）
6. 管理员邮箱不能已被其他企业绑定为管理员 → 返回 409

**响应 `data`:**
```json
{
  "id": 1,
  "name": "XX科技有限公司",
  "status": "active",
  "balance": "0.00000000",
  "admin_user_id": 42,
  "created_at": "2026-06-29T10:00:00Z"
}
```

**错误码:**

| code | message |
|:--:|------|
| 409 | 企业名称已存在 |
| 409 | 该邮箱已被其他企业绑定为管理员 |
| 422 | 管理员邮箱格式无效 |

---

#### 2.1.3 获取企业详情

```
GET /api/admin/enterprises/:id
```

**响应 `data`:** 同列表项，额外包含：
```json
{
  "...": "...",
  "subscriptions": [
    {
      "id": 1,
      "plan_name": "Claude 企业月包",
      "group_name": "Claude 分组",
      "status": "active",
      "daily_usage_usd": "12.34567890",
      "monthly_usage_usd": "234.56789001",
      "starts_at": "2026-06-01T00:00:00Z",
      "expires_at": "2026-07-01T00:00:00Z"
    }
  ]
}
```

---

#### 2.1.4 编辑企业信息

```
PUT /api/admin/enterprises/:id
```

**请求体 (全部可选，partial update):**
```json
{
  "name": "XX科技（更新）",
  "short_name": "XX",
  "contact_name": "李四",
  "contact_phone": "13900139000",
  "notes": "备注内容"
}
```

**可编辑字段**: name / short_name / credit_code / address / scale / industry / contact_name / contact_phone / contact_email / notes

**响应 `data`:** 更新后的完整企业对象

**错误码:**

| code | message |
|:--:|------|
| 409 | 企业名称已存在 |
| 404 | 企业不存在 |

---

#### 2.1.5 启停企业

```
POST /api/admin/enterprises/:id/toggle
```

**请求体:** (无)

**业务规则:**
- `active` → `disabled`: 级联禁用该企业所有 API Key（`api_keys.status = "disabled"`），企业管理员登录被拒绝
- `disabled` → `active`: 恢复 Key 原状态，管理员恢复登录权限

**响应 `data`:**
```json
{
  "id": 1,
  "status": "disabled",
  "affected_keys": 10
}
```

| 字段 | 说明 |
|------|------|
| affected_keys | 被级联禁用/启用的 Key 数量 |

---

#### 2.1.6 删除企业 🆕

```
DELETE /api/admin/enterprises/:id
```

**请求体:** (无)

**前置条件:** 企业状态必须为 `disabled`（必须先停用）

**业务规则:**
1. 校验 `status != "disabled"` → 返回 422 "请先停用企业"
2. 级联解绑所有 `enterprise_members`（`status="unbound"`, `unbound_at=NOW()`）
3. 级联禁用所有企业分配 Key
4. 企业套餐标记 `expired`
5. `enterprises.deleted_at = NOW()`（软删除）
6. 企业余额不退（由运营方人工处理）

**响应 `data`:**
```json
{
  "id": 1,
  "deleted": true,
  "unbound_members": 15
}
```

**错误码:**

| code | message |
|:--:|------|
| 422 | 请先停用企业 |
| 404 | 企业不存在 |

---

### 2.2 企业成员管理

> 前缀: `/api/admin/enterprises/:enterprise_id`

#### 2.2.1 获取成员列表

```
GET /api/admin/enterprises/:enterprise_id/members
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|:--:|------|
| page | int | 否 | |
| page_size | int | 否 | |
| search | string | 否 | 姓名/邮箱模糊搜索 |
| status | string | 否 | active / unbound |
| role | string | 否 | enterprise_admin / enterprise_member |
| department_id | int64 | 否 | 按部门筛选 |

**响应 `data.items`:**
```json
{
  "id": 1,
  "user_id": 100,
  "name": "张三",
  "email": "zhangsan@example.com",
  "role": "enterprise_member",
  "status": "active",
  "department_id": 5,
  "department_name": "研发部",
  "concurrency": 10,
  "rpm_limit": 100,
  "notes": "",
  "joined_at": "2026-06-15T09:00:00Z",
  "unbound_at": null,
  "last_active_at": "2026-06-28T18:30:00Z"
}
```

---

#### 2.2.2 创建企业用户（代注册）

```
POST /api/admin/enterprises/:enterprise_id/members
```

**请求体:**
```json
{
  "email": "zhangsan@example.com",
  "username": "张三",
  "password": "Abc@123456",
  "department_id": 5,
  "concurrency": 10,
  "rpm_limit": 100
}
```

| 字段 | 类型 | 必填 | 约束 |
|------|------|:--:|------|
| email | string(255) | ✅ | 不能已被任何 Sub2API 账号使用 |
| username | string(100) | 否 | |
| password | string(100) | 否 | 不提供则系统生成随机密码 |
| department_id | int64 | 否 | 必须属于该企业 |
| concurrency | int | 否 | 默认 0（无限制） |
| rpm_limit | int | 否 | 默认 0（无限制） |

**业务规则:**
1. 校验企业 `status = "active"`
2. 校验邮箱未被任何 `users` 记录使用
3. 校验 `department_id` 属于本企业（若提供）
4. 创建 `users` 记录 → 创建 `enterprise_members` 记录
5. 发送激活邮件（含登录链接和初始密码）

**响应 `data`:**
```json
{
  "id": 1,
  "user_id": 100,
  "name": "张三",
  "email": "zhangsan@example.com",
  "status": "pending",
  "created_at": "2026-06-29T10:00:00Z"
}
```

**错误码:**

| code | message |
|:--:|------|
| 409 | 该邮箱已注册 |
| 422 | 企业状态异常 |
| 422 | 部门不属于本企业 |
| 422 | 用户已是其他企业成员（1:1 约束） |

---

#### 2.2.3 编辑成员信息 🆕

```
PUT /api/admin/enterprises/:enterprise_id/members/:member_id
```

**请求体 (全部可选):**
```json
{
  "username": "张三(更新)",
  "password": "NewPass@123",
  "department_id": 6,
  "concurrency": 20,
  "rpm_limit": 200,
  "notes": "核心开发"
}
```

| 字段 | 说明 |
|------|------|
| password | 修改后发送密码变更通知邮件 |

---

#### 2.2.4 解绑成员

```
DELETE /api/admin/enterprises/:enterprise_id/members/:member_id
```

**业务规则:**
1. `enterprise_members.status → "unbound"`, `unbound_at = NOW()`
2. 禁用分配给该成员的所有 Key
3. 历史消费数据保留

**响应 `data`:**
```json
{
  "id": 1,
  "status": "unbound",
  "unbound_at": "2026-06-29T10:00:00Z",
  "disabled_keys": 3
}
```

---

### 2.3 部门管理

> 前缀: `/api/admin/departments` | 所有操作需指定 `enterprise_id`

#### 2.3.1 获取部门树

```
GET /api/admin/departments?enterprise_id=1
```

**响应 `data`:**
```json
[
  {
    "id": 1,
    "parent_id": 0,
    "name": "研发中心",
    "order_num": 0,
    "leader": "王五",
    "phone": "13700137000",
    "email": "wangwu@example.com",
    "status": "active",
    "member_count": 12,
    "created_at": "2026-06-01T00:00:00Z",
    "children": [
      {
        "id": 2,
        "parent_id": 1,
        "name": "前端组",
        "order_num": 1,
        "leader": "赵六",
        "member_count": 5,
        "children": []
      }
    ]
  }
]
```

---

#### 2.3.2 创建部门

```
POST /api/admin/departments
```

**请求体:**
```json
{
  "enterprise_id": 1,
  "parent_id": 0,
  "name": "研发中心",
  "order_num": 0,
  "leader": "王五",
  "phone": "13700137000",
  "email": "wangwu@example.com"
}
```

| 字段 | 类型 | 必填 | 约束 |
|------|------|:--:|------|
| enterprise_id | int64 | ✅ | |
| parent_id | int64 | 否 | 默认 0（顶级部门），必须属于同一企业 |
| name | string(100) | ✅ | 同一企业下唯一 |
| order_num | int | 否 | 排序，默认 0 |

**错误码:**

| code | message |
|:--:|------|
| 409 | 部门名称已存在 |
| 422 | 父部门不属于本企业 |

---

#### 2.3.3 编辑部门

```
PUT /api/admin/departments/:id
```

请求体同创建（全部可选）。

---

#### 2.3.4 删除部门

```
DELETE /api/admin/departments/:id
```

**业务规则:**
1. 检查是否有子部门 → 有则返回 422
2. 检查是否有成员 → 有则返回 422
3. 软删除 `deleted_at = NOW()`

**响应 `data`:**
```json
{
  "deleted": true
}
```

**错误码:**

| code | message |
|:--:|------|
| 422 | 存在子部门，请先删除子部门 |
| 422 | 存在关联成员，请先转移成员 |

---

## 三、企业管理员 API (enterprise)

> 前缀: `/api/enterprise` | 权限: `enterprise_members.role = "enterprise_admin"` | 企业状态: `active`

### 3.1 成员管理

企业管理员视角的成员管理，仅限本企业。

#### 3.1.1 成员列表

```
GET /api/enterprise/members
```

参数和响应格式同 [2.2.1](#221-获取成员列表)，但 `enterprise_id` 从登录上下文自动获取。

#### 3.1.2 创建成员

```
POST /api/enterprise/members
```

同 [2.2.2](#222-创建企业用户代注册)，`enterprise_id` 从上下文获取。

#### 3.1.3 编辑成员

```
PUT /api/enterprise/members/:member_id
```

同 [2.2.3](#223-编辑成员信息-)。

#### 3.1.4 解绑成员

```
DELETE /api/enterprise/members/:member_id
```

同 [2.2.4](#224-解绑成员)。

---

### 3.2 密钥管理

#### 3.2.1 企业密钥列表

```
GET /api/enterprise/keys
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|:--:|------|
| page | int | 否 | |
| page_size | int | 否 | |
| status | string | 否 | active / disabled |
| assigned_to | int64 | 否 | 按分配成员筛选（enterprise_members.id） |
| bound_tool | string | 否 | 按绑定工具筛选 |

**响应 `data.items`:**
```json
{
  "id": 1,
  "name": "张三的Claude Key",
  "key": "sk-abc...xyz",
  "key_prefix": "sk-abc***",
  "status": "active",
  "assigned_to": 5,
  "assigned_member_name": "张三",
  "assigned_member_email": "zhangsan@example.com",
  "groups": [
    { "id": 1, "name": "OpenAI 分组" },
    { "id": 2, "name": "Claude 分组" }
  ],
  "usage_purpose": "前端开发用",
  "bound_tool": "cursor",
  "quota": "100.00000000",
  "quota_used": "12.34567890",
  "expires_at": null,
  "created_at": "2026-06-15T09:00:00Z"
}
```

**数据范围:**
- 本企业管理员创建的所有 Key
- 包含 `assigned_to` 指向本企业成员的 Key
- 包含管理员自用 Key（`assigned_to = NULL`）

---

#### 3.2.2 创建密钥（含分配人员）

```
POST /api/enterprise/keys
```

**请求体:**
```json
{
  "name": "开发团队通用Key",
  "group_ids": [1, 2, 3],
  "assigned_to": 5,
  "usage_purpose": "前端项目开发",
  "bound_tool": "cursor",
  "quota": 100.0,
  "expires_at": "2026-12-31T23:59:59Z"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:--:|------|
| name | string(100) | ✅ | 密钥名称 |
| group_ids | []int64 | ✅ | 至少 1 个分组 ID（M:N 多选） |
| assigned_to | int64 | 否 | `enterprise_members.id`。NULL=管理员自用，有值=企业分配 |
| usage_purpose | string(200) | 否 | 用途说明（C1.5） |
| bound_tool | string(50) | 否 | cursor / trae / claude_code / codex / opencode / pixso / other |
| quota | float64 | 否 | Key 级配额（美元），不设则不限制 |
| expires_at | string | 否 | ISO 8601 时间，到期自动禁用 |

**业务规则:**
- `assigned_to` 对应的成员必须属于当前企业且 `status = "active"`
- `group_ids` 中的分组必须是有效分组
- 生成的 Key 格式: `sk-<random_string>`
- `assigned_to = NULL` → 管理员自用，消费走管理员个人池
- `assigned_to = member_id` → 企业分配，消费走企业资金池

**响应 `data`:**
```json
{
  "id": 10,
  "name": "开发团队通用Key",
  "key": "sk-a1b2c3d4e5f6...",
  "key_full": "sk-a1b2c3d4e5f6g7h8i9j0",
  "status": "active",
  "assigned_to": 5,
  "groups": [{ "id": 1, "name": "OpenAI" }, { "id": 2, "name": "Claude" }],
  "created_at": "2026-06-29T10:00:00Z"
}
```

> ⚠️ `key_full` 仅在创建响应中返回完整密钥，后续查询仅返回 `key_prefix`（脱敏）

**错误码:**

| code | message |
|:--:|------|
| 422 | 至少选择一个分组 |
| 422 | 指定的成员不属于本企业 |
| 404 | 分组不存在 |

---

#### 3.2.3 编辑密钥

```
PUT /api/enterprise/keys/:id
```

**请求体 (全部可选):**
```json
{
  "name": "更新后的名称",
  "group_ids": [1, 2],
  "assigned_to": 6,
  "usage_purpose": "新项目开发",
  "bound_tool": "trae",
  "quota": 200.0
}
```

**业务规则:**
- 只能编辑本企业下的密钥
- `assigned_to` 变更会改变计费归属（个人池 ↔ 企业池），写入审计日志

---

#### 3.2.4 删除密钥

```
DELETE /api/enterprise/keys/:id
```

软删除。已产生的消费记录保留。

---

#### 3.2.5 启停密钥

```
POST /api/enterprise/keys/:id/toggle
```

**请求体:** (无)

---

### 3.3 财务管理（企业资金池）

#### 3.3.1 企业财务汇总

```
GET /api/enterprise/finance
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|:--:|------|
| start_date | string | 否 | 用量统计起始日期，默认当月1日 |
| end_date | string | 否 | 用量统计截止日期，默认今日 |

**响应 `data`:**
```json
{
  "balance": "5000.00000000",
  "total_recharged": "10000.00000000",
  "subscriptions": [
    {
      "id": 1,
      "plan_name": "Claude 企业月包",
      "plan_id": 10,
      "group_name": "Claude 分组",
      "group_id": 2,
      "status": "active",
      "daily_usage_usd": "12.34567890",
      "weekly_usage_usd": "89.01234567",
      "monthly_usage_usd": "234.56789001",
      "starts_at": "2026-06-01T00:00:00Z",
      "expires_at": "2026-07-01T00:00:00Z"
    }
  ],
  "usage_summary": {
    "total_cost": "456.78901234",
    "total_calls": 12345,
    "by_member": [
      { "member_id": 5, "member_name": "张三", "cost": "100.00000000", "calls": 3000 },
      { "member_id": 6, "member_name": "李四", "cost": "200.00000000", "calls": 5000 }
    ],
    "by_model": [
      { "model": "gpt-4o", "cost": "150.00000000", "calls": 4000 },
      { "model": "claude-sonnet-4-20250514", "cost": "200.00000000", "calls": 5000 }
    ],
    "by_key": [
      { "key_id": 10, "key_name": "张三Key", "cost": "100.00000000", "calls": 3000 }
    ],
    "by_tool": [
      { "tool": "cursor", "cost": "150.00000000" },
      { "tool": "trae", "cost": "100.00000000" }
    ]
  }
}
```

**数据范围:** 仅统计 `usage_logs.pool_type = "enterprise"` 的记录。

---

#### 3.3.2 企业用量明细

```
GET /api/enterprise/usage
```

| 参数 | 类型 | 必填 | 说明 |
|------|------|:--:|------|
| page | int | 否 | |
| page_size | int | 否 | |
| start_date | string | 否 | |
| end_date | string | 否 | |
| member_id | int64 | 否 | 按成员筛选 |
| key_id | int64 | 否 | 按密钥筛选 |
| model | string | 否 | 按模型筛选 |
| bound_tool | string | 否 | 按工具筛选 |

**响应 `data.items`:**
```json
{
  "id": 1,
  "api_key_id": 10,
  "key_name": "张三Key",
  "user_id": 100,
  "user_name": "张三",
  "model": "claude-sonnet-4-20250514",
  "requested_model": "claude-sonnet-4-20250514",
  "total_cost": "0.01234567",
  "prompt_tokens": 500,
  "completion_tokens": 200,
  "pool_type": "enterprise",
  "created_at": "2026-06-29T10:30:00Z"
}
```

---

#### 3.3.3 企业充值 🆕

```
POST /api/enterprise/recharge
```

**请求体:**
```json
{
  "amount": 1000.00,
  "payment_method": "alipay"
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:--:|------|
| amount | float64 | ✅ | 充值金额（美元），必须 > 0 |
| payment_method | string | ✅ | alipay / wechat / stripe |

**业务规则:**
1. 复用现有支付系统，创建支付订单时标记 `order_type = "enterprise_recharge"` + `enterprise_id`
2. 支付回调成功 → `enterprises.balance += amount` + `enterprises.total_recharged += amount`
3. 支付回调失败 → 订单关闭，不影响企业余额

**响应 `data`:**
```json
{
  "order_id": "ORD20260629001",
  "amount": 1000.00,
  "payment_url": "https://pay.example.com/...",
  "qr_code": "data:image/png;base64,..."
}
```

**错误码:**

| code | message |
|:--:|------|
| 422 | 充值金额必须大于 0 |
| 422 | 不支持的支付方式 |

---

#### 3.3.4 企业购买套餐 🆕

```
POST /api/enterprise/subscribe
```

**请求体:**
```json
{
  "plan_id": 10,
  "group_id": 2
}
```

| 字段 | 类型 | 必填 | 说明 |
|------|------|:--:|------|
| plan_id | int64 | ✅ | 套餐计划 ID（subscription_plans.id） |
| group_id | int64 | ✅ | 分组 ID（本套餐适用的模型分组） |

**业务规则:**
1. 复用现有支付系统，标记 `order_type = "enterprise_subscription"` + `enterprise_id`
2. 支付成功 → 创建 `enterprise_subscriptions` 记录
3. 同一企业 + 同一分组可同时有多个有效套餐（叠加额度）

**响应 `data`:**
```json
{
  "id": 5,
  "enterprise_id": 1,
  "plan_name": "Claude 企业月包",
  "group_name": "Claude 分组",
  "starts_at": "2026-06-29T10:00:00Z",
  "expires_at": "2026-07-29T10:00:00Z",
  "status": "active"
}
```

---

### 3.4 部门管理

企业管理员视角，仅限操作本企业部门。接口同 [2.3 部门管理](#23-部门管理)，`enterprise_id` 从上下文获取。

```
GET    /api/enterprise/departments        # 部门树
POST   /api/enterprise/departments        # 创建部门
PUT    /api/enterprise/departments/:id    # 编辑部门
DELETE /api/enterprise/departments/:id    # 删除部门
```

---

### 3.5 企业设置

#### 3.5.1 获取企业设置

```
GET /api/enterprise/settings
```

**响应 `data`:**
```json
{
  "id": 1,
  "name": "XX科技有限公司",
  "short_name": "XX科技",
  "credit_code": "91110108XXXXXXXXXX",
  "address": "北京市海淀区...",
  "scale": "medium",
  "industry": "internet",
  "contact_name": "张三",
  "contact_phone": "13800138000",
  "contact_email": "admin@example.com",
  "notes": "",
  "status": "active",
  "balance": "5000.00000000",
  "member_count": 15,
  "created_at": "2026-06-01T10:00:00Z"
}
```

---

#### 3.5.2 更新企业设置

```
PUT /api/enterprise/settings
```

**请求体 (全部可选):**
```json
{
  "name": "XX科技（更新）",
  "short_name": "XX",
  "address": "新地址",
  "contact_name": "李四",
  "contact_phone": "13900139000",
  "contact_email": "new@example.com",
  "notes": "更新备注"
}
```

企业管理员只能编辑: name / short_name / address / contact_name / contact_phone / contact_email / notes

**不可编辑:** credit_code / scale / industry / status / balance / total_recharged（仅运营方可操作）

---

## 四、企业成员 API (member)

> 前缀: `/api/enterprise` | 权限: `enterprise_members.status = "active"`

### 4.1 企业 Profile

```
GET /api/enterprise/profile
```

**权限:** 企业管理员 + 企业普通成员

**响应 `data`:**
```json
{
  "enterprise": {
    "id": 1,
    "name": "XX科技有限公司",
    "short_name": "XX科技",
    "address": "北京市海淀区...",
    "scale": "medium",
    "industry": "internet",
    "contact_name": "张三",
    "contact_phone": "13800138000",
    "contact_email": "admin@example.com",
    "created_at": "2026-06-01T10:00:00Z"
  },
  "my_role": "enterprise_member",
  "my_department": "研发部",
  "my_joined_at": "2026-06-15T09:00:00Z",
  "monthly_usage": {
    "total_calls": 500,
    "total_cost": "12.34567890"
  }
}
```

> 管理员视角额外包含 `balance`、`total_recharged` 和 `subscriptions` 摘要。

### 4.2 部门列表（只读）

```
GET /api/enterprise/departments
```

返回本企业部门树（不含 leader/phone/email 联系人信息）。

---

## 五、个人用户 API 改造 (user)

> 前缀: `/api/user` | 现有 API，仅标注变更部分

### 5.1 个人密钥列表 ✏️ 改造

```
GET /api/user/keys
```

**变更:**
- 增加查询条件: `assigned_to` 指向当前用户的 `enterprise_members` 记录
- 响应增加字段:
  ```json
  {
    "id": 10,
    "name": "企业分配Key",
    "key_prefix": "sk-abc***",
    "source": "enterprise",           // 🆕 "personal" | "enterprise"
    "enterprise_name": "XX科技",       // 🆕 企业名称（若 source=enterprise）
    "assigned_member_name": "张三",    // 🆕 分配给的成员名
    "groups": [...],
    "usage_purpose": "前端开发",
    "bound_tool": "cursor",
    "status": "active",
    "created_at": "..."
  }
  ```
- 个人创建 Key → `source = "personal"`
- 企业分配 Key → `source = "enterprise"`, 不可编辑/删除

### 5.2 创建个人密钥 ✏️ 改造

```
POST /api/user/keys
```

**变更:** `group_id`（单选）→ `group_ids`（多选）

**请求体:**
```json
{
  "name": "我的开发Key",
  "group_ids": [1, 2],        // 🆕 多选分组（原 group_id: int）
  "usage_purpose": "个人项目",
  "bound_tool": "cursor"
}
```

> 个人用户创建的 Key，`assigned_to` 始终为 NULL（不可分配人员）。

### 5.3 个人用量 ✏️ 改造

```
GET /api/user/usage
```

**变更:**
- 查询范围增加: 企业分配给当前用户的 Key 产生的消费
- 响应增加 `pool_type` 字段标注每笔消费来源

---

## 六、通用错误码

| code | message | 触发场景 |
|:--:|------|------|
| 403 | 企业已停用 | 企业管理员或成员在企业被停用时访问企业 API |
| 403 | 需要企业管理员权限 | 非管理员访问 enterprise 路由 |
| 403 | 需要平台管理员权限 | 非 admin 访问 admin 路由 |
| 409 | 企业名称已存在 | 创建/更新企业时 |
| 409 | 该邮箱已注册 | B1 代注册时 |
| 409 | 部门名称已存在 | 创建部门时 |
| 422 | 用户已是其他企业成员 | 1:1 约束触发 |
| 422 | 指定的成员不属于本企业 | assigned_to 校验失败 |
| 422 | 密钥未授权使用该模型 | 分组路由匹配失败 |
| 422 | 企业额度不足 | 企业资金池消耗完 |
| 422 | 请先停用企业 | 删除未停用的企业 |
| 422 | 存在子部门 | 删除有子部门的部门 |
| 422 | 存在关联成员 | 删除有成员的部门 |

---

## 附录 A: 枚举值定义

| 枚举字段 | 可选值 |
|------|------|
| enterprise.scale | `micro` / `small` / `medium` / `large` |
| enterprise.industry | `internet` / `finance` / `education` / `healthcare` / `manufacturing` / `other` |
| enterprise.status | `active` / `disabled` |
| enterprise_member.role | `enterprise_admin` / `enterprise_member` |
| enterprise_member.status | `active` / `unbound` |
| enterprise_subscription.status | `active` / `expired` / `suspended` |
| api_key.bound_tool | `cursor` / `trae` / `claude_code` / `codex` / `opencode` / `pixso` / `other` |
| usage_log.pool_type | `personal` / `enterprise` |
| api_key.status | `active` / `disabled` |

## 附录 B: 接口变更影响范围

| 影响类型 | 涉及接口 | 说明 |
|:--:|------|------|
| 🆕 新增 | 全部 enterprise / member 路由 | 企业功能全新开发 |
| ✏️ 改造 | `POST /api/user/keys` | `group_id` → `group_ids`（单选→多选） |
| ✏️ 改造 | `GET /api/user/keys` | 增加企业来源 Key |
| ✏️ 改造 | `GET /api/user/usage` | 增加 pool_type 标注 |
| ✅ 不变 | 其他 `/api/user/*`、`/api/admin/users` 等 | 完全不变 |
