/**
 * 企业管理 Admin API
 * 对接后端 /api/admin/enterprises 真实接口
 */

import { apiClient } from '../client'
import type {
  Enterprise,
  EnterpriseDetail,
  CreateEnterpriseRequest,
  UpdateEnterpriseRequest,
  ToggleEnterpriseResponse,
  DeleteEnterpriseResponse,
  EnterpriseListParams,
  PaginatedResponse,
} from '@/types/enterprise'

export const enterpriseAPI = {
  /** 获取企业列表 */
  async list(params?: EnterpriseListParams): Promise<PaginatedResponse<Enterprise>> {
    return apiClient.get('/admin/enterprises', { params })
  },

  /** 获取企业详情 */
  async getById(id: number): Promise<EnterpriseDetail> {
    return apiClient.get(`/admin/enterprises/${id}`)
  },

  /** 创建企业 */
  async create(data: CreateEnterpriseRequest): Promise<Enterprise> {
    return apiClient.post('/admin/enterprises', data)
  },

  /** 编辑企业信息 */
  async update(id: number, data: UpdateEnterpriseRequest): Promise<Enterprise> {
    return apiClient.put(`/admin/enterprises/${id}`, data)
  },

  /** 启停企业 */
  async toggleStatus(id: number): Promise<ToggleEnterpriseResponse> {
    return apiClient.post(`/admin/enterprises/${id}/toggle`)
  },

  /** 删除企业 */
  async delete(id: number): Promise<DeleteEnterpriseResponse> {
    return apiClient.delete(`/admin/enterprises/${id}`)
  },
}

export default enterpriseAPI

// Re-export enterprise types for convenience
export type {
  Enterprise,
  EnterpriseDetail,
  EnterpriseStatus,
  EnterpriseScale,
  EnterpriseIndustry,
} from '@/types/enterprise'
