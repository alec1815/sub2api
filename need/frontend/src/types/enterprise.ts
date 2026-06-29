/**
 * 企业功能 — TypeScript 类型定义
 * 严格对照 05-API接口文档.md V1.0
 */

// ==================== 通用枚举 ====================

/** 企业规模 */
export type EnterpriseScale = 'micro' | 'small' | 'medium' | 'large'

/** 企业行业 */
export type EnterpriseIndustry =
  | 'internet'
  | 'finance'
  | 'education'
  | 'healthcare'
  | 'manufacturing'
  | 'other'

/** 企业状态 */
export type EnterpriseStatus = 'active' | 'disabled'

/** 企业成员角色 */
export type EnterpriseMemberRole = 'enterprise_admin' | 'enterprise_member'

/** 企业成员状态 */
export type EnterpriseMemberStatus = 'active' | 'pending' | 'unbound'

/** 套餐状态 */
export type SubscriptionStatus = 'active' | 'expired' | 'suspended'

/** API Key 绑定工具 */
export type BoundTool = 'cursor' | 'trae' | 'claude_code' | 'codex' | 'opencode' | 'pixso' | 'other'

/** 消费池类型 */
export type PoolType = 'personal' | 'enterprise'

/** 部门状态 */
export type DepartmentStatus = 'active' | 'disabled'

// ==================== 企业 (Enterprise) ====================

/** 企业列表项（GET /api/admin/enterprises） */
export interface Enterprise {
  id: number
  name: string
  short_name: string
  credit_code: string
  address: string
  scale: EnterpriseScale
  industry: EnterpriseIndustry
  parent_id: number
  status: EnterpriseStatus
  contact_name: string
  contact_phone: string
  contact_email: string
  notes: string
  balance: string // decimal(20,8) as string
  total_recharged: string
  admin_user_id: number
  admin_email: string
  member_count: number
  created_at: string
  updated_at: string
}

/** 企业订阅套餐 */
export interface EnterpriseSubscription {
  id: number
  plan_name: string
  plan_id: number
  group_name: string
  group_id: number
  status: SubscriptionStatus
  daily_usage_usd: string
  weekly_usage_usd?: string
  monthly_usage_usd: string
  starts_at: string
  expires_at: string
}

/** 企业详情（GET /api/admin/enterprises/:id） */
export interface EnterpriseDetail extends Enterprise {
  subscriptions: EnterpriseSubscription[]
}

/** 创建企业请求 */
export interface CreateEnterpriseRequest {
  name: string
  short_name?: string
  credit_code?: string
  address?: string
  scale?: EnterpriseScale
  industry?: EnterpriseIndustry
  contact_name?: string
  contact_phone?: string
  contact_email?: string
  admin_email: string
  admin_name?: string
  notes?: string
  parent_id?: number
}

/** 更新企业请求（全部可选，partial update） */
export interface UpdateEnterpriseRequest {
  name?: string
  short_name?: string
  credit_code?: string
  address?: string
  scale?: EnterpriseScale
  industry?: EnterpriseIndustry
  contact_name?: string
  contact_phone?: string
  contact_email?: string
  notes?: string
}

/** 启停企业响应 */
export interface ToggleEnterpriseResponse {
  id: number
  status: EnterpriseStatus
  affected_keys: number
}

/** 删除企业响应 */
export interface DeleteEnterpriseResponse {
  id: number
  deleted: boolean
  unbound_members: number
}

/** 企业列表查询参数 */
export interface EnterpriseListParams {
  page?: number
  page_size?: number
  search?: string
  status?: EnterpriseStatus
}

// ==================== 企业成员 (EnterpriseMember) ====================

/** 企业成员 */
export interface EnterpriseMember {
  id: number
  user_id: number
  name: string
  email: string
  role: EnterpriseMemberRole
  status: EnterpriseMemberStatus
  department_id: number | null
  department_name: string
  concurrency: number
  rpm_limit: number
  notes: string
  joined_at: string
  unbound_at: string | null
  last_active_at: string | null
}

/** 创建成员请求 */
export interface CreateMemberRequest {
  email: string
  username?: string
  password?: string
  department_id?: number
  concurrency?: number
  rpm_limit?: number
}

/** 更新成员请求（全部可选） */
export interface UpdateMemberRequest {
  username?: string
  password?: string
  department_id?: number
  concurrency?: number
  rpm_limit?: number
  notes?: string
}

/** 解绑成员响应 */
export interface UnbindMemberResponse {
  id: number
  status: 'unbound'
  unbound_at: string
  disabled_keys: number
}

/** 成员列表查询参数 */
export interface MemberListParams {
  page?: number
  page_size?: number
  search?: string
  status?: EnterpriseMemberStatus
  role?: EnterpriseMemberRole
  department_id?: number
}

// ==================== 部门 (Department) ====================

/** 部门节点（树形） */
export interface Department {
  id: number
  parent_id: number
  name: string
  order_num: number
  leader: string
  phone: string
  email: string
  status: DepartmentStatus
  member_count: number
  created_at: string
  children: Department[]
}

