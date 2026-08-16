import { motion } from "framer-motion"
import { Check } from "lucide-react"
import { Button } from "@/components/ui/button"

import type { Filters, Facets } from "@/types/logs"

interface FilterPanelProps {
  t: (key: string, fallback?: string) => string
  filters: Filters
  onChange: (filters: Filters) => void
  facets: Facets
}

export function FilterPanel({ t, filters, onChange, facets }: FilterPanelProps) {
  const normalizeFilterValue = (category: keyof Filters, value: string) =>
    category === "level" ? value.toLowerCase() : value

  const toggleFilter = (category: keyof Filters, value: string) => {
    const normalizedValue = normalizeFilterValue(category, value)
    const current = filters[category].map((entry) => normalizeFilterValue(category, entry))
    const updated = current.includes(normalizedValue)
      ? current.filter((entry) => entry !== normalizedValue)
      : [...current, normalizedValue]

    onChange({ ...filters, [category]: updated })
  }

  const clearAll = () => {
    onChange({ level: [], status: [], status_code: [], duration: [] })
  }

  const hasActiveFilters = Object.values(filters).some((group) => group.length > 0)

  const filterGroups = [
    { label: t("logs.filters.level"), key: "level" as keyof Filters, items: facets.levels },
    { label: t("logs.filters.status"), key: "status" as keyof Filters, items: facets.statuses },
    { label: t("logs.filters.statusCode"), key: "status_code" as keyof Filters, items: facets.status_codes.map(String) },
    { label: t("logs.filters.duration"), key: "duration" as keyof Filters, items: facets.durations },
  ]

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ delay: 0.05 }}
      className="flex h-full flex-col space-y-6 overflow-y-auto hide-scrollbar border border-border rounded-md p-4 bg-card"
    >
      <div className="flex items-center justify-between">
        <h3 className="text-md font-semibold text-foreground">
          {t("logs.filters.title")}
        </h3>
        {hasActiveFilters && (
          <Button variant="ghost" onClick={clearAll}>
            {t("logs.filters.clear")}
          </Button>
        )}
      </div>

      {filterGroups.map(({ label, key, items }) => (
        <div key={key} className="space-y-3">
          <p className="text-xs font-semibold uppercase tracking-wide text-muted-foreground">
            {label}
          </p>
          <div className="space-y-2">
            {items.map((item) => {
              const selected = filters[key]
                .map((entry) => normalizeFilterValue(key, entry))
                .includes(normalizeFilterValue(key, item))
              return (
                <motion.button
                  key={item}
                  type="button"
                  whileHover={{ x: 2 }}
                  onClick={() => toggleFilter(key, item)}
                  aria-pressed={selected}
                  className={`flex w-full items-center justify-between gap-2 border rounded-md px-3 py-1.5 text-sm transition-colors ${
                    selected
                      ? "border-primary bg-primary text-primary-foreground"
                      : "border-border text-muted-foreground hover:border-primary/40 hover:bg-primary/15"
                  }`}
                >
                  <span className="capitalize">{item}</span>
                  {selected && <Check className="h-3.5 w-3.5" />}
                </motion.button>
              )
            })}
          </div>
        </div>
      ))}
    </motion.div>
  )
}