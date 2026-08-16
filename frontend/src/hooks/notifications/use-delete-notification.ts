import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { NotificationService } from "@/api/services/notification-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import { useNotificationsSnapshot } from "./use-notifications-snapshot"
import { updateNotificationQueries } from "@/utils/notifications-utils"

export function useDeleteNotification(t: (key: string) => string) {
  const queryClient = useQueryClient()
  const { takeSnapshot, applySnapshot, applyUnreadCount } =
    useNotificationsSnapshot()

  return useMutation({
    mutationFn: NotificationService.delete,
    onMutate: async (id) => {
      const context = await takeSnapshot()
      const deletedItem = context.previousNotifications
        .flatMap(([, data]) => data?.items ?? [])
        .find((item) => item.id === id)

      if (!deletedItem) {
        return context
      }

      const nextNotifications = updateNotificationQueries(
        context.previousNotifications,
        (data, queryKey) => {
          const shouldAdjustTotal =
            queryKey[1] === "read" ? deletedItem.is_read : !deletedItem.is_read

          return {
            ...data,
            items: data.items.filter((item) => item.id !== id),
            counts: {
              all: Math.max(0, data.counts.all - 1),
              read: deletedItem.is_read
                ? Math.max(0, data.counts.read - 1)
                : data.counts.read,
              unread: deletedItem.is_read
                ? data.counts.unread
                : Math.max(0, data.counts.unread - 1),
            },
            pagination: {
              ...data.pagination,
              total: shouldAdjustTotal
                ? Math.max(0, data.pagination.total - 1)
                : data.pagination.total,
            },
          }
        }
      )

      nextNotifications.forEach(([queryKey, data]) => {
        queryClient.setQueryData(queryKey, data)
      })

      if (!deletedItem.is_read) {
        applyUnreadCount(Math.max(0, (context.previousUnreadCount ?? 0) - 1))
      }

      return context
    },
    onSuccess: () => {
      const toastId = toast.success(t("notifications.deleteSuccess"), {
        description: t("notifications.deleteSuccessDescription"),
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
