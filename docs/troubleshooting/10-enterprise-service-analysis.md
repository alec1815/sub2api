# 10 — 企业功能 Service 层分析记录

> **阶段**: P3 Service 业务逻辑层
> **日期**: 2026-06-29
> **范围**: 6 个 Service 模块、31+ 方法、架构模式、业务逻辑、依赖关系
> **相关**: 09-ent-service-pitfalls.md (排错)、07-编码实现分析.md (全栈分析)

---

## 一、模块总览

| # | Service | 文件 | 方法数 | 核心职责 |
|:--:|------|------|:--:|------|
| 1 | `EnterpriseService` | `enterprise_service.go` ✏️ | 7 | 企业 CRUD + 启停 + 删除 + 余额查询 |
| 2 | `EnterpriseMemberService` | `enterprise_member_service.go` 🆕 | 6 | 成员代注册/编辑/解绑/1:1约束 |
| 3 | `DepartmentService` | `department_service.go` ✏️ | 6 | 部门树形 CRUD + 删除前置校验 |
| 4 | `EnterpriseKeyService` | `enterprise_key_service.go` 🆕 | 4 | 企业密钥分配/列表/删除 |
| 5 | `EnterpriseBillingService` | `enterprise_billing_service.go` 🆕 | 3 | 财务汇总/套餐/充值校验 |
| 6 | `EnterpriseProfileService` | `enterprise_profile_service.go` 🆕 | 3 | Profile 多视角查询 |

---

## 二、架构模式分析

### 2.1 接口模式

每个 Service 定义自己需要的 Repository 接口（仅声明用到的子集）：

```
EnterpriseService
├── depends on → EnterpriseRepository (6 methods used)
├── depends on → EnterpriseMemberRepository (2 methods used)
├── depends on → EnterpriseSubscriptionRepository (1 method used)
├── depends on → UserRepository (2 methods used)
└── depends on → DepartmentRepository (0 → 构造接受但未使用)
```

**分析**：DepartmentRepository 被 `enterpriseService` 结构体持有但未被任何方法调用。
- 可能的原因：设计文档预期企业创建时可同时创建默认部门，但 P3 未实现
- 建议：P4 中实现或删除该依赖

### 2.2 错误定义模式

```go
// 继承 infraerrors 的 HTTP 状态码体系
var (
    ErrEnterpriseNameRequired    = infraerrors.BadRequest("ent.name_required", ...)
    ErrEnterpriseNotFound        = infraerrors.NotFound("ent.not_found", ...)
    ErrEnterpriseAlreadyDisabled = infraerrors.BadRequest("ent.already_disabled", ...)
    ErrDeleteRequiresDisabled    = infraerrors.BadRequest("ent.delete_requires_disabled", ...)
)
```

**分析**：
- 优点：错误码可国际化 (i18n key)，HTTP 状态码明确
- 风险：`infraerrors.BadRequest(...)` 作为 `var` 在 init 时求值，如果国际化文件未加载会得到空字符串
- 当前状态：项目已有 i18n 体系，安全

### 2.3 Repository 接口集中定义

所有 Repository 接口定义在对应的 Service 文件中：

| 接口 | 定义位置 | 引用者 |
|:--|:--|:--|
| `EnterpriseRepository` | `enterprise_service.go` | EnterpriseService, EnterpriseMemberService, DepartmentService, EnterpriseKeyService, EnterpriseBillingService, EnterpriseProfileService |
| `EnterpriseMemberRepository` | `enterprise_service.go` | EnterpriseService, EnterpriseMemberService, EnterpriseKeyService, EnterpriseBillingService, EnterpriseProfileService |
| `EnterpriseSubscriptionRepository` | `enterprise_service.go` | EnterpriseService, EnterpriseBillingService |
| `APIKeyGroupRepository` | `enterprise_service.go` | EnterpriseKeyService |
| `DepartmentRepository` | `department_service.go` | EnterpriseMemberService, EnterpriseProfileService |
| `UserRepository` | `user_service.go`（已有） | EnterpriseService, EnterpriseMemberService |
| `APIKeyRepository` | `api_key_service.go`（已有） | EnterpriseKeyService |

**分析**：接口定义分散但有规律：
- 核心企业接口（Enterprise/Member/Sub）集中在 `enterprise_service.go`
- 其他接口在各自 Service 文件中
- Go 的 package 内可见性使得跨文件引用无需 import

