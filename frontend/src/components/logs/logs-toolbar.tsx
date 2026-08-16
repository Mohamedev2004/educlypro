import { Download, Filter, Search } from "lucide-react"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import React from "react"
import { LOGS_PER_PAGE_OPTIONS } from "@/constants/logs"

interface LogsToolbarProps {
  t: (key: string, fallback?: string) => string
  searchQuery: string
  onSearchChange: (value: string) => void
  perPage: number
  onPerPageChange: (value: number) => void
  showFilters: boolean
  onToggleFilters: () => void
  activeFilters: number
  total: number
  shown: number
  isExporting: boolean
  onExport: () => void
  hasData: boolean
  summary: string
}

export function LogsToolbar({
  t,
  searchQuery,
  onSearchChange,
  perPage,
  onPerPageChange,
  showFilters,
  onToggleFilters,
  activeFilters,
  isExporting,
  onExport,
  hasData,
  summary,
}: LogsToolbarProps) {
  const [localSearch, setLocalSearch] = React.useState(searchQuery)

  React.useEffect(() => {
    const timeout = setTimeout(() => {
      onSearchChange(localSearch)
    }, 2000)

    return () => clearTimeout(timeout)
  }, [localSearch, onSearchChange])

  React.useEffect(() => {
    setLocalSearch(searchQuery)
  }, [searchQuery])

  return (
    <div className="flex flex-col gap-3">
      {/* Summary */}
      <p className="text-sm text-muted-foreground">
        {summary}
      </p>

      {/* Toolbar */}
      <div className="flex flex-wrap gap-2">
        
        {/* Search */}
        <div className="relative w-full sm:flex-1">
          <Search className="pointer-events-none absolute start-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            placeholder={t("logs.searchPlaceholder")}
            className="h-9 ps-9 text-sm w-full"
            value={localSearch}
            onChange={(e) => setLocalSearch(e.target.value)}
          />
        </div>

        {/* Per page */}
        <div className="w-full sm:w-auto">
          <Select
            value={String(perPage)}
            onValueChange={(value) => onPerPageChange(Number(value))}
          >
            <SelectTrigger className="!h-9 w-full sm:w-[130px]">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {LOGS_PER_PAGE_OPTIONS.map((option) => (
                <SelectItem key={option} value={String(option)}>
                  {option} {t("common.perPage")}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {/* Filters */}
        <Button
          variant={showFilters ? "default" : "outline"}
          onClick={onToggleFilters}
          className="relative h-9 w-full sm:w-auto"
        >
          <Filter className="h-4 w-4 me-2" />
          {t("logs.filters.button")}

          {activeFilters > 0 && (
            <Badge className="absolute bg-primary text-primary-foreground border border-primary-foreground border-2 -end-2 -top-2 h-5 w-5 p-0 text-xs">
              {activeFilters}
            </Badge>
          )}
        </Button>

        {/* Export */}
        {hasData && (
          <Button
            variant="default"
            onClick={onExport}
            disabled={isExporting}
            className="h-9 w-full sm:w-auto"
          >
            <Download className="h-4 w-4 me-2" />
            {isExporting ? t("logs.exporting") : t("logs.export")}
          </Button>
        )}
      </div>
    </div>
  )
}