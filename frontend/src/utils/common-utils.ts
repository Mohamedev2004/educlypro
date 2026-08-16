/**
 * Common pure utility functions.
 * 
 * Responsibility: Provide generic helper functions without side effects.
 * Layer: Utils
 */

/**
 * Replaces {key} placeholders in a template string with values from payload.
 * e.g. interpolate("Hello {username}!", { username: "Sara" }) → "Hello Sara!"
 */
export function interpolate(
  template: string,
  payload?: Record<string, unknown>
): string {
  if (!payload) return template
  return template.replace(/\{(\w+)\}/g, (_, key) =>
    key in payload ? String(payload[key]) : `{${key}}`
  )
}

/**
 * Builds an array of page numbers with ellipsis for pagination UI.
 */
export function buildPageItems(currentPage: number, totalPages: number) {
  if (totalPages <= 7) {
    return Array.from({ length: totalPages }, (_, index) => index + 1)
  }

  const pages = new Set([
    1,
    totalPages,
    currentPage - 1,
    currentPage,
    currentPage + 1,
  ])
  const safePages = [...pages]
    .filter((page) => page > 0 && page <= totalPages)
    .sort((left, right) => left - right)

  const result: Array<number | "ellipsis"> = []

  safePages.forEach((page, index) => {
    const previousPage = safePages[index - 1]
    if (previousPage && page - previousPage > 1) {
      result.push("ellipsis")
    }
    result.push(page)
  })

  return result
}
