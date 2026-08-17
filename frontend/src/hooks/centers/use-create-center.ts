import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { CentersService } from "@/api/services/centers-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for creating a center.
 *
 * Responsibility: trigger create + invalidate the centers list on success.
 * Layer: Hooks
 */
export function useCreateCenter(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: CentersService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["centers"] })
      const toastId = toast.success(t("centers.createSuccess"), {
        description: t("centers.createSuccessDescription"),
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
