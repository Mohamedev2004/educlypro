/**
 * Log related types.
 * 
 * Responsibility: Define UI-specific types for logs.
 * Layer: Types
 */

export type Filters = {
  level: string[]
  status: string[]
  status_code: string[]
  duration: string[]
}

export type Facets = {
  levels: string[]
  statuses: string[]
  status_codes: number[]
  durations: string[]
}
