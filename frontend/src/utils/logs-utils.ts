/**
 * Logs utility functions.
 * 
 * Responsibility: Provide helper functions for logs data processing.
 * Layer: Utils
 */

type LogForChart = {
  timestamp: string
}

export type LogsChartRange = "24h" | "7d"

export type LogsChartPoint = {
  key: string
  date: string
  label: string
  count: number
}

export function groupLogsByTime(
  logs: LogForChart[],
  range: LogsChartRange,
  locale: "ar" | "fr" | "en" | string
): LogsChartPoint[] {
  const now = new Date()
  const points: LogsChartPoint[] = []
  const keyToIndex = new Map<string, number>()
  const formatter = new Intl.DateTimeFormat(
    locale === "ar" ? "ar-MA" : locale === "fr" ? "fr-FR" : "en-US",
    range === "24h" ? { hour: "2-digit" } : { month: "short", day: "numeric" }
  )

  const steps = range === "24h" ? 24 : 7

  for (let i = steps - 1; i >= 0; i -= 1) {
    const bucket = new Date(now)
    if (range === "24h") {
      bucket.setMinutes(0, 0, 0)
      bucket.setHours(bucket.getHours() - i)
      const key = `${bucket.getFullYear()}-${bucket.getMonth()}-${bucket.getDate()}-${bucket.getHours()}`
      keyToIndex.set(key, points.length)
      points.push({
        key,
        date: bucket.toISOString(),
        label: formatter.format(bucket),
        count: 0,
      })
    } else {
      bucket.setHours(0, 0, 0, 0)
      bucket.setDate(bucket.getDate() - i)
      const key = `${bucket.getFullYear()}-${bucket.getMonth()}-${bucket.getDate()}`
      keyToIndex.set(key, points.length)
      points.push({
        key,
        date: bucket.toISOString(),
        label: formatter.format(bucket),
        count: 0,
      })
    }
  }

  for (const log of logs) {
    const logDate = new Date(log.timestamp)
    if (Number.isNaN(logDate.getTime()) || logDate > now) continue

    const key =
      range === "24h"
        ? `${logDate.getFullYear()}-${logDate.getMonth()}-${logDate.getDate()}-${logDate.getHours()}`
        : `${logDate.getFullYear()}-${logDate.getMonth()}-${logDate.getDate()}`

    const idx = keyToIndex.get(key)
    if (idx !== undefined) {
      points[idx].count += 1
    }
  }

  return points
}
