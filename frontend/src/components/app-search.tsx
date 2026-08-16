import { Search } from "lucide-react"
import { Kbd } from "@/components/ui/kbd"

export function AppSearch({ onClick, label }: { onClick: () => void; label: string }) {
  return (
    <div
      onClick={onClick}
      className="hidden h-9 w-[240px] cursor-pointer items-center gap-2 rounded-md border bg-background px-3 text-sm text-muted-foreground hover:bg-accent lg:flex"
    >
      <Search className="h-4 w-4" />
      <span className="flex-1 text-start">{label}</span>
      <Kbd>Ctrl K</Kbd>
    </div>
  )
}

export function SidebarSearch({ onClick, label }: { onClick: () => void; label: string }) {
  return (
    <div className="px-2 mt-2">
      <button
        onClick={onClick}
        className="
          group cursor-pointer flex w-full items-center gap-2 rounded-md border bg-background px-3 py-2 text-sm text-muted-foreground hover:bg-accent
          group-data-[collapsible=icon]:justify-center
          group-data-[collapsible=icon]:px-2
        "
      >
        <Search className="size-4 shrink-0" />

        {/* 👇 Hide text when collapsed */}
        <span className="flex-1 text-start group-data-[collapsible=icon]:hidden">
          {label}
        </span>

        {/* 👇 Hide shortcut when collapsed */}
        <Kbd className="group-data-[collapsible=icon]:hidden">Ctrl K</Kbd>
      </button>
    </div>
  )
}
