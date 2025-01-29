"use client"

import { Skeleton } from "@/components/ui/skeleton";
import React from "react";

export default function ChallengeSkeleton() {
    return (
        <div className="container p-4 justify-start h-full flex flex-col">

            <Skeleton className="mb-4 h-10 rounded-md" />

            <div className="flex flex-col w-full items-center sm:items-start flex-wrap gap-4 sm:flex-row">
                <Skeleton className="h-[200px] w-80 rounded-lg" />
                <Skeleton className="h-[200px] w-80 rounded-lg" />
                <Skeleton className="h-[200px] w-80 rounded-lg" />
                <Skeleton className="h-[200px] w-80 rounded-lg" />
                <Skeleton className="h-[200px] w-80 rounded-lg" />
                <Skeleton className="h-[200px] w-80 rounded-lg" />
            </div>
        </div>
    );
}