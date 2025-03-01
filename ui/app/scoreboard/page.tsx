"use client";

import React, { useState, useEffect } from "react"
import { Input } from "@/components/ui/input"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Button } from "@/components/ui/button";
import { Trophy, ChevronLeft, ChevronRight } from "lucide-react";
import { useScoreboardStore } from "@/store/scoreboardStore";
import { processScores } from "@/utils/processScores";
import type { TeamType, ScoreGraphEntryType } from "@/utils/types"
import { ScoreGraph } from "@/components/charts/ScoreGraph";
import { useMetadataStore } from "@/store/metadataStore";
import { ScoreGraphSkeleton } from "@/components/skeletons/scoreboard";

export default function Scoreboard() {
    const [searchQuery, setSearchQuery] = useState("")
    const [graphData, setGraphData] = useState<ScoreGraphEntryType[]>([])

    const { scores, totalPages, currentPage, fetchPage, graphLoading, topScores, fetchTopScores } = useScoreboardStore()
    const { eventStart } = useMetadataStore()

    useEffect(() => {
        fetchPage(currentPage)
        fetchTopScores()

        return () => { }
    }, [])

    useEffect(() => {
        const toProcess = topScores.map((team) => ({
            label: team.teamname,
            scores: team.submissions.map((sub) => ({
                timestamp: sub.timestamp,
                points: sub.points,
            })),
        }));

        const respon = processScores(toProcess, eventStart.toUTCString());
        setGraphData(respon);
    }, [topScores, eventStart])

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
            {graphLoading ?
                (<ScoreGraphSkeleton />) :
                topScores.length !== 0 && (<ScoreGraph plots={graphData} />)}

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

