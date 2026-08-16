import type { LogCounts } from "@/api/types/logs.types"
import { Badge } from "@/components/ui/badge"
import { useDirection } from "@/context/direction/direction-provider"
import { cn } from "@/utils/ui-utils"

type Props = {
  total: number
  counts: LogCounts
  activeLevel: string | null
  onLevelSelect: (level: keyof LogCounts | null) => void
}

function pct(n: number, total: number) {
  if (!total) return "0%"
  return `${((n / total) * 100).toFixed(1)}%`
}

const levelTextClass: Record<keyof LogCounts, string> = {
  info: "text-green-500 dark:text-green-400",
  warning: "text-yellow-500 dark:text-yellow-400",
  error: "text-red-500 dark:text-red-400",
}

const levelBadgeClass: Record<keyof LogCounts, string> = {
  info: "bg-green-500/10 text-green-600 dark:text-green-400",
  warning: "bg-yellow-500/10 text-yellow-600 dark:text-yellow-400",
  error: "bg-red-500/10 text-red-600 dark:text-red-400",
}

const levelActiveClass: Record<keyof LogCounts, string> = {
  info: "bg-green-500/10",
  warning: "bg-yellow-500/10",
  error: "bg-red-500/10",
}

const levels: (keyof LogCounts)[] = ["info", "warning", "error"]

export function LogCards({ total, counts, activeLevel, onLevelSelect }: Props) {
  const { t, locale, direction } = useDirection()
  const isRTL = direction === "rtl"
  const fmtLocale = locale === "ar" ? "ar-MA" : locale === "fr" ? "fr-FR" : "en-US"
  const fmt = (n: number) => n.toLocaleString(fmtLocale)

  return (
    <div className="flex flex-col h-full rounded-md border border-border bg-card overflow-hidden">

      {/* Total row */}
      <button
        type="button"
        onClick={() => onLevelSelect(null)}
        className={cn(
          "flex items-center justify-between px-5 py-4 transition-all cursor-pointer hover:bg-muted/40 flex-1",
          isRTL ? "flex-row-reverse text-left" : "text-left",
          !activeLevel && "bg-muted/20"
        )}
      >
        <div className={cn(isRTL && "text-left")}>
          <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-1">
            {t("logs.cards.total", "Total logs")}
          </p>
          <p className="text-2xl font-medium text-foreground">{fmt(total)}</p>
        </div>
        <p className={cn("text-xs text-muted-foreground", isRTL && "text-left")}>
          {t("logs.cards.allTime", "all time")}
        </p>
      </button>

      {/* Level rows */}
      {levels.map((level) => (
        <button
          key={level}
          type="button"
          onClick={() => onLevelSelect(level)}
          className={cn(
            "flex items-center justify-between px-5 py-4 transition-all cursor-pointer border-t border-border flex-1",
            isRTL ? "flex-row-reverse text-left" : "text-left",
            activeLevel === level
              ? levelActiveClass[level]
              : "hover:bg-muted/40"
          )}
        >
          <div className={cn(isRTL && "text-left")}>
            <Badge
              className={cn(
                "inline-flex items-center text-xs font-medium capitalize mb-1",
                levelBadgeClass[level]
              )}
            >
              {t(`logs.filters.${level}`, level)}
            </Badge>
            <p className={cn("text-2xl font-medium", levelTextClass[level])}>
              {fmt(counts[level])}
            </p>
          </div>
          <p className={cn("text-xs opacity-70", levelTextClass[level], isRTL ? "text-left" : "text-left")}>
            {pct(counts[level], total)}<br />
            <span className={cn("text-xs", levelTextClass[level])}>{t("logs.cards.ofTotal", "of total")}</span>
          </p>
        </button>
      ))}
    </div>
  )
}