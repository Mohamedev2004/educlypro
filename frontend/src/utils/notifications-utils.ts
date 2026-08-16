import type {
  NotificationFilter,
  NotificationsData,
} from "@/api/types/notification.types"

export type NotificationsQueryKey = [
  "notifications",
  NotificationFilter,
  number,
  number,
]

export type NotificationsSnapshot = Array<
  [readonly unknown[], NotificationsData | undefined]
>

export function isNotificationsListKey(
  queryKey: readonly unknown[]
): queryKey is NotificationsQueryKey {
  return (
    queryKey[0] === "notifications" &&
    (queryKey[1] === "read" || queryKey[1] === "unread") &&
    typeof queryKey[2] === "number" &&
    typeof queryKey[3] === "number"
  )
}

export function recalculatePagination(data: NotificationsData): NotificationsData {
  const totalPages = Math.max(
    1,
    Math.ceil(data.pagination.total / data.pagination.per_page)
  )

  return {
    ...data,
    pagination: {
      ...data.pagination,
      total_pages: totalPages,
      has_next: data.pagination.page < totalPages,
      has_prev: data.pagination.page > 1,
    },
  }
}

export function updateNotificationQueries(
  previousNotifications: NotificationsSnapshot,
  updater: (
    data: NotificationsData,
    queryKey: NotificationsQueryKey
  ) => NotificationsData
): NotificationsSnapshot {
  return previousNotifications.map(([queryKey, data]) => {
    if (!data || !isNotificationsListKey(queryKey)) {
      return [queryKey, data]
    }

    return [queryKey, recalculatePagination(updater(data, queryKey))]
  })
}
