export const logs = {
  title: "Activity logs",
  description:
    "Browse and inspect all system activity records. Filter by severity level, response status, or execution duration — or search by message. Adjust rows per page and paginate through results to audit specific time ranges or trace issues.",
  shownOfTotal: "{shown} of {total} logs",
  searchPlaceholder: "Search logs by message...",
  perPageOption: "{count} / page",
  chart: {
    title: "Logs over time",
    subtitle: "Trend of incoming logs",
    count: "Logs",
    rangeAria: "Select time range",
    range24h: "24h",
    range7d: "7d",
  },
  filters: {
    button: "Filter",
    title: "Filters",
    clear: "Clear",
    level: "Level",
    status: "Status",
    statusCode: "HTTP status",
    duration: "Duration",
  },
  empty: {
    title: "No logs found",
    description:
      "No logs match your current filters or search. Try adjusting or resetting them.",
    reset: "Reset filters",
  },
  row: {
    message: "Message",
    duration: "Duration",
    timestamp: "Timestamp",
  },
  cards: {
    total: "Total logs",
    allTime: "all time",
    ofTotal: "of total",
  },
  pagination: {
    showing: "Showing {count} logs",
    pageSummary: "Page {page} of {totalPages}",
  },
  export: "Export",
  exporting: "Exporting...",
  exportFailed: "Export failed. Please try again.",
  success: "Logs exported.",
  successDescription: "Your logs export is being prepared. You will receive a file download shortly.",
}
