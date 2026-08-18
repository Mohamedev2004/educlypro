import type { ApiEnvelope } from "./notification.types"

/**
 * Teachers API types.
 *
 * Responsibility: API-specific request/response shapes for a center
 * owner's teachers — plain records (not auth.User accounts), many-to-many
 * with both subjects and classes.
 * Layer: Types
 */

export type TeacherSubject = {
  id: number
  name: string
}

export type TeacherClass = {
  id: number
  name: string
}

export type Teacher = {
  id: number
  slug: string
  full_name: string
  email: string
  phone: string
  subjects: TeacherSubject[]
  classes: TeacherClass[]
  created_at: string
}

export type TeachersPagination = {
  page: number
  per_page: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

export type TeachersListData = {
  items: Teacher[]
  pagination: TeachersPagination
}

export type TeachersListParams = {
  page: number
  perPage: number
  q?: string
  sort?: "full_name" | "email" | "created_at"
  direction?: "asc" | "desc"
}

export type CreateTeacherPayload = {
  full_name: string
  email: string
  phone: string
}

export type UpdateTeacherPayload = CreateTeacherPayload

export type UpdateTeacherSubjectsPayload = {
  subject_ids: number[]
}

export type UpdateTeacherClassesPayload = {
  class_ids: number[]
}

export type TeachersListEnvelope = ApiEnvelope<TeachersListData>
export type TeacherEnvelope = ApiEnvelope<Teacher>
