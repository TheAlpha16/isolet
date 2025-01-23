"use client";

import React, { useState, useEffect } from "react"
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts"
import { Input } from "@/components/ui/input"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Trophy, ChevronLeft, ChevronRight } from "lucide-react"
import { format, isSameDay } from "date-fns"
import { Button } from "@/components/ui/button";
import { useScoreboardStore } from "@/store/scoreboardStore";
import { processTopScores } from "@/utils/processTopScores";
import type { TeamType } from "@/utils/types"

export default function Scoreboard() {
    const [searchQuery, setSearchQuery] = useState("")
    const [graphData, setGraphData] = useState<{
        timestamp: string;
        [key: string]: number | string;
    }[]>([])

    const { scores, totalPages, currentPage, loading, fetchPage, graphLoading, topScores, startTime, fetchTopScores } = useScoreboardStore()

    useEffect(() => {
        fetchPage(currentPage)
        fetchTopScores()

        return () => { }
    }, [])

    useEffect(() => {
        const respon = processTopScores(topScores, startTime);
        setGraphData(respon);
    }, [topScores, startTime])

    const formatXAxis = (tickItem: string, index: number) => {
        try {
            const date = new Date(tickItem)
            const prevDate = index > 0 ? new Date(graphData[index - 1].timestamp) : null
            const time = format(date, "HH:mm")

            if (!prevDate || !isSameDay(date, prevDate)) {
                return `${format(date, "MMM d")}`
            }

            return time
        } catch (e) {
            return ""
        }
    }

    const CustomTooltip = ({ active, payload, label }: any) => {
        if (active && payload && payload.length) {
            return (
                <div className="bg-background p-4 border rounded shadow">
                    <p className="font-bold">{(() => {
                        try {
                            return format(new Date(label), "MMM d, HH:mm")
                        } catch (e) {
                            return ""
                        }
                    })()}</p>
                    {payload.map((entry: any, index: number) => (
                        <p key={index} style={{ color: entry.color }}>
                            {entry.name}: {entry.value}
                        </p>
                    ))}
                </div>
            )
        }
        return null
    }

    const handlePageChange = async (newPage: number) => {
        if (newPage < 1 || newPage > totalPages) return
        await fetchPage(newPage)
    }

    const PageNavigation = () => (
        <div className="flex items-center justify-center gap-1">
            <Button
                onClick={() => handlePageChange(currentPage - 1)}
                disabled={currentPage === 1}
                className="p-2 rounded-full disabled:opacity-50 disabled:cursor-not-allowed"
                variant={"ghost"}
            >
                <ChevronLeft className="w-6 h-6" />
            </Button>
            <span className="text-sm text-gray-500">
                Page {currentPage} of {totalPages}
            </span>
            <Button
                onClick={() => handlePageChange(currentPage + 1)}
                disabled={currentPage === totalPages}
                className="p-2 rounded-full disabled:opacity-50 disabled:cursor-not-allowed"
                variant={"ghost"}
            >
                <ChevronRight className="w-6 h-6" />
            </Button>
        </div>
    )

    return (
        <div className="container mx-auto p-4 space-y-4">
            <h1 className="text-3xl font-bold mb-6">Scoreboard</h1>
            {!graphLoading && topScores.length !== 0 && (
                <Card>
                    <CardContent className="pt-6">

                        <ResponsiveContainer width="100%" height={400}>
                            <LineChart data={graphData} margin={{ top: 5, right: 30, left: 10, bottom: 5 }}>
                                <CartesianGrid strokeDasharray="3 3" />
                                <XAxis dataKey="timestamp" tickFormatter={formatXAxis} height={60} tick={{ fontSize: 12 }} />
                                <YAxis tick={{ fontSize: 12 }} width={40} />
                                <Tooltip content={<CustomTooltip />} />
                                <Legend />
                                {topScores.map((team, index) => (
                                    <Line
                                        key={team.teamname}
                                        type="monotone"
                                        dataKey={team.teamname}
                                        stroke={`hsl(${index * 36}, 70%, 50%)`}
                                        strokeWidth={2}
                                        dot={false}
                                        connectNulls={true}
                                    />
                                ))}
                            </LineChart>
                        </ResponsiveContainer>
                    </CardContent>
                </Card>)}

            <div className="flex flex-col sm:flex-row justify-between items-center space-y-4 sm:space-y-0">
                <Input
                    type="text"
                    placeholder="Search teams..."
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    className="max-w-sm"
                />
                <PageNavigation />
            </div>

            <Card>
                <CardContent className="p-0">
                    <Table>
                        <TableHeader>
                            <TableRow>
                                <TableHead className="w-24 text-center">Rank</TableHead>
                                <TableHead className="text-center">Team Name</TableHead>
                                <TableHead className="w-24 text-center">Score</TableHead>
                            </TableRow>
                        </TableHeader>
                        <TableBody>
                            {scores.map((team: TeamType) => (
                                <TableRow key={team.teamid}>
                                    <TableCell className="text-center">
                                        <div className="flex justify-center items-center">
                                            {team.rank <= 3 ? (
                                                <Trophy
                                                    className={`w-6 h-6 ${team.rank === 1 ? "text-yellow-500" : team.rank === 2 ? "text-gray-400" : "text-orange-500"
                                                        }`}
                                                />
                                            ) : (
                                                <Badge variant="secondary" className="w-8 flex justify-center">
                                                    #{team.rank}
                                                </Badge>
                                            )}
                                        </div>
                                    </TableCell>
                                    <TableCell className="font-medium text-center">
                                        <span className="truncate block max-w-xs mx-auto">{team.teamname}</span>
                                    </TableCell>
                                    <TableCell className="text-center">{team.score}</TableCell>
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </CardContent>
            </Card>

            <PageNavigation />
        </div>
    )
}

