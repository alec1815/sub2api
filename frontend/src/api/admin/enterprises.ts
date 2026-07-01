/**
 * Admin Enterprises API endpoints
 * Handles enterprise management for administrators
 */

import { apiClient } from '../client'
import type { BasePaginationResponse } from '@/types'
import type {
  Enterprise,
  EnterpriseDetail,
  CreateEnterpriseRequest,
  UpdateEnterpriseRequest,
  ToggleEnterpriseResponse,
  DeleteEnterpriseResponse,
  EnterpriseStatus,
} from '@/types/enterprise'

/**
 * List all enterprises with pagination
 * @param page - Page number (default: 1)
 * @param pageSize - Items per page (default: 20)
 * @param filters - Optional filters (search, status)
 * @param options - Optional request options (signal)
 * @returns Paginated list of enterprises
 */
export async function list(
  page: number = 1,
  pageSize: number = 20,
  filters?: {
    search?: string
    status?: EnterpriseStatus
  },
  options?: {
    signal?: AbortSignal
  }
): Promise<BasePaginationResponse<Enterprise>> {
  const params: Record<string, unknown> = {
    page,
    page_size: pageSize,
    ...filters,
  }
  const { data } = await apiClient.get<BasePaginationResponse<Enterprise>>(
    '/admin/enterprises',
    { params, signal: options?.signal }
  )
  return data
}

/**
 * Get enterprise detail by ID
 */
export async function getById(id: number): Promise<EnterpriseDetail> {
  const { data } = await apiClient.get<EnterpriseDetail>(`/admin/enterprises/${id}`)
  return data
}

/**
 * Create a new enterprise
 */
export async function create(
  enterprise: CreateEnterpriseRequest
): Promise<Enterprise> {
  const { data } = await apiClient.post<Enterprise>('/admin/enterprises', enterprise)
  return data
}

/**
 * Update an existing enterprise
 */
export async function update(
  id: number,
  updates: UpdateEnterpriseRequest
): Promise<Enterprise> {
  const { data } = await apiClient.put<Enterprise>(`/admin/enterprises/${id}`, updates)
  return data
}

/**
 * Toggle enterprise status (active/disabled)
 */
export async function toggleStatus(id: number): Promise<ToggleEnterpriseResponse> {
  const { data } = await apiClient.post<ToggleEnterpriseResponse>(
    `/admin/enterprises/${id}/toggle`
  )
  return data
}

/**
 * Delete an enterprise
 */
export async function deleteEnterprise(id: number): Promise<DeleteEnterpriseResponse> {
  const { data } = await apiClient.delete<DeleteEnterpriseResponse>(
    `/admin/enterprises/${id}`
  )
  return data
}

/**
 * Update enterprise balance (add / subtract / set)
 */
export async function updateBalance(
  id: number,
  balance: number,
  operation: 'set' | 'add' | 'subtract' = 'add',
  notes?: string
): Promise<Enterprise> {
  const { data } = await apiClient.post<Enterprise>(`/admin/enterprises/${id}/balance`, {
    balance,
    operation,
    notes: notes || ''
  })
  return data
}

/**
 * Get enterprise balance history
 */
export async function getBalanceHistory(
  id: number,
  page: number = 1,
  pageSize: number = 20
): Promise<BasePaginationResponse<{
  id: number
  amount: number
  operation: string
  notes: string
  created_at: string
}>> {
  const { data } = await apiClient.get(`/admin/enterprises/${id}/balance-history`, {
    params: { page, page_size: pageSize }
  })
  return data
}

/**
 * List API keys for an enterprise
 */
export async function getEnterpriseKeys(
  id: number,
  page: number = 1,
  pageSize: number = 20
): Promise<BasePaginationResponse<{ id: number; name: string; key: string; status: string; created_at: string }>> {
  const { data } = await apiClient.get(`/admin/enterprises/${id}/api-keys`, {
    params: { page, page_size: pageSize }
  })
  return data
}

const enterprisesAPI = {
  list,
  getById,
  create,
  update,
  toggleStatus,
  delete: deleteEnterprise,
  updateBalance,
  getBalanceHistory,
  getEnterpriseKeys,
}

export default enterprisesAPI
