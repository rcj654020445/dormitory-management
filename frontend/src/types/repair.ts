// Layer 0: Core types — no internal imports allowed
// Repair system types matching the backend Go types

export type RepairType =
  | 'facility'    // 设施维修（门/窗/床/柜）
  | 'plumbing'    // 水管/漏水
  | 'electrical'  // 电路/灯具/插座
  | 'network'     // 网络故障
  | 'cleaning'    // 保洁
  | 'other'       // 其他

export type RepairStatus =
  | 'pending'     // 待派单
  | 'assigned'    // 已派单
  | 'repairing'   // 维修中
  | 'completed'   // 已完成
  | 'cancelled'   // 已取消

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