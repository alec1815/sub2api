/**
 * Admin Enterprise Members API endpoints
 * Handles member management within an enterprise for administrators
 */

import { apiClient } from '../client'
import type { BasePaginationResponse } from '@/types'
import type {
  EnterpriseMember,
  CreateMemberRequest,
  UpdateMemberRequest,
  UnbindMemberResponse,
  MemberListParams,
} from '@/types/enterprise'

/**
 * List members of an enterprise
 * @param enterpriseId - Enterprise ID
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @param filters - Optional filters
 * @param options - Optional request options (signal)
 * @returns Paginated list of members
 */
export async function list(
  enterpriseId: number,
  page: number = 1,
  pageSize: number = 20,
  filters?: Omit<MemberListParams, 'page' | 'page_size'>,
  options?: {
    signal?: AbortSignal
  }
): Promise<BasePaginationResponse<EnterpriseMember>> {
  const params: Record<string, unknown> = {
    page,
    page_size: pageSize,
    ...filters,
  }
  const { data } = await apiClient.get<BasePaginationResponse<EnterpriseMember>>(
    `/admin/enterprises/${enterpriseId}/members`,
    { params, signal: options?.signal }
  )
  return data
}

/**
 * Create (register) a new member for an enterprise
 */
export async function create(
  enterpriseId: number,
  member: CreateMemberRequest
): Promise<EnterpriseMember> {
  const { data } = await apiClient.post<EnterpriseMember>(
    `/admin/enterprises/${enterpriseId}/members`,
    member
  )
  return data
}

/**
 * Update a member's information
 */
export async function update(
  enterpriseId: number,
  memberId: number,
  updates: UpdateMemberRequest
): Promise<EnterpriseMember> {
  const { data } = await apiClient.put<EnterpriseMember>(
    `/admin/enterprises/${enterpriseId}/members/${memberId}`,
    updates
  )
  return data
}

/**
 * Unbind (remove) a member from an enterprise
 */
export async function unbind(
  enterpriseId: number,
  memberId: number
): Promise<UnbindMemberResponse> {
  const { data } = await apiClient.delete<UnbindMemberResponse>(
    `/admin/enterprises/${enterpriseId}/members/${memberId}`
  )
  return data
}

const enterpriseMembersAPI = {
  list,
  create,
  update,
  unbind,
}

export default enterpriseMembersAPI
