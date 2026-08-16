import { useQuery } from "@tanstack/react-query"
import { NotificationService } from "@/api/services/notification-service"
import { useAuth } from "@/context/auth/auth-context"

/**
 * Hook to fetch and manage the unread notifications count.
 * 
 * Responsibility: Manage unread notification count state and side effects.
 * Layer: Hooks
 */
export function useUnreadNotificationsCount() {
  const { user } = useAuth()

  const { data: unreadCount = 0, ...query } = useQuery({
    queryKey: ["notifications", "unread-count"],
    queryFn: NotificationService.unreadCount,
    enabled: !!user,
    staleTime: 1000 * 60, // 1 minute
  })

  return {
    unreadCount,
    ...query,
  }
}
