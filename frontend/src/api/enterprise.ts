/**
 * Enterprise Self-Admin API endpoints
 * Handles enterprise-level operations for enterprise admins (not system admins)
 * Base path: /api/enterprise
 */

import { apiClient } from './client'
import type { BasePaginationResponse } from '@/types'
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
  UsageDetailItem,
  UsageListParams,
  EnterpriseRechargeRequest,
  EnterpriseRechargeResponse,
  EnterpriseSubscribeRequest,
  EnterpriseSubscribeResponse,
  Department,
  CreateDepartmentRequest,
  UpdateDepartmentRequest,
  EnterpriseSettings,
  UpdateEnterpriseSettingsRequest,
  EnterpriseProfile,
} from '@/types/enterprise'

// ==================== Members ====================

export async function listMembers(
  page: number = 1,
  pageSize: number = 20,
  filters?: Omit<MemberListParams, 'page' | 'page_size'>,
  options?: { signal?: AbortSignal }
): Promise<BasePaginationResponse<EnterpriseMember>> {
  const { data } = await apiClient.get<BasePaginationResponse<EnterpriseMember>>(
    '/enterprise/members',
    { params: { page, page_size: pageSize, ...filters }, signal: options?.signal }
  )
  return data
}

export async function createMember(member: CreateMemberRequest): Promise<EnterpriseMember> {
  const { data } = await apiClient.post<EnterpriseMember>('/enterprise/members', member)
  return data
}

export async function updateMember(
  memberId: number,
  updates: UpdateMemberRequest
): Promise<EnterpriseMember> {
  const { data } = await apiClient.put<EnterpriseMember>(
    `/enterprise/members/${memberId}`,
    updates
  )
  return data
}

export async function unbindMember(memberId: number): Promise<UnbindMemberResponse> {
  const { data } = await apiClient.delete<UnbindMemberResponse>(
    `/enterprise/members/${memberId}`
  )
  return data
}

// ==================== Keys ====================

export async function listKeys(
  page: number = 1,
  pageSize: number = 20,
  options?: { signal?: AbortSignal }
): Promise<BasePaginationResponse<EnterpriseKey>> {
  const { data } = await apiClient.get<BasePaginationResponse<EnterpriseKey>>(
    '/enterprise/keys',
    { params: { page, page_size: pageSize }, signal: options?.signal }
  )
  return data
}

export async function createKey(
  keyData: CreateEnterpriseKeyRequest
): Promise<CreateEnterpriseKeyResponse> {
  const { data } = await apiClient.post<CreateEnterpriseKeyResponse>(
    '/enterprise/keys',
    keyData
  )
  return data
}

export async function updateKey(
  keyId: number,
  updates: UpdateEnterpriseKeyRequest
): Promise<EnterpriseKey> {
  const { data } = await apiClient.put<EnterpriseKey>(
    `/enterprise/keys/${keyId}`,
    updates
  )
  return data
}

export async function deleteKey(keyId: number): Promise<void> {
  await apiClient.delete(`/enterprise/keys/${keyId}`)
}

export async function toggleKey(keyId: number): Promise<EnterpriseKey> {
  const { data } = await apiClient.post<EnterpriseKey>(`/enterprise/keys/${keyId}/toggle`)
  return data
}

// ==================== Finance ====================

export async function getFinance(): Promise<EnterpriseFinance> {
  const { data } = await apiClient.get<EnterpriseFinance>('/enterprise/finance')
  return data
}

export async function getUsage(
  page: number = 1,
  pageSize: number = 20,
  filters?: Omit<UsageListParams, 'page' | 'page_size'>,
  options?: { signal?: AbortSignal }
): Promise<BasePaginationResponse<UsageDetailItem>> {
  const { data } = await apiClient.get<BasePaginationResponse<UsageDetailItem>>(
    '/enterprise/usage',
    { params: { page, page_size: pageSize, ...filters }, signal: options?.signal }
  )
  return data
}

export async function recharge(
  req: EnterpriseRechargeRequest
): Promise<EnterpriseRechargeResponse> {
  const { data } = await apiClient.post<EnterpriseRechargeResponse>(
    '/enterprise/recharge',
    req
  )
  return data
}

export async function subscribe(
  req: EnterpriseSubscribeRequest
): Promise<EnterpriseSubscribeResponse> {
  const { data } = await apiClient.post<EnterpriseSubscribeResponse>(
    '/enterprise/subscribe',
    req
  )
  return data
}

// ==================== Departments ====================

export async function getDepartmentTree(): Promise<Department[]> {
  const { data } = await apiClient.get<Department[]>('/enterprise/departments')
  return data
}

export async function createDepartment(
  dept: CreateDepartmentRequest
): Promise<Department> {
  const { data } = await apiClient.post<Department>('/enterprise/departments', dept)
  return data
}

export async function updateDepartment(
  deptId: number,
  updates: UpdateDepartmentRequest
): Promise<Department> {
  const { data } = await apiClient.put<Department>(
    `/enterprise/departments/${deptId}`,
    updates
  )
  return data
}

export async function deleteDepartment(deptId: number): Promise<void> {
  await apiClient.delete(`/enterprise/departments/${deptId}`)
}

// ==================== Settings ====================

export async function getSettings(): Promise<EnterpriseSettings> {
  const { data } = await apiClient.get<EnterpriseSettings>('/enterprise/settings')
  return data
}

export async function updateSettings(
  updates: UpdateEnterpriseSettingsRequest
): Promise<EnterpriseSettings> {
  const { data } = await apiClient.put<EnterpriseSettings>(
    '/enterprise/settings',
    updates
  )
  return data
}

// ==================== Profile ====================

export async function getProfile(): Promise<EnterpriseProfile> {
  const { data } = await apiClient.get<EnterpriseProfile>('/enterprise/profile')
  return data
}

const enterpriseAdminAPI = {
  // Members
  listMembers,
  createMember,
  updateMember,
  unbindMember,
  // Keys
  listKeys,
  createKey,
  updateKey,
  deleteKey,
  toggleKey,
  // Finance
  getFinance,
  getUsage,
  recharge,
  subscribe,
  // Departments
  getDepartmentTree,
  createDepartment,
  updateDepartment,
  deleteDepartment,
  // Settings
  getSettings,
  updateSettings,
  // Profile
  getProfile,
}

export default enterpriseAdminAPI
