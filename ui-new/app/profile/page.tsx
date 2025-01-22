"use client"

import React, { useState } from "react"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs"
import { UserProfile } from "@/components/profile/UserProfile"
import { TeamProfile } from "@/components/profile/TeamProfile"
import { TeamManagement } from "@/components/profile/TeamManagement"
import { Submissions } from "@/components/profile/Submissions"
import { CategoryProgress } from "@/components/profile/CategoryProgress"
import {
    generateMockUser,
    generateMockTeam,
    generateMockSubmissions,
    generateMockCategoryProgress,
} from "@/utils/mockData"
import { useAuthStore } from "@/store/authStore"

export default function ProfilePage() {
    const [activeTab, setActiveTab] = useState("user")
    const { user } = useAuthStore()

    // Generate mock data
    const team = generateMockTeam("1", "1")
    const userSubmissions = generateMockSubmissions(20)
    const teamSubmissions = generateMockSubmissions(50)
    const userCategoryProgress = generateMockCategoryProgress()
    const teamCategoryProgress = generateMockCategoryProgress()

    const isUserCaptain = false;

    return (
        <div className="container mx-auto p-4 space-y-8">
            <Tabs value={activeTab} onValueChange={setActiveTab}>
                <TabsList className="grid w-full grid-cols-2">
                    <TabsTrigger value="user">User</TabsTrigger>
                    <TabsTrigger value="team">Team</TabsTrigger>
                </TabsList>
                <TabsContent value="user" className="space-y-4">
                    {user && (<UserProfile user={user} points={100} />)}
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Submissions submissions={userSubmissions} title="User Submissions" />
                        <CategoryProgress categories={userCategoryProgress} title="User Category Progress" />
                    </div>
                </TabsContent>
                <TabsContent value="team" className="space-y-4">
                    <TeamProfile team={team} />
                    <TeamManagement team={team} isUserCaptain={isUserCaptain} />
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <Submissions submissions={teamSubmissions} title="Team Submissions" />
                        <CategoryProgress categories={teamCategoryProgress} title="Team Category Progress" />
                    </div>
                </TabsContent>
            </Tabs>
        </div>
    )
}

