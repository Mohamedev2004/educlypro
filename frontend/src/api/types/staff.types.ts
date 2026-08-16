import type { ApiEnvelope } from "./notification.types"

/**
 * Staff (center_scanner / center_receptionist) API types.
 *
 * Responsibility: API-specific request/response shapes for staff management.
 * Layer: Types
 */

export type StaffRole = "center_scanner" | "center_receptionist"

export type Staff = {
  id: number
  username: string
  email: string
  role: StaffRole
  center_id: number
  created_at: string
}

export type StaffPagination = {
  page: number
  per_page: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

export type StaffListData = {
  items: Staff[]
  pagination: StaffPagination
}

export type StaffListParams = {
  page: number
  perPage: number
  role?: StaffRole
  q?: string
  sort?: "username" | "email" | "role" | "created_at"
  direction?: "asc" | "desc"
}

export type CreateStaffPayload = {
  username: string
  email: string
  password: string
  role: StaffRole
  center_id: number
}

export type UpdateStaffPayload = {
  username: string
  email: string
  role: StaffRole
  password?: string
  center_id: number
}

export type StaffListEnvelope = ApiEnvelope<StaffListData>
export type StaffEnvelope = ApiEnvelope<Staff>
