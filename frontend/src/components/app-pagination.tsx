import { motion } from "framer-motion"
import {
  Pagination,
  PaginationContent,
  PaginationItem,
  PaginationLink,
} from "@/components/ui/pagination"
import {
  ChevronLeft,
  ChevronRight,
  ChevronsLeft,
  ChevronsRight,
} from "lucide-react"

export type AppPaginationMeta = {
  page: number
  per_page: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

type AppPaginationProps = {
  pagination: AppPaginationMeta
  summaryTop?: React.ReactNode
  summaryBottom?: React.ReactNode
  onPageChange: (page: number) => void
}

export function AppPagination({
  pagination,
  summaryTop,
  summaryBottom,
  onPageChange,
}: AppPaginationProps) {
  return (
    <motion.div
      initial={{ opacity: 0, y: 8 }}
      animate={{ opacity: 1, y: 0 }}
      transition={{ duration: 0.2 }}
      className="flex flex-col gap-4 border-t border-border/70 pt-4 lg:flex-row lg:items-center lg:justify-between"
    >
      <div className="space-y-1 text-sm text-muted-foreground">
        {summaryTop ? <div>{summaryTop}</div> : null}
        {summaryBottom ? <div>{summaryBottom}</div> : null}
      </div>

      <Pagination className="mx-0 w-auto justify-end">
        <PaginationContent>
          <PaginationItem>
            <PaginationLink
              href="#"
              aria-disabled={pagination.page === 1}
              className={`flex items-center justify-center rounded-md border bg-background px-3 py-2 hover:bg-background border-border${
                pagination.page === 1 ? "pointer-events-none opacity-50" : ""
              }`}
              onClick={(e) => {
                e.preventDefault()
                if (pagination.page !== 1) onPageChange(1)
              }}
            >
              <ChevronsLeft className="size-4 rtl:rotate-180" />
            </PaginationLink>
          </PaginationItem>

          <PaginationItem>
            <PaginationLink
              href="#"
              aria-disabled={!pagination.has_prev}
              className={`flex items-center justify-center rounded-md border bg-background px-3 py-2 hover:bg-background border-border${
                !pagination.has_prev ? "pointer-events-none opacity-50" : ""
              }`}
              onClick={(e) => {
                e.preventDefault()
                if (pagination.has_prev) onPageChange(pagination.page - 1)
              }}
            >
              <ChevronLeft className="size-4 rtl:rotate-180" />
            </PaginationLink>
          </PaginationItem>

          <PaginationItem>
            <PaginationLink
              href="#"
              aria-disabled={!pagination.has_next}
              className={`flex items-center justify-center rounded-md border bg-background px-3 py-2 hover:bg-background border-border${
                !pagination.has_next ? "pointer-events-none opacity-50" : ""
              }`}
              onClick={(e) => {
                e.preventDefault()
                if (pagination.has_next) onPageChange(pagination.page + 1)
              }}
            >
              <ChevronRight className="size-4 rtl:rotate-180" />
            </PaginationLink>
          </PaginationItem>

          <PaginationItem>
            <PaginationLink
              href="#"
              aria-disabled={pagination.page === pagination.total_pages}
              className={`flex items-center justify-center rounded-md border border-border bg-background px-3 py-2 hover:bg-background ${
                pagination.page === pagination.total_pages
                  ? "pointer-events-none opacity-50"
                  : ""
              }`}
              onClick={(e) => {
                e.preventDefault()
                if (pagination.page !== pagination.total_pages) {
                  onPageChange(pagination.total_pages)
                }
              }}
            >
              <ChevronsRight className="size-4 rtl:rotate-180" />
            </PaginationLink>
          </PaginationItem>
        </PaginationContent>
      </Pagination>
    </motion.div>
  )
}

