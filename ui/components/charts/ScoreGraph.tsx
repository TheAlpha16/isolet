"use client";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  ChartConfig,
  ChartContainer,
  ChartLegend,
  ChartLegendContent,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart";
import { ChartData } from "@/models/score";
import { chartColor } from "@/utils/chartColors";
import { fromUnixSeconds } from "@/utils/helpers";
import { useMemo } from "react";
import { CartesianGrid, Line, LineChart, XAxis } from "recharts";

type ScoreGraphProps = {
  data: ChartData;
};

export function ScoreGraph({ data }: ScoreGraphProps) {
  const { labels, points } = data;

  const chartConfig = useMemo(
    () =>
      labels.reduce(
        (acc, label, index) => ({
          ...acc,
          [label]: {
            label,
            color: chartColor(index, labels.length),
          },
        }),
        {}
      ),
    [labels]
  ) satisfies ChartConfig;

  return (
    <Card className="hidden md:block">
      <CardHeader className="flex items-center gap-2 space-y-0 border-b py-5 sm:flex-row">
        <div className="grid flex-1 gap-1 text-center sm:text-left">
          <CardTitle>Scores</CardTitle>
          <CardDescription>Showing total points for each team over time</CardDescription>
        </div>
      </CardHeader>
      <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
        <ChartContainer config={chartConfig} className="aspect-auto h-[250px] w-full">
          <LineChart data={points}>
            <CartesianGrid vertical={false} />
            <XAxis dataKey="timestamp" tickLine={false} axisLine={false} tickFormatter={() => ""} />
            <ChartTooltip
              cursor={false}
              content={
                <ChartTooltipContent
                  labelFormatter={(value, payload) => {
                    // Get timestamp from the payload data
                    const timestamp = payload?.[0]?.payload?.timestamp;
                    if (!timestamp) return "";

                    const current_date = fromUnixSeconds(timestamp);
                    return current_date.toLocaleString("en-US", {
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
                stroke={chartColor(index, labels.length)}
                dot={false}
              />
            ))}
            <ChartLegend content={<ChartLegendContent />} />
          </LineChart>
        </ChartContainer>
      </CardContent>
    </Card>
  );
}
