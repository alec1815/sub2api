# 08 — P1/P2 企业功能开发排错

> **阶段**: P1 数据库 Schema & 迁移 / P2 Repository 数据访问层
> **日期**: 2026-06-29

---

## 问题清单

| # | 阶段 | 问题 | 根因 | 严重度 |
|---|:--:|------|------|:--:|
| 1 | P1 | Ent schema 命名 → Go 类型名不一致 | Ent 自动将 snake_case 转为 TitleCase，不理解规则会写错引用 | 高 |
| 2 | P1 | 部分唯一索引无法在 Ent schema 层定义 | PostgreSQL `WHERE deleted_at IS NULL` 部分索引只能通过迁移 SQL 实现 | 中 |
| 3 | P1 | `parent_id` 自引用不定义 edge | 根节点 parent_id=0 无对应记录，不适合做 FK | 中 |
| 4 | P1 | EnterpriseSubscription 不用 SoftDeleteMixin | 企业套餐通过 status 管理生命周期（active/expired/suspended） | 中 |
| 5 | P2 | OrderOption 类型使用 | Ent 生成的排序函数返回 `{Table}.OrderOption`，需传入 `sql.OrderDesc/Asc()` | 高 |
| 6 | P2 | Repository 接口放 service 包 | service 层定义接口、domain 层定义类型、repo 层实现接口 | 中 |
| 7 | P2 | Entity ↔ Service 映射无自动方案 | Ent 不提供 auto-mapping，需要手写转换函数 | 中 |
| 8 | P2 | `activeQuery()` 模式适配 | 有 SoftDeleteMixin 的 entity 需要 `DeletedAtIsNil()` 过滤 | 中 |
| 9 | P2 | 跨 entity 查询 | DepartmentRepo 查成员需要直接访问 `client.EnterpriseMember.Query()` | 低 |

---

## 1. Ent Schema 命名 → Go 类型名

### 现象

Ent 的 `go generate ./ent` 会根据 schema struct 名和字段名自动生成 Go 代码。命名规则：

```
Schema struct    →  Go 类型          →  数据库表名（@Table 注解）
──────────────────────────────────────────────────────────────
Enterprise       →  Enterprise       →  enterprises
EnterpriseMember →  EnterpriseMember →  enterprise_members
Department       →  Department       →  departments
APIKeyGroup      →  APIKeyGroup      →  api_key_groups
```

**字段方法名**：`field.String("contact_name")` → 生成 `SetContactName()` / `ContactName` 等。

**Ent 生成的包**：`backend/ent/enterprise`、`backend/ent/department` 等，每个 schema 一个子包，包含：
- 字段常量：`enterprise.FieldName`、`enterprise.FieldID`
- 排序函数：`enterprise.ByName()`、`enterprise.ByCreatedAt()`
- 条件谓词：`enterprise.NameEQ()`、`enterprise.NameContainsFold()`
- Schema 构造：`enterprise.IDEQ()`、`enterprise.DeletedAtIsNil()`

### 踩坑

❌ 错误用法：
```go
import "github.com/Wei-Shaw/sub2api/ent/enterprises"  // 不存在，包名是单数
```

✅ 正确用法：
```go
import (
    dbent "github.com/Wei-Shaw/sub2api/ent"           // 根包别名
    "github.com/Wei-Shaw/sub2api/ent/enterprise"      // 子包是单数 schema 名
)
```

### 教训

**永远先 `go generate ./ent` 再看生成代码确认 API**，不要凭猜测写引用。

---

## 2. 部分唯一索引（WHERE deleted_at IS NULL）

### 现象

Ent 原生不支持 `WHERE` 条件索引。企业名需要在未删除记录中唯一，但已软删除的记录可以重复。

### 方案

- Ent schema 层只定义普通唯一索引（`index.Fields("name")`）
- 真正的**部分唯一索引**在迁移 SQL 中实现

```sql
-- 154_create_enterprises.sql
CREATE UNIQUE INDEX idx_enterprises_name_active
    ON enterprises (name) WHERE deleted_at IS NULL;
```

### 教训

