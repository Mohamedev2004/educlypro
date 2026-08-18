import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { SubCentersService } from "@/api/services/subcenters-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import type { UpdateSubCenterPayload } from "@/api/types/subcenters.types"

/**
 * Hook for renaming a sub-center.
 *
 * Responsibility: trigger update + invalidate the sub-centers list on success.
 * Layer: Hooks
 */
export function useUpdateSubCenter(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({
      id,
      payload,
    }: {
      id: number
      payload: UpdateSubCenterPayload
    }) => SubCentersService.update(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["subcenters"] })
      const toastId = toast.success(t("subcenters.updateSuccess"), {
        description: t("subcenters.updateSuccessDescription"),
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
