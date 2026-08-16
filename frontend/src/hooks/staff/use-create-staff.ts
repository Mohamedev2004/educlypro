import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { StaffService } from "@/api/services/staff-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for creating a staff member.
 *
 * Responsibility: trigger create + invalidate the staff list on success.
 * Layer: Hooks
 */
export function useCreateStaff(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: StaffService.create,
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["staff"] })
      const toastId = toast.success(t("staff.createSuccess"), {
        description: t("staff.createSuccessDescription"),
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
