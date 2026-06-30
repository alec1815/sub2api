/**
 * 企业管理员 API（/api/enterprise/*）
 * 对接后端企业管理员视角的全部接口
 */

import { apiClient } from './client'
import type {
  EnterpriseMember,
  CreateMemberRequest,
  UpdateMemberRequest,
  UnbindMemberResponse,
  MemberListParams,
  EnterpriseKey,
  CreateEnterpriseKeyRequest,
  CreateEnterpriseKeyResponse,
  UpdateEnterpriseKeyRequest,
  EnterpriseFinance,
  EnterpriseRechargeRequest,
  EnterpriseRechargeResponse,
  EnterpriseSubscribeRequest,
  EnterpriseSubscribeResponse,
  EnterpriseSettings,
  UpdateEnterpriseSettingsRequest,
  Department,
  CreateDepartmentRequest,
  UpdateDepartmentRequest,
  EnterpriseProfile,
  UsageDetailItem,
  UsageListParams,
  PaginatedResponse,
} from '@/types/enterprise'

export const enterpriseAdminAPI = {
  // ==================== 成员管理 ====================

  /** 成员列表 */
  async listMembers(params?: MemberListParams): Promise<PaginatedResponse<EnterpriseMember>> {
    return apiClient.get('/enterprise/members', { params })
  },

  /** 创建成员（代注册） */
  async createMember(data: CreateMemberRequest): Promise<EnterpriseMember> {
    return apiClient.post('/enterprise/members', data)
  },

  /** 编辑成员 */
  async updateMember(memberId: number, data: UpdateMemberRequest): Promise<EnterpriseMember> {
    return apiClient.put(`/enterprise/members/${memberId}`, data)
  },

  /** 解绑成员 */
  async unbindMember(memberId: number): Promise<UnbindMemberResponse> {
    return apiClient.delete(`/enterprise/members/${memberId}`)
  },

  // ==================== 密钥管理 ====================

  /** 企业密钥列表 */
  async listKeys(params?: {
    page?: number
    page_size?: number
    status?: string
    assigned_to?: number
    bound_tool?: string
  }): Promise<PaginatedResponse<EnterpriseKey>> {
    return apiClient.get('/enterprise/keys', { params })
  },

  /** 创建密钥 */
  async createKey(data: CreateEnterpriseKeyRequest): Promise<CreateEnterpriseKeyResponse> {
    return apiClient.post('/enterprise/keys', data)
  },

  /** 编辑密钥 */
  async updateKey(keyId: number, data: UpdateEnterpriseKeyRequest): Promise<EnterpriseKey> {
    return apiClient.put(`/enterprise/keys/${keyId}`, data)
  },

  /** 删除密钥 */
  async deleteKey(keyId: number): Promise<{ deleted: boolean }> {
    return apiClient.delete(`/enterprise/keys/${keyId}`)
  },

  /** 启停密钥 */
  async toggleKey(keyId: number): Promise<EnterpriseKey> {
    return apiClient.post(`/enterprise/keys/${keyId}/toggle`)
  },

  // ==================== 财务管理 ====================

  /** 企业财务汇总 */
  async getFinance(params?: { start_date?: string; end_date?: string }): Promise<EnterpriseFinance> {
    return apiClient.get('/enterprise/finance', { params })
  },

  /** 企业用量明细 */
  async getUsage(params?: UsageListParams): Promise<PaginatedResponse<UsageDetailItem>> {
    return apiClient.get('/enterprise/usage', { params })
  },

  /** 企业充值 */
  async recharge(data: EnterpriseRechargeRequest): Promise<EnterpriseRechargeResponse> {
    return apiClient.post('/enterprise/recharge', data)
  },

  /** 企业购买套餐 */
  async subscribe(data: EnterpriseSubscribeRequest): Promise<EnterpriseSubscribeResponse> {
    return apiClient.post('/enterprise/subscribe', data)
  },

  // ==================== 部门管理 ====================

  /** 部门树 */
  async getDepartmentTree(): Promise<Department[]> {
    return apiClient.get('/enterprise/departments')
  },

  /** 创建部门 */
  async createDepartment(data: Omit<CreateDepartmentRequest, 'enterprise_id'>): Promise<Department> {
    return apiClient.post('/enterprise/departments', data)
  },

  /** 编辑部门 */
  async updateDepartment(id: number, data: UpdateDepartmentRequest): Promise<Department> {
    return apiClient.put(`/enterprise/departments/${id}`, data)
  },

  /** 删除部门 */
  async deleteDepartment(id: number): Promise<{ deleted: boolean }> {
    return apiClient.delete(`/enterprise/departments/${id}`)
  },

  // ==================== 企业设置 ====================

  /** 获取企业设置 */
  async getSettings(): Promise<EnterpriseSettings> {
    return apiClient.get('/enterprise/settings')
  },

  /** 更新企业设置 */
  async updateSettings(data: UpdateEnterpriseSettingsRequest): Promise<EnterpriseSettings> {
    return apiClient.put('/enterprise/settings', data)
  },

  // ==================== 企业 Profile ====================

  /** 获取企业 Profile */
  async getProfile(): Promise<EnterpriseProfile> {
    return apiClient.get('/enterprise/profile')
  },
}

export default enterpriseAdminAPI
