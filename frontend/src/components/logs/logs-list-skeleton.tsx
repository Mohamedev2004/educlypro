import { Skeleton } from "@/components/ui/skeleton"

type LogsListSkeletonProps = {
  rows?: number
}

export function LogsListSkeleton({ rows = 10 }: LogsListSkeletonProps) {
  return (
    <div className="divide-y">
      {Array.from({ length: rows }).map((_, idx) => (
        <div key={idx} className="w-full p-4">
          <div className="flex items-center gap-4">
            <Skeleton className="h-4 w-4 rounded-sm" />
            <Skeleton className="h-5 w-16 rounded-full" />
            <Skeleton className="h-4 w-20" />
            <Skeleton className="h-4 w-24" />
            <Skeleton className="h-4 flex-1" />
            <Skeleton className="h-4 w-10" />
            <Skeleton className="h-4 w-14" />
          </div>
        </div>
      ))}
    </div>
  )
}

