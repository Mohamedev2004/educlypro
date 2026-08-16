import { useQueryClient } from "@tanstack/react-query"
import type {
  NotificationFilter,
  NotificationsData,
} from "@/api/types/notification.types"
import {
  isNotificationsListKey,
  type NotificationsSnapshot,
} from "@/utils/notifications-utils"

export type MutationContext = {
  previousNotifications: NotificationsSnapshot
  previousUnreadCount: number | undefined
  previousActiveTab?: NotificationFilter
  previousPage?: number
}

export function useNotificationsSnapshot() {
  const queryClient = useQueryClient()

  async function takeSnapshot(): Promise<MutationContext> {
    return {
      previousNotifications: queryClient.getQueriesData<NotificationsData>({
        predicate: (query) => isNotificationsListKey(query.queryKey),
      }),
      previousUnreadCount: queryClient.getQueryData<number>([
        "notifications",
        "unread-count",
      ]),
    }
  }

  function applySnapshot(snapshot: NotificationsSnapshot) {
    snapshot.forEach(([queryKey, data]) => {
      queryClient.setQueryData(queryKey, data)
    })
  }

  function applyUnreadCount(nextCount: number | undefined) {
    queryClient.setQueryData(["notifications", "unread-count"], nextCount ?? 0)
  }

  return {
    takeSnapshot,
    applySnapshot,
    applyUnreadCount,
  }
}
