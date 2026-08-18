import * as React from "react"
import { PlusIcon } from "lucide-react"

import { Button } from "@/components/ui/button"
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandItem,
  CommandInput,
  CommandList,
} from "@/components/ui/command"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"

/**
 * Search-or-create picker used at every level of the onboarding wizard
 * (grade / major / subject) — type to filter suggestions, pick one, or
 * create a new entry from the typed query.
 *
 * Responsibility: pure UI + local filter state. Owns no server state —
 * `onAdd` is the only way data leaves this component.
 * Layer: Components (domain)
 */

type AddEntityPopoverProps = {
  t: (key: string) => string
  triggerLabel: string
  placeholder: string
  suggestions: readonly string[]
  existingNames: readonly string[]
  onAdd: (name: string) => void
  disabled?: boolean
  triggerVariant?: "default" | "outline" | "secondary"
  triggerSize?: "default" | "sm" | "xs"
}

export function AddEntityPopover({
  t,
  triggerLabel,
  placeholder,
  suggestions,
  existingNames,
  onAdd,
  disabled = false,
  triggerVariant = "outline",
  triggerSize = "sm",
}: AddEntityPopoverProps) {
  const [open, setOpen] = React.useState(false)
  const [query, setQuery] = React.useState("")

  const existingLower = React.useMemo(
    () => existingNames.map((name) => name.toLowerCase()),
    [existingNames]
  )

  const filteredSuggestions = React.useMemo(
    () =>
      suggestions.filter(
        (suggestion) =>
          !existingLower.includes(suggestion.toLowerCase()) &&
          suggestion.toLowerCase().includes(query.trim().toLowerCase())
      ),
    [suggestions, existingLower, query]
  )

  const trimmedQuery = query.trim()
  const isDuplicate = existingLower.includes(trimmedQuery.toLowerCase())
  const matchesSuggestion = suggestions.some(
    (suggestion) => suggestion.toLowerCase() === trimmedQuery.toLowerCase()
  )
  const canCreate =
    trimmedQuery.length > 0 && !isDuplicate && !matchesSuggestion

  function handleAdd(name: string) {
    const trimmed = name.trim()
    if (!trimmed || disabled) return
    if (existingLower.includes(trimmed.toLowerCase())) return

    onAdd(trimmed)
    setQuery("")
    setOpen(false)
  }

  return (
    <Popover
      open={open}
      onOpenChange={(nextOpen) => {
        if (disabled) return
        setOpen(nextOpen)
        if (!nextOpen) setQuery("")
      }}
    >
      <PopoverTrigger asChild>
        <Button variant={triggerVariant} size={triggerSize} disabled={disabled}>
          <PlusIcon />
          {triggerLabel}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-72 p-0" align="start">
        <Command shouldFilter={false}>
          <CommandInput
            value={query}
            onValueChange={setQuery}
            placeholder={placeholder}
          />
          <CommandList>
            {filteredSuggestions.length > 0 && (
              <CommandGroup heading={t("onboarding.suggestions")}>
                {filteredSuggestions.map((suggestion) => (
                  <CommandItem
                    key={suggestion}
                    value={suggestion}
                    onSelect={() => handleAdd(suggestion)}
                  >
                    {suggestion}
                  </CommandItem>
                ))}
              </CommandGroup>
            )}

            {canCreate && (
              <CommandGroup heading={t("onboarding.createNew")}>
                <CommandItem
                  value={`create-${trimmedQuery}`}
                  onSelect={() => handleAdd(trimmedQuery)}
                >
                  <PlusIcon />
                  {t("onboarding.createOption").replace("{name}", trimmedQuery)}
                </CommandItem>
              </CommandGroup>
            )}

            {isDuplicate && (
              <p className="px-3 py-2 text-xs text-muted-foreground">
                {t("onboarding.alreadyAdded").replace("{name}", trimmedQuery)}
              </p>
            )}

            {!trimmedQuery && filteredSuggestions.length === 0 && (
              <CommandEmpty>{t("onboarding.allSuggestionsAdded")}</CommandEmpty>
            )}
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  )
}
