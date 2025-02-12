import { Skeleton } from "@/components/ui/skeleton"

export function TabSkeleton() {
    return <Skeleton className="w-full h-10" />
}

export function ProfileSkeleton() {
    return <Skeleton className="w-full h-[200px]" />
}

export function ChartSkeleton() {
    return <Skeleton className="w-full h-[400px]" />
}

export function ProfilePageSkeleton() {
    return (
        <div className="container mx-auto p-4 space-y-8 w-full h-full">
            <TabSkeleton />
            <ProfileSkeleton />
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <ChartSkeleton />
                <ChartSkeleton />
            </div>
        </div>
    )
}
