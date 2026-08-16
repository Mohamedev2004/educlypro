import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { toast } from "sonner"
import { useAuth } from "@/context/auth/auth-context"
import { normalizeApiError, getFieldError, getApiMessage, getAuthFieldMessage } from "@/utils/error-utils"
import { getDashboardPath } from "@/utils/navigation-utils"
import type { LoginPayload } from "@/api/types/auth.types"

export function useLogin(t: (key: string) => string) {
  const { login } = useAuth()
  const navigate = useNavigate()
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [errors, setErrors] = useState<{
    email?: string
    password?: string
    general?: string
  }>({})

  const handleLogin = async (payload: LoginPayload) => {
    setIsSubmitting(true)
    setErrors({})

    try {
      const user = await login(payload)
      const toastId = toast.success(t("auth.loggedIn"), {
        description: t("auth.loggedInDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
      navigate(getDashboardPath(user.role), { replace: true })
    } catch (err) {
      const apiError = normalizeApiError(err)
      const emailError = getFieldError(apiError, "email")
      const passwordError = getFieldError(apiError, "password")

      setErrors({
        email: emailError
          ? getAuthFieldMessage(t, "auth", "email", emailError)
          : undefined,
        password: passwordError
          ? getAuthFieldMessage(t, "auth", "password", passwordError)
          : undefined,
        general:
          apiError.code === "invalid_credentials"
            ? getApiMessage(t, apiError)
            : !emailError && !passwordError
              ? getApiMessage(t, apiError)
              : undefined,
      })
    } finally {
      setIsSubmitting(false)
    }
  }

  return {
    isSubmitting,
    errors,
    handleLogin,
  }
}
