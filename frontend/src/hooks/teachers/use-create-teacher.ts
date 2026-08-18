import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { TeachersService } from "@/api/services/teachers-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for creating a teacher.
 *
 * Responsibility: trigger create + invalidate the teachers list on success.
 * Layer: Hooks
 */
export function useCreateTeacher(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: TeachersService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["teachers"] })
      const toastId = toast.success(t("teachers.createSuccess"), {
        description: t("teachers.createSuccessDescription"),
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
