import { useMemo, useState } from "react"
import type { PaginationState, SortingState } from "@tanstack/react-table"

import AppLayout from "@/layouts/app-layout"
import { useDirection } from "@/context/direction/direction-provider"
import { useCentersList } from "@/hooks/centers/use-centers-list"
import { createColumns } from "@/components/centers/columns"
import { DataTable } from "@/components/centers/data-table"
import { CenterFormDialog } from "@/components/centers/center-form-dialog"
import type { CentersListParams } from "@/api/types/centers.types"

const CENTERS_SORT_COLUMNS = new Set<CentersListParams["sort"]>([
  "name",
  "slug",
  "created_at",
])

export default function SuperAdminCenters() {
  const { t } = useDirection()

  const [pagination, setPagination] = useState<PaginationState>({ pageIndex: 0, pageSize: 10 })
  const [sorting, setSorting] = useState<SortingState>([{ id: "created_at", desc: true }])
  const [search, setSearch] = useState("")
  const [formOpen, setFormOpen] = useState(false)

  const sortColumn = sorting[0]?.id
  const listParams: CentersListParams = {
    page: pagination.pageIndex + 1,
    perPage: pagination.pageSize,
    q: search || undefined,
    sort: CENTERS_SORT_COLUMNS.has(sortColumn as CentersListParams["sort"])
      ? (sortColumn as CentersListParams["sort"])
      : "created_at",
    direction: sorting[0]?.desc ? "desc" : "asc",
  }

  const { data, isFetching } = useCentersList(listParams)

  const items = data?.items ?? []
  const pageCount = data?.pagination.total_pages ?? 1

  const columns = useMemo(() => createColumns({ t }), [t])

  return (
    <AppLayout
      breadcrumbs={[
        { label: t("roles.overview"), href: "/super-admin/dashboard" },
        { label: t("centers.title") },
      ]}
    >
      <div className="flex flex-1 flex-col gap-4">
        <div>
          <h1 className="text-2xl font-semibold text-foreground">{t("centers.title")}</h1>
          <p className="text-sm text-muted-foreground max-w-3xl text-justify">{t("centers.description")}</p>
        </div>

        <DataTable
          t={t}
          columns={columns}
          data={items}
          pagination={pagination}
          pageCount={pageCount}
          onPaginationChange={setPagination}
          sorting={sorting}
          onSortingChange={setSorting}
          search={search}
          onSearchChange={(value) => {
            setSearch(value)
            setPagination((prev) => ({ ...prev, pageIndex: 0 }))
          }}
          onAddClick={() => setFormOpen(true)}
        />

        {isFetching && (
          <p className="text-xs text-muted-foreground">{t("centers.loading")}</p>
        )}
      </div>

      <CenterFormDialog
        t={t}
        open={formOpen}
        onOpenChange={setFormOpen}
        onSuccess={() => {}}
      />
    </AppLayout>
  )
}
