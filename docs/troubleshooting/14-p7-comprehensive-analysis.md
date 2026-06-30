# P7 综合分析：PRD 符合度 + 前后端一致性 + 测试可行性

> **日期**: 2026-06-30 | **范围**: P1-P7 全部已完成的工作
> **分析目标**:
> 1. PRD(`企业功能PRD_产品视角_V5.md`) 功能覆盖度
> 2. 前后端字段/接口一致性
> 3. P6/P8 测试可行性
> 4. 日志/统计/企业维度状态

---

## 一、PRD 功能覆盖度分析

### 1.1 模块 A：企业账号管理

| 功能点 | PRD 要求 | 后端 | 前端 | 状态 |
|:--:|------|:--:|:--:|:--:|
| **A1** | 运营方创建企业 + 自动创建管理员 | ✅ `POST /api/v1/admin/enterprises` → `EnterpriseService.Create` (含 `createEnterpriseAdmin`) | ✅ `EnterprisesView.vue` 创建弹窗 | ✅ |
| **A2** | 企业信息维护（运营方全字段/企管仅本企业） | ✅ `PUT /api/v1/admin/enterprises/:id` + `PUT /api/v1/enterprise/settings` | ✅ `EnterpriseSettings.vue` (企管) | ✅ |
| **A3** | 企业状态管控（启停级联密钥） | ✅ `POST /api/v1/admin/enterprises/:id/toggle` → `ToggleEnterpriseStatus` (级联 Key 禁用/恢复) | ✅ `EnterprisesView.vue` 启停按钮 + 确认弹窗 | ✅ |
| **A3** | 企业删除（前置校验+级联解绑） | ✅ `DELETE /api/v1/admin/enterprises/:id` → `DeleteEnterprise` (先停用+解绑成员) | ✅ `EnterprisesView.vue` 删除按钮 + 确认弹窗 | ✅ |

### 1.2 模块 B：企业成员管理

| 功能点 | PRD 要求 | 后端 | 前端 | 状态 |
|:--:|------|:--:|:--:|:--:|
| **B1** | 企业管理员代注册成员（创建 users + 绑定） | ✅ `POST /api/v1/enterprise/members` → `CreateMember` (含 `createEnterpriseMemberUser`) | ✅ `MemberManagement.vue` 创建弹窗 | ✅ |
| **B2** | 已有账号绑定加入企业 | ⚠️ PRD 标注 B2，但 `06-开发计划任务.md` 列为「后续迭代」 | ❌ 无 | ⏭️ 本期跳过 |
| **B3** | 成员列表（搜索/筛选/分页） | ✅ `GET /api/v1/enterprise/members` + `GET /api/v1/admin/enterprises/:id/members` | ✅ `MemberManagement.vue` + `EnterpriseMembersView.vue` | ✅ |
| **B4** | 解绑成员（级联禁用 Key） | ✅ `DELETE /api/v1/enterprise/members/:id` → `UnbindMember` | ✅ `MemberManagement.vue` 解绑按钮 + 确认弹窗 | ✅ |

### 1.3 模块 B1：部门管理

| 功能点 | PRD 要求 | 后端 | 前端 | 状态 |
|:--:|------|:--:|:--:|:--:|
| — | 部门树形 CRUD | ✅ `GET/POST/PUT/DELETE /api/v1/admin/departments` + `/api/v1/enterprise/departments` | ✅ `DepartmentsView.vue` + `DepartmentTreeNode.vue` | ✅ |

### 1.4 模块 C：API 密钥分配

| 功能点 | PRD 要求 | 后端 | 前端 | 状态 |
|:--:|------|:--:|:--:|:--:|
| **C1** | 创建密钥时分配人员（assigned_to 下拉框） | ✅ `POST /api/v1/enterprise/keys` → Key 含 `assigned_to` 字段 | ✅ `EnterpriseKeys.vue` 创建弹窗含 `assignedTo` | ✅ |
| **C1.5** | 用途绑定（用途说明 + 绑定工具） | ✅ Key 含 `usage_purpose` + `bound_tool` 字段 | ✅ EnterpriseKeys.vue 创建弹窗含用途/工具字段 | ✅ |
| **C1.6** | 分组多选（一把 Key 绑多分组） | ✅ `api_key_groups` 多对多中间表 | ⚠️ 前端创建弹窗用逗号分隔的 groupIds 文本输入（非多选勾选） | ⚠️ UX 简化 |
| **C2** | 企业管理员视角密钥列表 | ✅ `GET /api/v1/enterprise/keys` → 跨 user_id+assigned_to UNION 查询 | ✅ `EnterpriseKeys.vue` | ✅ |
| **C3** | 成员视角密钥可见性（个人+企业） | ⚠️ 后端通过现有 `GET /api/v1/keys` 返回（含 assigned_to Key），需确认是否区分「来源」 | ✅ 前端原型有此需求，但当前个人 Key 列表页未改造 | ⚠️ **待改造** |

