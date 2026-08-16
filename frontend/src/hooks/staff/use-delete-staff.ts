import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { StaffService } from "@/api/services/staff-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"

/**
 * Hook for deleting a staff member.
 *
 * Responsibility: trigger delete + invalidate the staff list on success.
 * Layer: Hooks
 */
export function useDeleteStaff(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (id: number) => StaffService.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["staff"] })
      const toastId = toast.success(t("staff.deleteSuccess"), {
        description: t("staff.deleteSuccessDescription"),
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
