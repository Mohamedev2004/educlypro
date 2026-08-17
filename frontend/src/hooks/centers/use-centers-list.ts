import { useQuery } from "@tanstack/react-query"
import { CentersService } from "@/api/services/centers-service"
import type { CentersListParams } from "@/api/types/centers.types"

/**
 * Hook for fetching the centers list.
 *
 * Responsibility: manage centers fetching state with pagination/search/sort.
 * Layer: Hooks
 */
export function useCentersList(params: CentersListParams) {
  return useQuery({
    queryKey: [
      "centers",
      params.page,
      params.perPage,
      params.q,
      params.sort,
      params.direction,
    ],
    queryFn: () => CentersService.list(params),
    staleTime: 10_000,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    placeholderData: (prev) => prev,
  })
}
