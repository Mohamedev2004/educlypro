import type { ApiEnvelope } from "./notification.types"

export type LogsPagination = {
  page: number
  per_page: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

export type LogItem = {
  id: string
  timestamp: string
  level: "info" | "warning" | "error" | string
  service: string
  message: string
  duration: string
  status: string
  status_code: number
  request_id?: string
  actor_id?: string
  entity_id?: string
  payload?: unknown
}

export type LogFacets = {
  levels: string[]
  statuses: string[]
  status_codes: number[]
  durations: string[]
}

export type LogCounts = {
  info: number
  warning: number
  error: number
}

export type LogsData = {
  items: LogItem[]
  facets: LogFacets
  pagination: LogsPagination
  applied: Record<string, unknown>
  counts: LogCounts
}

export type LogsListParams = {
  page: number
  perPage: number
  q?: string
  level?: string[]
  status?: string[]
  status_code?: number[]
  duration?: string[]
  from?: string
  to?: string
}

export type LogsExportParams = {
  q?: string
  level?: string[]
  status?: string[]
  status_code?: number[]
  duration?: string[]
  from?: string
  to?: string
}

export type LogsChartRange = "24h" | "7d"

export type LogsChartPoint = {
  date: string
  count: number
}

export type LogsChartEnvelope = ApiEnvelope<LogsChartPoint[]>

export type LogsEnvelope = ApiEnvelope<LogsData>

