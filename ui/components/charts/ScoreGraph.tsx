"use client";

import * as React from "react";
import { isSameDay } from "date-fns";
import { Line, LineChart, CartesianGrid, XAxis } from "recharts";
import {
    Card,
    CardContent,
    CardDescription,
    CardHeader,
    CardTitle,
} from "@/components/ui/card";
import {
    ChartConfig,
    ChartContainer,
    ChartLegend,
    ChartLegendContent,
    ChartTooltip,
    ChartTooltipContent,
} from "@/components/ui/chart";
import type { ScoreGraphEntryType } from "@/utils/types";

type ScoreGraphProps = {
    plots: ScoreGraphEntryType[]
};

export function ScoreGraph({ plots }: ScoreGraphProps) {
    const labels = React.useMemo(() => {
        return plots.length > 0
            ? Object.keys(plots[0]).filter((key) => key !== "timestamp")
            : [];
    }, [plots]);

    const chartConfig = React.useMemo(
        () =>
            labels.reduce(
                (acc, label, index) => ({
                    ...acc,
                    [label]: {
                        label,
                        color: `hsl(${index * 36}, 70%, 50%)`,
                    },
                }),
                {}
            ),
        [labels]
    ) satisfies ChartConfig;

    // const formatXAxis = (tickItem: string, index: number) => {
    //     try {
    //         const date = new Date(tickItem);
    //         const prevDate = index > 0 ? new Date(plots[index - 1].timestamp) : null;

    //         if (!prevDate || !isSameDay(date, prevDate)) {
    //             return date.toLocaleDateString("en-US", {
    //                 month: "short",
    //                 day: "numeric",
    //             });
    //         }

    //         return date.toLocaleTimeString("en-US", {
    //             hour: "numeric",
    //             minute: "numeric",
    //         });

    //     } catch {
    //         return ""
    //     }
    // }

    return (
        <Card>
            <CardHeader className="flex items-center gap-2 space-y-0 border-b py-5 sm:flex-row">
                <div className="grid flex-1 gap-1 text-center sm:text-left">
                    <CardTitle>Scores</CardTitle>
                    <CardDescription>
                        Showing total points for each team over time
                    </CardDescription>
                </div>
            </CardHeader>
            <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
                <ChartContainer
                    config={chartConfig}
                    className="aspect-auto h-[250px] w-full"
                >
                    <LineChart data={plots}>
                        <CartesianGrid vertical={false} />
                        <XAxis
                            dataKey="timestamp"
                            tickLine={false}
                            axisLine={false}
                            // tickMargin={8}
                            // minTickGap={32}
                            tickFormatter={() => ""}
                        />
                        <ChartTooltip
                            cursor={false}
                            content={
                                <ChartTooltipContent
                                    labelFormatter={(value) => {
                                        return new Date(value).toLocaleString("en-US", {
                                            month: "short",
                                            day: "numeric",
                                            hour: "numeric",
                                            minute: "numeric",
                                        });
                                    }}
                                    indicator="dot"
                                />
                            }
                        />

                        {labels.map((label, index) => (
                            <Line
                                key={label}
                                dataKey={label}
                                type="monotone"
                                stroke={`hsl(${index * 36}, 70%, 50%)`}
                                dot={false} />
                        ))}
                        <ChartLegend content={<ChartLegendContent />} />
                    </LineChart>
                </ChartContainer>
            </CardContent>
        </Card>
    );
}
