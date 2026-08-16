import { axiosInstance } from "../axios"
import type {
  ApiEnvelope,
  NotificationListParams,
  NotificationsData,
} from "../types/notification.types"

export const NotificationService = {
  async list(params: NotificationListParams): Promise<NotificationsData> {
    const response = await axiosInstance.get<ApiEnvelope<NotificationsData>>(
      "/notifications",
      {
        params: {
          page: params.page,
          per_page: params.perPage,
          filter: params.filter,
        },
      }
    )

    return response.data.data
  },

  async markRead(id: number): Promise<void> {
    await axiosInstance.patch(`/notifications/${id}/read`)
  },

  async markAllRead(): Promise<void> {
    await axiosInstance.patch("/notifications/read-all")
  },

  async deleteAllRead(): Promise<void> {
    await axiosInstance.delete("/notifications/read")
  },

  async delete(id: number): Promise<void> {
    await axiosInstance.delete(`/notifications/${id}`)
  },

  async unreadCount(): Promise<number> {
    const response = await axiosInstance.get<ApiEnvelope<{ count: number }>>(
      "/notifications/unread-count"
    )

    return response.data.data.count
  },
}
