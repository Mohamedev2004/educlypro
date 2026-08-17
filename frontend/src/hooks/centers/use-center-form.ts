import { useState } from "react"
import {
  normalizeApiError,
  getFieldError,
  getApiMessage,
  getCentersFieldMessage,
} from "@/utils/error-utils"
import { useCreateCenter } from "./use-create-center"

/**
 * Hook for the create-center form.
 *
 * Responsibility: own the center form field state, map server validation
 * errors to fields, and submit to the create mutation. The slug is never a
 * user input — the backend always derives and de-duplicates it from name.
 * Layer: Hooks
 */
export function useCenterForm(t: (key: string) => string, onSuccess: () => void) {
  const [name, setName] = useState("")
  const [errors, setErrors] = useState<{
    name?: string
    general?: string
  }>({})

  const createCenter = useCreateCenter(t)
  const isSubmitting = createCenter.isPending

  const handleSubmit = async () => {
    setErrors({})

    try {
      await createCenter.mutateAsync({ name })
      setName("")
      onSuccess()
      return true
    } catch (err) {
      const apiError = normalizeApiError(err)
      const nameError = getFieldError(apiError, "name")

      setErrors({
        name: nameError ? getCentersFieldMessage(t, "name", nameError) : undefined,
        general: !nameError ? getApiMessage(t, apiError) : undefined,
      })
      return false
    }
  }

  return {
    name,
    setName,
    errors,
    isSubmitting,
    handleSubmit,
  }
}
