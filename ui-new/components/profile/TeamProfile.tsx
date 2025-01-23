import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import type { TeamType } from "@/utils/types"
import { Flag, Trophy, UsersRound } from "lucide-react"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"

interface TeamProfileProps {
    team: TeamType;
}

export function TeamProfile({ team }: TeamProfileProps) {
    return (
        <Card className="flex flex-wrap flex-row justify-between">
            <CardHeader className="flex flex-row items-center gap-4">
                <Avatar className="w-16 h-16">
                    <AvatarFallback>
                        <UsersRound size={30} />
                    </AvatarFallback>
                </Avatar>
                <div style={{ marginTop: "0px" }}>
                    <CardTitle className="text-2xl">{team.teamname}</CardTitle>
                    <p className="text-sm text-muted-foreground">{team.members.length} members</p>
                </div>
            </CardHeader>
            <CardContent className="sm:pt-6 sm:flex">
                <div className="flex items-center gap-2">
                    <TooltipProvider>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Badge variant="secondary" className="text-lg px-3 py-1 gap-1">
                                    <Trophy size={20} />
                                    {team.rank}
                                </Badge>
                            </TooltipTrigger>
                            <TooltipContent>
                                <p>
                                    {`Rank: ${team.rank}`}
                                </p>
                            </TooltipContent>
                        </Tooltip>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Badge variant="default" className="text-lg px-3 py-1 gap-1">
                                    <Flag size={20} />
                                    {team.score}
                                </Badge>
                            </TooltipTrigger>
                            <TooltipContent>
                                <p>
                                    {`Score: ${team.score}`}
                                </p>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider >
                </div>
            </CardContent>
        </Card>
    )
}

