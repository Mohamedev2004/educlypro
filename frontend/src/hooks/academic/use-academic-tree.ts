import { useQuery } from "@tanstack/react-query"
import { AcademicService } from "@/api/services/academic-service"

/**
 * Hook for fetching the center owner's grade/major/subject tree.
 *
 * Responsibility: fetch + cache the academic tree that drives the
 * onboarding wizard.
 * Layer: Hooks
 */
export function useAcademicTree() {
  return useQuery({
    queryKey: ["academic-tree"],
    queryFn: AcademicService.tree,
  })
}
