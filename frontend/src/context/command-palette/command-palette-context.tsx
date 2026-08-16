/* eslint-disable react-refresh/only-export-components */
import * as React from "react"

/**
 * Command palette (Ctrl/Cmd+K) open state.
 *
 * Responsibility: single source of truth for whether the command palette is
 * open, shared by every trigger (header search, sidebar search, keyboard
 * shortcut) so there is exactly one dialog instance and one state.
 * Layer: Context
 */

type CommandPaletteContextValue = {
  open: boolean
  setOpen: (open: boolean) => void
}

const CommandPaletteContext = React.createContext<
  CommandPaletteContextValue | undefined
>(undefined)

export function CommandPaletteProvider({
  children,
}: {
  children: React.ReactNode
}) {
  const [open, setOpen] = React.useState(false)

  React.useEffect(() => {
    const handler = (event: KeyboardEvent) => {
      if ((event.metaKey || event.ctrlKey) && event.key === "k") {
        event.preventDefault()
        setOpen((prev) => !prev)
      }
    }

    window.addEventListener("keydown", handler)
    return () => window.removeEventListener("keydown", handler)
  }, [])

  const value = React.useMemo(() => ({ open, setOpen }), [open])

  return (
    <CommandPaletteContext.Provider value={value}>
      {children}
    </CommandPaletteContext.Provider>
  )
}

export function useCommandPalette() {
  const context = React.useContext(CommandPaletteContext)
  if (!context) {
    throw new Error(
      "useCommandPalette must be used inside CommandPaletteProvider"
    )
  }
  return context
}
