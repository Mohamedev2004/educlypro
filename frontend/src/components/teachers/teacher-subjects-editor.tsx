import { useState } from "react"
import { BookOpen, LoaderCircle, SquarePen } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { SubjectMultiSelect } from "@/components/teachers/subject-multi-select"
import { useUpdateTeacherSubjects } from "@/hooks/teachers/use-update-teacher-subjects"
import type { TeacherSubject } from "@/api/types/teachers.types"

/**
 * Teacher detail page card for viewing and editing a teacher's subject
 * assignment in place — no separate dialog, since this is the only place
 * subjects are managed (the create/edit teacher form only covers identity
 * fields).
 *
 * Responsibility: own local edit-mode + selection state, submit via
 * `useUpdateTeacherSubjects`.
 * Layer: Components (domain)
 */

interface TeacherSubjectsEditorProps {
  t: (key: string) => string
  teacherId: number
  subjects: TeacherSubject[]
}

export function TeacherSubjectsEditor({
  t,
  teacherId,
  subjects,
}: TeacherSubjectsEditorProps) {
  const [editing, setEditing] = useState(false)
  const [selectedIds, setSelectedIds] = useState<number[]>([])
  const updateSubjects = useUpdateTeacherSubjects(t)

  function startEditing() {
    setSelectedIds(subjects.map((s) => s.id))
    setEditing(true)
  }

  async function handleSave() {
    await updateSubjects.mutateAsync({ id: teacherId, subjectIds: selectedIds })
    setEditing(false)
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <BookOpen className="h-4 w-4" />
          {t("teachers.detail.subjectsList")}
        </CardTitle>
        {!editing && (
          <CardAction>
            <Button variant="ghost" size="icon" onClick={startEditing}>
              <SquarePen className="h-4 w-4" />
              <span className="sr-only">{t("common.edit")}</span>
            </Button>
          </CardAction>
        )}
      </CardHeader>
      <CardContent className="space-y-3">
        {editing ? (
          <>
            <SubjectMultiSelect
              t={t}
              selectedIds={selectedIds}
              onChange={setSelectedIds}
              disabled={updateSubjects.isPending}
            />
            <div className="flex justify-end gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setEditing(false)}
                disabled={updateSubjects.isPending}
              >
                {t("common.cancel")}
              </Button>
              <Button
                size="sm"
                onClick={handleSave}
                disabled={updateSubjects.isPending}
              >
                {updateSubjects.isPending && (
                  <LoaderCircle className="mr-2 h-4 w-4 animate-spin" />
                )}
                {t("common.save")}
              </Button>
            </div>
          </>
        ) : subjects.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            {t("teachers.noSubjects")}
          </p>
        ) : (
          <div className="flex flex-wrap gap-1.5">
            {subjects.map((subject) => (
              <Badge key={subject.id} variant="secondary">
                {subject.name}
              </Badge>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
