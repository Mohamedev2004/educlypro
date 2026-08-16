import { AnimatePresence, motion } from "framer-motion"
import { ChevronDown } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import type { LogItem as Log } from "@/api/types/logs.types"
import type { Locale } from "@/context/direction/direction-provider"

type LogLevel = "info" | "warning" | "error"

const levelStyles: Record<LogLevel, string> = {
  info: "bg-green-500/10 text-green-600 dark:text-green-400",
  warning: "bg-yellow-500/10 text-yellow-600 dark:text-yellow-400",
  error: "bg-red-500/10 text-red-600 dark:text-red-400",
}

const statusStyles: Record<string, string> = {
  "200": "text-green-600 dark:text-green-400",
  "201": "text-green-600 dark:text-green-400",
  "401": "text-yellow-600 dark:text-yellow-400",
  "403": "text-yellow-600 dark:text-yellow-400",
  "429": "text-yellow-600 dark:text-yellow-400",
  "502": "text-red-600 dark:text-red-400",
  "503": "text-red-600 dark:text-red-400",
  warning: "text-yellow-600 dark:text-yellow-400",
}

interface LogRowProps {
  log: Log
  expanded: boolean
  onToggle: () => void
  locale: Locale
  t: (key: string, fallback?: string) => string
  tracedRequestId?: string | null
  onTraceRequest?: (requestId: string) => void
}

export function LogRow({
  log,
  expanded,
  onToggle,
  locale,
  t,
  tracedRequestId,
  onTraceRequest,
}: LogRowProps) {
  const statusLabel = log.status_code ? String(log.status_code) : log.status
  const formattedTime = new Date(log.timestamp).toLocaleTimeString(
    locale === "ar" ? "ar-MA" : locale === "fr" ? "fr-FR" : "en-US",
    {
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    }
  )

  const hasPayload =
    log.payload !== null &&
    log.payload !== undefined &&
    !(typeof log.payload === "object" && Object.keys(log.payload as object).length === 0)
  const isTraced = !!tracedRequestId && log.request_id === tracedRequestId

  return (
    <>
      <motion.div
        onClick={onToggle}
        role="button"
        tabIndex={0}
        onKeyDown={(event) => {
          if (event.key === "Enter" || event.key === " ") {
            event.preventDefault()
            onToggle()
          }
        }}
        className={`w-full p-4 text-left transition-colors hover:bg-muted/50 active:bg-muted/70 ${
          isTraced ? "bg-primary/5 ring-1 ring-primary/30" : ""
        }`}
        whileHover={{ backgroundColor: "rgba(0,0,0,0.02)" }}
      >
        <div className="flex items-center gap-4 cursor-pointer">
          <motion.div
            animate={{ rotate: expanded ? 180 : 0 }}
            transition={{ duration: 0.2 }}
            className="flex-shrink-0"
          >
            <ChevronDown className="h-4 w-4 text-muted-foreground" />
          </motion.div>

          <Badge
            variant="secondary"
            className={`flex-shrink-0 capitalize ${levelStyles[log.level as LogLevel]}`}
          >
            {log.level}
          </Badge>

          <time className="w-20 flex-shrink-0 whitespace-nowrap font-mono text-xs text-muted-foreground">
            {formattedTime}
          </time>

          <span className="flex-shrink-0 min-w-max text-sm font-medium text-foreground">
            {log.service}
          </span>

          {log.request_id && (
            <button
              type="button"
              onClick={(event) => {
                event.stopPropagation()
                onTraceRequest?.(log.request_id!)
              }}
              className={`hidden rounded border cursor-pointer px-2 py-0.5 font-mono text-xs md:inline-flex ${
                isTraced
                  ? "border-primary bg-primary/10 text-primary"
                  : "border-border text-muted-foreground hover:border-primary/40 hover:text-primary"
              }`}
            >
              {log.request_id}
            </button>
          )}

          <p className="flex-1 truncate text-sm text-muted-foreground">
            {log.message}
          </p>

          <span
            className={`flex-shrink-0 font-mono text-sm font-semibold ${
              statusStyles[statusLabel] ?? "text-muted-foreground"
            }`}
          >
            {statusLabel}
          </span>

          <span className="w-16 flex-shrink-0 text-right font-mono text-xs text-muted-foreground">
            {log.duration}
          </span>
        </div>
      </motion.div>

      <AnimatePresence initial={false}>
        {expanded && (
          <motion.div
            key="details"
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden bg-sidebar rounded-sm mb-2"
          >
            <div className="space-y-4 p-4">
              <div>
                <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                  {t("logs.row.message")}
                </p>
                <p className="rounded border border-border p-3 font-mono text-sm text-foreground">
                  {log.message}
                </p>
              </div>

              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                    {t("logs.row.duration")}
                  </p>
                  <p className="font-mono text-foreground">{log.duration}</p>
                </div>
                <div>
                  <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                    {t("logs.row.timestamp")}
                  </p>
                  <p className="font-mono text-xs text-foreground">
                    {log.timestamp}
                  </p>
                </div>
              </div>

              {log.request_id && (
                <div>
                  <p className="mb-1 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                    {t("logs.row.requestId", "Request ID")}
                  </p>
                  <button
                    type="button"
                    onClick={() => onTraceRequest?.(log.request_id!)}
                    className={`rounded border px-2 py-1 font-mono text-xs ${
                      isTraced
                        ? "border-primary bg-primary/10 text-primary"
                        : "border-border text-muted-foreground hover:border-primary/40 hover:text-primary"
                    }`}
                  >
                    {log.request_id}
                  </button>
                </div>
              )}

              {hasPayload && (
                <div>
                  <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-muted-foreground">
                    {t("logs.row.payload", "Payload")}
                  </p>
                  <pre className="rounded border border-border bg-muted/40 p-3 font-mono text-xs text-foreground overflow-x-auto whitespace-pre-wrap break-all">
                    {JSON.stringify(log.payload, null, 2)}
                  </pre>
                </div>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </>
  )
}