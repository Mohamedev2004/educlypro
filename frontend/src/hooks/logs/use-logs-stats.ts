import { useState, useEffect } from "react"
import type { LogCounts } from "@/api/types/logs.types"
import type { Filters } from "@/types/logs"

export function useLogsStats(
  isFetching: boolean,
  filteredTotal: number,
  counts: LogCounts,
  filters: Filters,
  searchQuery: string
) {
  const [globalTotal, setGlobalTotal] = useState<number | null>(null)
  const [globalCounts, setGlobalCounts] = useState<LogCounts | null>(null)

  useEffect(() => {
    if (
      !isFetching &&
      globalTotal === null &&
      filters.level.length === 0 &&
      filters.status.length === 0 &&
      filters.status_code.length === 0 &&
      filters.duration.length === 0 &&
      !searchQuery
    ) {
      setGlobalTotal(filteredTotal)
      setGlobalCounts(counts)
    }
  }, [isFetching, filteredTotal, counts, globalTotal, filters, searchQuery])

  return {
    globalTotal,
    globalCounts,
  }
}
