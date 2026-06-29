# 09 — P3 Service 层排错记录

> **阶段**: P3 Service 业务逻辑层
> **日期**: 2026-06-29
> **前置**: 08-ent-repo-pitfalls.md (P1/P2)

---

## 问题清单

| # | 问题 | 根因 | 严重度 |
|:--:|------|------|:--:|
| 1 | 同一文件创建两个 Service 导致依赖混淆 | EnterpriseService 和 EnterpriseMemberService 引用同一文件中的接口定义 | 中 |
| 2 | `APIKey` 结构体缺少企业新字段 | P1 在 Ent schema 加了字段但 Service 层 APIKey struct 没同步 | 高 |
| 3 | 分页计算中的类型不匹配 | `pagination.PaginationParams` 返回 `int` 但分页切割需要 `int64` | 低 |
| 4 | `nowFunc()` 未使用导致 lint 告警 | 月度用量计算是 TODO 占位，辅助函数无调用者 | 低 |
| 5 | DepartmentService 中变量重定义 | `ErrDepartmentHasChildren` 与 `ErrDepartmentHasMembers` 引用方式 | 低 |
| 6 | Repository 接口放在多个文件中 | 部分接口定义在 `enterprise_service.go` 但被其他 service 引用 | 中 |
| 7 | 1:1 约束校验在 Service 层而非 DB 层 | `enterprise_members` 已有部分唯一索引，但 Service 仍需先查再报 | 中 |
| 8 | Profile 月度用量为零值占位 | 网关计费层 P4 未实现，usage_logs.enterprise_id 无数据 | 中 |

---

## 1. 同一文件定义两个 Service 及其接口

### 现象

`enterprise_service.go` 同时定义了：
- `EnterpriseRepository` 接口
- `EnterpriseMemberRepository` 接口
- `EnterpriseSubscriptionRepository` 接口
- `APIKeyGroupRepository` 接口
- `EnterpriseService` 业务逻辑

导致其他文件（如 `enterprise_member_service.go`）引用 `EnterpriseMemberRepository` 时依赖 `enterprise_service.go` 编译，产生隐性耦合。

### 决策

**保持现状**。原因：
- 现有代码风格：`api_key_service.go` 同样定义 `APIKeyRepository` 接口
- 统一接口定义在 Service 层是项目约定（见 08 文档 #6）
- 拆分到独立文件会增加文件数量但不会减少耦合（仍 import 同一 package）

### 教训

如果未来 Repository 接口数量超过 10 个，可考虑独立 `repository_interfaces.go` 文件。

---

## 2. APIKey 结构体缺少企业新字段

### 现象

P1 在 `ent/schema/api_key.go` 添加了 `assigned_to`、`usage_purpose`、`bound_tool` 字段，Ent 生成代码已包含。但 `service/api_key.go` 中的 `APIKey` 结构体没有同步添加。

### 修复

在 `service/api_key.go` 的 `APIKey` struct 中新增：
```go
AssignedTo   *int64  // 分配给企业成员的 enterprise_members.id
UsagePurpose string  // 用途说明
BoundTool    string  // 绑定工具
```

### 教训

**Schema 加字段时，必须同步更新三层**：
1. `ent/schema/xxx.go` → Ent 生成代码
2. `service/xxx.go` → Service 层结构体
3. `repository/xxx_repo.go` → Entity↔Service 映射函数

缺失任一层会导致编译通过但运行时数据丢失。

---

## 3. 分页计算类型不匹配

### 现象

`EnterpriseKeyService.ListKeys` 中简易分页：
```go
offset := params.Offset()  // 返回 int
// ...
offset >= total            // int vs int64
```

### 修复

显式转换：
```go
offset := int64(params.Offset())
limit := int64(params.Limit())
```

---

## 4. nowFunc() 未使用导致 lint 告警

### 现象

`enterprise_profile_service.go` 定义了 `func nowFunc() time.Time { return time.Now() }` 但从未调用。Lint 报告 `unused function`。

### 根因

