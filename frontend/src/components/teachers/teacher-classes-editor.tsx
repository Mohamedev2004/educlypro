import { useState } from "react"
import { LoaderCircle, Layers, SquarePen } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { ClassMultiSelect } from "@/components/teachers/class-multi-select"
import { useUpdateTeacherClasses } from "@/hooks/teachers/use-update-teacher-classes"
import type { TeacherClass } from "@/api/types/teachers.types"

/**
 * Teacher detail page card for viewing and editing a teacher's class
 * assignment in place — see TeacherSubjectsEditor for the rationale (no
 * separate dialog, this is the only place classes are managed).
 *
 * Responsibility: own local edit-mode + selection state, submit via
 * `useUpdateTeacherClasses`.
 * Layer: Components (domain)
 */

interface TeacherClassesEditorProps {
  t: (key: string) => string
  teacherId: number
  classes: TeacherClass[]
}

export function TeacherClassesEditor({
  t,
  teacherId,
  classes,
}: TeacherClassesEditorProps) {
  const [editing, setEditing] = useState(false)
  const [selectedIds, setSelectedIds] = useState<number[]>([])
  const updateClasses = useUpdateTeacherClasses(t)

  function startEditing() {
    setSelectedIds(classes.map((c) => c.id))
    setEditing(true)
  }

  async function handleSave() {
    await updateClasses.mutateAsync({ id: teacherId, classIds: selectedIds })
    setEditing(false)
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          <Layers className="h-4 w-4" />
          {t("teachers.detail.classesList")}
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
            <ClassMultiSelect
              t={t}
              selectedIds={selectedIds}
              onChange={setSelectedIds}
              disabled={updateClasses.isPending}
            />
            <div className="flex justify-end gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setEditing(false)}
                disabled={updateClasses.isPending}
              >
                {t("common.cancel")}
              </Button>
              <Button
                size="sm"
                onClick={handleSave}
                disabled={updateClasses.isPending}
              >
                {updateClasses.isPending && (
                  <LoaderCircle className="mr-2 h-4 w-4 animate-spin" />
                )}
                {t("common.save")}
              </Button>
            </div>
          </>
        ) : classes.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            {t("teachers.noClasses")}
          </p>
        ) : (
          <div className="flex flex-wrap gap-1.5">
            {classes.map((cls) => (
              <Badge key={cls.id} variant="outline">
                {cls.name}
              </Badge>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  )
}
