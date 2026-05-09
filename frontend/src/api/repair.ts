// Layer 1: API calls — depends on types (Layer 0) only
import axios from 'axios'
import type { Repair, RepairStatus, RepairWithRoom, ApiResponse, PaginatedResult } from '@/types'

const client = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
})

// Add auth token to requests
client.interceptors.request.use((config) => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export interface CreateRepairDto {
  room_id: string
  type: Repair['type']
  description: string
  priority: Repair['priority']
}

export interface UpdateRepairDto {
  type?: Repair['type']
  description?: string
  priority?: Repair['priority']
  scheduled_at?: string | null
  remark?: string | null
}

export interface UpdateStatusDto {
  status: RepairStatus
  repairer_id?: string
  cost?: number | null
}

export interface RateRepairDto {
  rating: number
}

export interface ListRepairsParams {
  page?: number
  page_size?: number
  status?: RepairStatus
  room_id?: string
  priority?: Repair['priority']
}

export const repairApi = {
  list(params: ListRepairsParams = {}) {
    const { page = 1, page_size = 20, ...filters } = params
    return client.get<ApiResponse<PaginatedResult<RepairWithRoom>>>('/repairs', {
      params: { page, page_size, ...filters },
    })
  },

  get(id: string) {
    return client.get<ApiResponse<RepairWithRoom>>(`/repairs/${id}`)
  },

  create(data: CreateRepairDto) {
    return client.post<ApiResponse<RepairWithRoom>>('/repairs', data)
  },

  update(id: string, data: UpdateRepairDto) {
    return client.put<ApiResponse<RepairWithRoom>>(`/repairs/${id}`, data)
  },

  updateStatus(id: string, data: UpdateStatusDto) {
    return client.put<ApiResponse<RepairWithRoom>>(`/repairs/${id}/status`, data)
  },

  rate(id: string, data: RateRepairDto) {
    return client.put<ApiResponse<RepairWithRoom>>(`/repairs/${id}/rating`, data)
  },

  delete(id: string) {
    return client.delete<ApiResponse<null>>(`/repairs/${id}`)
  },
}