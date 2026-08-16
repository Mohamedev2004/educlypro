import { useQuery } from "@tanstack/react-query"
import { StaffService } from "@/api/services/staff-service"
import type { StaffListParams } from "@/api/types/staff.types"

/**
 * Hook for fetching the staff list.
 *
 * Responsibility: manage staff fetching state with pagination/filter/sort.
 * Layer: Hooks
 */
export function useStaffList(params: StaffListParams) {
  return useQuery({
    queryKey: [
      "staff",
      params.page,
      params.perPage,
      params.role,
      params.q,
      params.sort,
      params.direction,
    ],
    queryFn: () => StaffService.list(params),
    staleTime: 10_000,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    placeholderData: (prev) => prev,
  })
}
