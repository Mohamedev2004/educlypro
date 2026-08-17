import { useState } from "react"
import { Link, useParams } from "react-router-dom"
import { ArrowLeft, University, Plus, UserPen, Users } from "lucide-react"

import AppLayout from "@/layouts/app-layout"
import { useDirection } from "@/context/direction/direction-provider"
import { useCenterDetail } from "@/hooks/centers/use-center-detail"
import { AddOwnerDialog } from "@/components/centers/add-owner-dialog"
import { AddStaffDialog } from "@/components/centers/add-staff-dialog"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { roleLabels } from "@/constants/roles"

/**
 * Center detail page.
 *
 * Responsibility: render a single center's full summary — identity, owner,
 * and staff roster — with actions to provision an owner and staff members.
 * Layout mirrors the reference "entity detail" pattern (header + back link,
 * 2/3 content + 1/3 aside grid, Card-based info/action sections).
 * Layer: Pages
 */
export default function SuperAdminCenterDetail() {
  const { t } = useDirection()
  const { slug } = useParams<{ slug: string }>()
  const centerSlug = slug ?? ""

  const { data: center, isLoading, isError } = useCenterDetail(centerSlug)

  const [ownerDialogOpen, setOwnerDialogOpen] = useState(false)
  const [staffDialogOpen, setStaffDialogOpen] = useState(false)

  const breadcrumbs = [
    { label: t("roles.overview"), href: "/super-admin/dashboard" },
    { label: t("centers.title"), href: "/super-admin/centers" },
    { label: center?.name ?? "" },
  ]

  if (isLoading) {
    return (
      <AppLayout breadcrumbs={breadcrumbs}>
        <p className="text-sm text-muted-foreground">{t("centers.loading")}</p>
      </AppLayout>
    )
  }

  if (isError || !center) {
    return (
      <AppLayout breadcrumbs={breadcrumbs}>
        <p className="text-sm text-destructive">{t("centers.errors.notFound")}</p>
      </AppLayout>
    )
  }

  return (
    <AppLayout breadcrumbs={breadcrumbs}>
      <div className="flex h-full min-w-0 flex-1 flex-col gap-6">
        {/* Header */}
        <div className="flex flex-wrap items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <Button variant="outline" size="icon" asChild>
              <Link to="/super-admin/centers">
                <ArrowLeft className="h-4 w-4" />
                <span className="sr-only">{t("centers.detail.back")}</span>
              </Link>
            </Button>
            <div>
              <h1 className="text-2xl font-semibold">{center.name}</h1>
              <p className="text-sm text-muted-foreground">{center.slug}</p>
            </div>
          </div>
        </div>

        {/* Main grid */}
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          {/* 2/3 — Center info + staff */}
          <div className="space-y-6 lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <University className="h-4 w-4" />
                  {t("centers.detail.centerInfo")}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 text-sm">
                <InfoRow label={t("centers.fields.name")} value={center.name} />
                <InfoRow label={t("centers.fields.slug")} value={center.slug} />
                <InfoRow
                  label={t("centers.columns.createdAt")}
                  value={new Date(center.created_at).toLocaleDateString()}
                />
              </CardContent>
            </Card>

            <Card>
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle className="flex items-center gap-2">
                  <Users className="h-4 w-4" />
                  {t("centers.detail.staffList")}
                </CardTitle>
                <Button size="sm" variant="outline" onClick={() => setStaffDialogOpen(true)}>
                  <Plus className="h-4 w-4" />
                  {t("centers.detail.addStaffButton")}
                </Button>
              </CardHeader>
              <CardContent className="space-y-3">
                {center.staff.length === 0 ? (
                  <p className="text-sm text-muted-foreground">{t("centers.detail.noStaff")}</p>
                ) : (
                  center.staff.map((member) => (
                    <div key={member.id} className="rounded-lg border p-4 text-sm">
                      <div className="mb-2 flex items-start justify-between gap-2">
                        <div className="min-w-0 flex-1">
                          <p className="font-medium">{member.username}</p>
                          <p className="text-xs text-muted-foreground">{member.email}</p>
                        </div>
                        <Badge variant="secondary" className="shrink-0">
                          {t(roleLabels[member.role])}
                        </Badge>
                      </div>
                      <span className="text-xs text-muted-foreground">
                        {t("centers.columns.createdAt")}:{" "}
                        {new Date(member.created_at).toLocaleDateString()}
                      </span>
                    </div>
                  ))
                )}
              </CardContent>
            </Card>
          </div>

          {/* 1/3 — Owner (aside) */}
          <div className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <UserPen className="h-4 w-4" />
                  {t("centers.detail.owner")}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-4 text-sm">
                {center.owner ? (
                  <div className="space-y-3">
                    <InfoRow label={t("centers.fields.username")} value={center.owner.username} />
                    <InfoRow label={t("centers.fields.email")} value={center.owner.email} />
                    <InfoRow
                      label={t("centers.columns.createdAt")}
                      value={new Date(center.owner.created_at).toLocaleDateString()}
                    />
                  </div>
                ) : (
                  <>
                    <p className="text-muted-foreground">{t("centers.detail.noOwner")}</p>
                    <Button size="sm" className="w-full" onClick={() => setOwnerDialogOpen(true)}>
                      <Plus className="mr-1.5 h-3.5 w-3.5" />
                      {t("centers.detail.addOwnerButton")}
                    </Button>
                  </>
                )}
              </CardContent>
            </Card>
          </div>
        </div>
      </div>

      <AddOwnerDialog
        t={t}
        open={ownerDialogOpen}
        onOpenChange={setOwnerDialogOpen}
        centerSlug={centerSlug}
        onSuccess={() => {}}
      />

      <AddStaffDialog
        t={t}
        open={staffDialogOpen}
        onOpenChange={setStaffDialogOpen}
        centerSlug={centerSlug}
        onSuccess={() => {}}
      />
    </AppLayout>
  )
}

function InfoRow({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex items-start justify-between gap-4">
      <span className="min-w-[100px] shrink-0 text-muted-foreground">{label}</span>
      <span className="text-right">{value}</span>
    </div>
  )
}
