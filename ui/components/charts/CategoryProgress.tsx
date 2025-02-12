"use client"

import React, { useMemo } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import type { CategoryProgress as CategoryProgressType } from "@/utils/types"

interface CategoryProgressProps {
    categories: CategoryProgressType[]
    title: string
}

const generateColors = (count: number) => {
    const hueStep = 360 / count
    return Array.from({ length: count }, (_, i) => `hsl(${i * hueStep}, 70%, 50%)`)
}

export function CategoryProgress({ categories, title }: CategoryProgressProps) {
    const colors = useMemo(() => generateColors(categories.length), [categories.length])

    return (
        <Card>
            <CardHeader>
                <CardTitle>{title}</CardTitle>
            </CardHeader>
            <CardContent>
                {categories.map((category, index) => (
                    <div key={category.category} className="mb-4">
                        <div className="flex justify-between items-center mb-1">
                            <span className="font-medium">{category.category}</span>
                            <span className="text-sm text-muted-foreground">
                                {category.solved}/{category.total}
                            </span>
                        </div>
                        <Progress
                            value={(category.solved / category.total) * 100}
                            style={
                                {
                                    "--progress-background": colors[index],
                                } as React.CSSProperties
                            }
                        />
                    </div>
                ))}
            </CardContent>
        </Card>
    )
}

