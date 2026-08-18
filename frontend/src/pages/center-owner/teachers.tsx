import { useMemo, useState } from "react"
import type { PaginationState, SortingState } from "@tanstack/react-table"

import AppLayout from "@/layouts/app-layout"
import { useDirection } from "@/context/direction/direction-provider"
import { useTeachersList } from "@/hooks/teachers/use-teachers-list"
import { useDeleteTeacher } from "@/hooks/teachers/use-delete-teacher"
import { createColumns } from "@/components/teachers/columns"
import { DataTable } from "@/components/teachers/data-table"
import { TeacherFormDialog } from "@/components/teachers/teacher-form-dialog"
import { DeleteTeacherDialog } from "@/components/teachers/delete-teacher-dialog"
import type { Teacher, TeachersListParams } from "@/api/types/teachers.types"

const TEACHERS_SORT_COLUMNS = new Set<TeachersListParams["sort"]>([
  "full_name",
  "email",
  "created_at",
])

export default function CenterOwnerTeachers() {
  const { t } = useDirection()

  const [pagination, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  })
  const [sorting, setSorting] = useState<SortingState>([
    { id: "created_at", desc: true },
  ])
  const [search, setSearch] = useState("")

  const [formOpen, setFormOpen] = useState(false)
  const [editingTeacher, setEditingTeacher] = useState<Teacher | null>(null)
  const [teacherPendingDelete, setTeacherPendingDelete] =
    useState<Teacher | null>(null)

  const sortColumn = sorting[0]?.id
  const listParams: TeachersListParams = {
    page: pagination.pageIndex + 1,
    perPage: pagination.pageSize,
    q: search || undefined,
    sort: TEACHERS_SORT_COLUMNS.has(sortColumn as TeachersListParams["sort"])
      ? (sortColumn as TeachersListParams["sort"])
      : "created_at",
    direction: sorting[0]?.desc ? "desc" : "asc",
  }

  const { data, isFetching } = useTeachersList(listParams)
  const deleteTeacher = useDeleteTeacher(t)

  const items = data?.items ?? []
  const pageCount = data?.pagination.total_pages ?? 1

  const columns = useMemo(
    () =>
      createColumns({
        t,
        onEdit: (teacher) => {
          setEditingTeacher(teacher)
          setFormOpen(true)
        },
        onDelete: (teacher) => setTeacherPendingDelete(teacher),
      }),
    [t]
  )

  return (
    <AppLayout
      breadcrumbs={[
        { label: t("roles.overview"), href: "/center-owner/dashboard" },
        { label: t("teachers.title") },
      ]}
    >
      <div className="flex flex-1 flex-col gap-4">
        <div>
          <h1 className="text-lg font-semibold">{t("teachers.title")}</h1>
          <p className="text-sm text-muted-foreground">
            {t("teachers.description")}
          </p>
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
          onAddClick={() => {
            setEditingTeacher(null)
            setFormOpen(true)
          }}
        />

        {isFetching && (
          <p className="text-xs text-muted-foreground">
            {t("teachers.loading")}
          </p>
        )}
      </div>

      <TeacherFormDialog
        t={t}
        open={formOpen}
        onOpenChange={setFormOpen}
        teacher={editingTeacher}
        onSuccess={() => setEditingTeacher(null)}
      />

      <DeleteTeacherDialog
        t={t}
        teacher={teacherPendingDelete}
        onOpenChange={(open) => {
          if (!open) setTeacherPendingDelete(null)
        }}
        isDeleting={deleteTeacher.isPending}
        onConfirm={() => {
          if (!teacherPendingDelete) return
          deleteTeacher.mutate(
            { id: teacherPendingDelete.id, slug: teacherPendingDelete.slug },
            { onSuccess: () => setTeacherPendingDelete(null) }
          )
        }}
      />
    </AppLayout>
  )
}
