import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { SubCentersService } from "@/api/services/subcenters-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for deleting a sub-center.
 *
 * Responsibility: trigger delete + invalidate the sub-centers list on success.
 * Layer: Hooks
 */
export function useDeleteSubCenter(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => SubCentersService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["subcenters"] })
      const toastId = toast.success(t("subcenters.deleteSuccess"), {
        description: t("subcenters.deleteSuccessDescription"),
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
