/* eslint-disable @typescript-eslint/no-explicit-any */
import { useNavigate } from "react-router-dom"
import {
  Command,
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "../components/ui/command"

export function CommandPalette({
  open,
  setOpen,
  groups,
  t,
}: {
  open: boolean
  setOpen: (v: boolean) => void
  groups: any[]
  t: (key: string) => string
}) {
  const navigate = useNavigate()

  return (
    <CommandDialog open={open} onOpenChange={setOpen}>
      <Command>
        <CommandInput placeholder={t("shell.searchPagesPlaceholder")} />

        <CommandList className="max-h-[300px] overflow-y-auto">
          <CommandEmpty>{t("shell.noResults")}</CommandEmpty>

          {groups.map((group) => (
            <CommandGroup key={group.label} heading={group.label}>
              {group.items.map((item: any) => {
                const Icon = item.icon

                return (
                  <CommandItem
                    key={item.url}
                    onSelect={() => {
                      setOpen(false)
                      navigate(item.url)
                    }}
                    className="flex items-center gap-2"
                  >
                    <Icon className="size-4 text-muted-foreground" />
                    {item.title}
                  </CommandItem>
                )
              })}
            </CommandGroup>
          ))}
        </CommandList>
      </Command>
    </CommandDialog>
  )
}
