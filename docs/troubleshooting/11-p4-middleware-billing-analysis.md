# 11. P4 中间件 + 网关计费路由分析

> 分析日期: 2026-06-29 | 版本: P4 Phase | 状态: ✅ 已完成

## 一、P4 目标概述

P4 实现"隔离门"——**三个资金池的计费路径隔离**：

| 资金池 | APIKey.AssignedTo | 扣费来源 | 实现方式 |
|--------|-------------------|----------|----------|
| 个人池 (personal) | `nil` → 非企业admin | `userRepo.DeductBalance` | 原有逻辑 |
| 管理员自用池 | `nil` → 企业admin | `userRepo.DeductBalance` | 原有逻辑 |
| 企业池 (enterprise) | `!= nil` | `enterpriseRepo.DeductBalance` | **新增** |

## 二、涉及文件

| 文件 | 类别 | 说明 |
|------|------|------|
| `backend/internal/server/middleware/enterprise_auth.go` | **新建** | 企业认证中间件 |
| `backend/internal/service/enterprise_billing_router.go` | **新建** | 计费来源判定 |
| `backend/internal/service/enterprise_service.go` | 修改 | +DeductBalance 接口 |
| `backend/internal/repository/enterprise_repo.go` | 修改 | +DeductBalance 实现 |
| `backend/internal/service/gateway_service.go` | 修改 | 计费管道集成 |
| `backend/internal/service/usage_log.go` | 修改 | +EnterpriseID/PoolType |
| `backend/internal/repository/usage_log_repo.go` | 修改 | +INSERT 字段 |

---

## 三、模块①：企业认证中间件

### 3.1 项目现有代码风格

项目中已有的中间件模式（以 `JWTAuth` 为例）：

```go
// 1. ContextKey 使用自定义类型 middleware.ContextKey
const ContextKeyJWT middleware.ContextKey = "jwt_subject"

// 2. 取值: value, exists := c.Get(string(key)); type, ok := value.(T)
func GetAuthSubjectFromContext(c *gin.Context) (*AuthSubject, bool) {
    value, exists := c.Get(string(ContextKeyJWTSubject))
    if !exists { return nil, false }
    subject, ok := value.(*AuthSubject)
    return subject, ok
}

// 3. 中止请求: middleware.AbortWithError(c, statusCode, code, message)
// 4. gin.HandlerFunc 模式：return func(c *gin.Context) { ... }
// 5. Set context: c.Set(string(key), value)
```

### 3.2 P4 前逻辑

**不存在企业认证中间件**。所有路由只依赖 `JWTAuth` + 可能的 `RoleAuth`（admin/user 级别），不考虑企业成员身份。

### 3.3 P4 新逻辑

```go
// 两个中间件，在 JWTAuth 之后使用

// RequireEnterpriseAdmin:
//   JWTAuth → GetAuthSubjectFromContext → memberRepo.GetByUserID(userID)
//   → 校验 role="enterprise_admin" + status="active"
//   → entRepo.GetByID(member.EnterpriseID)
//   → 校验 enterprise.status="active"
//   → c.Set(enterprise_id, enterprise_member_id, enterprise_role)

// RequireEnterpriseMember:
//   JWTAuth → GetAuthSubjectFromContext → memberRepo.GetByUserID(userID)
//   → 校验 status="active" (不校验 role)
//   → entRepo.GetByID → 校验 enterprise.status="active"
//   → c.Set(...)
```

### 3.4 变更对照 — 新增 3 个 ContextKey + 2 个 HandlerFunc + 4 个 Accessor

```
P4 前: 无企业上下文、无企业权限中间件
P4 后:
  + ContextKeyEnterpriseID       = "enterprise_id"
  + ContextKeyEnterpriseMemberID = "enterprise_member_id"
  + ContextKeyEnterpriseRole     = "enterprise_role"
  + GetEnterpriseIDFromContext(c) → (int64, bool)
  + GetEnterpriseMemberIDFromContext(c) → (int64, bool)
  + GetEnterpriseRoleFromContext(c) → (string, bool)
  + RequireEnterpriseAdmin(memberRepo, entRepo) gin.HandlerFunc
  + RequireEnterpriseMember(memberRepo, entRepo) gin.HandlerFunc
```

### 3.5 校验点对比

| 校验 | RequireEnterpriseAdmin | RequireEnterpriseMember |
|------|----------------------|------------------------|
| JWT 认证 | ✅ (前置) | ✅ (前置) |
| 企业成员存在 | ✅ | ✅ |
| 成员角色检查 | `enterprise_admin` | 不限 |
| 成员状态检查 | `active` | `active` |
| 企业存在 | ✅ | ✅ |
| 企业状态检查 | `active` | `active` |