Ent schema 和迁移 SQL 是互补关系——schema 负责 Go 代码生成，迁移 SQL 负责数据库级约束。

---

## 3. parent_id 自引用：不定义 edge

### 现象

企业表有 `parent_id` 表示父企业，部门表有 `parent_id` 表示父部门。根节点的 `parent_id = 0`，没有对应记录。

### 决策

**不定义自引用 edge**。原因：
- `parent_id = 0` 是合法值，但不存在 id=0 的记录
- 如果定义 FK edge，Ent 会要求 parent 必须存在，导致根节点无法创建

### 教训

有"零值表示无父节点"语义的字段，不要定义自引用 edge，直接用字段。

---

## 4. EnterpriseSubscription 不用 SoftDeleteMixin

### 现象

企业套餐不是"软删除"语义，而是"生命周期状态"：`active → expired → suspended`。

### 对比

| 实体 | Mixin | 原因 |
|:--|:--|:--|
| Enterprise | TimeMixin + SoftDeleteMixin | 企业会删除，需要恢复能力 |
| EnterpriseMember | TimeMixin + SoftDeleteMixin | 成员解绑=软删除，保留历史 |
| Department | TimeMixin + SoftDeleteMixin | 部门删除后可能有历史数据引用 |
| EnterpriseSubscription | **仅 TimeMixin** | 通过 status 字段管理，不软删除 |
| APIKeyGroup | TimeMixin + SoftDeleteMixin | 中间表，解绑=删除 |

### 教训

不是所有 entity 都需要 SoftDeleteMixin——根据业务语义选择 mixin 组合。

---

## 5. OrderOption 类型使用

### 现象

Ent 生成的排序函数签名：
```go
// ent/enterprise/enterprise.go (生成代码)
func ByName(opts ...sql.OrderTermOption) OrderOption
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption
func ByID(opts ...sql.OrderTermOption) OrderOption
```

`OrderOption` 是函数类型，用于 `Query.Order(...)`。

### 动态排序实现

```go
func enterpriseListOrder(params pagination.PaginationParams) enterprise.OrderOption {
    isDesc := params.NormalizedSortOrder(pagination.SortOrderDesc) == pagination.SortOrderDesc
    switch strings.ToLower(strings.TrimSpace(params.SortBy)) {
    case "name":
        if isDesc {
            return enterprise.ByName(sql.OrderDesc())
        }
        return enterprise.ByName(sql.OrderAsc())
    case "", "created_at":
        if isDesc {
            return enterprise.ByCreatedAt(sql.OrderDesc())
        }
        return enterprise.ByCreatedAt(sql.OrderAsc())
    }
}
```

### 踩坑

❌ 错误写法：`return enterprise.ByName(sql.OrderDesc)` — 缺少函数调用括号

✅ 正确：`return enterprise.ByName(sql.OrderDesc())`

### 两种排序方式

Ent 支持两种排序语法：
1. **类型安全的 OrderOption**（推荐）：
   ```go
   .Order(enterprise.ByName(sql.OrderAsc()))
   ```
2. **字符串字段名**（适合固定排序）：
   ```go
   .Order(dbent.Asc(department.FieldOrderNum), dbent.Asc(department.FieldID))
   ```

---

## 6. 分层架构：接口放 service、类型放 domain、实现在 repo

### 架构约定

```
domain/enterprise.go          ← 领域类型 + 错误定义
service/enterprise_service.go ← Repository 接口 + Service 类型别名
repository/enterprise_repo.go ← 接口实现（私有结构体）
```

### 类型流向

```go
// domain: 定义原始类型
type Enterprise struct { ... }
var ErrEnterpriseNotFound = ...

// service: 别名暴露 + 接口定义
type Enterprise = domain.Enterprise       // 类型别名（不是新类型）
var ErrEnterpriseNotFound = domain.ErrEnterpriseNotFound

type EnterpriseRepository interface {
    Create(ctx context.Context, e *Enterprise) error
    // ...
}

// repository: 实现接口
type enterpriseRepository struct { client *dbent.Client }  // 私有结构体
func NewEnterpriseRepository(client *dbent.Client) service.EnterpriseRepository { ... }
```

