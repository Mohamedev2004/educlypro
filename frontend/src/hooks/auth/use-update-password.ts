import { useState } from "react"
import { toast } from "sonner"
import { useAuth } from "@/context/auth/auth-context"
import { normalizeApiError, getFieldError, getApiMessage, getAuthFieldMessage } from "@/utils/error-utils"
import type { UpdatePasswordPayload } from "@/api/types/auth.types"

export function useUpdatePassword(t: (key: string) => string) {
  const { updatePassword } = useAuth()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [errors, setErrors] = useState<{
    currentPassword?: string
    newPassword?: string
    confirmPassword?: string
    general?: string
  }>({})

  const handleUpdatePassword = async (payload: UpdatePasswordPayload & { confirmPassword: string }) => {
    const { currentPassword, newPassword, confirmPassword } = payload
    setErrors({})

    if (newPassword !== confirmPassword) {
      setErrors({ confirmPassword: t("settings.errors.passwordsMismatch") })
      return false
    }

    setIsSubmitting(true)

    try {
      await updatePassword({ currentPassword, newPassword })
      const toastId = toast.success(t("settings.passwordUpdated"), {
        description: t("settings.passwordUpdatedDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
      return true
    } catch (err) {
      const apiError = normalizeApiError(err)
      const currentPasswordError = getFieldError(apiError, "currentPassword")
      const newPasswordError = getFieldError(apiError, "newPassword")

      setErrors({
        currentPassword: currentPasswordError
          ? getAuthFieldMessage(t, "settings", "currentPassword", currentPasswordError)
          : undefined,
        newPassword: newPasswordError
          ? getAuthFieldMessage(t, "settings", "newPassword", newPasswordError)
          : undefined,
        general:
          !currentPasswordError && !newPasswordError
            ? getApiMessage(t, apiError)
            : undefined,
      })
      return false
    } finally {
      setIsSubmitting(false)
    }
  }

  return {
    isSubmitting,
    errors,
    handleUpdatePassword,
  }
}
