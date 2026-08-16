"use client"

import * as React from "react"
import { ChevronsUpDownIcon, GalleryVerticalEnd } from "lucide-react"

import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuBadge,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarGroup,
  SidebarGroupLabel,
} from "../components/ui/sidebar"
import { NavUser } from "../components/nav-user"
import { Link, useLocation } from "react-router-dom"
import { useDirection } from "../context/direction/direction-provider"
import { useAuth } from "../context/auth/auth-context"
import { useUnreadNotificationsCount } from "@/hooks/notifications/use-unread-notifications-count"
import { getDashboardPath } from "@/utils/navigation-utils"
import { roleLabels } from "@/constants/roles"
import { useCommandPalette } from "../context/command-palette/command-palette-context"
import { SidebarSearch } from "./app-search"
import type { RolePageGroup } from "@/types/navigation.types"
// import { Separator } from "./ui/separator"

type AppSidebarProps = React.ComponentProps<typeof Sidebar> & {
  groups: RolePageGroup[]
}

export function AppSidebar({ groups, ...props }: AppSidebarProps) {
  const location = useLocation()
  const { direction, t } = useDirection()
  const { user } = useAuth()
  const { setOpen } = useCommandPalette()

  const { unreadCount } = useUnreadNotificationsCount()

  const dashboardPath = React.useMemo(() => {
    if (!user) return "/"
    return getDashboardPath(user.role)
  }, [user])

  return (
    <Sidebar
      side={direction === "rtl" ? "right" : "left"}
      variant="inset"
      collapsible="icon"
      {...props}
    >
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton size="lg" asChild>
              <Link to={dashboardPath}>
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-primary text-primary-foreground">
                  <GalleryVerticalEnd className="size-4" />
                </div>
                <div className="grid flex-1 text-start text-sm leading-tight">
                  <span className="truncate font-medium">
                    {t("shell.appDashboardName")}
                  </span>
                  <span className="truncate text-xs text-muted-foreground">
                    {user ? t(roleLabels[user.role]) : t("shell.adminPanel")}
                  </span>
                </div>
                <ChevronsUpDownIcon className="ms-auto size-4" />
              </Link>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>

      {/* SEARCH */}
      <SidebarSearch onClick={() => setOpen(true)} label={t("shell.search")} />

      <SidebarContent
        className="overflow-y-auto max-h-[100vh] overscroll-contain"
      >
        <SidebarGroup>
          <SidebarMenu>
            {groups.map((group) => (
              <React.Fragment key={group.label}>
                <SidebarGroupLabel className="px-2 py-1 text-xs text-muted-foreground group-data-[collapsible=icon]:hidden">
                  {group.label}
                </SidebarGroupLabel>

                {group.items.map((item) => {
                  const isActive = location.pathname === item.url
                  const showUnreadBadge =
                    item.url.endsWith("/notifications") && unreadCount > 0

                  return (
                    <SidebarMenuItem key={item.title}>
                      <SidebarMenuButton
                        asChild
                        variant="sidebar"
                        data-active={isActive}
                        className="w-full"
                        tooltip={item.title}
                      >
                        <Link to={item.url} className="flex items-center gap-2">
                          <item.icon />
                          <span>{item.title}</span>
                        </Link>
                      </SidebarMenuButton>

                      {showUnreadBadge && (
                        <SidebarMenuBadge data-active={isActive}>
                          {unreadCount}
                        </SidebarMenuBadge>
                      )}
                    </SidebarMenuItem>
                  )
                })}
                {/* <Separator /> */}
              </React.Fragment>
            ))}
          </SidebarMenu>
        </SidebarGroup>
      </SidebarContent>

      <SidebarFooter>
        <NavUser
          user={{
            name: user?.username ?? "Guest",
            email: user?.email ?? "guest@example.com",
            avatar: "",
          }}
        />
      </SidebarFooter>
    </Sidebar>
  )
}