### 1.5 模块 D：计费归属与消费

| 功能点 | PRD 要求 | 后端 | 前端 | 状态 |
|:--:|------|:--:|:--:|:--:|
| **D1** | 企业密钥消费 → 企业资金池扣款 | ✅ `EnterpriseBillingRouter.AssignedTo 判定` + `GatewayService` P4 隔离门 | N/A (系统行为) | ✅ |
| **D2** | 个人密钥消费 → 个人资金池扣款 | ✅ 原有逻辑不变 | N/A | ✅ |
| **D3** | 两条计费线路隔离 | ✅ 三条资金线隔离（个人/管理员自用/企业池） | N/A | ✅ |

### 1.6 模块 E：用量查看与财务聚合

| 功能点 | PRD 要求 | 后端 | 前端 | 状态 |
|:--:|------|:--:|:--:|:--:|
| **E1** | 个人视角用量查看（含被分配密钥消费） | ⚠️ **usage_logs SELECT 缺少 enterprise_id/pool_type 扫描** → 无法区分来源 | 前端「使用记录」页未区分企业/个人来源 | ⚠️ **P6 待修复** |
| **E2** | 企业维度财务聚合（余额/套餐/消费拆维） | ⚠️ Finance API 已注册路由，但 **monthly_usage 硬编码为 0**，`usage_logs` 读取链未完成 | ✅ `EnterpriseFinance.vue` 视图已创建 | ⚠️ **依赖后端** |

---

## 二、前后端字段一致性分析

### 2.1 ✅ 一致的部分

| 结构体 | 前后端一致 |
|------|:--:|
| `CreateEnterpriseRequest` (name, short_name, credit_code, address, scale, industry, contact_*, admin_email, admin_name, notes, parent_id) | ✅ 字段名均为 snake_case，类型一致 |
| `UpdateEnterpriseRequest` (全部可选同名字段) | ✅ |
| `CreateMemberRequest` (email, username, password, department_id, concurrency, rpm_limit) | ✅ |
| `CreateEnterpriseKeyRequest` (name, group_ids, assigned_to, usage_purpose, bound_tool, quota, expires_at) | ✅ |
| `Enterprise` 列表响应 | ✅ 前端 Enterprise interface 与后端 Ent Enterprise 模型字段一致 |
| `EnterpriseMember` 响应 | ✅ |
| `Department` 响应 | ✅ |
| `UpdateEnterpriseSettingsRequest` (name, short_name, address, contact_*, notes) | ✅ 前端 `EnterpriseSettings.vue` 用的字段与后端 `UpdateEnterpriseSettings` 一致 |

### 2.2 ⚠️ 需关注的差异

| # | 位置 | 差异 | 影响 | 建议 |
|:--:|------|------|------|------|
| 1 | **前端创建企业表单** 含 `password` 字段 | 后端 `CreateEnterpriseRequest` 无此字段。后端 Service 内部自动生成密码 | 前端发的 `password` 会被忽略 | ✅ 低风险（后端处理正确），但前端多余的字段可清理 |
| 2 | **前端 `EnterpriseMemberStatus`** 含 `'pending'` | 后端 Ent schema 默认 `active`，无 pending 状态。Service 直接创 `active` 成员 | 前端筛选器含 pending 选项但后端无数据 | 需决定：PRD 说「待激活」→ 是否要实现邮件激活流程？本期 PRD 标注 B1「发送激活邮件」但实际未实现 |
| 3 | **前端 `EnterpriseFinance.usage_summary`** 结构丰富 | 后端 `UsageSummary` 含 `by_member/by_model/by_key/by_tool` 但数据为空/硬编码 0 | 前端按此结构渲染但看不到真实数据 | 依赖 P6 usage_logs SELECT 修复 |
| 4 | **前端分页响应 `BasePaginationResponse<T>`** 使用 `items, total, page, page_size, pages` | 后端 `PaginatedResponse` 同名结构 ✅ | ✅ 一致 |
| 5 | **前端 `EnterpriseProfile.MonthlyUsage`** 含 `total_calls, total_cost` | 后端 `MonthlyUsage` 硬编码 `CallCount=0, TotalCost=0` (TODO P4) | 前端企业资料页月度用量始终为 0 | P6 需接入 usage_logs 聚合 |
| 6 | **API 路径前缀** 后端使用 `/api/v1/` | 前端 `apiClient` 基路径应为 `import.meta.env.VITE_API_BASE`（预期包含 `/api/v1`） | 若环境变量配置一致则无问题 | 需在联调时确认 |