/** 创建部门请求 */
export interface CreateDepartmentRequest {
  enterprise_id: number
  parent_id?: number
  name: string
  order_num?: number
  leader?: string
  phone?: string
  email?: string
}

/** 更新部门请求 */
export interface UpdateDepartmentRequest {
  name?: string
  parent_id?: number
  order_num?: number
  leader?: string
  phone?: string
  email?: string
  status?: DepartmentStatus
}

// ==================== 企业财务 (Enterprise Finance) ====================

/** 用量按维度拆分项 */
export interface UsageByDimension {
  member_id?: number
  member_name?: string
  key_id?: number
  key_name?: string
  model?: string
  tool?: string
  cost: string
  calls: number
}

/** 企业财务汇总 */
export interface EnterpriseFinance {
  balance: string
  total_recharged: string
  subscriptions: EnterpriseSubscription[]
  usage_summary: {
    total_cost: string
    total_calls: number
    by_member: UsageByDimension[]
    by_model: UsageByDimension[]
    by_key: UsageByDimension[]
    by_tool: UsageByDimension[]
  }
}

/** 企业用量明细项 */
export interface UsageDetailItem {
  id: number
  api_key_id: number
  key_name: string
  user_id: number
  user_name: string
  model: string
  requested_model: string
  total_cost: string
  prompt_tokens: number
  completion_tokens: number
  pool_type: PoolType
  created_at: string
}

/** 用量明细查询参数 */
export interface UsageListParams {
  page?: number
  page_size?: number
  start_date?: string
  end_date?: string
  member_id?: number
  key_id?: number
  model?: string
  bound_tool?: string
}

// ==================== 企业充值 / 套餐 ====================

/** 企业充值请求 */
export interface EnterpriseRechargeRequest {
  amount: number
  payment_method: 'alipay' | 'wechat' | 'stripe'
}

/** 企业充值响应 */
export interface EnterpriseRechargeResponse {
  order_id: string
  amount: number
  payment_url: string
  qr_code?: string
}

/** 企业购买套餐请求 */
export interface EnterpriseSubscribeRequest {
  plan_id: number
  group_id: number
}

/** 企业购买套餐响应 */
export interface EnterpriseSubscribeResponse {
  id: number
  enterprise_id: number
  plan_name: string
  group_name: string
  starts_at: string
  expires_at: string
  status: SubscriptionStatus
}

// ==================== 企业设置 (Enterprise Settings) ====================

/** 企业设置（管理员视角） */
export interface EnterpriseSettings {
  id: number
  name: string
  short_name: string
  credit_code: string
  address: string
  scale: EnterpriseScale
  industry: EnterpriseIndustry
  contact_name: string
  contact_phone: string
  contact_email: string
  notes: string
  status: EnterpriseStatus
  balance: string
  member_count: number
  created_at: string
}

/** 更新企业设置请求 */
export interface UpdateEnterpriseSettingsRequest {
  name?: string
  short_name?: string
  address?: string
  contact_name?: string
  contact_phone?: string
  contact_email?: string
  notes?: string
}

// ==================== 企业 Profile ====================

/** 企业 Profile（成员视角） */
export interface EnterpriseProfile {
  enterprise: {
    id: number
    name: string
    short_name: string
    address: string
    scale: EnterpriseScale
    industry: EnterpriseIndustry
    contact_name: string
    contact_phone: string
    contact_email: string
    created_at: string
  }
  my_role: EnterpriseMemberRole
  my_department: string
  my_joined_at: string
  monthly_usage: {
    total_calls: number
    total_cost: string
  }
}

// ==================== 企业密钥 (Enterprise Key) ====================

/** 密钥分组 */
export interface KeyGroup {
  id: number
  name: string
}

/** 企业密钥列表项 */
export interface EnterpriseKey {
  id: number
  name: string
  key: string
  key_prefix: string
  status: 'active' | 'disabled'
  assigned_to: number | null
  assigned_member_name: string
  assigned_member_email: string
  groups: KeyGroup[]
  usage_purpose: string
  bound_tool: BoundTool
  quota: string
  quota_used: string
  expires_at: string | null
  created_at: string
}

/** 创建企业密钥请求 */
export interface CreateEnterpriseKeyRequest {
  name: string
  group_ids: number[]
  assigned_to?: number
  usage_purpose?: string
  bound_tool?: BoundTool
  quota?: number
  expires_at?: string
}

/** 创建密钥响应 */
export interface CreateEnterpriseKeyResponse {
  id: number
  name: string
  key: string
  key_full: string // 仅在创建时返回完整 key
  status: 'active'
  assigned_to: number | null
  groups: KeyGroup[]
  created_at: string
}

/** 更新密钥请求 */
export interface UpdateEnterpriseKeyRequest {
  name?: string
  group_ids?: number[]
  assigned_to?: number | null
  usage_purpose?: string
  bound_tool?: BoundTool
  quota?: number
}

// ==================== 通用分页 ====================

/** 分页响应 */
export interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  pages: number
}
