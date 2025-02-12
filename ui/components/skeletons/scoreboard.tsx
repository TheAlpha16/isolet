import { Skeleton } from "@/components/ui/skeleton";

export function ScoreGraphSkeleton() {
    return (
        <Skeleton className="w-full rounded-lg h-[380px]" />
    )
}

export function SearchSkeleton() {
    return (
        <div className="flex flex-col sm:flex-row justify-between items-center space-y-4 sm:space-y-0">
            <Skeleton className="w-80 min-h-[40px]" />
            <Skeleton className="w-[150px] min-h-[40px]" />
        </div>
    )
}

export function ScoresTableSkeleton() {
    return (
        <Skeleton className="flex w-full rounded-lg min-h-[200px]" />
    )
}

export function ScoreboardSkeleton() {
    return (
        <div className="container mx-auto p-4 space-y-4">
            <ScoreGraphSkeleton />
            <SearchSkeleton />
            <ScoresTableSkeleton />
        </div>
    )
}