### 优点

- 上层依赖接口，不依赖具体实现（方便测试 Mock）
- 类型定义集中在 domain，service 层通过别名引用
- Wire 依赖注入：接口 = service 层定义，实现 = repo 层构造函数

### 踩坑

P2 早期尝试把 Repository 接口放 `repository` 包，导致 service 层 import repo 层 → 循环依赖。

---

## 7. Entity ↔ Service 映射

### 现象

Ent 不提供自动 ORM 映射（类似 GORM 的 `Scan`），需要手写转换。

### 模式

```go
func enterpriseEntityToService(m *dbent.Enterprise) *service.Enterprise {
    if m == nil { return nil }
    return &service.Enterprise{
        ID: m.ID,
        Name: m.Name,
        // ... 逐字段映射
    }
}

func enterpriseEntitiesToService(models []*dbent.Enterprise) []service.Enterprise {
    out := make([]service.Enterprise, 0, len(models))
    for i := range models {
        if s := enterpriseEntityToService(models[i]); s != nil {
            out = append(out, *s)
        }
    }
    return out
}
```

### 注意

- 列表转换**返回值为值切片** `[]service.Enterprise`（不是指针切片），对齐现有代码风格
- 单个查询返回 `*service.Enterprise`（指针）
- `Create`/`Update` 方法回写 ID/时间戳到入参指针

---

## 8. activeQuery 模式

### 模式

```go
func (r *enterpriseRepository) activeQuery() *dbent.EnterpriseQuery {
    return r.client.Enterprise.Query().Where(enterprise.DeletedAtIsNil())
}
```

### 规则

| 场景 | 是否使用 activeQuery |
|:--|:--|
| 查询列表 | ✅ |
| 按 ID 查单个 | ✅ |
| 是否存在检查 | ✅ |
| 软删除操作 | ❌ 直接用 `client.Enterprise.DeleteOneID()` |

### 注意事项

- `EnterpriseSubscription` 没有 SoftDeleteMixin → 不需要 `activeQuery()`
- `HasMembers` 跨 entity 查询时，手动写 `enterprisemember.DeletedAtIsNil()`

---

## 9. 跨 entity 查询

### 现象

DepartmentRepo 的 `HasMembers()` 需要查 `enterprise_members` 表，但 department_repo 只有 `Department.Query()`。

### 方案

直接通过 Ent client 查询跨 entity：
```go
func (r *departmentRepository) HasMembers(ctx context.Context, id int64) (bool, error) {
    count, err := r.client.EnterpriseMember.Query().
        Where(
            enterprisemember.DepartmentIDEQ(id),
            enterprisemember.StatusEQ(domain.StatusActive),
            enterprisemember.DeletedAtIsNil(),   // 手动加软删除过滤
        ).Count(ctx)
    return count > 0, err
}
```

### 教训

Repository 不限于只访问自己的 entity——Ent client 可以查询任何 schema。跨 entity 查询只需 import 对应 schema 子包。

---

## 经验总结

### 1. Ent Schema → 代码生成是单向的

```
修改 schema.go → go generate ./ent → 新代码覆盖 ent/ 目录
                ⚠️ 不要手动修改 ent/ 下的生成代码！
```

### 2. 先读现有代码模式，再写新代码

P2 所有 Repo 遵循 `apiKeyRepository` 的模式：
- 私有结构体 `xxxRepository struct { client *dbent.Client }`
- 构造函数 `NewXxxRepository(client) service.XxxRepository`
- `clientFromContext()` 支持事务
- `translatePersistenceError()` 统一错误转换

### 3. 进程式开发：每完成一个 Repo 就 `go build`

不要等 5 个 Repo 全写完再编译。每完成一个就 `go build ./...` 验证，避免错误堆积。

### 4. Wire 注册要及时同步

新增 Repo 后立即在 `wire.go` 注册，否则 DI 容器编译失败：
```go
wire.Bind(new(service.EnterpriseRepository), new(*enterpriseRepository)),
```

---

> **下一步**：[P3 Service 层](../enterprise/v1/06-开发计划任务.md) · 依赖 P2 Repo 全部完成
