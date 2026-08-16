import { BellRing } from "lucide-react"

import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import type { Translate } from "@/types/notifications"

type NotificationsHeroProps = {
  unreadCount: number
  readCount: number
  t: Translate
  onTabChange: (tab: "unread" | "read") => void
}

export function NotificationsHero({
  unreadCount,
  readCount,
  t,
  onTabChange,
}: NotificationsHeroProps) {
  return (
    <div>
      <Card className="relative overflow-hidden border-none text-foreground">
        <CardHeader className="relative z-10">
          <div className="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
            {/* LEFT */}
            <div className="space-y-3">
              <div className="inline-flex items-center gap-2 rounded-md border border-foreground/10 bg-foreground/5 px-3 py-1 text-xs font-medium tracking-wider text-foreground/70 backdrop-blur">
                <BellRing className="size-3.5" />
                {t("notifications.title")}
              </div>

              <CardTitle className="text-2xl font-semibold tracking-tight text-foreground">
                {t("notifications.title")}
              </CardTitle>

              <CardDescription className="max-w-3xl text-sm text-muted-foreground text-justify">
                {t("notifications.description")}
              </CardDescription>
            </div>

            {/* RIGHT */}
            <div className="grid grid-cols-2 gap-4 sm:min-w-[260px]">
              {/* UNREAD */}
              <div
                onClick={() => onTabChange("unread")}
                className="group cursor-pointer relative overflow-hidden rounded-md border border-foreground/10 bg-foreground/5 p-5 transition hover:bg-foreground/10"
              >
                <div className="text-xs tracking-wider text-foreground/50 uppercase">
                  {t("notifications.unreadTab")}
                </div>

                <div className="mt-3 text-4xl font-bold text-foreground">
                  {unreadCount}
                </div>
              </div>

              {/* READ */}
              <div
                onClick={() => onTabChange("read")}
                className="group cursor-pointer relative overflow-hidden rounded-md border border-foreground/10 bg-foreground/5 p-5 transition hover:bg-foreground/10"
              >
                <div className="text-xs tracking-wider text-foreground/50 uppercase">
                  {t("notifications.readTab")}
                </div>

                <div className="mt-3 text-4xl font-bold text-foreground">
                  {readCount}
                </div>
              </div>
            </div>
          </div>
        </CardHeader>
      </Card>
    </div>
  )
}