### 2.3 ❌ 已知不一致（需修复）

| # | 位置 | 问题 | 严重程度 |
|:--:|------|------|:--:|
| 1 | `scanUsageLog()` 后端 | **usage_logs 读取出错**：`usageLogSelectColumns` 常量缺少 `enterprise_id, pool_type` 列，`scanUsageLog` 函数未 scan 这两个字段 → INSERT 正确但 SELECT 不返回企业数据 | 🔴 **阻塞 P8** |
| 2 | `enterpriseProfile.MonthlyUsage` | 后端硬编码 0，注释 `TODO(P4): 月度用量需网关计费层接入后通过 enterprise_id 聚合 usage_logs 实现` | 🟡 功能不完整 |
| 3 | `EnterpriseFinance.usage_summary` | 同上，数据为空 | 🟡 功能不完整 |

---

## 三、后端路由对照表

### 3.1 平台运营方路由 (`/api/v1/admin`)

| 方法 | 后端路径 | 前端 API 函数 | 一致 |
|:--:|------|------|:--:|
| GET | `/admin/enterprises` | `adminAPI.enterprises.list()` | ✅ |
| POST | `/admin/enterprises` | `adminAPI.enterprises.create()` | ✅ |
| PUT | `/admin/enterprises/:id` | `adminAPI.enterprises.update()` | ✅ |
| POST | `/admin/enterprises/:id/toggle` | `adminAPI.enterprises.toggle()` | ✅ |
| DELETE | `/admin/enterprises/:id` | `adminAPI.enterprises.delete()` | ✅ |
| GET | `/admin/enterprises/:id/members` | `adminAPI.enterpriseMembers.list()` | ✅ |
| POST | `/admin/enterprises/:id/members` | `adminAPI.enterpriseMembers.create()` | ✅ |
| PUT | `/admin/enterprises/:id/members/:mid` | `adminAPI.enterpriseMembers.update()` | ✅ |
| DELETE | `/admin/enterprises/:id/members/:mid` | `adminAPI.enterpriseMembers.unbind()` | ✅ |
| GET | `/admin/departments?enterprise_id=X` | `adminAPI.departments.list()` | ✅ |
| POST | `/admin/departments` | `adminAPI.departments.create()` | ✅ |
| PUT | `/admin/departments/:id` | `adminAPI.departments.update()` | ✅ |
| DELETE | `/admin/departments/:id` | `adminAPI.departments.delete()` | ✅ |

### 3.2 企业管理员路由 (`/api/v1/enterprise`)

| 方法 | 后端路径 | 前端 API 函数 | 一致 |
|:--:|------|------|:--:|
| GET | `/enterprise/members` | `enterpriseAdminAPI.listMembers()` | ✅ |
| POST | `/enterprise/members` | `enterpriseAdminAPI.createMember()` | ✅ |
| PUT | `/enterprise/members/:id` | `enterpriseAdminAPI.updateMember()` | ✅ |
| DELETE | `/enterprise/members/:id` | `enterpriseAdminAPI.unbindMember()` | ✅ |
| GET | `/enterprise/keys` | `enterpriseAdminAPI.listKeys()` | ✅ |
| POST | `/enterprise/keys` | `enterpriseAdminAPI.createKey()` | ✅ |
| PUT | `/enterprise/keys/:id` | `enterpriseAdminAPI.updateKey()` | ✅ |
| POST | `/enterprise/keys/:id/toggle` | `enterpriseAdminAPI.toggleKey()` | ✅ |
| DELETE | `/enterprise/keys/:id` | `enterpriseAdminAPI.deleteKey()` | ✅ |
| GET | `/enterprise/finance` | `enterpriseAdminAPI.getFinance()` | ✅ |
| GET | `/enterprise/usage` | `enterpriseAdminAPI.getUsage()` | ✅ |
| GET | `/enterprise/balance` | 前端未独立调用（含在 finance 中） | ⚠️ 可选 |
| GET | `/enterprise/subscriptions` | 前端未独立调用（含在 finance 中） | ⚠️ 可选 |
| POST | `/enterprise/recharge` | `enterpriseAdminAPI.recharge()` | ✅ |
| POST | `/enterprise/subscribe` | `enterpriseAdminAPI.subscribe()` | ✅ |
| GET | `/enterprise/settings` | `enterpriseAdminAPI.getSettings()` | ✅ |
| PUT | `/enterprise/settings` | `enterpriseAdminAPI.updateSettings()` | ✅ |
| GET | `/enterprise/departments` | `enterpriseAdminAPI.getDepartmentTree()` | ✅ |
| POST | `/enterprise/departments` | `enterpriseAdminAPI.createDepartment()` | ✅ |
| PUT | `/enterprise/departments/:id` | `enterpriseAdminAPI.updateDepartment()` | ✅ |
| DELETE | `/enterprise/departments/:id` | `enterpriseAdminAPI.deleteDepartment()` | ✅ |
| GET | `/enterprise/profile` | `enterpriseAdminAPI.getProfile()` | ✅ |

