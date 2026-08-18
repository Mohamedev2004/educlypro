import { useState } from "react"
import { useNavigate } from "react-router-dom"
import { useQueryClient } from "@tanstack/react-query"
import { ArrowRightIcon, LayersIcon, LoaderCircle } from "lucide-react"

import AuthLayout from "@/layouts/auth-layout"
import { Button } from "@/components/ui/button"
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import { Accordion } from "@/components/ui/accordion"
import { AddEntityPopover } from "@/components/onboarding/add-entity-popover"
import { GradeItem } from "@/components/onboarding/grade-item"
import { useDirection } from "@/context/direction/direction-provider"
import { useAuth } from "@/context/auth/auth-context"
import { useAcademicTree } from "@/hooks/academic/use-academic-tree"
import { useAddGrade } from "@/hooks/academic/use-add-grade"
import { useRemoveGrade } from "@/hooks/academic/use-remove-grade"
import { useAddMajor } from "@/hooks/academic/use-add-major"
import { useRemoveMajor } from "@/hooks/academic/use-remove-major"
import { useAddSubject } from "@/hooks/academic/use-add-subject"
import { useRemoveSubject } from "@/hooks/academic/use-remove-subject"
import { GRADE_SUGGESTIONS } from "@/constants/academic"
import { getDashboardPath } from "@/utils/navigation-utils"

/**
 * Center owner onboarding — the required first step before the dashboard
 * unlocks. Builds the center's grade -> major -> subject structure.
 *
 * Responsibility: compose the wizard from the academic-* hooks and
 * onboarding/* presentational components; own only the small bits of local
 * UI state (which row is mid-delete) the hooks don't already track.
 * Layer: Pages
 */
export default function Onboarding() {
  const { t } = useDirection()
  const { user } = useAuth()
  const navigate = useNavigate()
  const queryClient = useQueryClient()

  const { data, isLoading } = useAcademicTree()
  const grades = data?.grades ?? []

  // Mirrors the backend's own definition of "setup complete" (see
  // auth.Repository.CenterAcademicSetupComplete) — every grade needs at
  // least one major, and every major needs at least one subject, before
  // the dashboard unlocks.
  const hasIncompleteGrade = grades.some((grade) => grade.majors.length === 0)
  const hasIncompleteMajor = grades
    .flatMap((grade) => grade.majors)
    .some((major) => major.subjects.length === 0)
  const isStructureComplete =
    grades.length > 0 && !hasIncompleteGrade && !hasIncompleteMajor

  const addGrade = useAddGrade(t)
  const removeGrade = useRemoveGrade(t)
  const addMajor = useAddMajor(t)
  const removeMajor = useRemoveMajor(t)
  const addSubject = useAddSubject(t)
  const removeSubject = useRemoveSubject(t)

  const [removingGradeId, setRemovingGradeId] = useState<number | null>(null)
  const [removingMajorId, setRemovingMajorId] = useState<number | null>(null)
  const [removingSubjectId, setRemovingSubjectId] = useState<number | null>(
    null
  )
  const [isContinuing, setIsContinuing] = useState(false)

  async function handleContinue() {
    if (!user || !isStructureComplete) return
    setIsContinuing(true)
    // Refetch "me" first so the dashboard guard sees the just-updated
    // has_grades flag instead of racing the background invalidation from
    // the last add-grade mutation.
    await queryClient.invalidateQueries({ queryKey: ["me"] })
    navigate(getDashboardPath(user.role), { replace: true })
  }

  return (
    <AuthLayout contentClassName="max-w-xl">
      <div className="flex justify-center">
        <Card className="w-full">
          <CardHeader>
            <CardTitle>{t("onboarding.title")}</CardTitle>
            <CardDescription>{t("onboarding.description")}</CardDescription>
          </CardHeader>

          <CardContent className="flex flex-col gap-4">
            <div className="flex flex-wrap items-center justify-between gap-3">
              <AddEntityPopover
                t={t}
                triggerLabel={t("onboarding.addGrade")}
                placeholder={t("onboarding.addGradePlaceholder")}
                suggestions={GRADE_SUGGESTIONS}
                existingNames={grades.map((grade) => grade.name)}
                onAdd={(name) => addGrade.mutate({ name })}
                disabled={addGrade.isPending}
                triggerVariant="default"
              />

              <Button
                variant="outline"
                disabled={!isStructureComplete || isContinuing}
                onClick={handleContinue}
              >
                {isContinuing ? (
                  <LoaderCircle className="animate-spin" />
                ) : (
                  <>
                    {t("onboarding.continueButton")}
                    <ArrowRightIcon />
                  </>
                )}
              </Button>
            </div>

            {grades.length > 0 && !isStructureComplete && (
              <p className="text-xs text-muted-foreground">
                {t("onboarding.incompleteStructureHint")}
              </p>
            )}

            {isLoading ? (
              <p className="text-sm text-muted-foreground">
                {t("onboarding.loading")}
              </p>
            ) : grades.length > 0 ? (
              <Accordion multiple defaultValue={[String(grades[0]?.id ?? "")]}>
                {grades.map((grade) => (
                  <GradeItem
                    t={t}
                    key={grade.id}
                    grade={grade}
                    onAddMajor={(name) =>
                      addMajor.mutate({ gradeId: grade.id, payload: { name } })
                    }
                    onRemoveMajor={(majorId) => {
                      setRemovingMajorId(majorId)
                      removeMajor.mutate(majorId, {
                        onSettled: () => setRemovingMajorId(null),
                      })
                    }}
                    onAddSubject={(majorId, name) =>
                      addSubject.mutate({ majorId, payload: { name } })
                    }
                    onRemoveSubject={(subjectId) => {
                      setRemovingSubjectId(subjectId)
                      removeSubject.mutate(subjectId, {
                        onSettled: () => setRemovingSubjectId(null),
                      })
                    }}
                    onRemoveGrade={() => {
                      setRemovingGradeId(grade.id)
                      removeGrade.mutate(grade.id, {
                        onSettled: () => setRemovingGradeId(null),
                      })
                    }}
                    isAddingMajor={addMajor.isPending}
                    isAddingSubject={addSubject.isPending}
                    isRemovingGrade={removingGradeId === grade.id}
                    removingMajorId={removingMajorId}
                    removingSubjectId={removingSubjectId}
                  />
                ))}
              </Accordion>
            ) : (
              <div className="flex flex-col items-center gap-4 rounded-2xl border border-dashed py-16 text-center">
                <div className="flex size-12 items-center justify-center rounded-2xl bg-primary/10 text-primary">
                  <LayersIcon className="size-6" />
                </div>
                <div className="space-y-1">
                  <p className="font-medium">{t("onboarding.empty.title")}</p>
                  <p className="max-w-sm text-sm text-muted-foreground">
                    {t("onboarding.empty.description")}
                  </p>
                </div>
                <div className="flex flex-wrap justify-center gap-2">
                  {GRADE_SUGGESTIONS.slice(0, 4).map((suggestion) => (
                    <Button
                      key={suggestion}
                      variant="secondary"
                      onClick={() => addGrade.mutate({ name: suggestion })}
                      disabled={addGrade.isPending}
                    >
                      {suggestion}
                    </Button>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </AuthLayout>
  )
}
