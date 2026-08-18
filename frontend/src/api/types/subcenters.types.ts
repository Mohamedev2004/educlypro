import type { ApiEnvelope } from "./notification.types"

/**
 * Sub-centers API types.
 *
 * Responsibility: API-specific request/response shapes for a center
 * owner's sub-centers (operational locations within their own center).
 * Layer: Types
 */

export type SubCenter = {
  id: number
  name: string
  staff_count: number
  created_at: string
}

// Unpaginated — a center realistically has a handful of sub-centers, not
// pages of them.
export type SubCentersListData = {
  items: SubCenter[]
}

export type CreateSubCenterPayload = {
  name: string
}

export type UpdateSubCenterPayload = {
  name: string
}

export type SubCentersListEnvelope = ApiEnvelope<SubCentersListData>
export type SubCenterEnvelope = ApiEnvelope<SubCenter>
