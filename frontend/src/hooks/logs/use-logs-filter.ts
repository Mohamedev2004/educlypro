import { useState, useEffect } from "react"
import type { Filters } from "@/types/logs"

export function useLogsFilter() {
  const [searchQuery, setSearchQuery] = useState("")
  const [showFilters, setShowFilters] = useState(false)
  const [filters, setFilters] = useState<Filters>({
    level: [],
    status: [],
    status_code: [],
    duration: [],
  })
  const [page, setPage] = useState(1)
  const [perPage, setPerPage] = useState(10)

  useEffect(() => {
    setPage(1)
  }, [
    filters.level,
    filters.status,
    filters.status_code,
    filters.duration,
    searchQuery,
    perPage,
  ])

  const activeFilters =
    filters.level.length +
    filters.status.length +
    filters.status_code.length +
    filters.duration.length

  const applyFilters = (next: Filters) => {
    setPage(1)
    setFilters({
      ...next,
      level: next.level.map((l) => l.toLowerCase()),
    })
  }

  const handleReset = () => {
    applyFilters({ level: [], status: [], status_code: [], duration: [] })
    setSearchQuery("")
  }

  return {
    searchQuery,
    setSearchQuery,
    showFilters,
    setShowFilters,
    filters,
    setFilters,
    page,
    setPage,
    perPage,
    setPerPage,
    activeFilters,
    applyFilters,
    handleReset,
  }
}
