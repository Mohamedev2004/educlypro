import { useMemo } from "react"
import { BellRing, CheckCheck, Trash2 } from "lucide-react"

import type {
  NotificationFilter,
  NotificationItem,
} from "@/api/types/notification.types"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import { ar } from "@/locales/ar"
import { en } from "@/locales/en"
import { fr } from "@/locales/fr"
import type { NotificationLocale, Translate } from "@/types/notifications"

type NotificationListProps = {
  items: NotificationItem[]
  filter: NotificationFilter
  locale: NotificationLocale
  t: Translate
  isMutating: boolean
  onMarkRead: (id: number) => void
  onDelete: (id: number) => void
}

const messages = {
  ar,
  en,
  fr,
} as const

import { interpolate } from "@/utils/common-utils"

function resolveTopicText(
  topic: string | undefined,
  suffix: "" | ".body",
  locale: NotificationLocale,
  fallback: string,
  payload: Record<string, unknown> | undefined
): string {
  if (!topic) {
    return interpolate(fallback, payload)
  }

  const topicMap = messages[locale].notifications.topics as Record<
    string,
    string
  >
  const text = topicMap[`${topic}${suffix}`] ?? fallback
  return interpolate(text, payload)
}

export function NotificationList({
  items,
  filter,
  locale,
  t,
  isMutating,
  onMarkRead,
  onDelete,
}: NotificationListProps) {
  const formatter = useMemo(
    () =>
      new Intl.DateTimeFormat(
        locale === "ar" ? "ar-MA" : locale === "fr" ? "fr-FR" : "en-US",
        {
          dateStyle: "medium",
          timeStyle: "short",
        }
      ),
    [locale]
  )

  if (items.length === 0) {
    return (
      <div className="rounded-md border border-dashed border-border bg-muted/30 px-6 py-12 text-center">
        <div className="mx-auto mb-4 flex size-14 items-center justify-center rounded-md bg-background shadow-sm ring-1 ring-border">
          <BellRing className="size-6 text-primary" />
        </div>
        <h3 className="font-medium">
          {filter === "unread"
            ? t("notifications.allCaughtUp")
            : t("notifications.emptyRead")}
        </h3>
      </div>
    )
  }

  return (
    <div className="space-y-3">
      {items.map((notification) => {
        // payload comes from the backend as map[string]any — cast it safely
        const payload = notification.payload as
          | Record<string, unknown>
          | undefined

        const title = resolveTopicText(
          notification.topic,
          "",
          locale,
          notification.title,
          payload
        )
        const body = resolveTopicText(
          notification.topic,
          ".body",
          locale,
          notification.body,
          payload
        )

        return (
          <div
            key={notification.id}
            className="rounded-md border border-border bg-sidebar p-4 shadow-sm transform transition-all duration-300 ease-out hover:scale-[1.01] hover:shadow-md"          >
            <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
              <div className="space-y-3">
                <div className="flex flex-wrap items-center gap-2">
                  <Badge variant={notification.is_read ? "outline" : "default"}>
                    {notification.is_read
                      ? t("notifications.readStatus")
                      : t("notifications.unreadStatus")}
                  </Badge>
                  {/* <Badge variant="secondary">
                    {notification.channel === "email"
                      ? t("notifications.emailChannel")
                      : t("notifications.inAppChannel")}
                  </Badge> */}
                </div>

                <div className="space-y-1">
                  <h3 className="text-base font-medium text-foreground">
                    {title}
                  </h3>
                  <p className="max-w-3xl text-sm leading-6 text-muted-foreground">
                    {body}
                  </p>
                </div>

                <div className="flex flex-wrap items-center gap-2 text-xs text-muted-foreground">
                  <span>
                    {formatter.format(new Date(notification.created_at))}
                  </span>
                  {/* {notification.topic && (
                    <>
                      <span className="opacity-60">•</span>
                      <span className="rounded-full bg-muted px-2 py-1 font-medium">
                        {title}
                      </span>
                    </>
                  )} */}
                </div>
              </div>

              <div className="flex w-full shrink-0 items-center gap-2 sm:w-auto">
                {!notification.is_read && (
                  <Button
                    variant="default"
                    disabled={isMutating}
                    onClick={() => onMarkRead(notification.id)}
                    className="flex-1 sm:flex-none"
                  >
                    <CheckCheck className="me-2 size-4" />
                    {t("notifications.markAsRead")}
                  </Button>
                )}

                <Button
                  variant="destructive"
                  disabled={isMutating}
                  onClick={() => onDelete(notification.id)}
                  className="flex-1 sm:flex-none"
                >
                  <Trash2 className="me-2 size-4" />
                  {t("common.delete")}
                </Button>
              </div>
            </div>
          </div>
        )
      })}
    </div>
  )
}
