/**
 * Pure helper functions for handling and translating API errors.
 *
 * Responsibility: Normalize API errors and provide translated messages.
 * Layer: Utils
 */

export type ApiError = {
  message: string
  status?: number
  code?: string
  errors?: Record<string, string>
}

type Translate = (key: string, fallback?: string) => string

/**
 * Normalizes an unknown error into a standard ApiError object.
 */
export function normalizeApiError(error: unknown): ApiError {
  if (isApiError(error)) {
    return error
  }

  if (error instanceof Error) {
    return { message: error.message }
  }

  return { message: "Something went wrong. Please try again." }
}

/**
 * Checks if an unknown error matches the ApiError structure.
 */
function isApiError(error: unknown): error is ApiError {
  return Boolean(
    error &&
    typeof error === "object" &&
    "message" in error &&
    typeof (error as { message?: unknown }).message === "string"
  )
}

/**
 * Gets a specific field error from an ApiError object.
 */
export function getFieldError(error: ApiError, field: string) {
  return error.errors?.[field]
}

/**
 * Returns a translated error message for auth or settings fields.
 */
export function getAuthFieldMessage(
  t: Translate,
  section: "auth" | "settings",
  field: string,
  code?: string
) {
  if (!code) return undefined

  const keyMap: Record<string, string> = {
    [`${section}.email.validation.required`]: `${section}.errors.emailRequired`,
    [`${section}.email.validation.email`]: `${section}.errors.emailInvalid`,
    [`${section}.email.validation.taken`]: `${section}.errors.emailTaken`,
    [`${section}.password.validation.required`]: `${section}.errors.passwordRequired`,
    [`${section}.password.validation.min`]: `${section}.errors.passwordMin`,
    [`${section}.username.validation.required`]: `${section}.errors.usernameRequired`,
    [`${section}.username.validation.min`]: `${section}.errors.usernameMin`,
    "settings.currentPassword.validation.required":
      "settings.errors.currentPasswordRequired",
    "settings.currentPassword.auth.current_password_invalid":
      "settings.errors.currentPasswordInvalid",
    "settings.newPassword.validation.required":
      "settings.errors.newPasswordRequired",
    "settings.newPassword.validation.min": "settings.errors.newPasswordMin",
  }

  const key = keyMap[`${section}.${field}.${code}`]
  return key ? t(key) : undefined
}

/**
 * Returns a translated error message for the whole API error.
 */
export function getApiMessage(t: Translate, error: ApiError) {
  const codeMap: Record<string, string> = {
    invalid_credentials: "auth.errors.invalidCredentials",
    invalid_or_expired_token: "auth.errors.invalidResetToken",
    refresh_token_required: "api.sessionExpired",
    invalid_refresh_token: "api.sessionExpired",
    unauthorized: "api.sessionExpired",
    forbidden: "api.forbidden",
    registration_failed: "api.serverError",
    login_failed: "api.serverError",
    logout_failed: "api.serverError",
    profile_update_failed: "api.serverError",
    password_update_failed: "api.serverError",
    forgot_password_failed: "api.serverError",
    reset_password_failed: "api.serverError",
    invalid_notification_page: "notifications.errors.invalidPage",
    invalid_notification_per_page: "notifications.errors.invalidPerPage",
    invalid_notification_filter: "notifications.errors.invalidFilter",
    invalid_notification_id: "notifications.errors.invalidNotification",
    notifications_list_failed: "notifications.errors.loadFailed",
    notifications_unread_count_failed: "notifications.errors.loadFailed",
    notification_mark_read_failed: "notifications.errors.markReadFailed",
    notifications_mark_all_read_failed:
      "notifications.errors.markAllReadFailed",
    notification_delete_failed: "notifications.errors.deleteFailed",
    center_mismatch: "staff.errors.centerMismatch",
    staff_not_found: "staff.errors.notFound",
    staff_list_failed: "staff.errors.loadFailed",
    staff_create_failed: "staff.errors.createFailed",
    staff_update_failed: "staff.errors.updateFailed",
    staff_delete_failed: "staff.errors.deleteFailed",
    center_not_found: "centers.errors.notFound",
    owner_already_assigned: "centers.errors.ownerAlreadyAssigned",
    centers_list_failed: "centers.errors.loadFailed",
    centers_create_failed: "centers.errors.createFailed",
    centers_get_failed: "centers.errors.loadFailed",
    centers_add_owner_failed: "centers.errors.addOwnerFailed",
    centers_add_staff_failed: "centers.errors.addStaffFailed",
    centers_list_subcenters_failed: "subcenters.errors.loadFailed",
    subcenter_not_found: "subcenters.errors.notFound",
    subcenter_exists: "subcenters.errors.exists",
    subcenter_has_staff: "subcenters.errors.hasStaff",
    subcenter_mismatch: "subcenters.errors.mismatch",
    subcenters_list_failed: "subcenters.errors.loadFailed",
    subcenters_create_failed: "subcenters.errors.createFailed",
    subcenters_update_failed: "subcenters.errors.updateFailed",
    subcenters_delete_failed: "subcenters.errors.deleteFailed",
    grade_not_found: "onboarding.errors.gradeNotFound",
    grade_exists: "onboarding.errors.gradeExists",
    major_not_found: "onboarding.errors.majorNotFound",
    major_exists: "onboarding.errors.majorExists",
    subject_not_found: "onboarding.errors.subjectNotFound",
    subject_exists: "onboarding.errors.subjectExists",
    academic_tree_failed: "onboarding.errors.treeFailed",
    academic_add_grade_failed: "onboarding.errors.addGradeFailed",
    academic_remove_grade_failed: "onboarding.errors.removeGradeFailed",
    academic_add_major_failed: "onboarding.errors.addMajorFailed",
    academic_remove_major_failed: "onboarding.errors.removeMajorFailed",
    academic_add_subject_failed: "onboarding.errors.addSubjectFailed",
    academic_remove_subject_failed: "onboarding.errors.removeSubjectFailed",
    onboarding_required: "onboarding.errors.setupRequired",
    onboarding_check_failed: "api.serverError",
    teacher_not_found: "teachers.errors.notFound",
    invalid_subject_selection: "teachers.errors.invalidSubjectSelection",
    invalid_class_selection: "teachers.errors.invalidClassSelection",
    teacher_list_failed: "teachers.errors.loadFailed",
    teacher_create_failed: "teachers.errors.createFailed",
    teacher_update_failed: "teachers.errors.updateFailed",
    teacher_delete_failed: "teachers.errors.deleteFailed",
    network_error: "api.networkError",
  }

  const key = error.code ? codeMap[error.code] : undefined
  if (key) return t(key)
  if (error.message === "Network Error") return t("api.networkError")
  return error.message || t("api.defaultError")
}

