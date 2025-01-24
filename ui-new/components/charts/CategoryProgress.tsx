"use client"

import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import type { CategoryProgress as CategoryProgressType } from "@/utils/types"

interface CategoryProgressProps {
    categories: CategoryProgressType[]
    title: string
}

const categoryColors = {
    Web: "bg-blue-500",
    Crypto: "bg-green-500",
    Pwn: "bg-red-500",
    Reverse: "bg-purple-500",
    Forensics: "bg-yellow-500",
}

export function CategoryProgress({ categories, title }: CategoryProgressProps) {
    return (
        <Card>
            <CardHeader>
                <CardTitle>{title}</CardTitle>
            </CardHeader>
            <CardContent>
                {categories.map((category) => (
                    <div key={category.category} className="mb-4">
                        <div className="flex justify-between items-center mb-1">
                            <span className="font-medium">{category.category}</span>
                            <span className="text-sm text-muted-foreground">
                                {category.solved}/{category.total}
                            </span>
                        </div>
                        <Progress
                            value={(category.solved / category.total) * 100}
                            className={`${categoryColors[category.category as keyof typeof categoryColors] || "bg-gray-500"}`}
                        />
                    </div>
                ))}
            </CardContent>
        </Card>
    )
}

