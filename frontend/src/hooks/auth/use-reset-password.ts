import { useState, useEffect } from "react"
import { useNavigate, useSearchParams } from "react-router-dom"
import { toast } from "sonner"
import { AuthService } from "@/api/services/auth-service"
import { normalizeApiError, getFieldError, getApiMessage, getAuthFieldMessage } from "@/utils/error-utils"
import type { ResetPasswordPayload } from "@/api/types/auth.types"

export function useResetPassword(t: (key: string) => string) {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const token = searchParams.get("token")

  const [isSubmitting, setIsSubmitting] = useState(false)
  const [errors, setErrors] = useState<{
    password?: string
    confirmPassword?: string
    general?: string
  }>({})

  useEffect(() => {
    if (!token) {
      toast.error(t("api.defaultError"))
      navigate("/login")
    }
  }, [token, navigate, t])

  const handleResetPassword = async (payload: Omit<ResetPasswordPayload, "token"> & { confirmPassword: string }) => {
    const { password, confirmPassword } = payload
    setIsSubmitting(true)
    setErrors({})

    if (password !== confirmPassword) {
      setErrors({ confirmPassword: t("settings.errors.passwordsMismatch") })
      setIsSubmitting(false)
      return
    }

    if (!token) {
      setIsSubmitting(false)
      return
    }

    try {
      await AuthService.resetPassword({ token, password })
      const toastId = toast.success(t("auth.passwordResetSuccess"), {
        description: t("auth.passwordResetSuccessDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
      navigate("/login", { replace: true })
    } catch (err) {
      const apiError = normalizeApiError(err)
      const passwordError = getFieldError(apiError, "password")

      setErrors({
        password: passwordError
          ? getAuthFieldMessage(t, "auth", "password", passwordError)
          : undefined,
        general: !passwordError ? getApiMessage(t, apiError) : undefined,
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  return {
    isSubmitting,
    errors,
    handleResetPassword,
  }
}
