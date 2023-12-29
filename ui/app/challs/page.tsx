'use client'

import { useState, useEffect, useRef } from "react"
import Cookies from "js-cookie"
import { useRouter } from 'next/navigation'
import { challItem } from "@/components/Challenge"
import Challenge from "@/components/Challenge"
import User from "@/components/User"
import { toast } from "react-toastify"

interface accessCredsItem {
    userid: number
    level: number
    password: string
    port: string
    verified: boolean
    hostname: string
}

interface challList extends Array<challItem>{}
interface activatedItem {
    [level: number]: {password: string; port: string; verified: boolean; hostname: string};
}

function Challenges(){
    const { loggedin, respHook } = User()
    const visibleLevel = useRef("-1")
    const [ challenges, setChallenges ] = useState<challList>([])
    const [ activated, setActivated ] = useState<activatedItem>({ 100000: {
            "password": "fakepasswd",
            "port": "0",
            "verified": false,
            "hostname": ""
        }})
    const router = useRouter()

    const show = (status: string, message: string) => {
        switch (status) {
            case "success":
                toast.success(message, {
                    position: toast.POSITION.TOP_RIGHT,
                })
                break
            case "failure":
                toast.error(message, {
                    position: toast.POSITION.TOP_RIGHT,
                })
                break
            default:
                toast.warn(message, {
                    position: toast.POSITION.TOP_RIGHT,
                })
        }
    }

    const getChalls = async () => {

        const fetchTimeout = (url: string, ms: number, signal: AbortSignal, options = {}) => {
            const controller = new AbortController();
            const promise = fetch(url, { signal: controller.signal, ...options });
            if (signal) signal.addEventListener("abort", () => controller.abort());
            const timeout = setTimeout(() => controller.abort(), ms);
            return promise.finally(() => clearTimeout(timeout));
        }

        const controller = new AbortController()
		const { signal } = controller

        try {			
			const request = await fetchTimeout("/api/challs", 7000, signal, { 
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
		} catch (error: any) {
			if (error.name === "AbortError") {
				show("failure", "Request timed out! please reload")

			} else {
				show("failure", "Server not responding, contact admin")
			}
		}

        try {			
			const instancesRequest = await fetchTimeout("/api/status", 7000, signal, { 
                headers: {
                    "Authorization": `Bearer ${Cookies.get('token')}`
                }
			})
            const instancesStatus = await instancesRequest.status;
            if (instancesStatus != 200) {
                show("info", "User not logged in!")
                router.push("/logout");
            }
    
            const instancesJSON = await instancesRequest.json();
            let emptyDict = {} as activatedItem
    
            instancesJSON.map((item: accessCredsItem, index: number) => 
                emptyDict[item["level"]] = {"password": item["password"], "port": item["port"], "verified": item["verified"], "hostname": item["hostname"]}
            )
            setActivated(emptyDict)
		} catch (error: any) {
			if (error.name === "AbortError") {
				show("failure", "Request timed out! please reload")

			} else {
				show("failure", "Server not responding, contact admin")
			}
		}
    }

    useEffect(() => {
        if (respHook && !loggedin) {
            router.push("/login");
        } else {
            getChalls();
        }
    }, [respHook])

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
                        <Challenge key={ item.level } challObject={ item } isActive={ activated[item.level] != undefined } isVisible={ false } onClick={ handleVisibility } password={ activated[item.level] != undefined ? activated[item.level]["password"]: "" } port={activated[item.level] != undefined ? activated[item.level]["port"]: ""} hostname={ activated[item.level] != undefined ? activated[item.level]["hostname"]: "" }/>
                    )
                }
            </div>
        </>
    )
}

export default Challenges;