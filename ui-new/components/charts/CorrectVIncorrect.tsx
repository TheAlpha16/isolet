"use client"

import * as React from "react"
import { Label, Pie, PieChart } from "recharts"

import {
    Card,
    CardContent,
    CardDescription,
    CardFooter,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"
import {
    ChartConfig,
    ChartContainer,
    ChartTooltip,
    ChartTooltipContent,
} from "@/components/ui/chart"

const chartConfig = {
    subCount: {
        label: "Submissions",
    },
    correct: {
        label: "Correct",
        color: "hsl(var(--chart-2))",
    },
    incorrect: {
        label: "Incorrect",
        color: "hsl(var(--chart-5))",
    },
} satisfies ChartConfig

export function CorrectVIncorrect({ correct, incorrect }: { correct: number; incorrect: number }) {
    const totalSubmissions = correct + incorrect

    const submissionData = totalSubmissions > 0 ? [
        { subType: "correct", subCount: correct, fill: "var(--color-correct)" },
        { subType: "incorrect", subCount: incorrect, fill: "var(--color-incorrect)" },
    ] : [
        { subType: "", subCount: 1, fill: "hsl(var(--muted))" },
    ]

    return (
        <Card className="flex flex-col">
            <CardHeader className="items-center pb-0">
                <CardTitle>Submissions</CardTitle>
                <CardDescription>Correct vs Incorrect</CardDescription>
            </CardHeader>
            <CardContent className="flex-1 pb-0">
                <ChartContainer
                    config={chartConfig}
                    className="mx-auto aspect-square max-h-[250px]"
                >
                    <PieChart>
                        <ChartTooltip
                            cursor={false}
                            content={<ChartTooltipContent hideLabel />}
                        />
                        <Pie
                            data={submissionData}
                            dataKey="subCount"
                            nameKey="subType"
                            innerRadius={60}
                            strokeWidth={5}
                        >
                            <Label
                                content={({ viewBox }) => {
                                    if (viewBox && "cx" in viewBox && "cy" in viewBox) {
                                        return (
                                            <text
                                                x={viewBox.cx}
                                                y={viewBox.cy}
                                                textAnchor="middle"
                                                dominantBaseline="middle"
                                            >
                                                <tspan
                                                    x={viewBox.cx}
                                                    y={viewBox.cy}
                                                    className="fill-foreground text-3xl font-bold"
                                                >
                                                    {totalSubmissions.toLocaleString()}
                                                </tspan>
                                                <tspan
                                                    x={viewBox.cx}
                                                    y={(viewBox.cy || 0) + 24}
                                                    className="fill-muted-foreground"
                                                >
                                                    Submissions
                                                </tspan>
                                            </text>
                                        )
                                    }
                                }}
                            />
                        </Pie>
                    </PieChart>
                </ChartContainer>
            </CardContent>
            <CardFooter className="flex-col gap-2 text-sm">
                {totalSubmissions > 0 && (
                    <div className="flex items-center gap-2 font-medium leading-none">
                        Wow! {correct * 100 / totalSubmissions}% correct submissions
                    </div>)}
                <div className="leading-none text-muted-foreground">
                    {totalSubmissions > 0 ? "Showing total submissions till now" : "No submissions yet"}
                </div>
            </CardFooter>
        </Card>
    )
}