### 3.6 常量依赖

```go
// 中间件使用的是 service 包的常量 (从 domain 包 re-export)
service.EnterpriseRoleAdmin     = "enterprise_admin"
service.StatusActive            = "active"       // 通用 active 状态，非 member 专属
service.EnterpriseStatusActive  = "active"       // 企业 active 状态
```

⚠️ `service.StatusActive` 是 generic 的 active 常量（domain/constants.go），不是 member 专属的 `MemberStatusActive`。这是项目中已有的设计模式。

---

## 四、模块②：企业计费分流器 (EnterpriseBillingRouter)

### 4.1 项目现有代码风格

Service 层标准模式（以 `EnterpriseService` 为例）：

```go
// 1. 结构体 + 构造函数（私有字段，依赖接口）
type EnterpriseService struct {
    entRepo EnterpriseRepository    // 接口类型
    memberRepo EnterpriseMemberRepository
}

func NewEnterpriseService(entRepo EnterpriseRepository, ...) *EnterpriseService {
    return &EnterpriseService{entRepo: entRepo, ...}
}

// 2. 方法签名: (ctx context.Context, 业务参数...) (result, error)
func (s *EnterpriseService) CreateEnterprise(ctx context.Context, ...) (*Enterprise, error) { ... }
```

计费分流器同样遵循：

```go
type EnterpriseBillingRouter struct {
    memberRepo EnterpriseMemberRepository   // 接口注入
    entRepo    EnterpriseRepository         // 接口注入
}

func NewEnterpriseBillingRouter(memberRepo, entRepo) *EnterpriseBillingRouter { ... }
```

### 4.2 P4 前逻辑

**不存在计费分流概念**。网关层统一调用 `userRepo.DeductBalance(userID, cost)` 扣个人余额。以下两种情况一律扣 `apiKey.UserID` 的个人余额：

```
个人 Key (user_id=123, assigned_to=NULL)    → 扣 user 123 余额
企业管理员的个人 Key (user_id=456, assigned_to=NULL) → 扣 user 456 余额
```

### 4.3 P4 新逻辑

```
输入: apiKey (含 AssignedTo, UserID, Group, Status)

if apiKey.AssignedTo == nil:
    → BillingSource{PoolType: "personal", PayerUserID: apiKey.UserID}
    return   // 快速路径, 无 DB 查询

if apiKey.AssignedTo != nil:
    → memberRepo.GetByID(AssignedTo)           // DB 查询 1
    → entRepo.GetByID(member.EnterpriseID)     // DB 查询 2
    → 校验 enterprise.Status == "active"
    → BillingSource{
        PoolType:     "enterprise",
        EnterpriseID: &enterprise.ID,
        PayerUserID:  apiKey.UserID,          // 创建 Key 的管理员
      }
    return
```

### 4.4 变更对照

```
P4 前:
  无分流 → 所有调用统一扣个人余额

P4 后:
  + BillingSource 结构体 (PoolType, EnterpriseID, PayerUserID)
  + DetermineBillingSource(ctx, apiKey) (*BillingSource, error)
  + 个人 Key (AssignedTo==nil)        → PoolTypePersonal,   无 DB 查询
  + 企业 Key (AssignedTo!=nil)        → PoolTypeEnterprise,  2 次 DB 查询
  + 企业 inactive 时返回 error       → 阻止企业 Key 继续消费
```

### 4.5 性能分析

- 个人池走快速路径（`AssignedTo == nil` → 直接返回，0 DB 查询）
- 企业池需要 2 次 DB 查询（会被 handler 反复调用，存在 N+1 风险）
- 缺少缓存层 (TODO: P5 优化 — Redis 缓存 memberID→enterpriseID)

### 4.6 handler 集成点 (待 P5)

```
Gateway Handler
  → EnterpriseBillingRouter.DetermineBillingSource(ctx, apiKey)
  → RecordUsageInput{
      EnterpriseID:       source.EnterpriseID,
      EnterprisePoolType: source.PoolType,
    }
  → GatewayService.RecordUsage(ctx, input)
  → recordUsageCore(...)
  → postUsageBilling(...)
```

---

## 五、模块③：EnterpriseRepository.DeductBalance

### 5.1 项目现有代码风格

Repository 层特点（以 `UserRepository.DeductBalance` 为参考）：

