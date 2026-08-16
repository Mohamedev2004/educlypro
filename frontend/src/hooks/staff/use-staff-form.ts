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
 * Layer: Hooks
 */
export function useStaffForm(
  t: (key: string) => string,
  centerId: number,
  staff: Staff | null,
  onSuccess: () => void
) {
  const isEdit = !!staff

  const [username, setUsername] = useState(staff?.username ?? "")
  const [email, setEmail] = useState(staff?.email ?? "")
  const [role, setRole] = useState<StaffRole>(staff?.role ?? "center_scanner")
  const [password, setPassword] = useState("")
  const [errors, setErrors] = useState<{
    username?: string
    email?: string
    role?: string
    password?: string
    general?: string
  }>({})

  useEffect(() => {
    setUsername(staff?.username ?? "")
    setEmail(staff?.email ?? "")
    setRole(staff?.role ?? "center_scanner")
    setPassword("")
    setErrors({})
  }, [staff])

  const createStaff = useCreateStaff(t)
  const updateStaff = useUpdateStaff(t)
  const isSubmitting = createStaff.isPending || updateStaff.isPending

  const handleSubmit = async () => {
    setErrors({})

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
          },
        })
      } else {
        await createStaff.mutateAsync({
          username,
          email,
          password,
          role,
          center_id: centerId,
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
        email: emailError ? getStaffFieldMessage(t, "email", emailError) : undefined,
        role: roleError ? getStaffFieldMessage(t, "role", roleError) : undefined,
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
    errors,
    isSubmitting,
    handleSubmit,
  }
}
