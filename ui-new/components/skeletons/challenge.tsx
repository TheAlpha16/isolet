import { Skeleton } from "@/components/ui/skeleton";
import React from "react";

export function ChallengeSkeleton() {
    return (
        <div className="container p-4 justify-start h-full flex flex-col">

            <Skeleton className="mb-4 h-10 rounded-md" />

            <div className="flex flex-col w-full items-center sm:items-start flex-wrap gap-4 sm:flex-row">
                {Array.from({ length: 6 }).map((_, i) => (
                    <Skeleton key={i} className="h-[200px] w-[300px] rounded-md" />
                ))}
            </div>
        </div>
    )
}
