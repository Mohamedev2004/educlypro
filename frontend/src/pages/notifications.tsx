import { RefreshCcw } from "lucide-react"

import AppLayout from "@/layouts/app-layout"
import { useDirection } from "@/context/direction/direction-provider"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardHeader } from "@/components/ui/card"
import { Tabs, TabsContent } from "@/components/ui/tabs"
import { NotificationList } from "@/components/notifications/notification-list"
import { NotificationsHero } from "@/components/notifications/notifications-hero"
import { NotificationsPagination } from "@/components/notifications/notifications-pagination"
import { NotificationsToolbar } from "@/components/notifications/notifications-toolbar"
import { NotificationsListSkeleton } from "@/components/notifications/notifications-list-skeleton"
import type { NotificationsPageProps } from "@/types/notifications"
import { useNotificationsFilter } from "@/hooks/notifications/use-notifications-filter"
import { useNotificationsList } from "@/hooks/notifications/use-notifications-list"
import { useMarkNotificationRead } from "@/hooks/notifications/use-mark-notification-read"
import { useMarkAllNotificationsRead } from "@/hooks/notifications/use-mark-all-notifications-read"
import { useDeleteNotification } from "@/hooks/notifications/use-delete-notification"
import { useDeleteAllReadNotifications } from "@/hooks/notifications/use-delete-all-read-notifications"
import { useNotificationsErrors } from "@/hooks/notifications/use-notifications-errors"
import { useUnreadNotificationsCount } from "@/hooks/notifications/use-unread-notifications-count"

export default function NotificationsPage({
  dashboardHref = "/super-admin/dashboard",
  breadcrumbTaglineKey = "roles.overview",
  breadcrumbLabelKey = "roles.notifications",
}: NotificationsPageProps) {
  const { t, locale } = useDirection()

  const {
    activeTab,
    page,
    perPage,
    setPage,
    handleTabChange,
    handlePerPageChange,
    setActiveTab,
  } = useNotificationsFilter()

  const { unreadCount: globalUnreadCount } = useUnreadNotificationsCount()
  const notificationsQuery = useNotificationsList(activeTab, page, perPage)
  const { data } = notificationsQuery

  const markReadMutation = useMarkNotificationRead(t)
  const markAllReadMutation = useMarkAllNotificationsRead(t, setActiveTab, setPage)
  const deleteNotificationMutation = useDeleteNotification(t)
  const deleteAllReadMutation = useDeleteAllReadNotifications(t)

  const { pageError, perPageError, filterError, generalError } =
    useNotificationsErrors(t, notificationsQuery.error)

  const isMutating =
    markReadMutation.isPending ||
    markAllReadMutation.isPending ||
    deleteNotificationMutation.isPending ||
    deleteAllReadMutation.isPending

  const items = data?.items ?? []
  const counts = data?.counts
    ? {
        ...data.counts,
        unread: globalUnreadCount,
      }
    : undefined
  const pagination = data?.pagination
  const showSkeleton =
    notificationsQuery.isFetching || notificationsQuery.isLoading

  return (
    <AppLayout
      breadcrumbs={[
        { label: t(breadcrumbTaglineKey), href: dashboardHref },
        { label: t(breadcrumbLabelKey) },
      ]}
    >
      <div className="flex flex-1 flex-col gap-4">
        <NotificationsHero
          unreadCount={counts?.unread ?? 0}
          readCount={counts?.read ?? 0}
          t={t}
          onTabChange={(tab) => {
            setActiveTab(tab)
            setPage(1) // reset to first page when changing tabs
          }}
        />

        <div>
          <Card className="overflow-visible backdrop-blur-sm">
            <CardHeader className="gap-4 border-b border-border/70">
              <NotificationsToolbar
                activeTab={activeTab}
                counts={counts}
                perPage={perPage}
                isMutating={isMutating}
                t={t}
                onTabChange={handleTabChange}
                onPerPageChange={handlePerPageChange}
                onMarkAllRead={() => markAllReadMutation.mutate()}
                onDeleteAllRead={() => deleteAllReadMutation.mutate()}
              />

              {pageError && (
                <div className="text-sm font-medium text-destructive">
                  {pageError}
                </div>
              )}
              {perPageError && (
                <div className="text-sm font-medium text-destructive">
                  {perPageError}
                </div>
              )}
              {filterError && (
                <div className="text-sm font-medium text-destructive">
                  {filterError}
                </div>
              )}
              {generalError && (
                <div className="text-sm font-medium text-destructive">
                  {generalError}
                </div>
              )}
            </CardHeader>

            <CardContent className="space-y-6 pt-6">
              {showSkeleton ? (
                <NotificationsListSkeleton rows={perPage} />
              ) : (
                <div key={activeTab}>
                  <Tabs value={activeTab} onValueChange={handleTabChange}>
                    <TabsContent value="unread" className="space-y-4">
                      <div className="flex flex-wrap items-center justify-between gap-3">
                        <div>
                          <h2 className="font-medium">
                            {t("notifications.unreadTab")}
                          </h2>
                          <p className="text-sm text-muted-foreground">
                            {t("notifications.unreadDescription")}
                          </p>
                        </div>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => notificationsQuery.refetch()}
                          disabled={notificationsQuery.isFetching}
                        >
                          <RefreshCcw
                            className={
                              notificationsQuery.isFetching
                                ? "animate-spin"
                                : ""
                            }
                          />
                          {t("notifications.tryAgain")}
                        </Button>
                      </div>

                      <NotificationList
                        items={activeTab === "unread" ? items : []}
                        filter="unread"
                        locale={locale}
                        t={t}
                        isMutating={isMutating}
                        onMarkRead={(id) => markReadMutation.mutate(id)}
                        onDelete={(id) => deleteNotificationMutation.mutate(id)}
                      />
                    </TabsContent>

                    <TabsContent value="read" className="space-y-4">
                      <div>
                        <h2 className="font-medium">
                          {t("notifications.readTab")}
                        </h2>
                        <p className="text-sm text-muted-foreground">
                          {t("notifications.readDescription")}
                        </p>
                      </div>

                      <NotificationList
                        items={activeTab === "read" ? items : []}
                        filter="read"
                        locale={locale}
                        t={t}
                        isMutating={isMutating}
                        onMarkRead={(id) => markReadMutation.mutate(id)}
                        onDelete={(id) => deleteNotificationMutation.mutate(id)}
                      />
                    </TabsContent>
                  </Tabs>
                </div>
              )}

              {pagination && (
                <NotificationsPagination
                  pagination={pagination}
                  t={t}
                  onPageChange={setPage}
                  disabled={showSkeleton}
                />
              )}
            </CardContent>
          </Card>
        </div>
      </div>
    </AppLayout>
  )
}

