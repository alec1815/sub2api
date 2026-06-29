# ⑦ 平台添加分析 — DeepSeek 全链路源码追踪

> **核心问题**：分组页面的平台下拉框中没有 DeepSeek，要加需要改哪些地方？

**结论先行**：需要改动 **2 层、35+ 处**（前端 ~24 处 + 后端 ~13 处），分布在：
- **前端 UI 层**：类型定义 + 平台硬编码列表 + 颜色/图标 + i18n + channel monitor provider
- **后端核心链**：常量 + schema validator (×4) + 网关路由 + 渠道价格同步 + fallback model + mix-platform check

---

## 目录

1. [前端 — 平台类型与硬编码列表（~12 处）](#1-前端--平台类型与硬编码列表)
2. [前端 — Channel Monitor Provider（3 处 ⚠️ 新增）](#2-前端--channel-monitor-provider)
3. [前端 — 平台颜色/图标/样式（~5 处）](#3-前端--平台颜色图标样式)
4. [前端 — 账号创建/编辑（~3 处）](#4-前端--账号创建编辑)
5. [后端 — 常量与 Schema（8 处 ⚠️ 新增 4 处）](#5-后端--常量与-schema)
6. [后端 — 网关路由与调度（2 处）](#6-后端--网关路由与调度)
7. [后端 — 渠道价格同步与 Fallback Model（4 处 ⚠️ 新增 3 处）](#7-后端--渠道价格同步与-fallback-model)
8. [后端 — 计价/上游服务（2 处依赖检查）](#8-后端--计价上游服务)
9. [总结：三层递进实现方案](#总结三层递进实现方案)

---

## 1. 前端 — 平台类型与硬编码列表

这是**最多散落硬编码**的地方，每个管理页面几乎都有一份自己的平台列表。

### 1.1 核心类型定义（3 处）

| # | 文件 | 位置 | 当前值 | 改动 |
|---|------|------|--------|------|
| 1 | `frontend/src/types/index.ts` | L490 | `GroupPlatform = 'anthropic' \| 'openai' \| 'gemini' \| 'antigravity'` | 加 `\| 'deepseek'` |
| 2 | `frontend/src/types/index.ts` | L693 | `AccountPlatform = 'anthropic' \| 'openai' \| 'gemini' \| 'antigravity'` | 加 `\| 'deepseek'` |
| 3 | `frontend/src/api/admin/users.ts` | L310 | `PlatformQuotaPlatform = 'anthropic' \| 'openai' \| 'gemini' \| 'antigravity'` | 加 `\| 'deepseek'` |

这三个类型是所有其他平台引用**统一收敛的点**，改了它们，其他引用这两个类型的地方就能自动识别 `'deepseek'`。

### 1.2 系统设置 API（2 处 ⚠️ 新增）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 4 | `frontend/src/api/admin/settings.ts` | L20, L33 | `PlatformType` 类型 + `PLATFORMS` 数组 — 驱动 **SettingsView 中的默认平台限额模板**和认证源平台限额模板 |
| 5 | `frontend/src/api/admin/settings.ts` | L535-536, L787-790 | `SystemSettings` 和 `UpdateSettingsRequest` 中的 `fallback_model_anthropic/openai/gemini/antigravity` **4 个独立字段**，需加 `fallback_model_deepseek` |

### 1.3 分组管理页（2 处）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 6 | `frontend/src/views/admin/GroupsView.vue` | ~L3136 `platformOptions` | 创建/编辑分组时的平台选择下拉 |
| 7 | `frontend/src/views/admin/GroupsView.vue` | ~L3143 `platformFilterOptions` | 分组列表的筛选下拉 |

### 1.4 渠道管理页（1 处 ⚠️ 关键）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 8 | `frontend/src/views/admin/ChannelsView.vue` | ~L763 `platformOrder` | **这是你提到的地方！** 渠道创建对话框中"平台配置"的复选框列表 + per-platform 定价 Tab |

```typescript
// line 763
const platformOrder: GroupPlatform[] = ['anthropic', 'openai', 'gemini', 'antigravity']
```

这个列表还驱动了 `apiToForm()` 和 `formToAPI()` 的遍历逻辑，影响渠道数据在创建/编辑时的平台 Tab 生成。

### 1.5 账号管理（2 处）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 9 | `frontend/src/components/admin/account/AccountTableFilters.vue` | L28 `pOpts` | 账号列表的平台筛选下拉 |
| 10 | `frontend/src/components/account/CreateAccountModal.vue` | L70-150 | **硬编码的 4 个平台按钮**（Anthropic / OpenAI / Gemini / Antigravity），每个都是独立的 `<button>` 元素 |

### 1.6 订阅管理（1 处）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 11 | `frontend/src/views/admin/SubscriptionsView.vue` | ~L966 `platformFilterOptions` | 订阅列表筛选 |

### 1.7 运维面板（1 处）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 12 | `frontend/src/views/admin/ops/components/OpsDashboardHeader.vue` | ~L109 `platformOptions` | 运维仪表盘平台筛选 |

### 1.8 系统设置页面（2 处）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 13 | `frontend/src/views/admin/SettingsView.vue` | ~L3290 | 平台限额设置的默认模板 `['anthropic', 'openai', 'gemini', 'antigravity'] as const` |
| 14 | `frontend/src/views/admin/SettingsView.vue` | ~L3625 | 认证源平台限额的默认模板（同上） |

### 1.9 用户配额（2 处）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 15 | `frontend/src/components/admin/user/UserPlatformQuotaModal.vue` | L131 | `PLATFORMS` 数组 |
| 16 | `frontend/src/components/user/UserPlatformQuotaCell.vue` | L37 | `PLATFORM_ORDER` 数组 |

### 1.10 用户仪表盘（2 处）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 17 | `frontend/src/components/user/dashboard/UserDashboardStats.vue` | L281 | `PLATFORM_ORDER` 用于卡片排序 |
| 18 | `frontend/src/components/user/dashboard/UserDashboardStats.vue` | L250-254 | `PLATFORM_LABELS` 也需要加 `deepseek: 'DeepSeek'` |

### 1.11 用户使用分析（1 处 ⚠️ 新增）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 19 | `frontend/src/components/user/PlatformUsageBreakdown.vue` | L93-97 | `PLATFORM_LABELS` 用于平台标签展示 |

### 1.12 错误透传规则（1 处）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 20 | `frontend/src/components/admin/ErrorPassthroughRulesModal.vue` | L488 | `platformOptions` 硬编码数组 |

### 1.13 模型白名单（1 处）

| # | 文件 | 位置 | 说明 |
|---|------|------|------|
| 21 | `frontend/src/components/account/ModelWhitelistSelector.vue` | L184 | `upstreamSyncPlatforms` Set，控制"从上游同步"按钮可见性 |

---

## 2. 前端 — Channel Monitor Provider ⚠️ 新增

Channel Monitor 功能使用独立的 `Provider` 类型和常量，**不是复用 GroupPlatform/AccountPlatform**。

### 2.1 API 类型

| # | 文件 | 位置 | 改动 |
|---|------|------|------|
| 22 | `frontend/src/api/admin/channelMonitor.ts` | L8 | `Provider = 'openai' \| 'anthropic' \| 'gemini'` → 加 `\| 'deepseek'` |

### 2.2 常量定义

| # | 文件 | 位置 | 改动 |
|---|------|------|------|
| 23 | `frontend/src/constants/channelMonitor.ts` | L12-23 | 加 `PROVIDER_DEEPSEEK` 常量 + 加到 `PROVIDERS` 数组 |

### 2.3 格式化 composable

| # | 文件 | 位置 | 改动 |
|---|------|------|------|
| 24 | `frontend/src/composables/useChannelMonitorFormat.ts` | L59-102, L161-171 | 4 个函数需加 deepseek case：`providerLabel()`, `providerBadgeClass()`, `providerPickerClass()`, `providerGradient()` |

---

## 3. 前端 — 平台颜色/图标/样式

### 3.1 平台颜色定义（12 个颜色映射 + 2 个校验函数）

| # | 文件 | 改动 |
|---|------|------|
| 25 | `frontend/src/utils/platformColors.ts` | ① 类型 `Platform`（L8）加 `'deepseek'` |
| | | ② 所有 Record 加 `deepseek` 条目：`BADGE`, `BORDER`, `TEXT`, `ICON`, `BUTTON`, `BADGE_LIGHT`, `ACCENT_BAR`, `DISCOUNT`, `GRADIENT`, `GRADIENT_TEXT`, `GRADIENT_SUBTEXT` |
| | | ③ `isPlatform()` 函数（L109）加 `|| p === 'deepseek'` |
| | | ④ `platformLabel()` 函数（L157）加 `case 'deepseek': return 'DeepSeek'` |

### 3.2 渠道定价组件颜色

| # | 文件 | 改动 |
|---|------|------|
| 26 | `frontend/src/components/admin/channel/types.ts` | `getPlatformTagClass()`（L193）和 `getPlatformTextClass()`（L204）各加 `case 'deepseek'` |

### 3.3 平台图标

| # | 文件 | 改动 |
|---|------|------|
| 27 | `frontend/src/components/common/PlatformIcon.vue` | 在模板中添加 DeepSeek 的 SVG logo（在 antigravity 的 `<svg>` 之后加一个 `v-else-if="platform === 'deepseek'"`） |

### 3.4 平台 Badge

| # | 文件 | 改动 |
|---|------|------|
| 28 | `frontend/src/components/common/PlatformTypeBadge.vue` | ① `platformLabel` computed（L75）加 `if (props.platform === 'deepseek') return 'DeepSeek'` |
| | | ② `platformClass` computed（L119）加 deepseek 颜色分支 |

---

## 4. 前端 — 账号创建/编辑

### 4.1 CreateAccountModal（硬编码按钮）

| # | 文件 | 改动 |
|---|------|------|
| 29 | `frontend/src/components/account/CreateAccountModal.vue` | L70-150：在 Antigravity 按钮后面添加第 5 个 DeepSeek 平台按钮（`@click="form.platform = 'deepseek'"`） |

### 4.2 EditAccountModal（placeholder 映射）

| # | 文件 | 改动 |
|---|------|------|
| 30 | `frontend/src/components/account/EditAccountModal.vue` | L38-46 和 L60-67：baseUrl/ApiKey placeholder 需要加 deepseek 分支（如 `https://api.deepseek.com` / `sk-...`） |

### 4.3 测试文件（1 处 ⚠️ 新增）

| # | 文件 | 改动 |
|---|------|------|
| 31 | `frontend/src/views/admin/__tests__/SettingsView.spec.ts` | L1119：测试断言 `platforms = ["anthropic", "openai", "gemini", "antigravity"]`，需要加 `"deepseek"`；以及对 `quotas` 含 4 key 的断言需要更新为 5 key |

---

## 5. 后端 — 常量与 Schema

### 5.1 域常量定义（3 处）

| # | 文件 | 位置 | 改动 |
|---|------|------|------|
| 32 | `backend/internal/domain/constants.go` | L20-24 | 加 `PlatformDeepSeek = "deepseek"` |
| 33 | `backend/internal/service/domain_constants.go` | L39-44 | 加 `PlatformDeepSeek = domain.PlatformDeepSeek`（alias） |
| 34 | `backend/internal/service/domain_constants.go` | L49-54 `AllowedQuotaPlatforms` | 数组加 `PlatformDeepSeek` — **平台限额的唯一权威源** |

### 5.2 Ent Schema — 构建期校验（4 处 ⚠️ 新增 2 处）

| # | 文件 | 位置 | 改动 |
|---|------|------|------|
| 35 | `backend/ent/schema/group.go` | 枚举注释 | 在平台枚举注释中加 `deepseek`（注意：DB 字段是 varchar(50)，天然兼容） |
| 36 | `backend/ent/schema/user_platform_quota.go` | L41-44 validator | 白名单加 `case "deepseek": return nil` — **最容易遗漏！** |
| 37 | `backend/ent/schema/channel_monitor.go` | L38 `provider` Enum | `Values("openai", "anthropic", "gemini")` → **加 `"deepseek"`** |
| 38 | `backend/ent/schema/channel_monitor_request_template.go` | L42 `provider` Enum | 同上，**加 `"deepseek"`** |

> ⚠️ `channel_monitor` 的 provider 枚举是 Ent schema 的 `field.Enum().Values(...)`，新增值需要生成 migration（PostgreSQL 不支持 ALTER TYPE ADD VALUE 在事务内对已有 enum 生效，可能需要手动操作或 create new enum + swap）。

### 5.3 混合渠道检测（1 处 ⚠️ 新增）

| # | 文件 | 位置 | 改动 |
|---|------|------|------|
| 39 | `backend/internal/service/admin_service.go` | L3813-3821 `getAccountPlatform()` | switch 只覆盖 `PlatformAntigravity` 和 `PlatformAnthropic`/"claude"，需加 `case PlatformDeepSeek: return "DeepSeek"` |

这个函数在检查「账号的 platform 是否与已有分组 platform 冲突」时使用，不加的话 DeepSeek 账号的混合渠道检测会失效。

---

## 6. 后端 — 网关路由与调度

| # | 文件 | 改动 |
|---|------|------|
| 40 | `backend/internal/server/routes/gateway.go` | 路由分支判断中，将 DeepSeek 请求指向 OpenAI handler（因为 DeepSeek API 是 OpenAI 兼容的） |
| 41 | `backend/internal/service/proxy.go` 或相关文件 | 如果有 `IsOpenAI()` / `IsAnthropic()` 等方法，确认 DeepSeek 账号/分组能正确路由到 OpenAI 兼容处理路径 |

> **设计决策**：DeepSeek API 完全兼容 OpenAI Chat Completions 格式，推荐复用 OpenAI handler 路径。只需在路由判断中加 `|| platform == PlatformDeepSeek`。

---

## 7. 后端 — 渠道价格同步与 Fallback Model

### 7.1 渠道模型同步映射（1 处 ⚠️ 关键）

| # | 文件 | 位置 | 改动 |
|---|------|------|------|
| 42 | `backend/internal/handler/admin/channel_handler.go` | L507-512 | 加 `service.PlatformDeepSeek: "deepseek"` |

当前映射表：
```go
var platformToLiteLLMProvider = map[string]string{
    service.PlatformAnthropic:   "anthropic",
    service.PlatformOpenAI:      "openai",
    service.PlatformGemini:      "google",
    service.PlatformAntigravity: "anthropic",  // Antigravity 复用 Anthropic 定价
}
```

DeepSeek 在 LiteLLM 定价目录中的 provider 是 `"deepseek"`。需要添加：
```go
service.PlatformDeepSeek: "deepseek",
```

**作用**：渠道创建页面中，"同步最新模型"按钮会调用 `GET /api/v1/admin/channels/pricing/sync-models?platform=deepseek` → 后端通过这个映射找到 LiteLLM provider → `PricingService.ListModelNamesByProvider("deepseek")` → 返回 DeepSeek 模型列表。

### 7.2 Fallback Model 设置（3 处 ⚠️ 新增）

| # | 文件 | 位置 | 改动 |
|---|------|------|------|
| 43 | `backend/internal/service/domain_constants.go` | L312-315 | 加 `SettingKeyFallbackModelDeepSeek = "fallback_model_deepseek"` |
| 44 | `backend/internal/service/settings_view.go` | L160-163 | 加 `FallbackModelDeepSeek string` 字段 |
| 45 | `backend/internal/service/setting_service.go` | fallback model 读写 | 加 deepseek 分支的读写逻辑 |

当前后端 `SettingsView` 只有 4 个平台的 fallback model 字段，前端 `settings.ts` 的 `SystemSettings` 和 `UpdateSettingsRequest` 同理。如果不加，运维面板中 DeepSeek 平台就无法配置降级模型。

### 7.3 代理质量检测（1 处 ⚠️ 新增）

| # | 文件 | 位置 | 改动 |
|---|------|------|------|
| 46 | `backend/internal/service/admin_service.go` | L501-527 `proxyQualityTargets` | 加 DeepSeek 的 API endpoint（如 `https://api.deepseek.com/v1/models`，method GET，allowed status 401）用于代理质量检测 |

---

## 8. 后端 — 计价/上游服务

| # | 文件 | 说明 |
|---|------|------|
| 47 | `backend/internal/service/pricing_service.go` | LiteLLM 数据中的 `LiteLLMProvider: "deepseek"` 不需要改代码，但需要确保 JSON 数据源中包含 deepseek 的定价 |
| 48 | `backend/internal/service/upstream_models.go` | 上游模型同步：需确认 DeepSeek 模型如何加入 `ForwardAsRawChatCompletions` 的处理 |

---

## 9. 总结：三层递进实现方案

### 最小可行方案（MVP）— 让 DeepSeek 成为一等平台

改动清单（按优先级）：

```
优先级 P0（不改就没法用）：
  ① 后端 — domain/constants.go 加 PlatformDeepSeek
  ② 后端 — domain_constants.go 加 Alias + AllowedQuotaPlatforms
  ③ 后端 — 4 个 ent schema validator（group, user_platform_quota, channel_monitor × 2）
  ④ 前端 — 3 个类型定义（GroupPlatform, AccountPlatform, PlatformQuotaPlatform）
  ⑤ 前端 — 所有硬编码平台列表 × ~15 个页面/组件
  ⑥ 前端 — 颜色/图标/样式 × ~5 个文件
  ⑦ 后端 — 网关路由（让 DeepSeek 走 OpenAI handler）
  ⑧ 前端 — CreateAccountModal 加 DeepSeek 平台按钮

优先级 P1（功能完善）：
  ⑨ 后端 — channel_handler.go 的 platformToLiteLLMProvider 映射
  ⑩ 后端 — upstream_models.go 的模型同步
  ⑪ 后端 — SettingsView + settings.ts 的 fallback_model_deepseek 字段
  ⑫ 后端 — admin_service.go 的 getAccountPlatform() 混合渠道检测
  ⑬ 后端 — admin_service.go 的 proxyQualityTargets 代理质量检测
  ⑭ 前端 — Channel Monitor Provider 类型 + 常量 + 格式化函数
  ⑮ 前端 — EditAccountModal 的 placeholder 映射
  ⑯ 前端 — i18n 翻译（zh.ts / en.ts）
  ⑰ 前端 — 测试文件更新（SettingsView.spec.ts 等）

优先级 P2（可选优化）：
  ⑱ 后端 — 为 DeepSeek 添加专属路由分支（如需 reasoning/thinking 特殊处理）
  ⑲ 后端 — DeepSeek 专用的 warmup / 健康检查逻辑
  ⑳ 前端 — 合并所有 PLATFORM_ORDER/LABELS 为统一注册中心
```

### 关键架构观察

| 发现 | 影响 |
|------|------|
| **平台列表散落在 ~20 个位置** | 没有统一的"平台注册中心"，前端每页维护自己的列表 |
| **Channel Monitor 的 Provider 是独立类型** | 不同于 GroupPlatform/AccountPlatform，3 个 Provider 常量 + 1 个 API type + 1 个 composable（4 个 switch 函数） |
| **GroupPlatform 和 AccountPlatform 是两个独立类型** | 即使内容相同，也要分别修改 |
| **CreateAccountModal 用硬编码 button 而非 v-for** | 无法通过只改数组就加新平台，需要复制粘贴大量 HTML |
| **Fallback model 在前后端各 4 个独立字段** | 不是数组结构，每个平台一个字段（如 fallback_model_deepseek），需要前后端同步新增 |
| **Ent schema 有 4 个独立的 validator/enum** | user_platform_quota（validator switch）、channel_monitor（provider enum）、channel_monitor_request_template（provider enum）、group（注释） |
| **Antigravity 在 LiteLLM 映射中复用 Anthropic** | DeepSeek 与 OpenAI API 兼容，路由也建议复用 OpenAI handler |
| **DeepSeek 模型已存在于测试文件中** | `reasoning_content` 处理代码已有 `deepseek-*` 相关的测试

---

## 附录：完整文件改动清单

### 前端（~24 个文件）

```
# 类型定义（根源 — 3 处）
frontend/src/types/index.ts                            ← GroupPlatform (L490), AccountPlatform (L693)
frontend/src/api/admin/users.ts                        ← PlatformQuotaPlatform (L310)
frontend/src/api/admin/settings.ts                     ← PlatformType (L20), PLATFORMS (L33)
frontend/src/api/admin/settings.ts                     ← fallback_model_* 字段 × 2 (L535-536, L787-790)

# 平台硬编码下拉列表（~12 个页面/组件）
frontend/src/views/admin/GroupsView.vue                ← platformOptions, platformFilterOptions
frontend/src/views/admin/ChannelsView.vue              ← platformOrder (渠道配置, L763)
frontend/src/views/admin/SubscriptionsView.vue         ← platformFilterOptions
frontend/src/views/admin/ops/.../OpsDashboardHeader.vue ← platformOptions
frontend/src/views/admin/SettingsView.vue              ← inline 平台数组 × 2
frontend/src/components/admin/account/AccountTableFilters.vue ← pOpts
frontend/src/components/admin/user/UserPlatformQuotaModal.vue ← PLATFORMS (L131)
frontend/src/components/user/UserPlatformQuotaCell.vue ← PLATFORM_ORDER (L37)
frontend/src/components/user/dashboard/UserDashboardStats.vue ← PLATFORM_ORDER (L281) + PLATFORM_LABELS (L250)
frontend/src/components/user/PlatformUsageBreakdown.vue ← PLATFORM_LABELS (L93)
frontend/src/components/admin/ErrorPassthroughRulesModal.vue ← platformOptions
frontend/src/components/account/ModelWhitelistSelector.vue ← upstreamSyncPlatforms

# Channel Monitor Provider（3 处 ⚠️ 独立类型体系）
frontend/src/api/admin/channelMonitor.ts                ← Provider 类型 (L8)
frontend/src/constants/channelMonitor.ts                ← PROVIDER_DEEPSEEK + PROVIDERS 数组
frontend/src/composables/useChannelMonitorFormat.ts      ← 4 个 switch 函数

# 账号创建/编辑
frontend/src/components/account/CreateAccountModal.vue  ← 硬编码按钮 (L70-150)
frontend/src/components/account/EditAccountModal.vue    ← baseUrl/apiKey placeholder

# 颜色/图标/样式
frontend/src/utils/platformColors.ts                   ← Platform type + 12 个 Record
frontend/src/components/common/PlatformIcon.vue        ← SVG logo
frontend/src/components/common/PlatformTypeBadge.vue   ← label + class
frontend/src/components/admin/channel/types.ts         ← getPlatformTagClass/TextClass

# i18n
frontend/src/i18n/locales/zh.ts                        ← 平台标签、向导文本
frontend/src/i18n/locales/en.ts                        ← 同上

# 测试（需更新断言）
frontend/src/views/admin/__tests__/SettingsView.spec.ts ← 平台列表断言 (L1119)
frontend/src/components/user/__tests__/UserPlatformQuotaCell.spec.ts ← 平台数组断言
```

### 后端（~14 个文件）

```
# 域常量（源头）
backend/internal/domain/constants.go                   ← PlatformDeepSeek (L20-24)
backend/internal/service/domain_constants.go            ← Alias (L39-44) + AllowedQuotaPlatforms (L49-54) + SettingKey
backend/internal/service/domain_constants.go            ← SettingKeyFallbackModelDeepSeek (L312-315)

# Ent Schema 校验（4 处 — 需 migration）
backend/ent/schema/group.go                            ← 枚举注释
backend/ent/schema/user_platform_quota.go              ← validator switch (L41-44)
backend/ent/schema/channel_monitor.go                   ← provider enum (L38)
backend/ent/schema/channel_monitor_request_template.go  ← provider enum (L42)

# 路由/调度
backend/internal/server/routes/gateway.go              ← 路由分支判断
backend/internal/service/proxy.go                      ← 上游转发路径

# 渠道/定价
backend/internal/handler/admin/channel_handler.go      ← platformToLiteLLMProvider (L507)
backend/internal/service/pricing_service.go             ← 定价数据源依赖
backend/internal/service/upstream_models.go             ← 模型同步

# Fallback Model + Settings
backend/internal/service/settings_view.go               ← FallbackModelDeepSeek (L160-163)
backend/internal/service/setting_service.go             ← fallback model 读写
backend/internal/service/admin_service.go               ← getAccountPlatform() (L3813)
backend/internal/service/admin_service.go               ← proxyQualityTargets (L501)
```

---

> **下一版改进方向**：前端应考虑统一平台注册中心（类似 `platformColors.ts` 中已有 `Platform` 类型，但硬编码列表仍未收敛），将所有页面的下拉选项集中到一处管理。后端应考虑平台枚举统一生成（从 domain constants 自动生成 ent schema enum values），避免 4 个 schema 文件的手工同步。
