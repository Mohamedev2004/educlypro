// src/layouts/app-layout.tsx
import * as React from "react"
import { AppSidebar } from "../components/app-sidebar"
import { SidebarInset, SidebarProvider } from "../components/ui/sidebar"
import { AppHeader } from "../components/app-header"
import { CommandPalette } from "../components/command-palette"
import { useDirection } from "../context/direction/direction-provider"
import { useAuth } from "../context/auth/auth-context"
import { useCommandPalette } from "../context/command-palette/command-palette-context"
import { getRolePages } from "@/utils/navigation-utils"
import type { BreadcrumbType } from "@/components/app-breadcrumbs"

export default function AppLayout({
  children,
  breadcrumbs,
}: {
  children: React.ReactNode
  breadcrumbs: BreadcrumbType[]
}) {
  const { direction, t } = useDirection()
  const { user } = useAuth()
  const { open, setOpen } = useCommandPalette()
  const [defaultOpen] = React.useState(() => {
    const match = document.cookie.match(/sidebar_state=([^;]+)/)
    return match ? match[1] === "true" : true
  })

  const groups = React.useMemo(() => {
    if (!user) return []
    return getRolePages(user, t)
  }, [user, t])

  return (
    <SidebarProvider defaultOpen={defaultOpen} dir={direction}>
      <AppSidebar groups={groups} />
      {/* <SidebarInset className="overflow-hidden"> */}
      {/* IF YOU WANT TO ENABLE SCROLLING IN THE APP, WHITE AA STICKY HEADER */}
      <SidebarInset className="overflow-y-auto h-[calc(100vh-1rem)] hide-scrollbar">

        <AppHeader breadcrumbs={breadcrumbs} />
        <div className="p-6">{children}</div>
      </SidebarInset>

      <CommandPalette open={open} setOpen={setOpen} groups={groups} t={t} />
    </SidebarProvider>
  )
}
