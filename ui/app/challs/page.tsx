'use client'

import { useState, useEffect, useRef } from "react"
import Cookies from "js-cookie"
import { useRouter } from 'next/navigation'
import { challItem } from "@/components/Challenge"
import Challenge from "@/components/Challenge"
import User from "@/components/User"

interface accessCredsItem {
    userid: number
    level: number
    password: string
    port: string
    verified: boolean
}

interface challList extends Array<challItem>{}
interface activatedItem {
    [level: number]: {password: string; port: string; verified: boolean};
}

function Challenges(){
    const { loggedin, respHook } = User()
    const visibleLevel = useRef("-1")
    const [ challenges, setChallenges ] = useState<challList>([])
    const [ activated, setActivated ] = useState<activatedItem>({ 100000: {
            "password": "fakepasswd",
            "port": "0",
            "verified": false
        }})
    const router = useRouter()

    const getChalls = async () => {
        const request = await fetch("/api/challs", {
            headers: {
                "Authorization": `Bearer ${Cookies.get('token')}`
            }
        })

        const statusCode = await request.status;
        if (statusCode != 200) {
            router.push("/logout");
        }

        const challJson = await request.json();
        setChallenges(challJson);
        const instancesRequest = await fetch("/api/status", {
            headers: {
                "Authorization": `Bearer ${Cookies.get('token')}`
            }
        });

        const instancesStatus = await instancesRequest.status;
        if (instancesStatus != 200) {
            // toast.send()
            router.push("/logout");
        }

        const instancesJSON = await instancesRequest.json();
        let emptyDict = {} as activatedItem

        instancesJSON.map((item: accessCredsItem, index: number) => 
            emptyDict[item["level"]] = {"password": item["password"], "port": item["port"], "verified": item["verified"]}
        )
        setActivated(emptyDict)
    }

    useEffect(() => {
        if (respHook && !loggedin) {
            router.push("/login");
        } else {
            getChalls();
        }
    }, [respHook])

    // const nullFunction = () => {}

    const handleVisibility = (event: any) => {

        let level = event.target.dataset.level
        let newElement = document.getElementById(`submit-${level}`) as HTMLDivElement
        let prevElement = document.getElementById(`submit-${visibleLevel.current}`) as HTMLDivElement

        if (visibleLevel.current == level) {
            return
        }

        if (prevElement != null) {
            (document.getElementById(`level-${visibleLevel.current}`) as HTMLDivElement).classList.remove("border-palette-500")
            prevElement.classList.add("hidden")
        }

        if (newElement != null) {
            (document.getElementById(`level-${level}`) as HTMLDivElement).classList.add("border-palette-500")
            newElement.classList.remove("hidden")
        }
        visibleLevel.current = level
    }

    return (
        <>
            <div className={ `flex flex-col py-2 gap-2` }>
                {   
                    challenges.map((item: challItem, index) =>
                        <Challenge key={ item.level } challObject={ item } isActive={ activated[item.level] != undefined } isVisible={ false } onClick={ handleVisibility } password={ activated[item.level] != undefined ? activated[item.level]["password"]: "" } port={activated[item.level] != undefined ? activated[item.level]["port"]: ""}/>
                    )
                }
            </div>
        </>
    )
}

export default Challenges;