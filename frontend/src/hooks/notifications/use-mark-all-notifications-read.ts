import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { NotificationService } from "@/api/services/notification-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import { useNotificationsSnapshot } from "./use-notifications-snapshot"
import {
  isNotificationsListKey,
  updateNotificationQueries,
} from "@/utils/notifications-utils"
import type {
  NotificationFilter,
  NotificationItem,
  NotificationsData,
} from "@/api/types/notification.types"

export function useMarkAllNotificationsRead(
  t: (key: string) => string,
  setActiveTab: (tab: NotificationFilter) => void,
  setPage: (page: number) => void
) {
  const queryClient = useQueryClient()
  const { takeSnapshot, applySnapshot, applyUnreadCount } =
    useNotificationsSnapshot()

  return useMutation({
    mutationFn: NotificationService.markAllRead,
    onMutate: async () => {
      const context = await takeSnapshot()

      const promotedItems: NotificationItem[] = []

      const nextNotifications = updateNotificationQueries(
        context.previousNotifications,
        (data, queryKey) => {
          const newlyReadItems = data.items
            .filter((item) => !item.is_read)
            .map((item) => ({
              ...item,
              is_read: true,
              read_at: item.read_at ?? new Date().toISOString(),
            }))

          promotedItems.push(...newlyReadItems)

          if (queryKey[1] === "unread") {
            return {
              ...data,
              items: [],
              counts: {
                ...data.counts,
                read: data.counts.all,
                unread: 0,
              },
              pagination: {
                ...data.pagination,
                total: 0,
              },
            }
          }

          return {
            ...data,
            items: data.items.map((item) =>
              item.is_read
                ? item
                : {
                    ...item,
                    is_read: true,
                    read_at: new Date().toISOString(),
                  }
            ),
            counts: {
              ...data.counts,
              read: data.counts.all,
              unread: 0,
            },
          }
        }
      )

      nextNotifications.forEach(([queryKey, data]) => {
        if (!data || !isNotificationsListKey(queryKey)) {
          return
        }

        if (
          queryKey[1] === "read" &&
          queryKey[2] === 1 &&
          promotedItems.length
        ) {
          const uniqueItems = [
            ...promotedItems,
            ...data.items.filter(
              (item) =>
                !promotedItems.some((promoted) => promoted.id === item.id)
            ),
          ].slice(0, data.pagination.per_page)

          queryClient.setQueryData(queryKey, {
            ...data,
            items: uniqueItems,
            counts: {
              ...data.counts,
              read: data.counts.all,
              unread: 0,
            },
            pagination: {
              ...data.pagination,
              total: data.counts.all,
              total_pages: Math.max(
                1,
                Math.ceil(data.counts.all / data.pagination.per_page)
              ),
              has_next:
                data.pagination.page <
                Math.max(
                  1,
                  Math.ceil(data.counts.all / data.pagination.per_page)
                ),
            },
          } satisfies NotificationsData)
          return
        }

        queryClient.setQueryData(queryKey, data)
      })

      applyUnreadCount(0)
      setActiveTab("read")
      setPage(1)

      return context
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] })
      const toastId = toast.success(t("notifications.markAllReadSuccess"), {
        description: t("notifications.markAllReadSuccessDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
    },
    onError: (error, _variables, context) => {
      if (context) {
        applySnapshot(context.previousNotifications)
        applyUnreadCount(context.previousUnreadCount)
      }
      toast.error(getApiMessage(t, normalizeApiError(error)))
    },
  })
}
