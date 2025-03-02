"use client"

import React, { useEffect, useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { UserProfile } from "@/components/profile/UserProfile"
import { TeamProfile } from "@/components/profile/TeamProfile"
import { TeamManagement } from "@/components/profile/TeamManagement"
import { CategoryProgress } from "@/components/charts/CategoryProgress"
import { useAuthStore } from "@/store/authStore"
import { CorrectVIncorrect } from "@/components/charts/CorrectVIncorrect"
import { ScoreGraph } from "@/components/charts/ScoreGraph"
import { useProfileStore } from "@/store/profileStore"
import { ChartSkeleton } from "@/components/skeletons/profile"

export default function ProfilePage() {
    const [activeTab, setActiveTab] = useState("user")
    const { user } = useAuthStore()
    const { team, teamGraph, userCategoryProgress, teamCategoryProgress, userSubmissionsProgress, teamSubmissionsProgress, teamLoading, fetchSelfTeam } = useProfileStore()

    useEffect(() => {
        fetchSelfTeam()

        return () => { }
    }, [fetchSelfTeam])

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
                        {teamLoading ? <ChartSkeleton /> :
                            <CorrectVIncorrect correct={userSubmissionsProgress.correct} incorrect={userSubmissionsProgress.incorrect} />}
                        {teamLoading ? <ChartSkeleton /> :
                            <CategoryProgress categories={userCategoryProgress} title="Progress" />}
                    </div>
                </TabsContent>
                <TabsContent value="team" className="space-y-4">
                    <TeamProfile team={team} />
                    <TeamManagement team={team} user={user} />
                    {teamLoading ? <ChartSkeleton /> :
                        <ScoreGraph plots={teamGraph} />}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        {teamLoading ? <ChartSkeleton /> :
                            <CorrectVIncorrect correct={teamSubmissionsProgress.correct} incorrect={teamSubmissionsProgress.incorrect} />}
                        {teamLoading ? <ChartSkeleton /> :
                            <CategoryProgress categories={teamCategoryProgress} title="Progress" />}
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    )
}
