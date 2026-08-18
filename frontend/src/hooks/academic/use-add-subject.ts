import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { AcademicService } from "@/api/services/academic-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import type { AddSubjectPayload } from "@/api/types/academic.types"

/**
 * Hook for adding a subject to a major in the caller's own center.
 *
 * Responsibility: trigger the mutation + invalidate the academic tree.
 * Layer: Hooks
 */
export function useAddSubject(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      majorId,
      payload,
    }: {
      majorId: number
      payload: AddSubjectPayload
    }) => AcademicService.addSubject(majorId, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["academic-tree"] })
    },
    onError: (error) => {
      toast.error(getApiMessage(t, normalizeApiError(error)))
    },
  })
}