**潜在问题**：新增 Service 需要 `enterprise_service.go` 新增接口方法时，所有实现类都需要更新。Go 编译器会强制检查，不会遗漏。

---

## 三、关键业务逻辑分析

### 3.1 EnterpriseService

#### CreateEnterprise 流程
```
输入: name, admin_email, contact_phone?, contact_email?
  ↓
校验: name 非空, admin_email 非空
  ↓
创建用户: GetOrCreateUser(admin_email) → userID
  ↓
1:1 校验: GetMemberByUserID(userID) → 已有企业? → ErrAdminAlreadyMember
  ↓
创建企业: entRepo.Create(enterprise) → enterpriseID
  ↓
创建成员: memberRepo.Create(enterpriseID, userID, role=enterprise_admin, dept_id=0)
  ↓
输出: Enterprise + EnterpriseMember
```

**分析**：
- ✅ 管理员自动注册（PRD 需求）
- ✅ 1:1 约束在创建前校验（业务友好的错误消息）
- ⚠️ 事务问题：创建企业和创建成员是两个独立的 Repo 调用，不在一个数据库事务中。如果 CreateMember 失败，企业记录将孤立。
- **建议 P6**：包裹在 `ent.WithTx` 事务中

#### DeleteEnterprise 流程
```
输入: enterpriseID
  ↓
获取企业: GetByID → not found? → ErrEnterpriseNotFound
  ↓
前置校验: status != disabled? → ErrDeleteRequiresDisabled
  ↓
级联解绑: ListMembers(active) → for each → UnbindMember
  ↓
套餐标记: GetActiveSubscriptions → for each → UpdateStatus(expired)
  ↓
软删除: SoftDelete(enterpriseID)
  ↓
输出: success
```

**分析**：
- ✅ 必须先停用再删除（防止误删）
- ✅ 级联解绑所有成员（不残留数据）
- ⚠️ 未级联禁用 Key（标注 TODO P4）
- ⚠️ 无事务保护（同上）
- **建议 P6**：停用 → 删除两步操作应在一个 `ent.WithTx` 中

### 3.2 EnterpriseMemberService

#### CreateMember 流程
```
输入: enterpriseID, email, name?, department_id?
  ↓
企业状态校验: GetByID → status != active? → ErrEnterpriseNotActive
  ↓
1:1 约束校验: GetMemberByUserID(userID) → 已有企业? → ErrUserAlreadyMember
  ↓
部门校验: department_id > 0 → GetDepartment(department_id) → enterprise_id != 当前? → ErrDepartmentNotBelongToEnterprise
  ↓
创建/复用 users: GetOrCreateUser(email, name)
  ↓
创建成员: memberRepo.Create(enterpriseID, userID, dept_id, role=enterprise_member)
  ↓
输出: EnterpriseMember
```

**分析**：
- ✅ 三层校验（企业状态 → 1:1 → 部门归属）
- ✅ 代注册用户（不需要用户先自行注册）
- ✅ 部门校验确保跨企业越权
- ⚠️ 同上，无事务保护

#### UnbindMember 流程
```
输入: memberID, operatorUserID
  ↓
获取成员信息: GetByID(memberID)
  ↓
自解绑保护: member.Role == enterprise_admin && member.UserID == operatorUserID → ErrAdminCannotUnbindSelf
  ↓
解绑: memberRepo.Unbind(memberID) → status=unbound
  ↓
输出: success
```

**分析**：
- ✅ 企业管理员不能解绑自己（防止企业无管理员）
- ⚠️ 未级联禁用 Key（标注 TODO P4）
- ⚠️ 其他管理员成员解绑管理员未做保护（应至少保留 1 个管理员？设计文档中未明确）

### 3.3 DepartmentService

#### DeleteDepartment 流程
```
输入: departmentID
  ↓
获取部门: GetByID
  ↓
子部门检查: HasChildren(departmentID) → true? → ErrDepartmentHasChildren
  ↓
成员检查: HasMembers(departmentID) → true? → ErrDepartmentHasMembers
  ↓
软删除: SoftDelete(departmentID)
  ↓
输出: success
```

**分析**：
- ✅ 双前置检查（子部门 + 成员）
- ✅ 业务错误码会告知具体原因
- ⚠️ 无事务保护

#### Tree 查询
```
输入: enterpriseID
  ↓
全量查询: ListByEnterprise(enterpriseID) → []Department
  ↓
内存构造: parent_id → children map → 递归拼接成树
```

