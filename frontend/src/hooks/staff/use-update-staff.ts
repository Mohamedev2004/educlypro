import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { StaffService } from "@/api/services/staff-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import type { UpdateStaffPayload } from "@/api/types/staff.types"

/**
 * Hook for updating a staff member.
 *
 * Responsibility: trigger update + invalidate the staff list on success.
 * Layer: Hooks
 */
export function useUpdateStaff(t: (key: string) => string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: ({ id, payload }: { id: number; payload: UpdateStaffPayload }) =>
      StaffService.update(id, payload),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["staff"] })
      const toastId = toast.success(t("staff.updateSuccess"), {
        description: t("staff.updateSuccessDescription"),
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
