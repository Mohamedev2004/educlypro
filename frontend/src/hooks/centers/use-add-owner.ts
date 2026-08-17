import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { CentersService } from "@/api/services/centers-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import type { AddOwnerPayload } from "@/api/types/centers.types"

/**
 * Hook for assigning an owner to a center.
 *
 * Responsibility: trigger the add-owner mutation + invalidate the center
 * detail and centers list on success.
 * Layer: Hooks
 */
export function useAddOwner(t: (key: string) => string, centerSlug: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: AddOwnerPayload) => CentersService.addOwner(centerSlug, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["center", centerSlug] })
      queryClient.invalidateQueries({ queryKey: ["centers"] })
      const toastId = toast.success(t("centers.addOwnerSuccess"), {
        description: t("centers.addOwnerSuccessDescription"),
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
