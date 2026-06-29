/**
 * 企业成员管理 Admin API
 * 对接后端 /api/admin/enterprises/:id/members 真实接口
 */

import { apiClient } from '../client'
import type {
  EnterpriseMember,
  CreateMemberRequest,
  UpdateMemberRequest,
  UnbindMemberResponse,
  MemberListParams,
  PaginatedResponse,
} from '@/types/enterprise'

export const enterpriseMemberAPI = {
  /** 获取企业成员列表 */
  async list(enterpriseId: number, params?: MemberListParams): Promise<PaginatedResponse<EnterpriseMember>> {
    return apiClient.get(`/admin/enterprises/${enterpriseId}/members`, { params })
  },

  /** 创建企业成员（代注册） */
  async create(enterpriseId: number, data: CreateMemberRequest): Promise<EnterpriseMember> {
    return apiClient.post(`/admin/enterprises/${enterpriseId}/members`, data)
  },

  /** 编辑成员信息 */
  async update(enterpriseId: number, memberId: number, data: UpdateMemberRequest): Promise<EnterpriseMember> {
    return apiClient.put(`/admin/enterprises/${enterpriseId}/members/${memberId}`, data)
  },

  /** 解绑成员 */
  async unbind(enterpriseId: number, memberId: number): Promise<UnbindMemberResponse> {
    return apiClient.delete(`/admin/enterprises/${enterpriseId}/members/${memberId}`)
  },
}

export default enterpriseMemberAPI

export type {
  EnterpriseMember,
  EnterpriseMemberRole,
  EnterpriseMemberStatus,
} from '@/types/enterprise'
