import { useState } from "react"
import {
  normalizeApiError,
  getFieldError,
  getApiMessage,
  getCentersFieldMessage,
} from "@/utils/error-utils"
import { useAddOwner } from "./use-add-owner"

/**
 * Hook for the add-owner form.
 *
 * Responsibility: own the add-owner form field state, map server validation
 * errors to fields, and submit to the add-owner mutation.
 * Layer: Hooks
 */
export function useAddOwnerForm(
  t: (key: string) => string,
  centerSlug: string,
  onSuccess: () => void
) {
  const [username, setUsername] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [errors, setErrors] = useState<{
    username?: string
    email?: string
    password?: string
    general?: string
  }>({})

  const addOwner = useAddOwner(t, centerSlug)
  const isSubmitting = addOwner.isPending

  const handleSubmit = async () => {
    setErrors({})

    try {
      await addOwner.mutateAsync({ username, email, password })
      setUsername("")
      setEmail("")
      setPassword("")
      onSuccess()
      return true
    } catch (err) {
      const apiError = normalizeApiError(err)
      const usernameError = getFieldError(apiError, "username")
      const emailError = getFieldError(apiError, "email")
      const passwordError = getFieldError(apiError, "password")

      setErrors({
        username: usernameError
          ? getCentersFieldMessage(t, "username", usernameError)
          : undefined,
        email: emailError ? getCentersFieldMessage(t, "email", emailError) : undefined,
        password: passwordError
          ? getCentersFieldMessage(t, "password", passwordError)
          : undefined,
        general:
          !usernameError && !emailError && !passwordError
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
    password,
    setPassword,
    errors,
    isSubmitting,
    handleSubmit,
  }
}
