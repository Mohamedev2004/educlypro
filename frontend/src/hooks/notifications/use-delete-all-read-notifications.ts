import { useMutation, useQueryClient } from "@tanstack/react-query"
import { toast } from "sonner"
import { NotificationService } from "@/api/services/notification-service"
import { normalizeApiError, getApiMessage } from "@/utils/error-utils"
import { useNotificationsSnapshot } from "./use-notifications-snapshot"
import { updateNotificationQueries } from "@/utils/notifications-utils"

export function useDeleteAllReadNotifications(t: (key: string) => string) {
  const queryClient = useQueryClient()
  const { takeSnapshot, applySnapshot, applyUnreadCount } =
    useNotificationsSnapshot()

  return useMutation({
    mutationFn: NotificationService.deleteAllRead,
    onMutate: async () => {
      const context = await takeSnapshot()

      const nextNotifications = updateNotificationQueries(
        context.previousNotifications,
        (data, queryKey) => {
          const isReadTab = queryKey[1] === "read"

          return {
            ...data,
            items: isReadTab ? [] : data.items,
            counts: {
              ...data.counts,
              all: Math.max(0, data.counts.all - data.counts.read),
              read: 0,
            },
            pagination: {
              ...data.pagination,
              total: isReadTab ? 0 : data.pagination.total,
            },
          }
        }
      )

      nextNotifications.forEach(([queryKey, data]) => {
        queryClient.setQueryData(queryKey, data)
      })

      // Read notifications don't affect unread count, no need to adjust it

      return context
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["notifications"] })
      const toastId = toast.success(t("notifications.deleteAllReadSuccess"), {
        description: t("notifications.deleteAllReadSuccessDescription"),
        action: {
          label: t("common.close"),
          onClick: () => toast.dismiss(toastId),
        },
      })
    },
    onError: (error, _vars, context) => {
      if (context) {
        applySnapshot(context.previousNotifications)
        applyUnreadCount(context.previousUnreadCount)
      }
      toast.error(getApiMessage(t, normalizeApiError(error)))
    },
  })
}