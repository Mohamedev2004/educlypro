import type { NotificationPagination } from "@/api/types/notification.types"
import { AppPagination } from "@/components/app-pagination"
import type { Translate } from "@/types/notifications"
import { interpolate } from "@/utils/common-utils"

type NotificationsPaginationProps = {
  pagination: NotificationPagination
  t: Translate
  onPageChange: (page: number) => void
  disabled?: boolean
}

export function NotificationsPagination({
  pagination,
  t,
  onPageChange,
  disabled = false,
}: NotificationsPaginationProps) {
  // const pageItems = buildPageItems(pagination.page, pagination.total_pages)

  return (
    <div className={disabled ? "pointer-events-none opacity-60" : ""}>
      <AppPagination
        pagination={pagination}
        summaryTop={interpolate(t("notifications.showing"), {
          count: pagination.total,
        })}
        summaryBottom={interpolate(t("notifications.pageSummary"), {
          page: pagination.page,
          totalPages: pagination.total_pages,
        })}
        onPageChange={onPageChange}
      />
    </div>
  )
}
