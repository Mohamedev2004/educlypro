import { axiosInstance } from "../axios"
import type {
  LogsChartEnvelope,
  LogsChartPoint,
  LogsChartRange,
  LogsData,
  LogsEnvelope,
  LogsExportParams,
  LogsListParams,
} from "../types/logs.types"

export const LogsService = {
  async list(params: LogsListParams): Promise<LogsData> {
    const response = await axiosInstance.get<LogsEnvelope>("/logs", {
      params: {
        page: params.page,
        per_page: params.perPage,
        q: params.q,
        level: params.level,
        status: params.status,
        status_code: params.status_code,
        duration: params.duration,
        from: params.from,
        to: params.to,
      },
    })

    return response.data.data
  },

  async chart(range: LogsChartRange): Promise<LogsChartPoint[]> {
    const response = await axiosInstance.get<LogsChartEnvelope>("/logs/chart", {
      params: { range },
    })
    return response.data.data
  },

  async export(params: LogsExportParams): Promise<void> {
    const response = await axiosInstance.get("/logs/export", {
      params: {
        q: params.q,
        level: params.level,
        status: params.status,
        status_code: params.status_code,
        duration: params.duration,
        from: params.from,
        to: params.to,
      },
      responseType: "blob",
    })

    const url = URL.createObjectURL(new Blob([response.data]))
    const link = document.createElement("a")
    link.href = url
    link.download = `audit-logs-${new Date().toISOString().slice(0, 10)}.xlsx`
    document.body.appendChild(link)
    link.click()
    link.remove()
    URL.revokeObjectURL(url)
  },
}



