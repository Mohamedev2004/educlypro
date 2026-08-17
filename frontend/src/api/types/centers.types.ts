import type { ApiEnvelope } from "./notification.types"
import type { Staff, StaffRole } from "./staff.types"

/**
 * Centers (super admin overview) API types.
 *
 * Responsibility: API-specific request/response shapes for the centers
 * list, detail view, and provisioning actions (create center, add owner,
 * add staff).
 * Layer: Types
 */

export type Center = {
  id: number
  name: string
  slug: string
  owner_username: string | null
  owner_email: string | null
  staff_count: number
  created_at: string
}

export type CentersPagination = {
  page: number
  per_page: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

export type CentersListData = {
  items: Center[]
  pagination: CentersPagination
}

export type CentersListParams = {
  page: number
  perPage: number
  q?: string
  sort?: "name" | "slug" | "created_at"
  direction?: "asc" | "desc"
}

export type CenterOwner = {
  id: number
  username: string
  email: string
  created_at: string
}

export type CenterStaffMember = {
  id: number
  username: string
  email: string
  role: StaffRole
  created_at: string
}

export type CenterDetail = {
  id: number
  name: string
  slug: string
  owner: CenterOwner | null
  staff: CenterStaffMember[]
  created_at: string
}

// The slug is never client-supplied — the backend always derives and
// de-duplicates it from name.
export type CreateCenterPayload = {
  name: string
}

export type AddOwnerPayload = {
  username: string
  email: string
  password: string
}

export type AddStaffPayload = {
  username: string
  email: string
  password: string
  role: StaffRole
}

export type CentersListEnvelope = ApiEnvelope<CentersListData>
export type CenterEnvelope = ApiEnvelope<Center>
export type CenterDetailEnvelope = ApiEnvelope<CenterDetail>
export type CenterOwnerEnvelope = ApiEnvelope<CenterOwner>
export type AddStaffEnvelope = ApiEnvelope<Staff>
