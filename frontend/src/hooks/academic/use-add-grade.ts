import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { AcademicService } from "@/api/services/academic-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import type { AddGradePayload } from "@/api/types/academic.types"

/**
 * Hook for adding a grade to the caller's own center.
 *
 * Responsibility: trigger the mutation + invalidate the academic tree and
 * the cached user — the user's has_grades flag flips once the first grade
 * exists, which is what unlocks the dashboard.
 * Layer: Hooks
 */
export function useAddGrade(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: AddGradePayload) => AcademicService.addGrade(payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["academic-tree"] })
      queryClient.invalidateQueries({ queryKey: ["me"] })
      const toastId = toast.success(t("onboarding.gradeAddedSuccess"), {
        description: t("onboarding.gradeAddedSuccessDescription"),
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
