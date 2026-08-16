import { useState } from "react"
import { toast } from "sonner"
import { AuthService } from "@/api/services/auth-service"
import { normalizeApiError, getFieldError, getApiMessage, getAuthFieldMessage } from "@/utils/error-utils"
import type { ForgotPasswordPayload } from "@/api/types/auth.types"

export function useForgotPassword(t: (key: string) => string) {
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [errors, setErrors] = useState<{ email?: string; general?: string }>({})

  const handleForgotPassword = async (payload: ForgotPasswordPayload) => {
    setIsSubmitting(true)
    setErrors({})

    try {
      await AuthService.forgotPassword(payload)
      const toastId = toast.success(t("auth.resetLinkSent"), {
        description: t("auth.resetLinkSentDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
    } catch (err) {
      const apiError = normalizeApiError(err)
      const emailError = getFieldError(apiError, "email")

      setErrors({
        email: emailError
          ? getAuthFieldMessage(t, "auth", "email", emailError)
          : undefined,
        general: !emailError ? getApiMessage(t, apiError) : undefined,
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  return {
    isSubmitting,
    errors,
    handleForgotPassword,
  }
}
