import { useQuery } from "@tanstack/react-query"
import { TeachersService } from "@/api/services/teachers-service"

/**
 * Hook for fetching a single teacher's detail view by slug.
 *
 * Responsibility: manage teacher detail fetching state.
 * Layer: Hooks
 */
export function useTeacherDetail(slug: string) {
  return useQuery({
    queryKey: ["teacher", slug],
    queryFn: () => TeachersService.getBySlug(slug),
    enabled: slug.length > 0,
  })
}
