// src/layouts/settings-layout.tsx

import type { PropsWithChildren } from "react"
import { useLocation, Link } from "react-router-dom"
import { Button } from "../components/ui/button"
import { Separator } from "../components/ui/separator"
import { useDirection } from "../context/direction/direction-provider"

type NavItem = {
  title: string
  href: string
}

export default function SettingLayout({ children }: PropsWithChildren) {
  const location = useLocation()
  const { direction, t } = useDirection()

  const sidebarNavItems: NavItem[] = [
    { title: t("nav.profile"), href: "/profile" },
    { title: t("nav.password"), href: "/password" },
    { title: t("nav.appearance"), href: "/appearance" },
  ]

  return (
    <div dir={direction} className="px-4">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-xl font-semibold">{t("settings.title")}</h1>
        <p className="text-sm text-muted-foreground">
          {t("settings.description")}
        </p>
      </div>

      <div className="flex flex-col gap-6 lg:flex-row lg:gap-12">
        {/* Sidebar */}
        <aside className="w-full max-w-xl lg:w-48">
          <nav className="flex flex-col space-y-1">
            {sidebarNavItems.map((item, index) => {
              const isActive = location.pathname === item.href

              return (
                <Button
                  key={index}
                  variant="sidebar"
                  size="sm"
                  asChild
                  data-active={isActive}
                  className="w-full justify-start text-start"
                >
                  <Link to={item.href}>{item.title}</Link>
                </Button>
              )
            })}
          </nav>
        </aside>

        {/* Mobile separator */}
        <Separator className="my-6 lg:hidden" />

        {/* Content */}
        <div className="flex-1 md:max-w-xl">
          <section className="space-y-6">{children}</section>
        </div>
      </div>
    </div>
  )
}
