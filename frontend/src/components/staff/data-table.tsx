import * as React from "react"
import {
  type ColumnDef,
  type PaginationState,
  type SortingState,
  type VisibilityState,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table"
import { Plus } from "lucide-react"

import { DataTablePagination } from "@/components/data-table-pagination"
import { DataTableViewOptions } from "@/components/data-table-view-options"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { STAFF_PER_PAGE_OPTIONS } from "@/constants/staff"
import type { Staff, StaffRole } from "@/api/types/staff.types"

/**
 * Staff data-table composition: search + role filter + Add button, wired to
 * a server-driven (manual) @tanstack/react-table instance.
 *
 * Responsibility: orchestrate the generic table primitives for the staff
 * domain. Pure UI — all data fetching lives in the page/hooks.
 * Layer: Components (domain)
 */

interface DataTableProps {
  t: (key: string) => string
  columns: ColumnDef<Staff>[]
  data: Staff[]
  pagination: PaginationState
  pageCount: number
  onPaginationChange: (pagination: PaginationState) => void
  sorting: SortingState
  onSortingChange: (sorting: SortingState) => void
  search: string
  onSearchChange: (value: string) => void
  role: StaffRole | "all"
  onRoleChange: (value: StaffRole | "all") => void
  onAddClick: () => void
}

export function DataTable({
  t,
  columns,
  data,
  pagination,
  pageCount,
  onPaginationChange,
  sorting,
  onSortingChange,
  search,
  onSearchChange,
  role,
  onRoleChange,
  onAddClick,
}: DataTableProps) {
  const [columnVisibility, setColumnVisibility] = React.useState<VisibilityState>({})

  const table = useReactTable({
    data,
    columns,
    pageCount,
    getCoreRowModel: getCoreRowModel(),
    manualPagination: true,
    manualSorting: true,
    manualFiltering: true,
    onPaginationChange: (updater) => {
      const next = typeof updater === "function" ? updater(pagination) : updater
      onPaginationChange(next)
    },
    onSortingChange: (updater) => {
      const next = typeof updater === "function" ? updater(sorting) : updater
      onSortingChange(next)
    },
    onColumnVisibilityChange: setColumnVisibility,
    state: {
      pagination,
      sorting,
      columnVisibility,
    },
  })

  return (
    <div className="w-full min-w-0">
      <div className="flex flex-col justify-between gap-4 py-4 lg:flex-row lg:items-end">
        <div className="flex flex-col flex-wrap gap-3 sm:flex-row sm:items-center">
          <div className="flex flex-col gap-2">
            <span className="text-xs text-muted-foreground">{t("staff.filters.search")}</span>
            <Input
              placeholder={t("staff.filters.searchPlaceholder")}
              value={search}
              onChange={(event) => onSearchChange(event.target.value)}
              className="w-full sm:w-64"
            />
          </div>

          <div className="flex flex-col gap-2">
            <span className="text-xs text-muted-foreground">{t("staff.filters.role")}</span>
            <Select value={role} onValueChange={(value) => onRoleChange(value as StaffRole | "all")}>
              <SelectTrigger className="w-full sm:w-44">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">{t("staff.filters.allRoles")}</SelectItem>
                <SelectItem value="center_scanner">{t("roles.center_scanner")}</SelectItem>
                <SelectItem value="center_receptionist">{t("roles.center_receptionist")}</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        <div className="flex flex-wrap items-center justify-end gap-2">
          <Button onClick={onAddClick}>
            <Plus className="mr-2 h-4 w-4" />
            {t("staff.addButton")}
          </Button>
          <DataTableViewOptions table={table} t={t} />
        </div>
      </div>

      <div className="min-w-0 overflow-hidden rounded-md border">
        <Table className="min-w-full">
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead key={header.id} className="px-3 first:ps-6">
                    {header.isPlaceholder
                      ? null
                      : flexRender(header.column.columnDef.header, header.getContext())}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {table.getRowModel().rows.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow key={row.id}>
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id} className="px-3 py-2 first:ps-6">
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={columns.length} className="h-24 text-center text-muted-foreground">
                  {t("staff.empty")}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      <DataTablePagination table={table} t={t} pageSizeOptions={STAFF_PER_PAGE_OPTIONS} />
    </div>
  )
}
