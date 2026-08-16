/* eslint-disable @typescript-eslint/no-explicit-any */
"use client"

import { AnimatePresence, motion } from "framer-motion"
import { useMemo, useState } from "react"

import AppLayout from "@/layouts/app-layout"
import { useDirection } from "@/context/direction/direction-provider"
import { AppPagination } from "@/components/app-pagination"

import { LogRow } from "@/components/logs/log-row"
import { FilterPanel } from "@/components/logs/logs-filter"
import { LogsToolbar } from "@/components/logs/logs-toolbar"
import { LogsEmptyState } from "@/components/logs/logs-empty"
import { LogsListSkeleton } from "@/components/logs/logs-list-skeleton"
import { interpolate } from "@/utils/common-utils"
import type { LogItem as Log, LogCounts, LogsChartRange } from "@/api/types/logs.types"
import type { Facets } from "@/types/logs"
import { LogCards } from "@/components/logs/log-cards"
import { LogsChart } from "@/components/logs/logs-chart"
import { Button } from "@/components/ui/button"
import { LogsHero } from "@/components/logs/logs-hero"
import { useLogsFilter } from "@/hooks/logs/use-logs-filter"
import { useLogsList } from "@/hooks/logs/use-logs-list"
import { useLogsChart } from "@/hooks/logs/use-logs-chart"
import { useLogsTracing } from "@/hooks/logs/use-logs-tracing"
import { useLogsExport } from "@/hooks/logs/use-logs-export"
import { useLogsStats } from "@/hooks/logs/use-logs-stats"

