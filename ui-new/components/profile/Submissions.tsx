import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { PieChart, Pie, Cell, ResponsiveContainer, Tooltip } from "recharts"
import type { SubmissionType } from "@/utils/types"
import { format } from "date-fns"

interface SubmissionsProps {
    submissions: SubmissionType[]
    title: string
}

export function Submissions({ submissions, title }: SubmissionsProps) {
    const correctSubmissions = submissions.filter((s) => s.correct).length
    const incorrectSubmissions = submissions.length - correctSubmissions

    const pieData = [
        { name: "Correct", value: correctSubmissions },
        { name: "Incorrect", value: incorrectSubmissions },
    ]

    const COLORS = ["#00C49F", "#FF8042"]

    return (
        <Card>
            <CardHeader>
                <CardTitle>{title}</CardTitle>
            </CardHeader>
            <CardContent>
                <div className="h-64">
                    <ResponsiveContainer width="100%" height="100%">
                        <PieChart>
                            <Pie data={pieData} cx="50%" cy="50%" innerRadius={60} outerRadius={80} paddingAngle={5} dataKey="value">
                                {pieData.map((entry, index) => (
                                    <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                ))}
                            </Pie>
                            <Tooltip />
                        </PieChart>
                    </ResponsiveContainer>
                </div>
                <h3 className="text-lg font-semibold mt-4 mb-2">Recent Submissions</h3>
                <ul className="space-y-2">
                    {submissions.slice(0, 5).map((submission) => (
                        <li key={submission.sid} className="flex justify-between items-center">
                            <span>{submission.chall_name}</span>
                            <span>{format(new Date(submission.timestamp), "MMM dd, HH:mm")}</span>
                        </li>
                    ))}
                </ul>
            </CardContent>
        </Card>
    )
}

