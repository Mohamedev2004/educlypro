export type NotificationFilter = "read" | "unread"

export type NotificationItem = {
  id: number
  request_id: string
  user_id: number
  topic: string
  title: string
  body: string
  payload: Record<string, unknown> | null
  channel: string
  is_read: boolean
  read_at: string | null
  created_at: string
}

export type NotificationCounts = {
  all: number
  read: number
  unread: number
}

export type NotificationPagination = {
  page: number
  per_page: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

export type NotificationsData = {
  items: NotificationItem[]
  filter: NotificationFilter
  counts: NotificationCounts
  pagination: NotificationPagination
}

export type NotificationListParams = {
  page: number
  perPage: number
  filter: NotificationFilter
}

export type ApiEnvelope<T> = {
  success: boolean
  message: string
  data: T
}
