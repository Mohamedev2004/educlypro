import { motion } from "framer-motion"
import { Area, AreaChart, CartesianGrid, XAxis, YAxis } from "recharts"

import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
  type ChartConfig,
} from "@/components/ui/chart"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import type { LogsChartPoint, LogsChartRange } from "@/api/types/logs.types"

function resolveLocale(locale: string) {
  if (locale === "ar") return "ar-MA"
  if (locale === "fr") return "fr-FR"
  return "en-US"
}

type Props = {
  data: LogsChartPoint[]
  range: LogsChartRange
  onRangeChange: (range: LogsChartRange) => void
  locale: string
  t: (key: string, fallback?: string) => string
}

export function LogsChart({ data, range, onRangeChange, locale, t }: Props) {
  const resolvedLocale = resolveLocale(locale)
  const chartConfig = {
    count: {
      label: t("logs.chart.count"),
      color: "var(--chart-1)",
    },
  } satisfies ChartConfig

  return (
    <motion.div
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.25 }}
    >
      <Card className="pt-0">
        <CardHeader className="flex items-center gap-2 space-y-0 border-b py-5 sm:flex-row">
          <div className="grid flex-1 gap-1">
            <CardTitle>{t("logs.chart.title")}</CardTitle>
            <CardDescription>{t("logs.chart.subtitle")}</CardDescription>
          </div>
          <Select value={range} onValueChange={(value) => onRangeChange(value as LogsChartRange)}>
            <SelectTrigger
              className="w-[120px] rounded-lg sm:ml-auto"
              aria-label={t("logs.chart.rangeAria")}
            >
              <SelectValue />
            </SelectTrigger>
            <SelectContent className="rounded-xl">
              <SelectItem value="24h" className="rounded-lg">
                {t("logs.chart.range24h")}
              </SelectItem>
              <SelectItem value="7d" className="rounded-lg">
                {t("logs.chart.range7d")}
              </SelectItem>
            </SelectContent>
          </Select>
        </CardHeader>
        <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
          <ChartContainer config={chartConfig} className="aspect-auto h-[250px] w-full">
          <AreaChart data={data}>
            <defs>
              <linearGradient id="logsGradient" x1="0" y1="0" x2="0" y2="1">
                <stop offset="5%" stopColor="var(--color-count)" stopOpacity={0.8} />
                <stop offset="95%" stopColor="var(--color-count)" stopOpacity={0.1} />
              </linearGradient>
            </defs>
            <CartesianGrid vertical={false} />
            <XAxis
              dataKey="date"
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              minTickGap={24}
              tickFormatter={(value) => {
                const date = new Date(value)
                return range === "24h"
                  ? date.toLocaleTimeString(resolvedLocale, { hour: "2-digit" })
                  : date.toLocaleDateString(resolvedLocale, { month: "short", day: "numeric" })
              }}
            />
            <YAxis
              allowDecimals={false}
              tickLine={false}
              axisLine={false}
              tickMargin={8}
              width={32}
            />
            <ChartTooltip
              cursor={false}
              content={
                <ChartTooltipContent
                  indicator="dot"
                  labelFormatter={(value) => {
                    const date = new Date(String(value))
                    return range === "24h"
                      ? date.toLocaleString(resolvedLocale, { hour: "2-digit", minute: "2-digit" })
                      : date.toLocaleDateString(resolvedLocale, { month: "short", day: "numeric" })
                  }}
                />
              }
            />
            <Area
              dataKey="count"
              type="natural"
              fill="url(#logsGradient)"
              stroke="var(--color-count)"
              strokeWidth={2}
            />
          </AreaChart>
          </ChartContainer>
        </CardContent>
      </Card>
    </motion.div>
  )
}
