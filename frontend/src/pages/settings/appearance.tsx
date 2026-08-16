// src/pages/appearance.tsx
import { useDirection } from "@/context/direction/direction-provider"
import { useTheme } from "../../components/theme-provider"
import AppLayout from "../../layouts/app-layout"
import SettingLayout from "../../layouts/setting-layout"
import { cn } from "@/utils/ui-utils"
import { Monitor, Moon, Sun } from "lucide-react"

type Theme = "light" | "dark" | "system"

const themes: { value: Theme; labelKey: string; icon: React.ElementType; descriptionKey: string }[] = [
  {
    value: "light",
    labelKey: "settings.themes.light",
    icon: Sun,
    descriptionKey: "settings.themes.lightDescription",
  },
  {
    value: "dark",
    labelKey: "settings.themes.dark",
    descriptionKey: "settings.themes.darkDescription",
    icon: Moon,
  },
  {
    value: "system",
    labelKey: "settings.themes.system",
    icon: Monitor,
    descriptionKey: "settings.themes.systemDescription",
  },
]

export default function Appearance() {
  const { theme, setTheme } = useTheme()
  const { t } = useDirection()

  return (
    <AppLayout breadcrumbs={[{ label: t("shell.settings") }, { label: t("nav.appearance") }]}>
      <SettingLayout>
        <div className="space-y-1">
          <h2 className="text-base font-medium">{t("settings.appearanceTitle")}</h2>
          <p className="text-sm text-muted-foreground">
            {t("settings.appearanceDescription")}
          </p>
        </div>

        <div className="grid grid-cols-3 gap-3 pt-2">
          {themes.map(({ value, labelKey, icon: Icon, descriptionKey }) => {
            const isActive = theme === value
            return (
              <button
                key={value}
                onClick={() => setTheme(value)}
                className={cn(
                  "relative flex flex-col items-center gap-3 rounded-xl border-2 p-4 text-center transition-all duration-150 hover:bg-muted/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring",
                  isActive
                    ? "border-primary bg-primary/5"
                    : "border-border bg-transparent"
                )}
              >
                {/* Theme preview */}
                <div
                  className={cn(
                    "flex h-16 w-full items-center justify-center rounded-lg",
                    value === "light" && "bg-white ring-1 ring-border",
                    value === "dark" && "bg-zinc-900 ring-1 ring-zinc-700",
                    value === "system" &&
                      "bg-gradient-to-br from-white to-zinc-900 ring-1 ring-border"
                  )}
                >
                  <Icon
                    className={cn(
                      "size-6",
                      value === "light" && "text-zinc-800",
                      value === "dark" && "text-zinc-200",
                      value === "system" && "text-zinc-500"
                    )}
                  />
                </div>

                <div className="space-y-0.5">
                  <p className="text-sm font-medium">{t(labelKey)}</p>
                  <p className="text-xs text-muted-foreground">{t(descriptionKey)}</p>
                </div>

                {/* Active dot */}
                {isActive && (
                  <span className="absolute right-2 top-2 size-2 rounded-full bg-primary" />
                )}
              </button>
            )
          })}
        </div>
      </SettingLayout>
    </AppLayout>
  )
}