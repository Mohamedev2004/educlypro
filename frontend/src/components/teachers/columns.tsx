import type { ColumnDef } from "@tanstack/react-table"
import { Link } from "react-router-dom"
import { MoreHorizontal, SquarePen, Trash2 } from "lucide-react"

import { DataTableColumnHeader } from "@/components/data-table-column-header"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import type { Teacher } from "@/api/types/teachers.types"

/**
 * Column definitions for the teachers data table.
 *
 * Responsibility: describe how each Teacher field renders as a table
 * column. Kept to compact overview fields — the full subjects/classes
 * lists (with add/edit) live on the teacher detail page
 * (`/center-owner/teachers/:slug`, linked from the name cell here); this
 * table only shows their counts.
 * Layer: Components (domain)
 */

interface ColumnHandlers {
  t: (key: string) => string
  onEdit: (teacher: Teacher) => void
  onDelete: (teacher: Teacher) => void
}

export const createColumns = ({
  t,
  onEdit,
  onDelete,
}: ColumnHandlers): ColumnDef<Teacher>[] => [
  {
    accessorKey: "full_name",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("teachers.columns.fullName")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <Link
        to={`/center-owner/teachers/${row.original.slug}`}
        className="font-medium hover:underline"
      >
        {row.original.full_name}
      </Link>
    ),
  },
  {
    accessorKey: "email",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("teachers.columns.email")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <span className="text-muted-foreground">{row.original.email}</span>
    ),
  },
  {
    accessorKey: "phone",
    enableSorting: false,
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("teachers.columns.phone")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <span className="text-muted-foreground">{row.original.phone}</span>
    ),
  },
  {
    id: "subjectsCount",
    accessorFn: (row) => row.subjects.length,
    enableSorting: false,
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("teachers.columns.subjects")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <span className="text-muted-foreground">
        {row.original.subjects.length}
      </span>
    ),
  },
  {
    id: "classesCount",
    accessorFn: (row) => row.classes.length,
    enableSorting: false,
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("teachers.columns.classes")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <span className="text-muted-foreground">
        {row.original.classes.length}
      </span>
    ),
  },
  {
    accessorKey: "created_at",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("teachers.columns.createdAt")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <span className="text-muted-foreground">
        {new Date(row.original.created_at).toLocaleDateString()}
      </span>
    ),
  },
  {
    id: "actions",
    enableHiding: false,
    enableSorting: false,
    cell: ({ row }) => {
      const teacher = row.original

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">{t("teachers.columns.openMenu")}</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>
              {t("teachers.columns.actions")}
            </DropdownMenuLabel>
            <DropdownMenuItem onClick={() => onEdit(teacher)}>
              <SquarePen className="mr-2 h-4 w-4" />
              {t("common.edit")}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              variant="destructive"
              onClick={() => onDelete(teacher)}
            >
              <Trash2 className="mr-2 h-4 w-4" />
              {t("common.delete")}
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      )
    },
  },
]
