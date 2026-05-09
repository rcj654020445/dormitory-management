// Layer 0: Core types — no internal imports allowed
// This module defines shared TypeScript interfaces for the dormitory system.

export interface Student {
  id: string
  student_id: string
  name: string
  gender: 'male' | 'female'
  phone: string
  email: string
  major: string
  grade: number
  room_id?: string
  check_in_at?: string
  status: 'pending' | 'checked_in' | 'graduated' | 'suspended'
  created_at: string
  updated_at: string
}

export interface Building {
  id: string
  name: string
  gender: 'male' | 'female'
  floor_count: number
  room_per_floor: number
  description: string
  status: 'active' | 'maintenance' | 'retired'
  created_at: string
  updated_at: string
}

export interface Room {
  id: string
  building_id: string
  number: string
  floor: number
  type: 'standard' | 'suite' | 'triple'
  beds_total: number
  beds_used: number
  has_bathroom: boolean
  has_ac: boolean
  status: 'available' | 'full' | 'maintenance'
  created_at: string
  updated_at: string
}

export interface Allocation {
  id: string
  student_id: string
  room_id: string
  bed_number: number
  status: 'active' | 'checked_out'
  check_in_at: string
  check_out_at?: string
  reason?: string
  created_at: string
}

export interface Violation {
  id: string
  student_id: string
  type: 'late_return' | 'noise' | 'damage' | 'property_violation' | 'other'
  description: string
  points: number
  handled_by: string
  handled_at: string
  status: 'pending' | 'resolved'
  created_at: string
}

export type RepairType =
  | 'facility'
  | 'plumbing'
  | 'electrical'
  | 'network'
  | 'cleaning'
  | 'other'

export type RepairStatus =
  | 'pending'
  | 'assigned'
  | 'repairing'
  | 'completed'
  | 'cancelled'

export type RepairPriority = 'urgent' | 'normal' | 'low'

export interface Repair {
  id: string
  room_id: string
  reporter_id: string
  repairer_id: string | null
  type: RepairType
  description: string
  status: RepairStatus
  priority: RepairPriority
  scheduled_at: string | null
  completed_at: string | null
  cost: number | null
  rating: number | null
  remark: string | null
  created_at: string
  updated_at: string
}

export interface RepairWithRoom extends Repair {
  room_number: string
  building_name: string
  reporter_name: string
  repairer_name: string | null
}

export interface Pagination {
  page: number
  page_size: number
  total_items: number
  total_pages: number
}

export interface PaginatedResult<T> {
  data: T[]
  pagination: Pagination
}

export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: {
    code: number
    message: string
    details?: string
  }
}
