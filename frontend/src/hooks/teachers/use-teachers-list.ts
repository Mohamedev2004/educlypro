import { useQuery } from "@tanstack/react-query"
import { TeachersService } from "@/api/services/teachers-service"
import type { TeachersListParams } from "@/api/types/teachers.types"

/**
 * Hook for fetching the teachers list.
 *
 * Responsibility: manage teacher fetching state with pagination/search/sort.
 * Layer: Hooks
 */
export function useTeachersList(params: TeachersListParams) {
  return useQuery({
    queryKey: [
      "teachers",
      params.page,
      params.perPage,
      params.q,
      params.sort,
      params.direction,
    ],
    queryFn: () => TeachersService.list(params),
    staleTime: 10_000,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    placeholderData: (prev) => prev,
  })
}