```go
// 扣减模式：原子 UPDATE 带 WHERE 防卫
// userRepo 使用 Ent ORM:
user, err := client.User.UpdateOneID(id).AddBalance(-amount).Where(...).Save(ctx)

// enterpriseRepo 使用 Raw SQL:
_, err := r.db.ExecContext(ctx, `UPDATE enterprises SET balance = balance - $1 WHERE id = $2 AND balance >= $1`, amount, id)
```

项目同时使用 Ent ORM（CRUD + 简单查询）和 Raw SQL（复杂查询/原子操作）。

### 5.2 P4 前逻辑

企业 Repository **只有 `GetBalance`**（查询余额），没有扣减能力：

```go
// P4 前 EnterpriseRepository
type EnterpriseRepository interface {
    Create(ctx, *Enterprise) error
    GetByID(ctx, id) (*Enterprise, error)
    List(ctx, params, filters) ([]Enterprise, *PaginationResult, error)
    Update(ctx, *Enterprise) error
    SoftDelete(ctx, id) error
    GetBalance(ctx, id) (float64, error)   // ← 只读
}
```

### 5.3 P4 新逻辑

```go
// P4 后 EnterpriseRepository
type EnterpriseRepository interface {
    // ... 原有 6 个方法 ...
    DeductBalance(ctx context.Context, id int64, amount float64) (float64, error)  // +P4
}

// 实现 (enterprise_repo.go):
func (r *enterpriseRepository) DeductBalance(ctx, id, amount) (float64, error) {
    // Step 1: 原子扣减 (WHERE balance >= amount 防超扣)
    result, err := r.db.ExecContext(ctx, `
        UPDATE enterprises
        SET balance = balance - $1
        WHERE id = $2 AND balance >= $1
    `, amount, id)
    
    if rowsAffected == 0 {
        return 0, ErrEnterpriseInsufficientBalance
    }
    
    // Step 2: 查询新余额返回
    return 计算新余额, nil
}
```

### 5.4 变更对照

```
P4 前:
  EnterpriseRepository: 6 方法 (CRUD + GetBalance)
  企业余额只能查询，不可扣减

P4 后:
  + DeductBalance(ctx, id, amount) (float64, error)
  + ErrEnterpriseInsufficientBalance 错误常量
  + 原子 UPDATE: WHERE balance >= amount (防超扣)
  + 返回新余额供调用方使用
```

---

## 六、模块④：UsageLog 结构体三层同步

### 6.1 项目现有代码风格

此项目的数据层级同步模式（三层）：
```
Ent Schema (ent/schema/usage_log.go)
    → `go generate` 生成 Ent 类型 (ent/usage_log.go)
    → Service 层结构体 (service/usage_log.go) ← 手动同步
    → Repository 层 INSERT/SELECT (repository/usage_log_repo.go) ← 手动同步
```

三层之间**没有自动同步机制**，完全依赖人工保持一致。这是已知风险点（P1/P2/P3 中多次排错记录于此）。

### 6.2 P4 前逻辑

P1 已在 `ent/schema/usage_log.go` 添加了 `enterprise_id` 和 `pool_type` 字段，但 **Service 层和 Repository 层未同步**：

```
P1 Schema ✅   usage_log → enterprise_id (int64), pool_type (string)
P2-P3 ❌       service/usage_log.go → UsageLog 结构体缺少这两个字段
P2-P3 ❌       repository/usage_log_repo.go → INSERT 缺少这两列
```

这意味着 P3 阶段写入 UsageLog 时无法记录 `enterprise_id` 和 `pool_type`，即使网关知道了计费来源。

### 6.3 P4 新逻辑

**三层同步补全**：

```go
// 1. service/usage_log.go — 结构体新增
type UsageLog struct {
    // ... 原有 ~40 个字段 ...
    EnterpriseID *int64  // P4 sync from P1 schema
    PoolType     string  // P4 sync from P1 schema
}

// 2. repository/usage_log_repo.go — INSERT 新增
var usageLogInsertArgTypes = []string{
    // ... 原有 49 个 arg type ...
    "bigint",    // $51 enterprise_id    ← +P4
    "text",      // $52 pool_type        ← +P4
    "timestamptz",
}

// execUsageLogInsertNoResult:
//   INSERT columns: ..., enterprise_id, pool_type, created_at)
//   VALUES: $1..$50, $51, $52)         ← 从 50 变 52

// prepareUsageLogInsert:
//   args[50] = nullInt64(log.EnterpriseID)   ← +P4
//   args[51] = log.PoolType                  ← +P4
```

