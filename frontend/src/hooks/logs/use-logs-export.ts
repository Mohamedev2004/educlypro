import { useState } from "react"
import { toast } from "sonner"
import { LogsService } from "@/api/services/logs-service"
import type { Filters } from "@/types/logs"

export function useLogsExport(t: (key: string) => string) {
  const [isExporting, setIsExporting] = useState(false)

  const handleExport = async (searchQuery: string, filters: Filters) => {
    setIsExporting(true)
    try {
      await LogsService.export({
        q: searchQuery || undefined,
        level: filters.level,
        status: filters.status,
        status_code: filters.status_code.map((v) => Number(v)),
        duration: filters.duration,
      })

      const toastId = toast.success(t("logs.success"), {
        description: t("logs.successDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
    } catch {
      toast.error(t("logs.exportFailed"))
    } finally {
      setIsExporting(false)
    }
  }

  return {
    isExporting,
    handleExport,
  }
}
