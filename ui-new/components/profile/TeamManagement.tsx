import React, { useState } from "react"
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { type Team, TeamMember } from "@/utils/types"
import { Trash2, Copy } from "lucide-react"

interface TeamManagementProps {
    team: Team
    isUserCaptain: boolean
}

export function TeamManagement({ team, isUserCaptain }: TeamManagementProps) {
    const [inviteToken, setInviteToken] = useState("")

    const generateInviteToken = () => {
        const token = Math.random().toString(36).substring(2, 15)
        setInviteToken(token)
    }

    const copyInviteLink = () => {
        navigator.clipboard.writeText(`https://yourctfplatform.com/join-team?token=${inviteToken}`)
        // You might want to show a toast notification here
    }

    const removeMember = (memberId: string) => {
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
                        <li key={member.id} className="flex justify-between items-center">
                            <span>
                                {member.username} ({member.role})
                            </span>
                            {isUserCaptain && member.role !== "captain" && (
                                <Button variant="ghost" size="sm" onClick={() => removeMember(member.id)}>
                                    <Trash2 size={16} />
                                </Button>
                            )}
                        </li>
                    ))}
                </ul>
                {isUserCaptain && (
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

