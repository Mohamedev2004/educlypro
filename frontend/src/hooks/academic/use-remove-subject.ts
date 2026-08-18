import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { AcademicService } from "@/api/services/academic-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for removing a subject from a major in the caller's own center.
 *
 * Responsibility: trigger the mutation + invalidate the academic tree.
 * Layer: Hooks
 */
export function useRemoveSubject(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (subjectId: number) => AcademicService.removeSubject(subjectId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["academic-tree"] })
    },
    onError: (error) => {
      toast.error(getApiMessage(t, normalizeApiError(error)))
    },
  })
}
