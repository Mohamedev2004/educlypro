import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { AcademicService } from "@/api/services/academic-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for removing a major (and, transitively, its subjects) from a grade
 * in the caller's own center.
 *
 * Responsibility: trigger the mutation + invalidate the academic tree.
 * Layer: Hooks
 */
export function useRemoveMajor(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (majorId: number) => AcademicService.removeMajor(majorId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["academic-tree"] })
      const toastId = toast.success(t("onboarding.majorRemovedSuccess"), {
        description: t("onboarding.majorRemovedSuccessDescription"),
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
