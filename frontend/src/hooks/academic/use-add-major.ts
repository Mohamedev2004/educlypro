import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { AcademicService } from "@/api/services/academic-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import type { AddMajorPayload } from "@/api/types/academic.types"

/**
 * Hook for adding a major to a grade in the caller's own center.
 *
 * Responsibility: trigger the mutation + invalidate the academic tree.
 * Layer: Hooks
 */
export function useAddMajor(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      gradeId,
      payload,
    }: {
      gradeId: number
      payload: AddMajorPayload
    }) => AcademicService.addMajor(gradeId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["academic-tree"] })
      const toastId = toast.success(t("onboarding.majorAddedSuccess"), {
        description: t("onboarding.majorAddedSuccessDescription"),
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
