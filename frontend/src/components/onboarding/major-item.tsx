import { Trash2Icon } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import {
  AccordionItem,
  AccordionTrigger,
  AccordionContent,
} from "@/components/ui/accordion"
import { SubjectChips } from "@/components/onboarding/subject-chips"
import type { Major } from "@/api/types/academic.types"

type MajorItemProps = {
  t: (key: string) => string
  major: Major
  onAddSubject: (name: string) => void
  onRemoveSubject: (subjectId: number) => void
  onRemoveMajor: () => void
  isAddingSubject?: boolean
  isRemoving?: boolean
  removingSubjectId?: number | null
}

export function MajorItem({
  t,
  major,
  onAddSubject,
  onRemoveSubject,
  onRemoveMajor,
  isAddingSubject,
  isRemoving,
  removingSubjectId,
}: MajorItemProps) {
  return (
    <AccordionItem value={String(major.id)}>
      <div className="flex items-center gap-1 pe-2">
        <AccordionTrigger className="py-3 ps-3">
          <span className="flex items-center gap-2 font-normal">
            {major.name}
            <Badge variant="outline">{major.subjects.length}</Badge>
          </span>
        </AccordionTrigger>
        <Button
          type="button"
          variant="ghost"
          size="icon-sm"
          className="text-muted-foreground hover:text-destructive"
          onClick={onRemoveMajor}
          disabled={isRemoving}
        >
          <Trash2Icon />
          <span className="sr-only">
            {t("onboarding.removeLabel")} {major.name}
          </span>
        </Button>
      </div>

      <AccordionContent className="ps-3">
        <SubjectChips
          t={t}
          subjects={major.subjects}
          onAdd={onAddSubject}
          onRemove={onRemoveSubject}
          isAdding={isAddingSubject}
          removingId={removingSubjectId}
        />
      </AccordionContent>
    </AccordionItem>
  )
}
