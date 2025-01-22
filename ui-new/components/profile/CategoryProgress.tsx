import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Progress } from "@/components/ui/progress"
import type { CategoryProgress as CategoryProgressType } from "@/utils/types"

interface CategoryProgressProps {
    categories: CategoryProgressType[]
    title: string
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
                            <span>{category.category}</span>
                            <span>
                                {category.solved}/{category.total}
                            </span>
                        </div>
                        <Progress value={(category.solved / category.total) * 100} />
                    </div>
                ))}
            </CardContent>
        </Card>
    )
}

