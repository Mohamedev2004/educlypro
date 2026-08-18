import { useQuery } from "@tanstack/react-query"
import { SubCentersService } from "@/api/services/subcenters-service"

/**
 * Hook for fetching the caller's own center's sub-centers.
 *
 * Responsibility: fetch + cache the sub-centers list.
 * Layer: Hooks
 */
export function useSubCentersList() {
  return useQuery({
    queryKey: ["subcenters"],
    queryFn: SubCentersService.list,
  })
}