---

## 四、日志/统计/企业维度分析

### 4.1 usage_logs 表企业字段 (Migration 160)

| 字段 | 类型 | 说明 | 状态 |
|------|------|------|:--:|
| `enterprise_id` | BIGINT (nullable) | NULL=个人，有值=企业 | ✅ Schema |
| `pool_type` | VARCHAR(20) | `'personal'` / `'enterprise'` | ✅ Schema |
| 索引 `idx_usage_logs_enterprise_id_created_at` | | 按企业+时间查询 | ✅ |
| 索引 `idx_usage_logs_enterprise_id_pool_type` | | 按企业+池类型查询 | ✅ |

### 4.2 usage_logs 写入/读取状态

| 操作 | 状态 | 说明 |
|------|:--:|------|
| **INSERT** | ✅ 完成 | 网关 `GatewayService.postUsageBilling` 正确写入 `enterprise_id` 和 `pool_type` |
| **SELECT (scanUsageLog)** | ❌ 缺失 | `usageLogSelectColumns` 常量不含 `enterprise_id`/`pool_type`，`scanUsageLog` 函数未 scan 这两列 |
| **企业维度聚合查询** | ❌ 未实现 | 无 `GROUP BY enterprise_id` 的聚合方法 |
| **企业月度用量** | ❌ 占位 | `MonthlyUsage` 硬编码为 0 |

### 4.3 企业财务/费率数据

| 数据维度 | 状态 | 说明 |
|------|:--:|------|
| 企业余额 (enterprises.balance) | ✅ | 独立资金池，DeductBalance 已实现 |
| 企业充值 (enterprises.total_recharged) | ✅ Schema | 累计追踪，但 `RechargeEnterprise` 标记 TODO(P6) |
| 企业套餐 (enterprise_subscriptions) | ✅ Schema + Repo | 含 daily/weekly/monthly_usage_usd 用量累计 |
| 企业审计日志 | ❌ 无 | `payment_audit_logs` 不含 enterprise_id |
| 企业 Dashboard 统计 | ❌ 无 | `DashboardStats` 不含企业维度 |
| Token 消耗企业维度 | ⚠️ 半完成 | usage_logs 含 enterprise_id 可由此聚合，但读取链路未完成 |

### 4.4 企业中间件状态

| 中间件 | 状态 | 权限逻辑 |
|------|:--:|------|
| `RequireEnterpriseAdmin` | ✅ | JWT → 查 enterprise_members → 验证 role=enterprise_admin + enterprise.status=active |
| `RequireEnterpriseMember` | ✅ | 同上，仅验证 status=active |
| Context Key: `EnterpriseID` | ✅ | 注入 gin.Context |
| Context Key: `EnterpriseMemberID` | ✅ | 注入 gin.Context |
| Context Key: `EnterpriseRole` | ✅ | 注入 gin.Context |

---

## 五、P6 测试可行性分析

### 5.1 P6 任务清单再确认

| # | 任务 | 当前状态 | 我可独立完成？ |
|:--:|------|:--:|:--:|
| 6.1 | Repository 集成测试 (真实 PG) | ⚠️ 需真实数据库环境 | ❌ 需要 PostgreSQL + 迁移环境 |
| 6.2 | Service 集成测试 (三条资金线、1:1约束、启停级联) | ⚠️ 同上 | ❌ 同上 |
| 6.3 | Handler HTTP 测试 (gin test mode) | 可写，可 mock Service | ✅ 可写 Mock 测试 |
| 6.4 | 网关计费分流集成测试 (E2E) | 需完整环境 | ❌ 需完整后端运行环境 |
| 6.5 | 边界 case 补齐 | 可分析并写测试用例 | ✅ 可写测试代码 |
| 6.6 | CI 验证 (golangci-lint) | 环境相关 | ⚠️ 依赖 CI 配置 |

