import { useQuery } from "@tanstack/react-query"
import { CentersService } from "@/api/services/centers-service"

/**
 * Hook for fetching a single center's detail view.
 *
 * Responsibility: manage center detail fetching state.
 * Layer: Hooks
 */
export function useCenterDetail(slug: string) {
  return useQuery({
    queryKey: ["center", slug],
    queryFn: () => CentersService.getBySlug(slug),
    enabled: slug.length > 0,
  })
}
