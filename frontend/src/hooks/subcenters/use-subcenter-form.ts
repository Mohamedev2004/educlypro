import { useEffect, useState } from "react"
import {
  normalizeApiError,
  getFieldError,
  getApiMessage,
  getSubCentersFieldMessage,
} from "@/utils/error-utils"
import { useCreateSubCenter } from "./use-create-subcenter"
import { useUpdateSubCenter } from "./use-update-subcenter"
import type { SubCenter } from "@/api/types/subcenters.types"

/**
 * Hook for the sub-center create/edit form.
 *
 * Responsibility: own the sub-center form field state, map server
 * validation errors to fields, and submit to the create or update mutation
 * depending on whether an existing sub-center was passed in. Fields reset
 * every time the dialog opens (not just when `subCenter` changes identity)
 * — otherwise reopening in "add" mode back-to-back would keep showing
 * whatever was typed (or just created) last time, since `subCenter` stays
 * `null` across two adds and never re-triggers a reset on its own.
 * Layer: Hooks
 */
export function useSubCenterForm(
  t: (key: string) => string,
  open: boolean,
  subCenter: SubCenter | null,
  onSuccess: () => void
) {
  const isEdit = !!subCenter

  const [name, setName] = useState(subCenter?.name ?? "")
  const [errors, setErrors] = useState<{
    name?: string
    general?: string
  }>({})

  useEffect(() => {
    if (!open) return
    setName(subCenter?.name ?? "")
    setErrors({})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  const createSubCenter = useCreateSubCenter(t)
  const updateSubCenter = useUpdateSubCenter(t)
  const isSubmitting = createSubCenter.isPending || updateSubCenter.isPending

  const handleSubmit = async () => {
    setErrors({})

    try {
      if (isEdit) {
        await updateSubCenter.mutateAsync({
          id: subCenter.id,
          payload: { name },
        })
      } else {
        await createSubCenter.mutateAsync({ name })
      }
      onSuccess()
      return true
    } catch (err) {
      const apiError = normalizeApiError(err)
      const nameError = getFieldError(apiError, "name")

      setErrors({
        name: nameError
          ? getSubCentersFieldMessage(t, "name", nameError)
          : undefined,
        general: !nameError ? getApiMessage(t, apiError) : undefined,
      })
      return false
    }
  }

  return {
    isEdit,
    name,
    setName,
    errors,
    isSubmitting,
    handleSubmit,
  }
}
