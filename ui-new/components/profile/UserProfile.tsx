import React from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Badge } from "@/components/ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar"
import type { User } from "@/store/profileStore"
import { Flag, UsersRound } from "lucide-react"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"

interface UserProfileProps {
    user: User
    points: number
}

export function UserProfile({ user, points }: UserProfileProps) {
    return (
        <Card className="flex flex-wrap flex-row justify-between">
            <CardHeader className="flex flex-row items-center gap-4">
                <Avatar className="w-16 h-16">
                    <AvatarFallback>
                        <UsersRound size={30} />
                    </AvatarFallback>
                </Avatar>
                <div style={{ marginTop: "0px" }}>
                    <CardTitle className="text-2xl">{user.username}</CardTitle>
                    <p className="text-sm text-muted-foreground truncate">{user.email}</p>
                </div>
            </CardHeader>
            <CardContent className="sm:pt-6 sm:flex">
                <div className="flex items-center">
                    <TooltipProvider>
                        <Tooltip>
                            <TooltipTrigger asChild>
                                <Badge variant="default" className="px-3 py-1 text-lg gap-1">
                                    <Flag size={20} />
                                    {points}
                                </Badge>
                            </TooltipTrigger>
                            <TooltipContent>
                                <p>
                                    {`Points: ${points}`}
                                </p>
                            </TooltipContent>
                        </Tooltip>
                    </TooltipProvider>
                </div>
            </CardContent>
        </Card>
    )
}