### 6.4 变更对照

```
P4 前:
  UsageLog 结构体: 无 EnterpriseID, PoolType
  INSERT SQL:      50 个占位符 ($1..$50)
  argTypes:        不包含 "bigint"/"text"
  prepare args:    不包含 enterprise fields

P4 后:
  UsageLog 结构体: +EnterpriseID *int64, +PoolType string
  INSERT SQL:      52 个占位符 ($1..$52)
  argTypes:        +"bigint" ($51), +"text" ($52)
  prepare args:    +nullInt64(log.EnterpriseID) ($51), +log.PoolType ($52)
```

⚠️ 这种 Raw SQL position-based 模式风险很高：新增字段必须在 `created_at` 之前追加，否则参数位置对不上。

---

## 七、模块⑤：网关计费管道集成

### 7.1 项目现有代码风格

Gateway 计费采用"参数体 + 依赖体"模式：

```go
// 参数体：打包所有计费所需数据
type postUsageBillingParams struct {
    Cost               *CostBreakdown
    User               *User
    APIKey             *APIKey
    Account            *Account
    Subscription       *UserSubscription
    RequestPayloadHash string
    IsSubscriptionBill bool
    // ...
}

// 依赖体：打包所有计费所需的服务
type billingDeps struct {
    accountRepo         AccountRepository
    userRepo            UserRepository
    userSubRepo         UserSubscriptionRepository
    billingCacheService *BillingCacheService
    // ...
}

// 调用链: recordUsageCore → applyUsageBilling → postUsageBilling (legacy) 或 repo.Apply (生产)
```

### 7.2 P4 前逻辑

`postUsageBilling` 只处理两种计费路径：

```go
func postUsageBilling(ctx, p, deps) {
    if p.IsSubscriptionBill {
        // 订阅计费: deps.userSubRepo.IncrementUsage
    } else {
        // 余额计费: deps.userRepo.DeductBalance  ← 一律扣个人余额
    }
    // + API Key 配额、速率限制、Account 配额 (与谁扣费无关)
}
```

**无论 API Key 是否分配给企业成员，一律扣 `p.User.ID`（创建者）的个人余额。**

### 7.3 P4 新逻辑

```go
func postUsageBilling(ctx, p, deps) {
    if p.IsSubscriptionBill {
        // 订阅计费: 不变
    } else {
        // P4 隔离门: 企业池 vs 个人池
        if p.EnterprisePoolType == "enterprise" && p.EnterpriseID != nil && deps.enterpriseRepo != nil {
            // 企业池: deps.enterpriseRepo.DeductBalance
        } else {
            // 个人池: deps.userRepo.DeductBalance (原有逻辑)
        }
    }
    // API Key 配额、速率、Account 配额 → 不变，双方都正常扣除
}
```

### 7.4 完整数据流 (上下文传递链)

```
GatewayService.RecordUsage(ctx, &RecordUsageInput{
    EnterpriseID:       *int64,      ← P4 handler 传入
    EnterprisePoolType: string,      ← P4 handler 传入
})
    ↓
recordUsageCore(ctx, &recordUsageCoreInput{
    EnterpriseID:       *int64,      ← 从 RecordUsageInput 透传
    EnterprisePoolType: string,      ← 从 RecordUsageInput 透传
})
    ↓ 两条分支
    ├─ buildRecordUsageLog(...) → UsageLog{
    │    EnterpriseID: *int64,       ← 写入数据库
    │    PoolType:     string,       ← 写入数据库
    │  }
    └─ applyUsageBilling(..., &postUsageBillingParams{
         EnterpriseID:       *int64, ← 传给计费引擎
         EnterprisePoolType: string, ← 传给计费引擎
       }, deps, repo)
           ↓
       postUsageBilling(ctx, p, deps)
           → p.EnterprisePoolType == "enterprise"
             → deps.enterpriseRepo.DeductBalance  ← 扣企业余额
           → else
             → deps.userRepo.DeductBalance        ← 扣个人余额 (原有)
```

### 7.5 变更对照

