import type { ColumnDef } from "@tanstack/react-table"
import { Link } from "react-router-dom"

import { DataTableColumnHeader } from "@/components/data-table-column-header"
import { Badge } from "@/components/ui/badge"
import type { Center } from "@/api/types/centers.types"

/**
 * Column definitions for the centers data table.
 *
 * Responsibility: describe how each Center field renders as a table column.
 * Layer: Components (domain)
 */

interface ColumnHandlers {
  t: (key: string) => string
}

export const createColumns = ({ t }: ColumnHandlers): ColumnDef<Center>[] => [
  {
    accessorKey: "name",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title={t("centers.columns.name")} t={t} />
    ),
    cell: ({ row }) => (
      <Link
        to={`/super-admin/centers/${row.original.slug}`}
        className="font-medium hover:underline"
      >
        {row.original.name}
      </Link>
    ),
  },
  {
    accessorKey: "slug",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title={t("centers.columns.slug")} t={t} />
    ),
    cell: ({ row }) => (
      <span className="text-muted-foreground">{row.original.slug}</span>
    ),
  },
  {
    id: "owner",
    accessorFn: (row) => row.owner_username ?? "",
    enableSorting: false,
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title={t("centers.columns.owner")} t={t} />
    ),
    cell: ({ row }) => {
      const { owner_username, owner_email } = row.original
      if (!owner_username) {
        return <span className="text-muted-foreground">{t("centers.noOwner")}</span>
      }
      return (
        <div className="flex flex-col">
          <span className="font-medium">{owner_username}</span>
          <span className="text-xs text-muted-foreground">{owner_email}</span>
        </div>
      )
    },
  },
  {
    accessorKey: "staff_count",
    enableSorting: false,
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title={t("centers.columns.staffCount")} t={t} />
    ),
    cell: ({ row }) => <Badge variant="secondary">{row.original.staff_count}</Badge>,
  },
  {
    accessorKey: "created_at",
    header: ({ column }) => (
      <DataTableColumnHeader column={column} title={t("centers.columns.createdAt")} t={t} />
    ),
    cell: ({ row }) => (
      <span className="text-muted-foreground">
        {new Date(row.original.created_at).toLocaleDateString()}
      </span>
    ),
  },
]
