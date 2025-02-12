import React from "react"
import { Skeleton } from "@/components/ui/skeleton"

export function FormSkeleton() {
    return (
        <div className="container flex flex-col items-center justify-center h-full">
            <Skeleton className="rounded-lg w-[350px] h-[450px]" />
        </div>
    )
}
