import { Trash2Icon } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
  Accordion,
  AccordionItem,
  AccordionTrigger,
  AccordionContent,
} from "@/components/ui/accordion"
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog"
import { AddEntityPopover } from "@/components/onboarding/add-entity-popover"
import { MajorItem } from "@/components/onboarding/major-item"
import { MAJOR_SUGGESTIONS } from "@/constants/academic"
import type { Grade } from "@/api/types/academic.types"

type GradeItemProps = {
  t: (key: string) => string
  grade: Grade
  onAddMajor: (name: string) => void
  onRemoveMajor: (majorId: number) => void
  onAddSubject: (majorId: number, name: string) => void
  onRemoveSubject: (subjectId: number) => void
  onRemoveGrade: () => void
  isAddingMajor?: boolean
  isAddingSubject?: boolean
  isRemovingGrade?: boolean
  removingMajorId?: number | null
  removingSubjectId?: number | null
}

export function GradeItem({
  t,
  grade,
  onAddMajor,
  onRemoveMajor,
  onAddSubject,
  onRemoveSubject,
  onRemoveGrade,
  isAddingMajor,
  isAddingSubject,
  isRemovingGrade,
  removingMajorId,
  removingSubjectId,
}: GradeItemProps) {
  return (
    <AccordionItem value={String(grade.id)}>
      <div className="flex items-center gap-1 pe-3">
        <AccordionTrigger>
          <span className="flex items-center gap-2">
            {grade.name}
            <Badge variant="secondary">{grade.majors.length}</Badge>
          </span>
        </AccordionTrigger>

        <AlertDialog>
          <AlertDialogTrigger asChild>
            <Button
              type="button"
              variant="ghost"
              size="icon-sm"
              className="text-muted-foreground hover:text-destructive"
              disabled={isRemovingGrade}
            >
              <Trash2Icon />
              <span className="sr-only">
                {t("onboarding.removeLabel")} {grade.name}
              </span>
            </Button>
          </AlertDialogTrigger>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>
                {t("onboarding.removeGrade.title").replace(
                  "{name}",
                  grade.name
                )}
              </AlertDialogTitle>
              <AlertDialogDescription>
                {t("onboarding.removeGrade.description")}
              </AlertDialogDescription>
            </AlertDialogHeader>
            <AlertDialogFooter>
              <AlertDialogCancel disabled={isRemovingGrade}>
                {t("common.cancel")}
              </AlertDialogCancel>
              <AlertDialogAction
                className="bg-destructive text-white hover:bg-destructive/90"
                onClick={onRemoveGrade}
                disabled={isRemovingGrade}
              >
                {t("onboarding.removeGrade.confirm")}
              </AlertDialogAction>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </div>

      <AccordionContent>
        <div className="flex flex-col gap-3">
          <div>
            <AddEntityPopover
              t={t}
              triggerLabel={t("onboarding.addMajor")}
              placeholder={t("onboarding.addMajorPlaceholder")}
              suggestions={MAJOR_SUGGESTIONS}
              existingNames={grade.majors.map((major) => major.name)}
              onAdd={onAddMajor}
              disabled={isAddingMajor}
            />
          </div>

          {grade.majors.length > 0 ? (
            <Accordion multiple className="rounded-xl">
              {grade.majors.map((major) => (
                <MajorItem
                  t={t}
                  key={major.id}
                  major={major}
                  onAddSubject={(name) => onAddSubject(major.id, name)}
                  onRemoveSubject={onRemoveSubject}
                  onRemoveMajor={() => onRemoveMajor(major.id)}
                  isAddingSubject={isAddingSubject}
                  isRemoving={removingMajorId === major.id}
                  removingSubjectId={removingSubjectId}
                />
              ))}
            </Accordion>
          ) : (
            <p className="text-sm text-muted-foreground">
              {t("onboarding.noMajors")}
            </p>
          )}
        </div>
      </AccordionContent>
    </AccordionItem>
  )
}
