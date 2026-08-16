import { useState } from "react"
import type { NotificationFilter } from "@/api/types/notification.types"

export function useNotificationsFilter() {
  const [activeTab, setActiveTab] = useState<NotificationFilter>("unread")
  const [page, setPage] = useState(1)
  const [perPage, setPerPage] = useState<number>(10)

  function handleTabChange(value: string) {
    setActiveTab(value as NotificationFilter)
    setPage(1)
  }

  function handlePerPageChange(value: string) {
    setPerPage(Number(value))
    setPage(1)
  }

  return {
    activeTab,
    page,
    perPage,
    setPage,
    setActiveTab,
    handleTabChange,
    handlePerPageChange,
  }
}
