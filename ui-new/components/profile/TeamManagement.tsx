import React, { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import type { TeamType, UserType } from "@/utils/types"
import { UserPlus2, UserMinus2, Crown } from "lucide-react"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"
import { TeamInvite } from "@/components/profile/TeamInvite"
import { useMetadataStore } from "@/store/metadataStore"

interface TeamManagementProps {
    user: UserType
    team: TeamType
}

export function TeamManagement({ team, user }: TeamManagementProps) {
    const { teamLen } = useMetadataStore()
    const [inviteToken, setInviteToken] = useState("")
    const [showInvite, setShowInvite] = useState(false)

    const generateInviteToken = () => {
        const token = Math.random().toString(36).substring(2, 15)
        setInviteToken(token)
    }

    const getRandomPastelColor = (index: number) => {
        const hue = Math.floor(Math.random() * 360)
        return `hsl(${360 / teamLen * index}, 50%, 60%)`
    }

    return (
        <>
            <Card>
                <CardHeader className="flex flex-row items-center justify-between">
                    <CardTitle className="flex items-center gap-2">
                        Members
                    </CardTitle>
                    {user.rank === 2 && (
                        <TooltipProvider>
                            <Tooltip>
                                <TooltipTrigger asChild>
                                    <Button
                                        variant="outline"
                                        size="icon"
                                        onClick={() => setShowInvite(!showInvite)}
                                        className="flex items-center gap-1"
                                    >
                                        <UserPlus2 size={16} />
                                    </Button>
                                </TooltipTrigger>
                                <TooltipContent>
                                    <p>Invite new members</p>
                                </TooltipContent>
                            </Tooltip>
                        </TooltipProvider>
                    )}
                </CardHeader>
                <CardContent>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
                        {team.members.map((member, index) => (
                            <Card key={member.userid} className="bg-card">
                                <CardContent className="p-4">
                                    <div className="flex items-center space-x-4">
                                        <div
                                            className="min-w-10 h-10 rounded-full flex items-center justify-center text-lg font-semibold"
                                            style={{ backgroundColor: getRandomPastelColor(index) }}
                                        >
                                            {member.username.charAt(0).toUpperCase()}
                                        </div>
                                        <div className="flex-grow">
                                            <p className="font-medium">{member.username}</p>
                                            <p className="text-sm text-muted-foreground">{member.email}</p>
                                        </div>
                                        <div className="flex items-center">
                                            {member.rank === 2 && (
                                                <span className="bg-yellow-500/20 dark:bg-yellow-600/30 text-yellow-700 dark:text-yellow-300 text-xs font-medium px-2.5 py-0.5 rounded-full flex items-center gap-1">
                                                    <Crown size={12} className="" />
                                                    <div className="hidden sm:block">Captain</div>
                                                </span>
                                                // ) : (
                                                // user.rank === 2 && (
                                                //     <TooltipProvider>
                                                //         <Tooltip>
                                                //             <TooltipTrigger asChild>
                                                //                 <Button variant="ghost" size="sm" onClick={() => removeMember(member.userid)}>
                                                //                     <UserMinus2 size={16} />
                                                //                 </Button>
                                                //             </TooltipTrigger>
                                                //             <TooltipContent>
                                                //                 <p>Remove member</p>
                                                //             </TooltipContent>
                                                //         </Tooltip>
                                                //     </TooltipProvider>
                                                // )
                                            )}
                                        </div>
                                    </div>
                                </CardContent>
                            </Card>
                        ))}
                    </div>
                </CardContent>
            </Card>
            <TeamInvite
                isOpen={showInvite}
                onClose={() => setShowInvite(false)}
                onGenerate={generateInviteToken}
                token={inviteToken}
            />
        </>
    )
}

