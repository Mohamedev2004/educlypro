import { useQuery } from "@tanstack/react-query"
import { LogsService } from "@/api/services/logs-service"
import type { Filters } from "@/types/logs"

/**
 * Hook for fetching logs list.
 * 
 * Responsibility: Manage logs fetching state with pagination and filters.
 * Layer: Hooks
 */
export function useLogsList(
  page: number,
  perPage: number,
  searchQuery: string,
  filters: Filters
) {
  return useQuery({
    queryKey: [
      "logs",
      page,
      perPage,
      searchQuery,
      filters.level,
      filters.status,
      filters.status_code,
      filters.duration,
    ],
    queryFn: () =>
      LogsService.list({
        page,
        perPage,
        q: searchQuery || undefined,
        level: filters.level.length ? filters.level : undefined,
        status: filters.status.length ? filters.status : undefined,
        status_code: filters.status_code.length
          ? filters.status_code.map((v) => Number(v))
          : undefined,
        duration: filters.duration.length ? filters.duration : undefined,
      }),
    staleTime: 30_000,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    placeholderData: (prev) => prev,
  })
}
