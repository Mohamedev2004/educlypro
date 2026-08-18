import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { TeachersService } from "@/api/services/teachers-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for deleting a teacher.
 *
 * Responsibility: trigger delete + invalidate the teachers list and this
 * teacher's detail view on success.
 * Layer: Hooks
 */
export function useDeleteTeacher(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id }: { id: number; slug: string }) =>
      TeachersService.delete(id),
    onSuccess: (_data, { slug }) => {
      queryClient.invalidateQueries({ queryKey: ["teachers"] })
      queryClient.invalidateQueries({ queryKey: ["teacher", slug] })
      const toastId = toast.success(t("teachers.deleteSuccess"), {
        description: t("teachers.deleteSuccessDescription"),
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