/**
 * Returns a translated error message for notification-specific fields.
 */
export function getNotificationFieldMessage(
  t: Translate,
  field: string,
  code?: string
) {
  if (!code) return undefined

  const keyMap: Record<string, string> = {
    "page.validation.invalid": "notifications.errors.invalidPage",
    "per_page.validation.invalid": "notifications.errors.invalidPerPage",
    "filter.validation.invalid": "notifications.errors.invalidFilter",
  }

  const key = keyMap[`${field}.${code}`]
  return key ? t(key) : undefined
}

/**
 * Returns a translated error message for centers-specific fields (create
 * center, add owner, add staff forms).
 */
export function getCentersFieldMessage(
  t: Translate,
  field: string,
  code?: string
) {
  if (!code) return undefined

  const keyMap: Record<string, string> = {
    "name.validation.required": "centers.errors.nameRequired",
    "name.validation.min": "centers.errors.nameMin",
    "username.validation.required": "centers.errors.usernameRequired",
    "username.validation.min": "centers.errors.usernameMin",
    "email.validation.required": "centers.errors.emailRequired",
    "email.validation.email": "centers.errors.emailInvalid",
    "email.validation.taken": "centers.errors.emailTaken",
    "password.validation.required": "centers.errors.passwordRequired",
    "password.validation.min": "centers.errors.passwordMin",
    "role.validation.invalid": "centers.errors.roleInvalid",
  }

  const key = keyMap[`${field}.${code}`]
  return key ? t(key) : undefined
}

/**
 * Returns a translated error message for sub-center-specific fields (create
 * / rename forms).
 */
export function getSubCentersFieldMessage(
  t: Translate,
  field: string,
  code?: string
) {
  if (!code) return undefined

  const keyMap: Record<string, string> = {
    "name.validation.required": "subcenters.errors.nameRequired",
    "name.validation.min": "subcenters.errors.nameMin",
    "name.validation.taken": "subcenters.errors.exists",
  }

  const key = keyMap[`${field}.${code}`]
  return key ? t(key) : undefined
}

/**
 * Returns a translated error message for teacher-specific fields.
 */
export function getTeachersFieldMessage(
  t: Translate,
  field: string,
  code?: string
) {
  if (!code) return undefined

  const keyMap: Record<string, string> = {
    "fullName.validation.required": "teachers.errors.fullNameRequired",
    "fullName.validation.min": "teachers.errors.fullNameMin",
    "email.validation.required": "teachers.errors.emailRequired",
    "email.validation.email": "teachers.errors.emailInvalid",
    "email.validation.taken": "teachers.errors.emailTaken",
    "phone.validation.required": "teachers.errors.phoneRequired",
    "phone.validation.min": "teachers.errors.phoneMin",
  }

  const key = keyMap[`${field}.${code}`]
  return key ? t(key) : undefined
}

/**
 * Returns a translated error message for staff-specific fields.
 */
export function getStaffFieldMessage(
  t: Translate,
  field: string,
  code?: string
) {
  if (!code) return undefined

  const keyMap: Record<string, string> = {
    "username.validation.required": "staff.errors.usernameRequired",
    "username.validation.min": "staff.errors.usernameMin",
    "email.validation.required": "staff.errors.emailRequired",
    "email.validation.email": "staff.errors.emailInvalid",
    "email.validation.taken": "staff.errors.emailTaken",
    "password.validation.required": "staff.errors.passwordRequired",
    "password.validation.min": "staff.errors.passwordMin",
    "role.validation.invalid": "staff.errors.roleInvalid",
  }

  const key = keyMap[`${field}.${code}`]
  return key ? t(key) : undefined
}
