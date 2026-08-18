import { useState } from "react"
import { Link, useParams } from "react-router-dom"
import { ArrowLeft, ArrowRight, IdCard, SquarePen } from "lucide-react"

import AppLayout from "@/layouts/app-layout"
import { useDirection } from "@/context/direction/direction-provider"
import { useTeacherDetail } from "@/hooks/teachers/use-teacher-detail"
import { TeacherFormDialog } from "@/components/teachers/teacher-form-dialog"
import { TeacherSubjectsEditor } from "@/components/teachers/teacher-subjects-editor"
import { TeacherClassesEditor } from "@/components/teachers/teacher-classes-editor"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"

/**
 * Teacher detail page.
 *
 * Responsibility: render a single teacher's full summary — identity,
 * subjects, and classes — with an edit action for identity fields, and
 * in-place editors for subjects/classes. Layout mirrors the reference
 * "entity detail" pattern used by the center detail page (header + back
 * link, 2/3 content + 1/3 aside grid, Card-based info sections). Looked up
 * by slug (not id) so links to a teacher stay stable across renames.
 * Layer: Pages
 */
export default function CenterOwnerTeacherDetail() {
  const { t, direction } = useDirection()
  const { slug } = useParams<{ slug: string }>()

  const { data: teacher, isLoading, isError } = useTeacherDetail(slug ?? "")
  const [editOpen, setEditOpen] = useState(false)

  const breadcrumbs = [
    { label: t("roles.overview"), href: "/center-owner/dashboard" },
    { label: t("teachers.title"), href: "/center-owner/teachers" },
    { label: teacher?.full_name ?? "" },
  ]

  if (isLoading) {
    return (
      <AppLayout breadcrumbs={breadcrumbs}>
        <p className="text-sm text-muted-foreground">{t("teachers.loading")}</p>
      </AppLayout>
    )
  }

  if (isError || !teacher) {
    return (
      <AppLayout breadcrumbs={breadcrumbs}>
        <p className="text-sm text-destructive">
          {t("teachers.errors.notFound")}
        </p>
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
              <Link to="/center-owner/teachers">
                {direction === "rtl" ? (
                  <ArrowRight className="h-4 w-4" />
                ) : (
                  <ArrowLeft className="h-4 w-4" />
                )}
                <span className="sr-only">{t("teachers.detail.back")}</span>
              </Link>
            </Button>
            <div>
              <h1 className="text-2xl font-semibold">{teacher.full_name}</h1>
              <p className="text-sm text-muted-foreground">{teacher.email}</p>
            </div>
          </div>

          <Button variant="outline" onClick={() => setEditOpen(true)}>
            <SquarePen className="h-4 w-4" />
            {t("common.edit")}
          </Button>
        </div>

        {/* Main grid */}
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          {/* 2/3 — Teacher info + subjects */}
          <div className="space-y-6 lg:col-span-2">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <IdCard className="h-4 w-4" />
                  {t("teachers.detail.teacherInfo")}
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3 text-sm">
                <InfoRow
                  label={t("teachers.fields.fullName")}
                  value={teacher.full_name}
                />
                <InfoRow
                  label={t("teachers.fields.email")}
                  value={teacher.email}
                />
                <InfoRow
                  label={t("teachers.fields.phone")}
                  value={teacher.phone}
                />
                <InfoRow
                  label={t("teachers.columns.createdAt")}
                  value={new Date(teacher.created_at).toLocaleDateString()}
                />
              </CardContent>
            </Card>

            <TeacherSubjectsEditor
              t={t}
              teacherId={teacher.id}
              subjects={teacher.subjects}
            />
          </div>

          {/* 1/3 — Classes (aside) */}
          <div className="space-y-6">
            <TeacherClassesEditor
              t={t}
              teacherId={teacher.id}
              classes={teacher.classes}
            />
          </div>
        </div>
      </div>

      <TeacherFormDialog
        t={t}
        open={editOpen}
        onOpenChange={setEditOpen}
        teacher={teacher}
        onSuccess={() => {}}
      />
    </AppLayout>
  )
}

function InfoRow({ label, value }: { label: string; value: React.ReactNode }) {
  return (
    <div className="flex items-start justify-between gap-4">
      <span className="min-w-[100px] shrink-0 text-muted-foreground">
        {label}
      </span>
      <span className="text-right">{value}</span>
    </div>
  )
}
