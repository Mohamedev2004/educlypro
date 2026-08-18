import type { ApiEnvelope } from "./notification.types"

/**
 * Academic setup (center_owner onboarding) API types.
 *
 * Responsibility: API-specific request/response shapes for the
 * grade -> major -> subject tree a center owner builds during onboarding.
 * Layer: Types
 */

export type Subject = {
  id: number
  name: string
}

export type Class = {
  id: number
  name: string
}

export type Major = {
  id: number
  name: string
  subjects: Subject[]
  class: Class | null
}

export type Grade = {
  id: number
  name: string
  majors: Major[]
}

export type AcademicTree = {
  grades: Grade[]
}

export type AddGradePayload = {
  name: string
}

export type AddMajorPayload = {
  name: string
}

export type AddSubjectPayload = {
  name: string
}

export type AcademicTreeEnvelope = ApiEnvelope<AcademicTree>
export type GradeEnvelope = ApiEnvelope<Grade>
export type MajorEnvelope = ApiEnvelope<Major>
export type SubjectEnvelope = ApiEnvelope<Subject>
