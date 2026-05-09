// Layer 1: API calls — depends on types (Layer 0) only
import axios from 'axios'
import type { Room, ApiResponse, PaginatedResult } from '@/types'

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

export interface CreateRoomDto {
  building_id: string
  number: string
  floor: number
  type: 'standard' | 'suite' | 'triple'
  capacity: number
  has_bathroom: boolean
  has_ac: boolean
}

export interface UpdateRoomDto {
  number?: string
  floor?: number
  type?: 'standard' | 'suite' | 'triple'
  capacity?: number
  has_bathroom?: boolean
  has_ac?: boolean
  status?: 'available' | 'full' | 'maintenance'
}

export const roomApi = {
  list(page = 1, pageSize = 20) {
    return client.get<ApiResponse<PaginatedResult<Room>>>('/rooms', {
      params: { page, page_size: pageSize },
    })
  },

  get(id: string) {
    return client.get<ApiResponse<Room>>(`/rooms/${id}`)
  },

  create(data: CreateRoomDto) {
    return client.post<ApiResponse<Room>>('/rooms', data)
  },

  update(id: string, data: UpdateRoomDto) {
    return client.put<ApiResponse<Room>>(`/rooms/${id}`, data)
  },

  delete(id: string) {
    return client.delete<ApiResponse<null>>(`/rooms/${id}`)
  },

  listByBuilding(buildingId: string) {
    return client.get<ApiResponse<Room[]>>(`/buildings/${buildingId}/rooms`)
  },

  listAvailable(gender?: string) {
    return client.get<ApiResponse<Room[]>>('/rooms/available', {
      params: gender ? { gender } : {},
    })
  },
}