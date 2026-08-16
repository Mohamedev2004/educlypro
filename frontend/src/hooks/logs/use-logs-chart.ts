import { useQuery } from "@tanstack/react-query"
import { LogsService } from "@/api/services/logs-service"
import type { LogsChartRange } from "@/api/types/logs.types"

/**
 * Hook for fetching logs chart data.
 * 
 * Responsibility: Manage logs chart data fetching state.
 * Layer: Hooks
 */
export function useLogsChart(chartRange: LogsChartRange) {
  return useQuery({
    queryKey: ["logs-chart", chartRange],
    queryFn: () => LogsService.chart(chartRange),
    staleTime: 30_000,
    refetchOnMount: false,
    refetchOnWindowFocus: false,
    placeholderData: (prev) => prev,
  })
}
