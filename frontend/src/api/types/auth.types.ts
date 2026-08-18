export type Role =
  | "super_admin"
  | "center_owner"
  | "center_scanner"
  | "center_receptionist"

export interface Center {
  id: number
  name: string
  slug: string
}

export interface User {
  id: number
  username: string
  email: string
  role: Role
  center?: Center
  // True once the user's center has at least one grade. Only meaningful
  // for center-scoped roles — drives the center_owner onboarding gate.
  has_grades: boolean
}

export interface LoginPayload {
  email: string
  password: string
}

export interface UpdateProfilePayload {
  username: string
  email: string
}

export interface UpdatePasswordPayload {
  currentPassword: string
  newPassword: string
}

export interface ForgotPasswordPayload {
  email: string
}

export interface ResetPasswordPayload {
  token: string
  password: string
}

export interface AuthEnvelope<T> {
  message: string
  data: T
}

export interface AuthData {
  user: User
  token?: string
}
