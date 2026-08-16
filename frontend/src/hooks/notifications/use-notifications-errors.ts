import { normalizeApiError, getApiMessage, getNotificationFieldMessage } from "@/utils/error-utils"

export function useNotificationsErrors(t: (key: string) => string, error: unknown) {
  const apiError = error ? normalizeApiError(error) : null

  const pageError = apiError
    ? getNotificationFieldMessage(t, "page", apiError.errors?.page)
    : undefined
  const perPageError = apiError
    ? getNotificationFieldMessage(t, "per_page", apiError.errors?.per_page)
    : undefined
  const filterError = apiError
    ? getNotificationFieldMessage(t, "filter", apiError.errors?.filter)
    : undefined
  const generalError =
    apiError && !pageError && !perPageError && !filterError
      ? getApiMessage(t, apiError)
      : undefined

  return {
    pageError,
    perPageError,
    filterError,
    generalError,
  }
}
