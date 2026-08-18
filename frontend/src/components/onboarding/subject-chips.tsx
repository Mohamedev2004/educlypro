import { XIcon } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { AddEntityPopover } from "@/components/onboarding/add-entity-popover"
import { SUBJECT_SUGGESTIONS } from "@/constants/academic"
import type { Subject } from "@/api/types/academic.types"

type SubjectChipsProps = {
  t: (key: string) => string
  subjects: Subject[]
  onAdd: (name: string) => void
  onRemove: (subjectId: number) => void
  isAdding?: boolean
  removingId?: number | null
}

export function SubjectChips({
  t,
  subjects,
  onAdd,
  onRemove,
  isAdding,
  removingId,
}: SubjectChipsProps) {
  return (
    <div className="flex flex-wrap items-center gap-1.5">
      {subjects.map((subject) => (
        <Badge key={subject.id} variant="secondary" className="gap-1 pe-1">
          {subject.name}
          <button
            type="button"
            onClick={() => onRemove(subject.id)}
            disabled={removingId === subject.id}
            className="rounded-full p-0.5 hover:bg-foreground/10 disabled:pointer-events-none disabled:opacity-50"
            aria-label={`${t("onboarding.removeLabel")} ${subject.name}`}
          >
            <XIcon className="size-3" />
          </button>
        </Badge>
      ))}

      {subjects.length === 0 && (
        <p className="text-xs text-muted-foreground">
          {t("onboarding.noSubjects")}
        </p>
      )}

      <AddEntityPopover
        t={t}
        triggerLabel={t("onboarding.addSubject")}
        placeholder={t("onboarding.addSubjectPlaceholder")}
        suggestions={SUBJECT_SUGGESTIONS}
        existingNames={subjects.map((subject) => subject.name)}
        onAdd={onAdd}
        disabled={isAdding}
        triggerVariant="secondary"
        triggerSize="xs"
      />
    </div>
  )
}
