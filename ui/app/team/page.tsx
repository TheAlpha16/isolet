'use client'

import { useEffect, useState } from 'react'
import Cookies from 'js-cookie'
import { useRouter } from 'next/navigation'
import LoginStatus from "@/components/User";
import { JoinTeam, TeamPage } from '@/components/Team'

function Team() {
    const { loggedin, respHook } = LoginStatus();
    const [teamid, setTeamID] = useState(0);
    const router = useRouter();

    useEffect(() => {
        if (!respHook) {
            router.push("/");
        } else if (respHook && !loggedin) {
            router.push("/login");
        } else {
            var token = Cookies.get("token");
            if (token) {
                const teamid = JSON.parse(atob(token.split(".")[1])).teamid;
                setTeamID(teamid);
            } else {
                router.push("/login");
            }
        }
    }, [respHook]);
    
    return (
        <div className="flex flex-col gap-2">
            {
                teamid == -1 ? (
                    <JoinTeam router={router} />
                ) : (
                        // <TeamPage teamid={teamid} />
                        <div>Team Page</div>
                )
            }
        </div>
    )
}

export default Team