import { keepPreviousData, useQuery } from "@tanstack/react-query"
import { NotificationService } from "@/api/services/notification-service"
import type { NotificationFilter } from "@/api/types/notification.types"

export function useNotificationsList(
  activeTab: NotificationFilter,
  page: number,
  perPage: number
) {
  return useQuery({
    queryKey: ["notifications", activeTab, page, perPage],
    queryFn: () =>
      NotificationService.list({
        page,
        perPage,
        filter: activeTab,
      }),
    placeholderData: keepPreviousData,
  })
}
