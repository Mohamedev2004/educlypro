import { useQuery } from "@tanstack/react-query"
import { CentersService } from "@/api/services/centers-service"

/**
 * Hook for fetching an arbitrary center's sub-centers (super_admin), used to
 * populate the sub-center picker in the center detail page's add-staff form.
 *
 * Responsibility: fetch + cache the sub-centers list for one center.
 * Layer: Hooks
 */
export function useCenterSubCenters(centerSlug: string) {
  return useQuery({
    queryKey: ["center", centerSlug, "subcenters"],
    queryFn: () => CentersService.listSubCenters(centerSlug),
    enabled: !!centerSlug,
  })
}
