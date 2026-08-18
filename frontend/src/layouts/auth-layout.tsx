/* eslint-disable @typescript-eslint/no-explicit-any */
import type { ReactNode } from "react"
import { Link } from "react-router-dom"
import { GalleryVerticalEnd, Moon, Sun } from "lucide-react"
import { useDirection } from "../context/direction/direction-provider"

import { Button } from "@/components/ui/button"
import { AppFullscreen } from "@/components/app-fullscreen"
import { useTheme } from "@/components/theme-provider"
import { AppLanguage } from "@/components/app-language"
import { cn } from "@/utils/ui-utils"

type AuthLayoutProps = {
  children: ReactNode
  // Overrides the content column's max width — defaults to the narrow
  // max-w-sm form column every auth page uses. Wider flows (e.g. the
  // onboarding wizard) can widen it without affecting any other page.
  contentClassName?: string
}

export default function AuthLayout({
  children,
  contentClassName,
}: AuthLayoutProps) {
  const { direction, locale, setLocale, t } = useDirection()
  const { theme, setTheme } = useTheme()

  const isDark =
    theme === "dark" ||
    (theme === "system" &&
      window.matchMedia("(prefers-color-scheme: dark)").matches)

  const toggleTheme = () => setTheme(isDark ? "light" : "dark")

  return (
    <div dir={direction} className="grid min-h-svh lg:grid-cols-2">
      {/* Controls */}
      <div
        dir="ltr"
        className="fixed top-4 right-1/2 z-50 flex items-center justify-end gap-2"
      >
        <AppLanguage locale={locale} setLocale={(v) => setLocale(v as any)} />

        <AppFullscreen />

        <Button
          variant="default"
          size="icon"
          onClick={toggleTheme}
          className="relative size-9"
        >
          <Sun className="size-4 transition-all dark:scale-0 dark:-rotate-90" />
          <Moon className="absolute size-4 scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0" />
        </Button>
      </div>

      {/* COLORED SECTION → RIGHT */}
      <div
        className={[
          "relative hidden p-4 lg:flex lg:flex-col",
          direction === "rtl" ? "lg:order-1" : "lg:order-2",
        ].join(" ")}
      >
        <div className="relative flex-1 overflow-hidden rounded-xl bg-primary/15">
          {/* Your content / image can go here */}
        </div>
      </div>

      {/* FORM / CONTENT → LEFT */}
      <div
        className={[
          "relative flex items-center justify-center py-6",
          direction === "rtl" ? "lg:order-2" : "lg:order-1",
        ].join(" ")}
      >
        <div className={cn("w-full space-y-6", contentClassName ?? "max-w-sm")}>
          {/* Logo */}
          <div className="flex items-center justify-center gap-2 self-center font-medium">
            <div className="flex size-8 items-center justify-center rounded-md bg-primary text-primary-foreground">
              <GalleryVerticalEnd className="size-4" />
            </div>
            {t("shell.appName")}.
          </div>

          {children}

          <p className="text-center text-xs text-muted-foreground">
            By clicking continue, you agree to our{" "}
            <Link
              to="/terms"
              className="underline underline-offset-4 hover:text-primary"
            >
              Terms of Service
            </Link>{" "}
            and{" "}
            <Link
              to="/privacy"
              className="underline underline-offset-4 hover:text-primary"
            >
              Privacy Policy
            </Link>
            .
          </p>
        </div>
      </div>
    </div>
  )
}