**分析**：
- ✅ 单次查询获取全部（非 N+1）
- ✅ 内存中构建树，复杂度 O(n)
- 适配场景：部门数通常 < 100，全量查询无性能问题

### 3.4 EnterpriseKeyService

#### CreateEnterpriseKey 流程
```
输入: memberID (当前用户), targetMemberID (分配给谁), key_name, tool, purpose, groupIDs?
  ↓
管理员校验: GetMemberByUserID(operator) → role != enterprise_admin? → ErrNotEnterpriseAdmin
  ↓
目标成员校验: GetByID(targetMemberID) → enterprise_id != 操作者? → ErrMemberNotInEnterprise
  ↓
工具校验: BoundTool 必须是白名单中的枚举值
  ↓
用途校验: len(UsagePurpose) ≤ 200
  ↓
创建 Key: apiKeyRepo.CreateKey(assignedTo=targetMemberID, usagePurpose, boundTool)
  ↓
分组关联: for each groupID → keyGroupRepo.SetGroups(keyID, groupIDs)
```

**分析**：
- ✅ 操作者权限 + 目标成员归属双重校验
- ✅ 工具枚举白名单（防止非法工具名）
- ⚠️ 未校验目标成员是否 active（解绑的成员不应被分配 Key）
- ⚠️ 分组成员未校验所属企业（跨企业分组越权可能）

#### ListEnterpriseKeys — 当前实现
```
当前: 逐用户查询
GetMembers(enterpriseID) → for each member → GetAPIKeys(userID) → 合并
```

**分析**：
- ⚠️ N+1 查询问题：100 个成员 = 101 次数据库查询
- **规划和 P4 改为 UNION**：
  ```sql
  SELECT * FROM api_keys WHERE user_id = ? AND assigned_to IS NULL
  UNION ALL
  SELECT * FROM api_keys WHERE assigned_to IN (SELECT id FROM enterprise_members WHERE enterprise_id = ?)
  ```

### 3.5 EnterpriseBillingService

#### GetFinanceOverview 流程
```
输入: enterpriseID
  ↓
获取企业: GetByID → balance, total_recharged
  ↓
获取套餐: GetActiveSubscriptions → 最近套餐的 plan_id, expires_at
  ↓
获取成员数: ListMembers → count
  ↓
输出: EnterpriseFinance{Balance, TotalRecharged, ActivePlan, MemberCount, MonthlyUsage=0}
```

**分析**：
- ✅ 余额/充值总额/套餐/成员数直接查询
- ⚠️ 月度用量 = 0（P4 占位）
- ⚠️ 最近套餐只取第一条（多套餐场景不准确）

### 3.6 EnterpriseProfileService

#### GetProfile — 三个视角
```
GetProfile(userID):
  → GetMemberByUserID → 获取用户的企业成员身份
  → 无身份 → ErrNotEnterpriseMember
  → GetEnterprise + GetMember + GetDepartment → Profile 组装

GetProfileByEnterprise(enterpriseID):
  → GetByID → 获取企业基本信息
  → 无成员上下文 → Profile（无 role/dept 信息）

GetProfileWithMember(memberID):
  → GetByID → GetMember → GetEnterprise → GetDepartment
  → 从成员 ID 反查完整 Profile
```

**分析**：
- ✅ 三个入口覆盖：用户视角 / 企业视角 / 特定成员视角
- ⚠️ 月度用量 = 0（三个方法统一占位）
- ⚠️ GetProfileByEnterprise 未校验调用者是否有权限查看

---

## 四、依赖关系图

```
enterprise_service.go (定义核心接口)
  ├── EnterpriseRepository interface         ← 5 个 Service 引用
  ├── EnterpriseMemberRepository interface   ← 5 个 Service 引用
  ├── EnterpriseSubscriptionRepository interface ← 2 个 Service 引用
  ├── APIKeyGroupRepository interface        ← 1 个 Service 引用
  └── EnterpriseService struct               ← 唯一实现者

enterprise_member_service.go
  └── EnterpriseMemberService ← 依赖 ENT + Member + User + Dept

department_service.go
  ├── DepartmentRepository interface
  └── DepartmentService ← 依赖 Dept + ENT

enterprise_key_service.go
  └── EnterpriseKeyService ← 依赖 APIKey + Member + ENT + KeyGroup

enterprise_billing_service.go
  └── EnterpriseBillingService ← 依赖 ENT + Sub + Member

enterprise_profile_service.go
  └── EnterpriseProfileService ← 依赖 ENT + Member + Dept
```

