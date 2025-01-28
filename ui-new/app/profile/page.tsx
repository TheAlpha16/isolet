"use client"

import React, { useEffect, useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { UserProfile } from "@/components/profile/UserProfile"
import { TeamProfile } from "@/components/profile/TeamProfile"
import { TeamManagement } from "@/components/profile/TeamManagement"
import { CategoryProgress } from "@/components/charts/CategoryProgress"
import {
    generateMockTeam,
    generateMockSubmissions,
    generateMockCategoryProgress,
} from "@/utils/mockData"
import { useAuthStore } from "@/store/authStore"
import { CorrectVIncorrect } from "@/components/charts/CorrectVIncorrect"
import { ScoreGraph } from "@/components/charts/ScoreGraph"
import { processScores } from "@/utils/processScores";
import { ScoreGraphEntryType } from "@/utils/types"

const userSubmissions = generateMockSubmissions(20)
const teamSubmissions = generateMockSubmissions(50)
const userCategoryProgress = generateMockCategoryProgress()
const teamCategoryProgress = generateMockCategoryProgress()

export default function ProfilePage() {
    const [activeTab, setActiveTab] = useState("team")
    const { user } = useAuthStore()
    const [plots, setPlots] = useState<ScoreGraphEntryType[]>([{ timestamp: "", value: 0 }])

    // Generate mock data
    const team = generateMockTeam(user.teamid || 1, user.userid || 1)

    useEffect(() => {
        const teamData = {
            label: team.teamname,
            scores: teamSubmissions.filter((sub) => sub.correct === true).map((sub) => ({
                timestamp: sub.timestamp,
                points: sub.points
            }))
        }

        setPlots(processScores([teamData], "2021-08-01T00:00:00.000Z"))
    }, [teamSubmissions])

    return (
        <div className="container mx-auto p-4 space-y-8">
            <Tabs value={activeTab} onValueChange={setActiveTab}>
                <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="user">User</TabsTrigger>
                    <TabsTrigger value="team">Team</TabsTrigger>
                </TabsList>
                <TabsContent value="user" className="space-y-4">
                    <UserProfile user={user} team={team} />
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <CorrectVIncorrect correct={userSubmissions.filter((sub) => sub.correct == true).length} incorrect={userSubmissions.filter((sub) => sub.correct == false).length} />
                        <CategoryProgress categories={userCategoryProgress} title="Challenges" />
                    </div>
                </TabsContent>
                <TabsContent value="team" className="space-y-4">
                    <TeamProfile team={team} />
                    <TeamManagement team={team} user={user} />
                    <ScoreGraph plots={plots} />
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <CorrectVIncorrect correct={teamSubmissions.filter((sub) => sub.correct == true).length} incorrect={teamSubmissions.filter((sub) => sub.correct == false).length} />
                        <CategoryProgress categories={teamCategoryProgress} title="Team Category Progress" />
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    )
}
