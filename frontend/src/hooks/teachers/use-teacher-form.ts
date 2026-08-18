import { useEffect, useState } from "react"
import {
  normalizeApiError,
  getFieldError,
  getApiMessage,
  getTeachersFieldMessage,
} from "@/utils/error-utils"
import { useCreateTeacher } from "./use-create-teacher"
import { useUpdateTeacher } from "./use-update-teacher"
import type { Teacher } from "@/api/types/teachers.types"

/**
 * Hook for the teacher create/edit form.
 *
 * Responsibility: own the teacher form field state (identity fields only —
 * subject/class assignment lives on the teacher detail page instead), map
 * server validation errors to fields, and submit to the create or update
 * mutation depending on whether an existing teacher was passed in. Fields
 * reset every time the dialog opens (not just when `teacher` changes
 * identity) — otherwise reopening in "add" mode back-to-back would keep
 * showing whatever was typed (or just created) last time, since `teacher`
 * stays `null` across two adds and never re-triggers a reset on its own.
 * Layer: Hooks
 */
export function useTeacherForm(
  t: (key: string) => string,
  open: boolean,
  teacher: Teacher | null,
  onSuccess: () => void
) {
  const isEdit = !!teacher

  const [fullName, setFullName] = useState(teacher?.full_name ?? "")
  const [email, setEmail] = useState(teacher?.email ?? "")
  const [phone, setPhone] = useState(teacher?.phone ?? "")
  const [errors, setErrors] = useState<{
    fullName?: string
    email?: string
    phone?: string
    general?: string
  }>({})

  useEffect(() => {
    if (!open) return
    setFullName(teacher?.full_name ?? "")
    setEmail(teacher?.email ?? "")
    setPhone(teacher?.phone ?? "")
    setErrors({})
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open])

  const createTeacher = useCreateTeacher(t)
  const updateTeacher = useUpdateTeacher(t)
  const isSubmitting = createTeacher.isPending || updateTeacher.isPending

  const handleSubmit = async () => {
    setErrors({})

    try {
      const payload = { full_name: fullName, email, phone }
      if (isEdit) {
        await updateTeacher.mutateAsync({ id: teacher.id, payload })
      } else {
        await createTeacher.mutateAsync(payload)
      }
      onSuccess()
      return true
    } catch (err) {
      const apiError = normalizeApiError(err)
      const fullNameError = getFieldError(apiError, "fullName")
      const emailError = getFieldError(apiError, "email")
      const phoneError = getFieldError(apiError, "phone")

      setErrors({
        fullName: fullNameError
          ? getTeachersFieldMessage(t, "fullName", fullNameError)
          : undefined,
        email: emailError
          ? getTeachersFieldMessage(t, "email", emailError)
          : undefined,
        phone: phoneError
          ? getTeachersFieldMessage(t, "phone", phoneError)
          : undefined,
        general:
          !fullNameError && !emailError && !phoneError
            ? getApiMessage(t, apiError)
            : undefined,
      })
      return false
    }
  }

  return {
    isEdit,
    fullName,
    setFullName,
    email,
    setEmail,
    phone,
    setPhone,
    errors,
    isSubmitting,
    handleSubmit,
  }
}
