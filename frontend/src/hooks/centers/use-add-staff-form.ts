import { useState } from "react"
import {
  normalizeApiError,
  getFieldError,
  getApiMessage,
  getCentersFieldMessage,
} from "@/utils/error-utils"
import { useAddCenterStaff } from "./use-add-staff"
import type { StaffRole } from "@/api/types/staff.types"

/**
 * Hook for the center detail page's add-staff form.
 *
 * Responsibility: own the add-staff form field state, map server validation
 * errors to fields, and submit to the add-staff mutation.
 * Layer: Hooks
 */
export function useAddStaffForm(
  t: (key: string) => string,
  centerSlug: string,
  onSuccess: () => void
) {
  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [role, setRole] = useState<StaffRole>("center_scanner")
  const [password, setPassword] = useState("")
  const [errors, setErrors] = useState<{
    username?: string
    email?: string
    role?: string
    password?: string
    general?: string
  }>({})

  const addStaff = useAddCenterStaff(t, centerSlug)
  const isSubmitting = addStaff.isPending

  const handleSubmit = async () => {
    setErrors({})

    try {
      await addStaff.mutateAsync({ username, email, password, role })
      setUsername("")
      setEmail("")
      setRole("center_scanner")
      setPassword("")
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
          ? getCentersFieldMessage(t, "username", usernameError)
          : undefined,
        email: emailError ? getCentersFieldMessage(t, "email", emailError) : undefined,
        role: roleError ? getCentersFieldMessage(t, "role", roleError) : undefined,
        password: passwordError
          ? getCentersFieldMessage(t, "password", passwordError)
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
