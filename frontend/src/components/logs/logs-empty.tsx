import { motion } from "framer-motion"
import { Filter } from "lucide-react"
import { Button } from "@/components/ui/button"

interface LogsEmptyStateProps {
  t: (key: string, fallback?: string) => string
  hasActiveFilters: boolean
  onReset: () => void
}

export function LogsEmptyState({
  t,
  hasActiveFilters,
  onReset,
}: LogsEmptyStateProps) {
  return (
    <motion.div
      key="empty-state"
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="flex flex-col items-center justify-center py-16 text-center"
    >
      <div className="mb-4 rounded-md bg-primary p-3">
        <Filter className="h-5 w-5 text-primary-foreground" />
      </div>

      <h3 className="text-lg font-semibold text-foreground">
        {t("logs.empty.title")}
      </h3>

      <p className="mt-1 text-sm text-muted-foreground max-w-sm">
        {t("logs.empty.description")}
      </p>

      {hasActiveFilters && (
        <Button variant="outline" className="mt-4" onClick={onReset}>
          {t("logs.empty.reset")}
        </Button>
      )}
    </motion.div>
  )
}