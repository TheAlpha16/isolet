"use client";

import React, { useState, useEffect } from "react"
import { generateFakeData, getPagedTeams, searchTeams, type Team, type ScoreHistory } from "../../utils/fakeData"
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from "recharts"
import { Input } from "@/components/ui/input"
import { Card, CardContent } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table"
import { Trophy, ChevronLeft, ChevronRight } from "lucide-react"
import { format, isSameDay } from "date-fns"
import { Button } from "@/components/ui/button";

const PAGE_SIZE = 10

export default function Scoreboard() {
    const [teams, setTeams] = useState<Team[]>([])
    const [scoreHistory, setScoreHistory] = useState<ScoreHistory[]>([])
    const [currentPage, setCurrentPage] = useState(1)
    const [searchQuery, setSearchQuery] = useState("")
    const [filteredTeams, setFilteredTeams] = useState<Team[]>([])

    useEffect(() => {
        const { teams, scoreHistory } = generateFakeData(50, 3)
        setTeams(teams)
        setScoreHistory(scoreHistory)
        setFilteredTeams(teams)
    }, [])

    useEffect(() => {
        const filtered = searchTeams(teams, searchQuery)
        setFilteredTeams(filtered)
        setCurrentPage(1)
    }, [searchQuery, teams])

    const pageCount = Math.ceil(filteredTeams.length / PAGE_SIZE)
    const pagedTeams = getPagedTeams(filteredTeams, currentPage, PAGE_SIZE)

    const top10Teams = teams.slice(0, 10)
    const graphData = scoreHistory.map((entry) => {
        const graphEntry: { [key: string]: any } = { timestamp: entry.timestamp }
        top10Teams.forEach((team) => {
            graphEntry[team.name] = entry[team.name]
        })
        return graphEntry
    })

    const formatXAxis = (tickItem: string, index: number) => {
        const date = new Date(tickItem)
        const prevDate = index > 0 ? new Date(graphData[index - 1].timestamp) : null
        const time = format(date, "HH:mm")

        if (!prevDate || !isSameDay(date, prevDate)) {
            return `${format(date, "MMM d")}`
        }
        return time
    }

    const CustomTooltip = ({ active, payload, label }: any) => {
        if (active && payload && payload.length) {
            return (
                <div className="bg-background p-4 border rounded shadow">
                    <p className="font-bold">{format(new Date(label), "MMM d, HH:mm")}</p>
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

    const handlePageChange = (newPage: number) => {
        setCurrentPage(Math.max(1, Math.min(newPage, pageCount)))
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
                Page {currentPage} of {pageCount}
            </span>
            <Button
                onClick={() => handlePageChange(currentPage + 1)}
                disabled={currentPage === pageCount}
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
            <Card>
                <CardContent className="pt-6">
                    <ResponsiveContainer width="100%" height={400}>
                        <LineChart data={graphData} margin={{ top: 5, right: 30, left: 10, bottom: 5 }}>
                            <CartesianGrid strokeDasharray="3 3" />
                            <XAxis dataKey="timestamp" tickFormatter={formatXAxis} height={60} tick={{ fontSize: 12 }} />
                            <YAxis tick={{ fontSize: 12 }} width={40} />
                            <Tooltip content={<CustomTooltip />} />
                            <Legend />
                            {top10Teams.map((team, index) => (
                                <Line
                                    key={team.id}
                                    type="monotone"
                                    dataKey={team.name}
                                    stroke={`hsl(${index * 36}, 70%, 50%)`}
                                    strokeWidth={2}
                                    dot={false}
                                />
                            ))}
                        </LineChart>
                    </ResponsiveContainer>
                </CardContent>
            </Card>

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
                            {pagedTeams.map((team) => (
                                <TableRow key={team.id}>
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
                                        <span className="truncate block max-w-xs mx-auto">{team.name}</span>
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

