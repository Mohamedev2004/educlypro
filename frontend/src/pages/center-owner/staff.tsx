import { useMemo, useState } from "react"
import type { PaginationState, SortingState } from "@tanstack/react-table"

import AppLayout from "@/layouts/app-layout"
import { useDirection } from "@/context/direction/direction-provider"
import { useAuth } from "@/context/auth/auth-context"
import { useStaffList } from "@/hooks/staff/use-staff-list"
import { useDeleteStaff } from "@/hooks/staff/use-delete-staff"
import { createColumns } from "@/components/staff/columns"
import { DataTable } from "@/components/staff/data-table"
import { StaffFormDialog } from "@/components/staff/staff-form-dialog"
import { DeleteStaffDialog } from "@/components/staff/delete-staff-dialog"
import type { Staff, StaffListParams, StaffRole } from "@/api/types/staff.types"

const STAFF_SORT_COLUMNS = new Set<StaffListParams["sort"]>([
  "username",
  "email",
  "role",
  "created_at",
])

export default function CenterOwnerStaff() {
  const { t } = useDirection()
  const { user } = useAuth()
  const centerId = user?.center?.id

  const [pagination, setPagination] = useState<PaginationState>({ pageIndex: 0, pageSize: 10 })
  const [sorting, setSorting] = useState<SortingState>([{ id: "created_at", desc: true }])
  const [search, setSearch] = useState("")
  const [role, setRole] = useState<StaffRole | "all">("all")

  const [formOpen, setFormOpen] = useState(false)
  const [editingStaff, setEditingStaff] = useState<Staff | null>(null)
  const [staffPendingDelete, setStaffPendingDelete] = useState<Staff | null>(null)

  const sortColumn = sorting[0]?.id
  const listParams: StaffListParams = {
    page: pagination.pageIndex + 1,
    perPage: pagination.pageSize,
    role: role === "all" ? undefined : role,
    q: search || undefined,
    sort: STAFF_SORT_COLUMNS.has(sortColumn as StaffListParams["sort"])
      ? (sortColumn as StaffListParams["sort"])
      : "created_at",
    direction: sorting[0]?.desc ? "desc" : "asc",
  }

  const { data, isFetching } = useStaffList(listParams)
  const deleteStaff = useDeleteStaff(t)

  const items = data?.items ?? []
  const pageCount = data?.pagination.total_pages ?? 1

  const columns = useMemo(
    () =>
      createColumns({
        t,
        onEdit: (staff) => {
          setEditingStaff(staff)
          setFormOpen(true)
        },
        onDelete: (staff) => setStaffPendingDelete(staff),
      }),
    [t]
  )

  if (!centerId) {
    return null
  }

  return (
    <AppLayout
      breadcrumbs={[
        { label: t("roles.overview"), href: "/center-owner/dashboard" },
        { label: t("staff.title") },
      ]}
    >
      <div className="flex flex-1 flex-col gap-4">
        <div>
          <h1 className="text-lg font-semibold">{t("staff.title")}</h1>
          <p className="text-sm text-muted-foreground">{t("staff.description")}</p>
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
          role={role}
          onRoleChange={(value) => {
            setRole(value)
            setPagination((prev) => ({ ...prev, pageIndex: 0 }))
          }}
          onAddClick={() => {
            setEditingStaff(null)
            setFormOpen(true)
          }}
        />

        {isFetching && (
          <p className="text-xs text-muted-foreground">{t("staff.loading")}</p>
        )}
      </div>

      <StaffFormDialog
        t={t}
        open={formOpen}
        onOpenChange={setFormOpen}
        centerId={centerId}
        staff={editingStaff}
        onSuccess={() => setEditingStaff(null)}
      />

      <DeleteStaffDialog
        t={t}
        staff={staffPendingDelete}
        onOpenChange={(open) => {
          if (!open) setStaffPendingDelete(null)
        }}
        isDeleting={deleteStaff.isPending}
        onConfirm={() => {
          if (!staffPendingDelete) return
          deleteStaff.mutate(staffPendingDelete.id, {
            onSuccess: () => setStaffPendingDelete(null),
          })
        }}
      />
    </AppLayout>
  )
}
