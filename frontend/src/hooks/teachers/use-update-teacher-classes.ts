import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { TeachersService } from "@/api/services/teachers-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for replacing a teacher's full class assignment, from the teacher
 * detail page's classes editor.
 *
 * Responsibility: trigger the update + invalidate the teachers list and
 * this teacher's detail view on success.
 * Layer: Hooks
 */
export function useUpdateTeacherClasses(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, classIds }: { id: number; classIds: number[] }) =>
      TeachersService.updateClasses(id, { class_ids: classIds }),
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
