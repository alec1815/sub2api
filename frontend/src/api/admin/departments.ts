/**
 * Admin Departments API endpoints
 * Handles department tree management for administrators
 */

import { apiClient } from '../client'
import type {
  Department,
  CreateDepartmentRequest,
  UpdateDepartmentRequest,
} from '@/types/enterprise'

/**
 * Get department tree for an enterprise
 * @param enterpriseId - Enterprise ID
 * @returns Department tree with children
 */
export async function getTree(enterpriseId: number): Promise<Department[]> {
  const { data } = await apiClient.get<Department[]>('/admin/departments', {
    params: { enterprise_id: enterpriseId },
  })
  return data
}

/**
 * Create a new department
 */
export async function create(dept: CreateDepartmentRequest): Promise<Department> {
  const { data } = await apiClient.post<Department>('/admin/departments', dept)
  return data
}

/**
 * Update a department
 */
export async function update(
  id: number,
  updates: UpdateDepartmentRequest
): Promise<Department> {
  const { data } = await apiClient.put<Department>(`/admin/departments/${id}`, updates)
  return data
}

/**
 * Delete a department
 */
export async function deleteDept(id: number): Promise<void> {
  await apiClient.delete(`/admin/departments/${id}`)
}

const departmentsAPI = {
  getTree,
  create,
  update,
  delete: deleteDept,
}

export default departmentsAPI