export default function SuperAdminLogs() {
  const { t, locale } = useDirection()

  const [expandedId, setExpandedId] = useState<string | null>(null)
  const [chartRange, setChartRange] = useState<LogsChartRange>("24h")

  const {
    searchQuery,
    setSearchQuery,
    showFilters,
    setShowFilters,
    filters,
    page,
    setPage,
    perPage,
    setPerPage,
    activeFilters,
    applyFilters,
    handleReset: resetFilters,
  } = useLogsFilter()

  const logsQuery = useLogsList(page, perPage, searchQuery, filters)
  const chartQuery = useLogsChart(chartRange)

  const logs = useMemo(
    () => (logsQuery.data?.items ?? []) as Log[],
    [logsQuery.data?.items]
  )

  const { tracedRequestId, setTracedRequestId, handleTraceRequest, visibleLogs } =
    useLogsTracing(logs)

  const { isExporting, handleExport } = useLogsExport(t)

  const filteredTotal = logsQuery.data?.pagination.total ?? 0
  const counts = useMemo<LogCounts>(
    () => logsQuery.data?.counts ?? { info: 0, warning: 0, error: 0 },
    [logsQuery.data?.counts]
  )

  const { globalTotal, globalCounts } = useLogsStats(
    logsQuery.isFetching,
    filteredTotal,
    counts,
    filters,
    searchQuery
  )

  const pagination = logsQuery.data?.pagination ?? null
  const facets: Facets = logsQuery.data?.facets ?? {
    levels: [],
    statuses: [],
    status_codes: [],
    durations: [],
  }

  const showSkeleton = logsQuery.isLoading && !logsQuery.data

  const handleCardLevelSelect = (level: "info" | "warning" | "error" | null) => {
    if (level === null) {
      applyFilters({ level: [], status: [], status_code: [], duration: [] })
      setTracedRequestId(null)
      return
    }

    applyFilters({
      ...filters,
      level:
        filters.level.length === 1 && filters.level[0] === level ? [] : [level],
    })
  }

  const handleReset = () => {
    resetFilters()
    setTracedRequestId(null)
  }

  return (
    <AppLayout
      breadcrumbs={[
        { label: t("roles.overview"), href: "/super-admin/dashboard" },
        { label: t("roles.logs") },
      ]}
    >
      <main className="w-full bg-background">
        <div className="flex h-full flex-col gap-4">
          <LogsHero t={t} />

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 items-stretch">
            {/* LEFT: CHART */}
            <div className="lg:col-span-2 rounded-lg">
              <LogsChart
                data={chartQuery.data ?? []}
                range={chartRange}
                onRangeChange={setChartRange}
                locale={locale}
                t={t}
              />
            </div>

            {/* RIGHT: CARDS */}
            <div className="h-full">
              <LogCards
                total={globalTotal ?? filteredTotal}
                counts={globalCounts ?? counts}
                activeLevel={
                  filters.level.length === 1 ? (filters.level[0] as any) : null
                }
                onLevelSelect={handleCardLevelSelect}
              />
            </div>
          </div>

          <LogsToolbar
            t={t}
            searchQuery={searchQuery}
            onSearchChange={setSearchQuery}
            perPage={perPage}
            onPerPageChange={setPerPage}
            showFilters={showFilters}
            onToggleFilters={() => setShowFilters((v) => !v)}
            activeFilters={activeFilters}
            total={filteredTotal}
            shown={visibleLogs.length}
            isExporting={isExporting}
            onExport={() => handleExport(searchQuery, filters)}
            hasData={visibleLogs.length > 0}
            summary={interpolate(t("logs.shownOfTotal"), {
              shown: visibleLogs.length,
              total: filteredTotal,
            })}
          />

          <AnimatePresence>
            {tracedRequestId && (
              <motion.div
                initial={{ opacity: 0, y: -6 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -6 }}
                className="flex items-center justify-between rounded-md border border-primary/30 bg-primary/10 px-4 py-2"
              >
                <p className="text-sm text-foreground">
                  {t("logs.trace.active", "Showing logs for request")}{" "}
                  <span className="font-mono">{tracedRequestId}</span>
                </p>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setTracedRequestId(null)}
                >
                  {t("logs.trace.clear", "Clear tracing")}
                </Button>
              </motion.div>
            )}
          </AnimatePresence>

          <div className="flex flex-1 mt-4 overflow-hidden">
            <AnimatePresence initial={false}>
              {showFilters && (
                <motion.div
                  key="filters"
                  initial={{ width: 0, opacity: 0 }}
                  animate={{ width: 280, opacity: 1 }}
                  exit={{ width: 0, opacity: 0 }}
                  transition={{ duration: 0.2 }}
                  className="overflow-hidden"
                >
                  <FilterPanel
                    t={t}
                    filters={filters}
                    onChange={(next) => {
                      setExpandedId(null)
                      applyFilters(next)
                    }}
                    facets={facets}
                  />
                </motion.div>
              )}
            </AnimatePresence>

            {/* RIGHT SIDE (LOGS + PAGINATION) */}
            <div className="flex flex-1 flex-col overflow-hidden">
              {/* SCROLLABLE LOGS */}
              <div className="flex-1 overflow-auto hide-scrollbar">
                <div className="min-w-[900px]">
                  {showSkeleton ? (
                    <LogsListSkeleton rows={perPage} />
                  ) : visibleLogs.length > 0 ? (
                    visibleLogs.map((log, index) => (
                      <motion.div
                        key={log.id}
                        data-log-id={log.id}
                        initial={{ opacity: 0, y: -10 }}
                        animate={{ opacity: 1, y: 0 }}
                        exit={{ opacity: 0, y: -10 }}
                        transition={{ duration: 0.2, delay: index * 0.02 }}
                      >
                        <LogRow
                          log={log}
                          expanded={expandedId === log.id}
                          locale={locale}
                          t={t}
                          tracedRequestId={tracedRequestId}
                          onTraceRequest={(rid) => {
                            setExpandedId(null)
                            setPage(1)
                            handleTraceRequest(rid)
                          }}
                          onToggle={() =>
                            setExpandedId((c) => (c === log.id ? null : log.id))
                          }
                        />
                      </motion.div>
                    ))
                  ) : (
                    <LogsEmptyState
                      t={t}
                      hasActiveFilters={
                        activeFilters > 0 || !!searchQuery || !!tracedRequestId
                      }
                      onReset={handleReset}
                    />
                  )}
                </div>
              </div>

              {/* FIXED PAGINATION */}
              {pagination && (
                <div className="bg-background px-2 py-3 sticky bottom-0 z-10 shadow-[0_-2px_8px_rgba(0,0,0,0.05)]">
                  <AppPagination
                    pagination={pagination}
                    summaryTop={interpolate(t("logs.pagination.showing"), {
                      count: pagination.total,
                    })}
                    summaryBottom={interpolate(
                      t("logs.pagination.pageSummary"),
                      {
                        page: pagination.page,
                        totalPages: pagination.total_pages,
                      }
                    )}
                    onPageChange={(next) => {
                      setExpandedId(null)
                      setPage(next)
                    }}
                  />
                </div>
              )}
            </div>
          </div>
        </div>
      </main>
    </AppLayout>
  )
}
