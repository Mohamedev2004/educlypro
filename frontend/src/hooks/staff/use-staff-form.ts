import { useEffect, useState } from "react"
import {
  normalizeApiError,
  getFieldError,
  getApiMessage,
  getStaffFieldMessage,
} from "@/utils/error-utils"
import { useCreateStaff } from "./use-create-staff"
import { useUpdateStaff } from "./use-update-staff"
import type { Staff, StaffRole } from "@/api/types/staff.types"

/**
 * Hook for the staff create/edit form.
 *
 * Responsibility: own the staff form field state, map server validation
 * errors to fields, and submit to the create or update mutation depending
 * on whether an existing staff record was passed in. `centerId` is a fixed
 * value supplied by the caller (the signed-in owner's own center) — it is
 * never exposed as an editable field, only included in the submitted payload.
 * Fields reset every time the dialog opens (not just when `staff` changes
 * identity) — otherwise reopening in "add" mode back-to-back would keep
 * showing whatever was typed (or just created) last time, since `staff`
 * stays `null` across two adds and never re-triggers a reset on its own.
 * Layer: Hooks
 */
export function useStaffForm(
  t: (key: string) => string,
  open: boolean,
  centerId: number,
  staff: Staff | null,
  onSuccess: () => void
) {
  const isEdit = !!staff

  const [username, setUsername] = useState(staff?.username ?? "")
  const [email, setEmail] = useState(staff?.email ?? "")
  const [role, setRole] = useState<StaffRole>(staff?.role ?? "center_scanner")
  const [password, setPassword] = useState("")
  const [subCenterId, setSubCenterId] = useState(staff?.sub_center_id ?? 0)
  const [errors, setErrors] = useState<{
    username?: string
    email?: string
    role?: string
    password?: string
    subCenterId?: string
    general?: string
  }>({})

  useEffect(() => {
    if (!open) return
    setUsername(staff?.username ?? "")
    setEmail(staff?.email ?? "")
    setRole(staff?.role ?? "center_scanner")
    setPassword("")
    setSubCenterId(staff?.sub_center_id ?? 0)
    setErrors({})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  const createStaff = useCreateStaff(t)
  const updateStaff = useUpdateStaff(t)
  const isSubmitting = createStaff.isPending || updateStaff.isPending

  const handleSubmit = async () => {
    setErrors({})

    if (!subCenterId) {
      setErrors({ subCenterId: t("staff.errors.subCenterRequired") })
      return false
    }

    try {
      if (isEdit) {
        await updateStaff.mutateAsync({
          id: staff.id,
          payload: {
            username,
            email,
            role,
            password: password || undefined,
            center_id: centerId,
            sub_center_id: subCenterId,
          },
        })
      } else {
        await createStaff.mutateAsync({
          username,
          email,
          password,
          role,
          center_id: centerId,
          sub_center_id: subCenterId,
        })
      }
      onSuccess()
      return true
    } catch (err) {
      const apiError = normalizeApiError(err)
      const usernameError = getFieldError(apiError, "username")
      const emailError = getFieldError(apiError, "email")
      const roleError = getFieldError(apiError, "role")
      const passwordError = getFieldError(apiError, "password")

      setErrors({
        username: usernameError
          ? getStaffFieldMessage(t, "username", usernameError)
          : undefined,
        email: emailError
          ? getStaffFieldMessage(t, "email", emailError)
          : undefined,
        role: roleError
          ? getStaffFieldMessage(t, "role", roleError)
          : undefined,
        password: passwordError
          ? getStaffFieldMessage(t, "password", passwordError)
          : undefined,
        general:
          !usernameError && !emailError && !roleError && !passwordError
            ? getApiMessage(t, apiError)
            : undefined,
      })
      return false
    }
  }

  return {
    isEdit,
    username,
    setUsername,
    email,
    setEmail,
    role,
    setRole,
    password,
    setPassword,
    subCenterId,
    setSubCenterId,
    errors,
    isSubmitting,
    handleSubmit,
  }
}
