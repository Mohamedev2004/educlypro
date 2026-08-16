/**
 * Notification related types.
 * 
 * Responsibility: Define UI-specific types for notifications.
 * Layer: Types
 */

export type NotificationLocale = "ar" | "en" | "fr"

export type Translate = (key: string, fallback?: string) => string

export type NotificationsPageProps = {
  dashboardHref?: string
  breadcrumbTaglineKey?: string
  breadcrumbLabelKey?: string
}
