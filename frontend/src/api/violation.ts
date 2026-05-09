// Layer 1: API calls — depends on types (Layer 0) only
import axios from 'axios'
import type { Violation, ApiResponse, PaginatedResult } from '@/types'

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

export interface CreateViolationDto {
  student_id: string
  type: 'late_return' | 'noise' | 'damage' | 'property_violation' | 'other'
  description: string
  points: number
  handled_by: string
}

export interface UpdateViolationDto {
  type?: 'late_return' | 'noise' | 'damage' | 'property_violation' | 'other'
  description?: string
  points?: number
  handled_by?: string
  status?: 'pending' | 'resolved'
}

export interface ResolveViolationDto {
  status: 'resolved'
}

export interface ListViolationQuery {
  student_id?: string
  type?: string
  status?: string
  page?: number
  page_size?: number
}

export const violationApi = {
  list(query: ListViolationQuery = {}) {
    return client.get<ApiResponse<PaginatedResult<Violation>>>('/violations', {
      params: {
        student_id: query.student_id,
        type: query.type,
        status: query.status,
        page: query.page ?? 1,
        page_size: query.page_size ?? 20,
      },
    })
  },

  get(id: string) {
    return client.get<ApiResponse<Violation>>(`/violations/${id}`)
  },

  create(data: CreateViolationDto) {
    return client.post<ApiResponse<Violation>>('/violations', data)
  },

  update(id: string, data: UpdateViolationDto) {
    return client.put<ApiResponse<Violation>>(`/violations/${id}`, data)
  },

  delete(id: string) {
    return client.delete<ApiResponse<null>>(`/violations/${id}`)
  },

  resolve(id: string, data: ResolveViolationDto) {
    return client.put<ApiResponse<Violation>>(`/violations/${id}/resolve`, data)
  },
}