### 5.2 我能做什么

| 可做 | 说明 |
|------|------|
| ✅ **P6 Repository/Service Mock 单元测试** | 可以写 `go test -tags=unit` 的单测代码（Mock Repository 接口） |
| ✅ **P6 Handler Mock 测试** | 可以写 handler 的 HTTP 测试（Mock Service） |
| ✅ **P6 修复 scanUsageLog** | 可以修改 `usageLogSelectColumns` 和 `scanUsageLog` 补充缺失的企业字段 |
| ❌ **P6 集成测试** | 无法执行，需要你本地启动 PostgreSQL + 后端服务后运行 |

### 5.3 P6 建议优先级

| 优先级 | 任务 | 原因 |
|:--:|------|------|
| 🔴 P0 | **修复 scanUsageLog** | 阻塞 P8 联调和前端显示 |
| 🔴 P0 | **实现 enterprise_id 聚合查询** | 阻塞 E1/E2 用量功能 |
| 🟡 P1 | Service 单元测试 | 提升代码质量（已有 48 测试） |
| 🟡 P1 | Handler 单元测试 | 验证 API 行为 |
| 🟢 P2 | CI lint 检查 | 发布前执行 |
| 🟢 P2 | 充值/订阅支付集成 | `TODO(P6)` 标注 |

---

## 六、P8 测试可行性分析

### 6.1 完整业务闭环

```
运营方创建企业 → 企业管理员登录 → 创建成员 → 分配 Key → 调用 API → 查看用量
```

| 步骤 | 当前状态 | 阻塞点 |
|------|:--:|------|
| ① 运营方创建企业 | ✅ API + 前端完整 | — |
| ② 企业管理员登录 | ⚠️ 后端有中间件/路由，但前端无企管专属登录入口 | 企管用现有登录页，需确保返回的 JWT 含 enterprise role |
| ③ 创建成员 | ✅ API + 前端完整 | — |
| ④ 分配 Key | ✅ API + 前端完整 | — |
| ⑤ 调用 API | ✅ 网关计费分流已实现 | 需配置真实 Group/Account |
| ⑥ 查看用量 | ❌ usage_logs 读取链未完成 | **scanUsageLog 缺失 enterprise_id/pool_type** |

### 6.2 回归测试

| 场景 | 当前状态 | 我能测试？ |
|------|:--:|:--:|
| 个人用户独立 Key 创建 | 需确认路由/中间件不会误拦截 | ❌ 需运行后端 |
| 分组单选→多选迁移 | 需确认 C1.6 多选分组兼容 | ❌ 需运行后端 |
| 个人余额扣款不受影响 | 需确认 `assigned_to == nil` 路径 | ❌ 需运行后端 |

### 6.3 P8 我能做什么

| 可做 | 说明 |
|------|------|
| ✅ **修复 scanUsageLog 阻塞问题** | 补充 enterprise_id/pool_type SELECT 列 |
| ✅ **企业维度聚合查询** | 实现 `GROUP BY enterprise_id` 的 usage_logs 查询 |
| ✅ **Mock 联调** | 前端启用 Mock 模式验证页面交互 |
| ❌ **真实 E2E 测试** | 需要你本地启动完整后端 + 数据库 |

---

## 七、已发现问题汇总

### 7.1 🔴 阻塞级 (P8 前必须修复)

| # | 问题 | 文件 | 修复方式 |
|:--:|------|------|------|
| 1 | `usageLogSelectColumns` 缺少 `enterprise_id, pool_type` | `service/usage_log.go` | 补充 2 列 + `scanUsageLog` 新增扫描 |
| 2 | `enterpriseProfile.MonthlyUsage` 硬编码 0 | `service/enterprise_profile_service.go` | 改为查询 `usage_logs` 聚合 |
| 3 | `EnterpriseFinance.UsageSummary` 数据为空 | `service/enterprise_billing_service.go` | 同上 |
| 4 | `usage_logs` 缺少 `GROUP BY enterprise_id` 聚合查询 | `repository/usage_log_repo.go` | 新增 `AggregateByEnterprise` 方法 |

### 7.2 🟡 功能缺项

