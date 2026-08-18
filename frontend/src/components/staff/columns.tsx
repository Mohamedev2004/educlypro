import type { ColumnDef } from "@tanstack/react-table"
import { MoreHorizontal, SquarePen, Trash2 } from "lucide-react"

import { DataTableColumnHeader } from "@/components/data-table-column-header"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { roleLabels } from "@/constants/roles"
import type { Staff } from "@/api/types/staff.types"

/**
 * Column definitions for the staff data table.
 *
 * Responsibility: describe how each Staff field renders as a table column.
 * Layer: Components (domain)
 */

interface ColumnHandlers {
  t: (key: string) => string
  onEdit: (staff: Staff) => void
  onDelete: (staff: Staff) => void
}

export const createColumns = ({
  t,
  onEdit,
  onDelete,
}: ColumnHandlers): ColumnDef<Staff>[] => [
  {
    accessorKey: "username",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("staff.columns.username")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <span className="font-medium">{row.original.username}</span>
    ),
  },
  {
    accessorKey: "email",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("staff.columns.email")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <span className="text-muted-foreground">{row.original.email}</span>
    ),
  },
  {
    accessorKey: "role",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("staff.columns.role")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <Badge variant="secondary">{t(roleLabels[row.original.role])}</Badge>
    ),
  },
  {
    accessorKey: "sub_center_name",
    enableSorting: false,
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("staff.columns.subCenter")}
        t={t}
      />
    ),
    cell: ({ row }) => (
      <span className="text-muted-foreground">
        {row.original.sub_center_name}
      </span>
    ),
  },
  {
    accessorKey: "created_at",
    header: ({ column }) => (
      <DataTableColumnHeader
        column={column}
        title={t("staff.columns.createdAt")}
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
      const staff = row.original

      return (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="h-8 w-8 p-0">
              <span className="sr-only">{t("staff.columns.openMenu")}</span>
              <MoreHorizontal className="h-4 w-4" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end">
            <DropdownMenuLabel>{t("staff.columns.actions")}</DropdownMenuLabel>
            <DropdownMenuItem onClick={() => onEdit(staff)}>
              <SquarePen className="mr-2 h-4 w-4" />
              {t("common.edit")}
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              variant="destructive"
              onClick={() => onDelete(staff)}
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