```
P4 前:
  postUsageBillingParams: 11 个字段 (Cost/User/APIKey/Account/...)
  billingDeps:            7 个依赖 (accountRepo/userRepo/...)
  postUsageBilling:       一路扣 deps.userRepo.DeductBalance
  GatewayService:         无 enterpriseRepo

P4 后:
  postUsageBillingParams: +EnterpriseID (*int64), +EnterprisePoolType (string)
  billingDeps:            +enterpriseRepo (EnterpriseRepository)
  postUsageBilling:       if-else 分支 → enterprise vs user DeductBalance
  GatewayService:         +enterpriseRepo (通过 SetEnterpriseRepo 注入)
  + SetEnterpriseRepo()   避免修改 NewGatewayService 构造函数签名
  RecordUsageInput:       +EnterpriseID, +EnterprisePoolType (透传企业信息)
  recordUsageCoreInput:   +EnterpriseID, +EnterprisePoolType (内部分发)
```

### 7.6 注入模式

采用 **Setter 注入** 而非 **构造器注入**：

```go
// 构造器不改变 (NewGatewayService 有 30+ 参数, 4 个调用点)
gatewayService := NewGatewayService(accountRepo, groupRepo, ..., userPlatformQuotaRepo)

// P4 通过 Setter 注入
gatewayService.SetEnterpriseRepo(enterpriseRepo)
```

原因：变更 `NewGatewayService` 签名会破坏 `wire_gen.go` + 3 个测试文件的调用点。

### 7.7 隔离保证

| 保证 | 实现方式 |
|------|----------|
| 企业 Key 不扣个人余额 | `EnterprisePoolType=="enterprise"` → 走 `enterpriseRepo.DeductBalance` |
| 个人 Key 不扣企业余额 | `EnterprisePoolType!="enterprise"` → 走 `userRepo.DeductBalance` |
| API Key 配额 双方正常扣除 | 计费后的 `shouldDeductAPIKeyQuota()` 不受 pool 影响 |
| Account 配额 双方正常扣除 | `shouldUpdateAccountQuota()` 不受 pool 影响 |
| enterpriseRepo 为 nil 时降级 | `deps.enterpriseRepo != nil` 守卫 → 回退到个人扣费 |
| 企业 inactive 阻止消费 | `DetermineBillingSource` 返回 error → handler 拒绝请求 |

---

## 八、常量引用图

```
domain/constants.go
    ├─ PoolTypePersonal   = "personal"
    ├─ PoolTypeEnterprise = "enterprise"
    ├─ StatusActive       = "active"
    └─ EnterpriseStatusActive = "active"
         ↓ (re-exported)
service/domain_constants.go
         ↓
    enterprise_service.go     ← 使用 PoolType*, EnterpriseStatusActive
    enterprise_billing_router ← 使用 PoolType*
    middleware/enterprise_auth ← 使用 StatusActive, EnterpriseRoleAdmin, EnterpriseStatusActive
```

---

## 九、TODO 清单

| 优先级 | 项 | 说明 |
|--------|-----|------|
| P5 | handler 集成 | 在 GatewayHandler 中调用 DetermineBillingSource + 传递 EnterpriseID |
| P5 | EnterpriseBillingRouter 缓存 | 减少 DB 查询 (可选: Redis 缓存 memberID→enterpriseID 映射) |
| P5 | usageBillingRepo.Apply 企业支持 | 生产路径也需企业余额原子扣减 |
| P5 | 企业余额缓存 | `billingCacheService.InvalidateEnterpriseBalance()` |
| P5 | API Key CRUD 企业字段 | AssignedTo / KeyType 的 CRUD 接口 |
| P6 | 企业财务报表 | UsageLog 按 EnterpriseID+PoolType 聚合统计 |

---

## 十、设计决策记录

| 决策 | 原因 |
|------|------|
| Setter 模式 (SetEnterpriseRepo) | 避免修改 NewGatewayService 签名（破坏 4+ 调用点） |
| billingDeps 包含 enterpriseRepo | 与现有 deps 模式一致，通过 billingDeps() 统一构造 |
| postUsageBilling 走 if-else 而非策略模式 | 只有 2 个分支，策略模式过度设计 |
| DeductBalance 使用 AddBalance(-amount) | 复用已有 GTE WHERE 防卫，保证原子性 |
| enterpriseRepo == nil 时降级到个人扣费 | 兼容现有测试（无需 mock enterpriseRepo） |
| ContextKey 复用 middleware.ContextKey 类型 | 与项目已有的 JWTAuth context keys 一致 |

---

## 十一、验证门

| 门 | 状态 | 说明 |
|----|------|------|
| 编译门 | ✅ 通过 | `go build ./...` 全量编译成功 |
| 测试门 | ⏳ TODO | 需要编写 middleware + billing router 单元测试 |
| 隔离门 | ✅ 代码已实现 | 三个池通过 EnterprisePoolType 分支隔离 |