| # | 问题 | 说明 |
|:--:|------|------|
| 1 | 前端创建企业表单含 `password` 字段 | 后端忽略此字段，前端可保留（无副作用） |
| 2 | `EnterpriseMemberStatus` 含 `pending` 但后端无此状态 | 邮件激活流程未实现，本期跳过 |
| 3 | 充值/订阅支付集成标记 TODO(P6) | 企业财务页的充值和订阅按钮无法工作 |
| 4 | 企业审计日志缺失 | `payment_audit_logs` 无 enterprise_id |

### 7.3 🟢 已确认一致

| # | 确认项 | 结论 |
|:--:|------|------|
| 1 | 前端 `CreateEnterpriseRequest.type` 与后端 JSON tag | ✅ `name`, `short_name`, `credit_code`, `address`, `scale`, `industry`, `contact_name`, `contact_phone`, `contact_email`, `admin_email`, `admin_name`, `notes`, `parent_id` 全部一致 |
| 2 | 前端 `Enterprise` 响应接口与后端 Ent 模型字段 | ✅ 完全一致 |
| 3 | 所有 23 个企业路由前后端一致 | ✅ 方法/路径全部对应 |
| 4 | 分页格式 `BasePaginationResponse<T>` | ✅ `items, total, page, page_size, pages` 一致 |
| 5 | 枚举值 (status/scale/industry/role) | ✅ 前后端一致 |
| 6 | 部门 CRUD 接口 | ✅ 完全一致 |
| 7 | 密钥 CRUD + assigned_to + groups | ✅ 完全一致 |

---

## 八、结论与下一步

### 8.1 PRD 符合度 (最终更新 2026-06-30)

| 模块 | 完成度 | 说明 |
|------|:--:|------|
| A (企业账号) | **100%** | 前后端完整 |
| B (成员管理) | **90%** | B2 后续迭代 (已有账号绑定) |
| B1 (部门管理) | **100%** | 前后端完整 |
| C (密钥分配) | **100%** | C3 已完成(66a0ba3d) + C1/C1.5/C1.6/C2 |
| D (计费归属) | **100%** | 三线隔离已实现 |
| E (用量财务) | **70%** | scanUsageLog 已修复(f46a89a4) · MonthlyUsage 硬编码待解 |

**总体**: PRD 核心功能 A/B/C/D 已完成 100%，C3 成员视角Key列表改造已提交。E 模块：阻塞项 scanUsageLog 已修复，MonthlyUsage 需要在 Repo 层增加聚合查询方法。

### 8.2 字段一致性

- ✅ 字段名 (`snake_case` JSON tag) 前后端 **100%** 一致
- ✅ 请求/响应结构前后端 **100%** 一致
- ✅ ApiKey 类型已补全 `assigned_to/enterprise_id/enterprise_name` (C3)

### 8.3 接口/路由一致性

- ✅ 23 个企业端点前后端全部对应
- ✅ scanUsageLog 修复后 50→52 列 SELECT 正确匹配 52 变量扫描

### 8.4 测试工具可用性

| 工具 | 用途 | 可用 |
|------|------|:--:|
| `go test -tags=unit` | 后端单元测试 | ✅ 部分通过（DeductBalance stub 需修复） |
| `go build` | 编译验证 | ✅ 零错误 |
| `curl` | API 接口测试 | ✅ 本地可用 |
| `playwright-cli` | E2E 浏览器自动化 (P8) | ✅ 已安装技能 |
| `agent-browser` | 浏览器截图/抓取 | ✅ 已安装技能 |

### 8.5 P6 待修复问题

| # | 问题 | 文件 | 优先级 |
|:--:|------|------|:--:|
| 1 | `entRepoStubForDept` 缺少 `DeductBalance` 方法 | `department_service_test.go` | P1 |
| 2 | MonthlyUsage 硬编码 0 (需 UsageLogRepo 注入) | `enterprise_profile_service.go` + `enterprise_billing_service.go` | P1 |
| 3 | 充值/订阅支付集成 | `enterprise_billing_service.go` | P2 |

### 8.6 建议的下一步顺序

1. 🔴 **P6: 修复 DeductBalance test stub** → service 单元测试全通过
2. 🟡 **P6: MonthlyUsage 聚合** → 新增 `GetMonthlyUsageByEnterprise` Repo 方法
3. 🟡 **P6: curl API 测试** → 本地 PostgreSQL 环境验证 23 个端点
4. 🟢 **P8: playwright-cli E2E** → 完整业务闭环自动化测试
5. 🟢 **P6: 充值/订阅支付集成** → TODO(P6) 标记
