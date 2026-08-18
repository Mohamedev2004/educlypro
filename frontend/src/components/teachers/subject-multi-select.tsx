import * as React from "react"
import { PlusIcon, XIcon } from "lucide-react"

import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { useAcademicTree } from "@/hooks/academic/use-academic-tree"

/**
 * Multi-select for the subjects a teacher teaches, sourced from the
 * center's own grade -> major -> subject tree (flattened client-side — no
 * dedicated "list all subjects" endpoint needed). Same subject name under
 * different majors is disambiguated by showing the major alongside it.
 *
 * Responsibility: pure UI + local flatten/filter. Owns no server state —
 * `onChange` is the only way the selection leaves this component.
 * Layer: Components (domain)
 */

type SubjectOption = {
  id: number
  label: string
}

type SubjectMultiSelectProps = {
  t: (key: string) => string
  selectedIds: number[]
  onChange: (ids: number[]) => void
  disabled?: boolean
}

export function SubjectMultiSelect({
  t,
  selectedIds,
  onChange,
  disabled,
}: SubjectMultiSelectProps) {
  const { data } = useAcademicTree()
  const [open, setOpen] = React.useState(false)

  const options: SubjectOption[] = React.useMemo(
    () =>
      (data?.grades ?? []).flatMap((grade) =>
        grade.majors.flatMap((major) =>
          major.subjects.map((subject) => ({
            id: subject.id,
            label: `${subject.name} — ${major.name}`,
          }))
        )
      ),
    [data]
  )

  const selectedOptions = options.filter((option) =>
    selectedIds.includes(option.id)
  )

  function toggle(id: number) {
    if (selectedIds.includes(id)) {
      onChange(selectedIds.filter((existing) => existing !== id))
    } else {
      onChange([...selectedIds, id])
    }
  }

  return (
    <div className="flex flex-col gap-2">
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            type="button"
            variant="outline"
            className="w-full justify-start"
            disabled={disabled}
          >
            <PlusIcon />
            {t("teachers.fields.subjectsButton")}
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-80 p-0" align="start">
          <Command>
            <CommandInput
              placeholder={t("teachers.fields.subjectsSearchPlaceholder")}
            />
            <CommandList>
              <CommandEmpty>
                {t("teachers.fields.noSubjectsAvailable")}
              </CommandEmpty>
              <CommandGroup>
                {options.map((option) => {
                  const checked = selectedIds.includes(option.id)
                  return (
                    <CommandItem
                      key={option.id}
                      value={option.label}
                      data-checked={checked}
                      onSelect={() => toggle(option.id)}
                    >
                      {option.label}
                    </CommandItem>
                  )
                })}
              </CommandGroup>
            </CommandList>
          </Command>
        </PopoverContent>
      </Popover>

      {selectedOptions.length > 0 && (
        <div className="flex flex-wrap items-center gap-1.5">
          {selectedOptions.map((option) => (
            <Badge key={option.id} variant="secondary" className="gap-1 pe-1">
              {option.label}
              <button
                type="button"
                onClick={() => toggle(option.id)}
                disabled={disabled}
                className="rounded-full p-0.5 hover:bg-foreground/10 disabled:pointer-events-none disabled:opacity-50"
                aria-label={`${t("common.delete")} ${option.label}`}
              >
                <XIcon className="size-3" />
              </button>
            </Badge>
          ))}
        </div>
      )}
    </div>
  )
}