月度用量计算是 P4 阶段的 TODO。`GetProfile` 返回的 `MonthlyUsageOverview` 当前为零值占位，`nowFunc()` 是本意用于计算"当月"范围的时间工具函数。

### 修复

删除 `nowFunc()`，在三处 `GetProfile`/`GetProfileByEnterprise`/`GetProfileWithMember` 中添加 `TODO(P4)` 注释。月度用量零值占位不影响其他逻辑。

### 教训

辅助函数应与调用方**同批次实现**，不应预留给后续阶段。如果必须预留，使用 `_ = nowFunc` 消除 lint 警告并添加 TODO。

---

## 5. DepartmentService 变量重定义

### 现象

```go
var (
    ErrDepartmentCircularRef = infraerrors.BadRequest(...)
    ErrDepartmentHasChildrenErr = ErrDepartmentHasChildren  // 引用 domain 中已定义变量
    ErrDepartmentHasMembersErr  = ErrDepartmentHasMembers
)
```

`ErrDepartmentHasChildren` 和 `ErrDepartmentHasMembers` 已在 `domain/enterprise.go` 定义，这里重新创建别名是为了在 service 包内直接使用。

### 决策

保持别名定义。原因：Service 层所有错误在同一个 package 内引用，避免跨 package 查找。

---

## 6. Repository 接口集中定义

### 现象

- `EnterpriseRepository` → `enterprise_service.go`
- `EnterpriseMemberRepository` → `enterprise_service.go`
- `DepartmentRepository` → `department_service.go`
- `APIKeyRepository` → `api_key_service.go`

接口分散在多个文件中，新增 Service 时需要确认接口定义位置。

### 决策

保持现有模式。因为 Go 的 package 内引用不需要 import，接口定义在哪一个文件在同一个 package 中都可见。

---

## 7. 1:1 约束的双重校验

### 现象

`enterprise_members` 表有 `UNIQUE (user_id) WHERE status='active' AND deleted_at IS NULL` 的部分唯一索引。但 `EnterpriseMemberService.CreateMember` 中仍然手动查询该用户是否已是其他企业 active 成员。

### 决策

**保留 Service 层校验**。原因：
- 索引冲突返回的是数据库错误（`pq: duplicate key`），错误消息对前端不友好
- Service 层提前校验可以返回 `ErrMemberUserAlreadyMember` 业务错误，语义更清晰
- 数据库索引作为**最后防线**防止并发竞争

---

## 8. Profile 月度用量占位

### 现象

`EnterpriseProfileService` 三个方法返回的 `MonthlyUsageOverview` 均为 `{CallCount: 0, TotalCost: 0}`。

### 根因

月度用量统计依赖 `usage_logs.enterprise_id` 字段聚合，该字段由 P4 网关计费层写入。P3 阶段网关层尚未实现，usage_logs 中没有 enterprise_id 数据。

### 影响

前端 Profile 页展示"本月 API 调用次数: 0 / 总花费: $0.00"，不影响其他功能。

### 后续计划

P4 网关计费分流实现后，在 Profile 查询中接入：
```sql
SELECT COUNT(*), SUM(total_cost)
FROM usage_logs
WHERE enterprise_id = ? AND pool_type = 'enterprise'
  AND created_at BETWEEN month_start AND month_end
```

---

## 经验总结

### 1. Schema 变更 → 三层同步

每次修改 Ent schema，除了 `go generate ./ent`，还必须检查：
- `service/` 中的结构体定义
- `repository/` 中的映射函数

### 2. 类型安全

Go 的 `int` vs `int64` 在 32/64 位系统上行为不同。Service 层与 Repository 层之间统一使用 `int64`（数据库主键类型）。

### 3. TODO 标记规范

跨阶段预留的 TODO 必须有明确的阶段标记：`TODO(P4): xxx`，便于搜索和进度追踪。

### 4. 编译即门控

P3 每完成一个 Service 就 `go build ./...` 编译。不要等所有 Service 写完再编译。

---

> **下一步**: [P4 中间件 + 网关计费路由](../enterprise/v1/06-开发计划任务.md#p4中间件--网关计费路由)
