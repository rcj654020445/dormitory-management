// Layer 1: API calls — depends on types (Layer 0) only
import axios from 'axios'
import type { Student, ApiResponse, PaginatedResult } from '@/types'

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

export interface CreateStudentDto {
  student_id: string
  name: string
  gender: 'male' | 'female'
  phone: string
  email: string
  major: string
  grade: number
}

export interface UpdateStudentDto {
  name?: string
  gender?: 'male' | 'female'
  phone?: string
  email?: string
  major?: string
  grade?: number
  status?: string
}

export interface AllocateRoomDto {
  room_id: string
  bed_number: number
}

export const studentApi = {
  list(page = 1, pageSize = 20) {
    return client.get<ApiResponse<PaginatedResult<Student>>>('/students', {
      params: { page, page_size: pageSize },
    })
  },

  get(id: string) {
    return client.get<ApiResponse<Student>>(`/students/${id}`)
  },

  create(data: CreateStudentDto) {
    return client.post<ApiResponse<Student>>('/students', data)
  },

  update(id: string, data: UpdateStudentDto) {
    return client.put<ApiResponse<Student>>(`/students/${id}`, data)
  },

  delete(id: string) {
    return client.delete<ApiResponse<null>>(`/students/${id}`)
  },

  allocateRoom(id: string, data: AllocateRoomDto) {
    return client.post<ApiResponse<{ message: string }>>(`/students/${id}/allocate`, data)
  },

  vacate(id: string) {
    return client.post<ApiResponse<{ message: string }>>(`/students/${id}/vacate`)
  },
}
