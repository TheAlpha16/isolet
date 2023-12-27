'use client'

import { useEffect, useState } from "react"
import Cookies from "js-cookie"
import { useRouter } from "next/navigation"
import User from "@/components/User"

function Scoreboard() {
    const { loggedin, respHook } = User()
    const router = useRouter()
    const [ scores, setScores ] = useState([{"username": "", "score": 0}])

    const getScores = async () => {
        const request = await fetch("/api/scoreboard", {
            headers: {
                "Authorization": `Bearer ${Cookies.get("token")}`
            }
        })

        const status = await request.status
        if (status != 200) {
            router.push("/logout")
        }

        const scoreJSON = await request.json()
        setScores(scoreJSON)
    }

    useEffect(() => {
        if (respHook && !loggedin) {
            router.push("/login")
        } else {
            getScores()
        }
    }, [respHook])

    return (
        <div className="flex flex-col w-full items-center p-2 pt-5">
            <div className="flex p-3 border-b w-11/12 sm:w-1/2 bg-palette-600 text-2xl justify-between font-medium">
                <div>
                    Rank
                </div>
                <div>
                    Player
                </div>
                <div>
                    Score
                </div>
            </div>
            {   
                scores.map((item, index) =>
                    <div key={ index } className="flex p-3 border-b w-11/12 sm:w-1/2 bg-palette-600 text-xl justify-between font-Roboto">
                        <div>
                            #{ index + 1 }
                        </div>
                        <div>
                            { item.username }
                        </div>
                        <div>
                            { item.score }
                        </div>
                    </div>
                )
            }
        </div>
    )
}

export default Scoreboard;