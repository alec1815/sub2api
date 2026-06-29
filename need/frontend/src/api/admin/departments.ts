/**
 * 部门管理 Admin API
 * 对接后端 /api/admin/departments 真实接口
 */

import { apiClient } from '../client'
import type {
  Department,
  CreateDepartmentRequest,
  UpdateDepartmentRequest,
} from '@/types/enterprise'

export const departmentAPI = {
  /** 获取部门树 */
  async getTree(enterpriseId: number): Promise<Department[]> {
    return apiClient.get('/admin/departments', {
      params: { enterprise_id: enterpriseId },
    })
  },

  /** 创建部门 */
  async create(data: CreateDepartmentRequest): Promise<Department> {
    return apiClient.post('/admin/departments', data)
  },

  /** 编辑部门 */
  async update(id: number, data: UpdateDepartmentRequest): Promise<Department> {
    return apiClient.put(`/admin/departments/${id}`, data)
  },

  /** 删除部门 */
  async delete(id: number): Promise<{ deleted: boolean }> {
    return apiClient.delete(`/admin/departments/${id}`)
  },
}

export default departmentAPI

export type {
  Department,
  DepartmentStatus,
} from '@/types/enterprise'
