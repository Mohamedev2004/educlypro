import { useState, useEffect } from "react"
import { toast } from "sonner"
import { useAuth } from "@/context/auth/auth-context"
import { normalizeApiError, getFieldError, getApiMessage, getAuthFieldMessage } from "@/utils/error-utils"
import type { UpdateProfilePayload } from "@/api/types/auth.types"

export function useUpdateProfile(t: (key: string) => string) {
  const { user, updateProfile } = useAuth()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [username, setUsername] = useState(user?.username || "")
  const [email, setEmail] = useState(user?.email || "")
  const [errors, setErrors] = useState<{
    username?: string
    email?: string
    general?: string
  }>({})

  useEffect(() => {
    if (user) {
      setUsername(user.username)
      setEmail(user.email)
    }
  }, [user])

  const handleUpdateProfile = async (payload: UpdateProfilePayload) => {
    setIsSubmitting(true)
    setErrors({})

    try {
      await updateProfile(payload)
      const toastId = toast.success(t("settings.profileUpdated"), {
        description: t("settings.profileUpdatedDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
      return true
    } catch (err) {
      const apiError = normalizeApiError(err)
      const usernameError = getFieldError(apiError, "username")
      const emailError = getFieldError(apiError, "email")

      setErrors({
        username: usernameError
          ? getAuthFieldMessage(t, "settings", "username", usernameError)
          : undefined,
        email: emailError
          ? getAuthFieldMessage(t, "settings", "email", emailError)
          : undefined,
        general:
          !usernameError && !emailError
            ? getApiMessage(t, apiError)
            : undefined,
      })
      return false
    } finally {
      setIsSubmitting(false)
    }
  }

  return {
    username,
    setUsername,
    email,
    setEmail,
    isSubmitting,
    errors,
    handleUpdateProfile,
  }
}
