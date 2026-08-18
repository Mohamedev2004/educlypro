import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { TeachersService } from "@/api/services/teachers-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import type { UpdateTeacherPayload } from "@/api/types/teachers.types"

/**
 * Hook for updating a teacher.
 *
 * Responsibility: trigger update + invalidate the teachers list and this
 * teacher's detail view on success.
 * Layer: Hooks
 */
export function useUpdateTeacher(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: number
      payload: UpdateTeacherPayload
    }) => TeachersService.update(id, payload),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ["teachers"] })
      queryClient.invalidateQueries({ queryKey: ["teacher", data.slug] })
      const toastId = toast.success(t("teachers.updateSuccess"), {
        description: t("teachers.updateSuccessDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
    },
    onError: (error) => {
      toast.error(getApiMessage(t, normalizeApiError(error)))
    },
  })
}
