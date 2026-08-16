import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { NotificationService } from "@/api/services/notification-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import { useNotificationsSnapshot } from "./use-notifications-snapshot"
import { updateNotificationQueries } from "@/utils/notifications-utils"

export function useMarkNotificationRead(t: (key: string) => string) {
  const queryClient = useQueryClient()
  const { takeSnapshot, applySnapshot, applyUnreadCount } = useNotificationsSnapshot()

  return useMutation({
    mutationFn: NotificationService.markRead,
    onMutate: async (id) => {
      const context = await takeSnapshot()
      const now = new Date().toISOString()
      const changedNotification = context.previousNotifications
        .flatMap(([, data]) => data?.items ?? [])
        .find((item) => item.id === id && !item.is_read)

      if (!changedNotification) {
        return context
      }

      const nextNotifications = updateNotificationQueries(
        context.previousNotifications,
        (data, queryKey) => {
          if (queryKey[1] === "unread") {
            return {
              ...data,
              items: data.items.filter((item) => item.id !== id),
              counts: {
                ...data.counts,
                read: data.counts.read + 1,
                unread: Math.max(0, data.counts.unread - 1),
              },
              pagination: {
                ...data.pagination,
                total: Math.max(0, data.pagination.total - 1),
              },
            }
          }

          const updatedItem = {
            ...changedNotification,
            is_read: true,
            read_at: now,
          }
          const alreadyPresent = data.items.some((item) => item.id === id)
          const shouldInsertOnFirstPage = queryKey[2] === 1 && !alreadyPresent
          const nextItems = shouldInsertOnFirstPage
            ? [updatedItem, ...data.items].slice(0, data.pagination.per_page)
            : data.items.map((item) =>
                item.id === id ? { ...item, is_read: true, read_at: now } : item
              )

          return {
            ...data,
            items: nextItems,
            counts: {
              ...data.counts,
              read: data.counts.read + 1,
              unread: Math.max(0, data.counts.unread - 1),
            },
            pagination: {
              ...data.pagination,
              total: data.pagination.total + 1,
            },
          }
        }
      )

      nextNotifications.forEach(([queryKey, data]) => {
        queryClient.setQueryData(queryKey, data)
      })

      applyUnreadCount(Math.max(0, (context.previousUnreadCount ?? 0) - 1))

      return context
    },
    onSuccess: () => {
      const toastId = toast.success(t("notifications.markReadSuccess"), {
        description: t("notifications.markReadSuccessDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
    },
    onError: (error, _id, context) => {
      if (context) {
        applySnapshot(context.previousNotifications)
        applyUnreadCount(context.previousUnreadCount)
      }
      toast.error(getApiMessage(t, normalizeApiError(error)))
    },
  })
}
