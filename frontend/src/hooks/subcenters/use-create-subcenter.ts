import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { SubCentersService } from "@/api/services/subcenters-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for creating a sub-center.
 *
 * Responsibility: trigger create + invalidate the sub-centers list on success.
 * Layer: Hooks
 */
export function useCreateSubCenter(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: SubCentersService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["subcenters"] })
      const toastId = toast.success(t("subcenters.createSuccess"), {
        description: t("subcenters.createSuccessDescription"),
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
