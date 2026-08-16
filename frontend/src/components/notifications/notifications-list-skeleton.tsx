import { Skeleton } from "@/components/ui/skeleton"

type NotificationsListSkeletonProps = {
  rows?: number
}

export function NotificationsListSkeleton({
  rows = 6,
}: NotificationsListSkeletonProps) {
  return (
    <div className="space-y-3">
      {Array.from({ length: rows }).map((_, idx) => (
        <div
          key={idx}
          className="rounded-md border border-border/70 bg-background/90 p-4 shadow-sm"
        >
          <div className="flex flex-col gap-4 lg:flex-row lg:items-start lg:justify-between">
            <div className="space-y-3">
              <div className="flex flex-wrap items-center gap-2">
                <Skeleton className="h-5 w-20 rounded-full" />
              </div>

              <div className="space-y-2">
                <Skeleton className="h-5 w-[280px] max-w-full" />
                <Skeleton className="h-4 w-[520px] max-w-full" />
                <Skeleton className="h-4 w-[460px] max-w-full" />
              </div>

              <Skeleton className="h-3.5 w-40" />
            </div>

            <div className="flex shrink-0 items-center gap-2">
              <Skeleton className="h-9 w-28 rounded-md" />
              <Skeleton className="h-9 w-28 rounded-md" />
            </div>
          </div>
        </div>
      ))}
    </div>
  )
}

