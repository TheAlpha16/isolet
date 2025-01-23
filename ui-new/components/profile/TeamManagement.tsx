import React, { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import type { TeamType, UserType } from "@/utils/types"
import { Trash2, Copy } from "lucide-react"

interface TeamManagementProps {
    user: UserType
    team: TeamType
}

export function TeamManagement({ team, user }: TeamManagementProps) {
    const [inviteToken, setInviteToken] = useState("")

    const generateInviteToken = () => {
        const token = Math.random().toString(36).substring(2, 15)
        setInviteToken(token)
    }

    const copyInviteLink = () => {
        navigator.clipboard.writeText(`https://yourctfplatform.com/join-team?token=${inviteToken}`)
        // You might want to show a toast notification here
    }

    const removeMember = (memberId: number) => {
        // Implement member removal logic here
        console.log(`Removing member with ID: ${memberId}`)
    }

    return (
        <Card>
            <CardHeader>
                <CardTitle>Team Management</CardTitle>
            </CardHeader>
            <CardContent>
                <h3 className="text-lg font-semibold mb-2">Team Members</h3>
                <ul className="space-y-2 mb-4">
                    {team.members.map((member) => (
                        <li key={member.userid} className="flex justify-between items-center">
                            <span>
                                {member.username} ({member.rank})
                            </span>
                            {user.rank === 2 && (
                                <Button variant="ghost" size="sm" onClick={() => removeMember(member.userid)}>
                                    <Trash2 size={16} />
                                </Button>
                            )}
                        </li>
                    ))}
                </ul>
                {user.rank === 2 && (
                    <div className="space-y-2">
                        <Button onClick={generateInviteToken}>Generate Invite Token</Button>
                        {inviteToken && (
                            <div className="flex items-center space-x-2">
                                <Input value={inviteToken} readOnly />
                                <Button variant="outline" onClick={copyInviteLink}>
                                    <Copy size={16} />
                                </Button>
                            </div>
                        )}
                    </div>
                )}
            </CardContent>
        </Card>
    )
}

