import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { AcademicService } from "@/api/services/academic-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for removing a grade (and, transitively, its majors and subjects)
 * from the caller's own center.
 *
 * Responsibility: trigger the mutation + invalidate the academic tree and
 * the cached user — removing the last grade flips has_grades back to false.
 * Layer: Hooks
 */
export function useRemoveGrade(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (gradeId: number) => AcademicService.removeGrade(gradeId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["academic-tree"] })
      queryClient.invalidateQueries({ queryKey: ["me"] })
      const toastId = toast.success(t("onboarding.gradeRemovedSuccess"), {
        description: t("onboarding.gradeRemovedSuccessDescription"),
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
