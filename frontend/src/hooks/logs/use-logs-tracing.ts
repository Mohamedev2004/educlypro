import { useState, useEffect } from "react"
import type { LogItem as Log } from "@/api/types/logs.types"

export function useLogsTracing(logs: Log[]) {
  const [tracedRequestId, setTracedRequestId] = useState<string | null>(null)

  const handleTraceRequest = (requestId: string) => {
    setTracedRequestId(requestId)
  }

  const visibleLogs = tracedRequestId
    ? logs.filter((log) => log.request_id === tracedRequestId)
    : logs

  useEffect(() => {
    if (!tracedRequestId || visibleLogs.length === 0) return
    const firstRow = document.querySelector<HTMLElement>(
      `[data-log-id="${visibleLogs[0].id}"]`
    )
    firstRow?.scrollIntoView({ behavior: "smooth", block: "center" })
  }, [tracedRequestId, visibleLogs])

  return {
    tracedRequestId,
    setTracedRequestId,
    handleTraceRequest,
    visibleLogs,
  }
}
