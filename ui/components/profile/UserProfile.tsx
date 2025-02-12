"use client"

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Avatar, AvatarFallback } from "@/components/ui/avatar"
import { Trophy, UsersRound, Stars } from "lucide-react"
import type { TeamType, UserType } from "@/utils/types"

interface UserProfileProps {
    user: UserType
    team: TeamType
}

export function UserProfile({ user, team }: UserProfileProps) {
    return (
        <Card className="flex flex-col">
            <CardHeader className="flex flex-row items-center gap-4">
                <Avatar className="h-12 w-12">
                    <AvatarFallback className="text-lg">{user.username.slice(0, 2).toUpperCase()}</AvatarFallback>
                </Avatar>
                <div>
                    <CardTitle>{user.username}</CardTitle>
                    <p className="text-sm text-muted-foreground">{user.email}</p>
                </div>
            </CardHeader>
            <CardContent className="sm:flex w-full">
                <div className="flex flex-wrap justify-between items-center w-full gap-4">
                    <div className="flex items-center gap-2">
                        <UsersRound className="h-6 w-6 text-green-500" />
                        <div>
                            <p className="text-sm font-medium">Team</p>
                            <p className="text-3xl font-bold">{team.teamname}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-2">
                        <Stars className="h-6 w-6 text-blue-500" />
                        <div>
                            <p className="text-sm font-medium">Score</p>
                            <p className="text-3xl font-bold">{user.score}</p>
                        </div>
                    </div>
                    <div className="flex items-center gap-2">
                        <Trophy className="h-6 w-6 text-yellow-500" />
                        <div>
                            <p className="text-sm font-medium">Rank</p>
                            <p className="text-3xl font-bold">{team.rank}</p>
                        </div>
                    </div>
                </div>
            </CardContent>
        </Card>
    )
}