---

## 五、与设计文档对照

| 设计模块 (04-概要设计) | Service 方法 | 对照结果 |
|:--|:--|:--:|
| A1 企业创建 | `EnterpriseService.CreateEnterprise` | ✅ 含管理员自动注册 + 1:1 校验 |
| A2 企业编辑 | `EnterpriseService.UpdateEnterprise` | ✅ 可选字段指针 |
| A3 企业启停 | `EnterpriseService.ToggleStatus` | ✅ 级联 Key 标注 P4 |
| A3.5 企业删除 | `EnterpriseService.DeleteEnterprise` | ✅ 前置 disabled + 级联解绑 |
| B1 成员创建 | `EnterpriseMemberService.CreateMember` | ✅ 代注册 + 1:1 + 部门校验 |
| B2 成员列表 | `EnterpriseMemberService.ListMembers` | ✅ 分页 |
| B4 成员解绑 | `EnterpriseMemberService.UnbindMember` | ✅ 含自解绑拦截 |
| B5 成员编辑 | `EnterpriseMemberService.UpdateMember` | ✅ 含部门校验 |
| C1 企业 Key | `EnterpriseKeyService.CreateKey` | ⚠️ 工具枚举白名单完整，但 UNION 待 P4 |
| C2 Key 列表 | `EnterpriseKeyService.ListKeys` | ⚠️ N+1 查询，待 P4 UNION 优化 |
| D1 财务汇总 | `EnterpriseBillingService.GetFinanceOverview` | ✅ 余额/充值/套餐/成员，月度0占位 |
| E1 套餐管理 | `EnterpriseBillingService.GetSubscriptions` | ✅ 直接查询 |
| F1 部门树 | `DepartmentService.GetTree` | ✅ 全量查询 + 内存构建 |
| F2 部门 CRUD | `DepartmentService.Create/Update/Delete` | ✅ 含循环引用 + 子节点 + 成员前置校验 |
| G1 Profile | `EnterpriseProfileService` | ✅ 三视角查询，月度0占位 |

**遗漏项**：
1. ⚠️ **计费分流**：P4 网关层实现，P3 未涉及
2. ⚠️ **充值/支付**：复用现有支付系统，P3 只做 `ValidateRecharge` 校验
3. ⚠️ **通知/告警**：余额不足通知等，PRD 提到但 P3 未实现

---

## 六、跨阶段 TODO 追踪

| 编号 | 位置 | 内容 | 阶段 | 影响 |
|:--:|:--|------|:--:|:--|
| T1 | `EnterpriseKeyService` | UNION 查询替代 N+1 | P4 | 性能 |
| T2 | `EnterpriseProfileService` ×3 | 月度用量真实聚合 | P4 | 功能 |
| T3 | `EnterpriseMemberService.Unbind` | 级联禁用 Key | P4 | 安全 |
| T4 | `EnterpriseService.ToggleStatus` | 级联禁用 Key | P4 | 安全 |
| T5 | `EnterpriseService.CreateEnterprise` | 事务包裹 | P6 | 数据一致性 |
| T6 | `EnterpriseService.DeleteEnterprise` | 事务包裹 | P6 | 数据一致性 |
| T7 | `EnterpriseKeyService.CreateKey` | 目标成员 active 校验 | P6 | 安全 |

---

## 七、编码习惯总结

| 习惯 | 说明 | 评价 |
|:--|------|:--:|
| Repository 接口在 Service 层定义 | 依赖反转，Service 控制接口 | ✅ 标准 Go |
| 错误用 infraerrors + i18n key | HTTP 状态码 + 国际化的错误体系 | ✅ 一致 |
| Request 可选字段用指针 | `*string` / `*int64` 区分"未传"和"零值" | ✅ 清晰 |
| 树形结构单次全量查 | 避免 N+1，内存 O(n) 构建 | ✅ 高效 |
| 跨层 TODO 标记阶段 | `TODO(P4): xxx` 格式区分 | ✅ 可追踪 |
| 部分唯一索引 | PGSQL 层约束 + Service 层友好错误 | ✅ 双重防线 |

---

> **下一步**: [P4 中间件 + 网关计费路由](../enterprise/v1/06-开发计划任务.md#p4中间件--网关计费路由)
