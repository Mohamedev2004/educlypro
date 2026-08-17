import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { CentersService } from "@/api/services/centers-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import type { AddStaffPayload } from "@/api/types/centers.types"

/**
 * Hook for adding a staff member to a center.
 *
 * Responsibility: trigger the add-staff mutation + invalidate the center
 * detail and centers list (staff_count) on success.
 * Layer: Hooks
 */
export function useAddCenterStaff(t: (key: string) => string, centerSlug: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (payload: AddStaffPayload) => CentersService.addStaff(centerSlug, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["center", centerSlug] })
      queryClient.invalidateQueries({ queryKey: ["centers"] })
      const toastId = toast.success(t("centers.addStaffSuccess"), {
        description: t("centers.addStaffSuccessDescription"),
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
