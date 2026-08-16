import { CheckCheck, Trash2 } from "lucide-react"

import type {
  NotificationCounts,
  NotificationFilter,
} from "@/api/types/notification.types"
import { Badge } from "@/components/ui/badge"
import { Button } from "@/components/ui/button"
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { NOTIFICATIONS_PER_PAGE_OPTIONS } from "@/constants/notifications"
import type { Translate } from "@/types/notifications"
import { useAuth } from "@/context/auth/auth-context"

type NotificationsToolbarProps = {
  activeTab: NotificationFilter
  counts?: NotificationCounts
  perPage: number
  isMutating: boolean
  t: Translate
  onTabChange: (value: string) => void
  onPerPageChange: (value: string) => void
  onMarkAllRead: () => void
  onDeleteAllRead: () => void
}

export function NotificationsToolbar({
  activeTab,
  counts,
  perPage,
  isMutating,
  t,
  onTabChange,
  onPerPageChange,
  onMarkAllRead,
  onDeleteAllRead,
}: NotificationsToolbarProps) {
  const { user } = useAuth()
  const is_super_admin = user?.role === "super_admin"
  return (
    <div
      className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between"
    >
      <Tabs value={activeTab} onValueChange={onTabChange} className="gap-4">
        <TabsList variant="line" className="gap-2 rounded-md p-0">
          <TabsTrigger value="unread" className="gap-2 px-1">
            {t("notifications.unreadTab")}
            <Badge variant="ghost">{counts?.unread ?? 0}</Badge>
          </TabsTrigger>
          <TabsTrigger value="read" className="gap-2 px-1">
            {t("notifications.readTab")}
            <Badge variant="ghost">{counts?.read ?? 0}</Badge>
          </TabsTrigger>
        </TabsList>
      </Tabs>

      <div className="flex flex-col items-end gap-3 sm:flex-row">
        <div className="w-full space-y-1 sm:w-auto">
          <div className="text-xs font-medium text-muted-foreground">
            {t("notifications.perPageLabel")}
          </div>

          <Select value={String(perPage)} onValueChange={onPerPageChange}>
            <SelectTrigger className="w-full sm:w-30">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {NOTIFICATIONS_PER_PAGE_OPTIONS.map((option) => (
                <SelectItem key={option} value={String(option)}>
                  {option} {t("common.perPage")}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>

        {activeTab === "unread" && (
          <Button
            variant="default"
            disabled={isMutating || (counts?.unread ?? 0) === 0}
            onClick={onMarkAllRead}
            className="w-full sm:w-auto"
          >
            <CheckCheck className="me-2 size-4" />
            {t("notifications.markAllRead")}
          </Button>
        )}
        {activeTab === "read" && is_super_admin && (
          <Button
            variant="destructive"
            disabled={isMutating || (counts?.read ?? 0) === 0}
            onClick={onDeleteAllRead}
            className="w-full sm:w-auto"
          >
            <Trash2 className="me-2 size-4" />
            {t("notifications.deleteAllRead")}
          </Button>
        )}
      </div>
    </div>
  )
}
