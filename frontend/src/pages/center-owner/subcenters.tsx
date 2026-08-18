import { useState } from "react"
import { MoreHorizontal, Plus, SquarePen, Trash2 } from "lucide-react"

import AppLayout from "@/layouts/app-layout"
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table"
import { SubCenterFormDialog } from "@/components/subcenters/subcenter-form-dialog"
import { DeleteSubCenterDialog } from "@/components/subcenters/delete-subcenter-dialog"
import { useDirection } from "@/context/direction/direction-provider"
import { useSubCentersList } from "@/hooks/subcenters/use-subcenters-list"
import { useDeleteSubCenter } from "@/hooks/subcenters/use-delete-subcenter"
import type { SubCenter } from "@/api/types/subcenters.types"

export default function CenterOwnerSubCenters() {
  const { t } = useDirection()

  const { data, isFetching } = useSubCentersList()
  const deleteSubCenter = useDeleteSubCenter(t)

  const [formOpen, setFormOpen] = useState(false)
  const [editingSubCenter, setEditingSubCenter] = useState<SubCenter | null>(
    null
  )
  const [subCenterPendingDelete, setSubCenterPendingDelete] =
    useState<SubCenter | null>(null)

  const items = data?.items ?? []

  return (
    <AppLayout
      breadcrumbs={[
        { label: t("roles.overview"), href: "/center-owner/dashboard" },
        { label: t("subcenters.title") },
      ]}
    >
      <div className="flex flex-1 flex-col gap-4">
        <div className="flex flex-wrap items-center justify-between gap-3">
          <div>
            <h1 className="text-lg font-semibold">{t("subcenters.title")}</h1>
            <p className="text-sm text-muted-foreground">
              {t("subcenters.description")}
            </p>
          </div>

          <Button
            onClick={() => {
              setEditingSubCenter(null)
              setFormOpen(true)
            }}
          >
            <Plus className="mr-2 h-4 w-4" />
            {t("subcenters.addButton")}
          </Button>
        </div>

        <div className="overflow-hidden rounded-md border">
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="px-3 first:ps-6">
                  {t("subcenters.columns.name")}
                </TableHead>
                <TableHead className="px-3">
                  {t("subcenters.columns.staffCount")}
                </TableHead>
                <TableHead className="px-3">
                  {t("subcenters.columns.createdAt")}
                </TableHead>
                <TableHead className="px-3" />
              </TableRow>
            </TableHeader>
            <TableBody>
              {items.length ? (
                items.map((subCenter) => (
                  <TableRow key={subCenter.id}>
                    <TableCell className="px-3 py-2 font-medium first:ps-6">
                      {subCenter.name}
                    </TableCell>
                    <TableCell className="px-3 py-2">
                      <Badge variant="secondary">{subCenter.staff_count}</Badge>
                    </TableCell>
                    <TableCell className="px-3 py-2 text-muted-foreground">
                      {new Date(subCenter.created_at).toLocaleDateString()}
                    </TableCell>
                    <TableCell className="px-3 py-2">
                      <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                          <Button variant="ghost" className="h-8 w-8 p-0">
                            <span className="sr-only">
                              {t("subcenters.columns.openMenu")}
                            </span>
                            <MoreHorizontal className="h-4 w-4" />
                          </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent align="end">
                          <DropdownMenuLabel>
                            {t("subcenters.columns.actions")}
                          </DropdownMenuLabel>
                          <DropdownMenuItem
                            onClick={() => {
                              setEditingSubCenter(subCenter)
                              setFormOpen(true)
                            }}
                          >
                            <SquarePen className="mr-2 h-4 w-4" />
                            {t("common.edit")}
                          </DropdownMenuItem>
                          <DropdownMenuSeparator />
                          <DropdownMenuItem
                            variant="destructive"
                            onClick={() => setSubCenterPendingDelete(subCenter)}
                          >
                            <Trash2 className="mr-2 h-4 w-4" />
                            {t("common.delete")}
                          </DropdownMenuItem>
                        </DropdownMenuContent>
                      </DropdownMenu>
                    </TableCell>
                  </TableRow>
                ))
              ) : (
                <TableRow>
                  <TableCell
                    colSpan={4}
                    className="h-24 text-center text-muted-foreground"
                  >
                    {t("subcenters.empty")}
                  </TableCell>
                </TableRow>
              )}
            </TableBody>
          </Table>
        </div>

        {isFetching && (
          <p className="text-xs text-muted-foreground">
            {t("subcenters.loading")}
          </p>
        )}
      </div>

      <SubCenterFormDialog
        t={t}
        open={formOpen}
        onOpenChange={setFormOpen}
        subCenter={editingSubCenter}
        onSuccess={() => setEditingSubCenter(null)}
      />

      <DeleteSubCenterDialog
        t={t}
        subCenter={subCenterPendingDelete}
        onOpenChange={(open) => {
          if (!open) setSubCenterPendingDelete(null)
        }}
        isDeleting={deleteSubCenter.isPending}
        onConfirm={() => {
          if (!subCenterPendingDelete) return
          deleteSubCenter.mutate(subCenterPendingDelete.id, {
            onSuccess: () => setSubCenterPendingDelete(null),
          })
        }}
      />
    </AppLayout>
  )
}
