'use client'

import Cookies from "js-cookie"
import { useRouter } from "next/navigation"
import { useEffect, useState, useRef } from "react"
import { toast } from "react-toastify"
import "react-toastify/dist/ReactToastify.css"

export interface challItem {
    level: number
    name: string
    prompt: string
    solves: number
    tags: Array<string>
}

interface Props {
    challObject: challItem
    isActive: boolean
    isVisible: boolean
    onClick: any
    password: string
    port: string
    hostname: string
}

function Challenge(props: Props) {
    const [isActive, setActive] = useState(props.isActive)
    const port = useRef(props.port)
    const password = useRef(props.password)
    const hostname = useRef(props.hostname)

    useEffect(() => {
        setActive(props.isActive)
        port.current = props.port
        password.current = props.password
        hostname.current = props.hostname
    }, [props])

    const fetchTimeout = (url: string, ms: number, signal: AbortSignal, options = {}) => {
        const controller = new AbortController();
        const promise = fetch(url, { signal: controller.signal, ...options });
        if (signal) signal.addEventListener("abort", () => controller.abort());
        const timeout = setTimeout(() => controller.abort(), ms);
        return promise.finally(() => clearTimeout(timeout));
    }

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

    const router = useRouter()

    const handleSubmit = async () => {
        const data = new FormData()
        const flag = document.getElementById(`flag-${props.challObject.level}`) as HTMLInputElement

        if (flag.value == "") {
            show("failure", "empty string is not the flag!")
            return
        }

        data.append("flag", flag.value)
        data.append("level", `${props.challObject.level}`)

        const controller = new AbortController()
		const { signal } = controller

        try {			
			const request = await fetchTimeout("/api/submit", 5000, signal, { 
                method: "POST",
                headers: {
                    "Authorization": `Bearer ${Cookies.get('token')}`
                },
                body: data
			})
            const status = await request.status
            if (status == 401) {
                show("failure", "not logged in :/")
                router.push("/logout")
                return
            }
            const submitJSON = await request.json()
            show(submitJSON.status, submitJSON.message)
            if (submitJSON.status == "success") {
                flag.value = ""
            }
		} catch (error: any) {
			if (error.name === "AbortError") {
				show("failure", "Request timed out! please reload")
			} else {
				show("failure", "Server not responding, contact admin")
			}
		}
    }

    const changeBtn = (btn: HTMLButtonElement, status: string) => {
        const launchButton = document.getElementById(`launch-${props.challObject.level}`) as HTMLButtonElement
        
        switch(status) {
            case "stopped":
                btn.classList.remove("bg-rose-500", "bg-amber-300", "text-black", "text-palette-100")
                btn.classList.add("bg-palette-500", "text-black")
                btn.innerText = "Start"
                launchButton.addEventListener("click", eventListen)
                break
            case "running":
                btn.classList.add("bg-rose-500", "text-palette-100")
                btn.classList.remove("bg-palette-500", "bg-amber-300", "text-black", "text-palette-100")
                btn.innerText = "Stop"
                launchButton.addEventListener("click", eventListen)
                break
            case "starting":
                btn.classList.remove("bg-rose-500", "bg-palette-500", "text-black", "text-palette-100")
                btn.classList.add("bg-amber-300", "text-black")
                btn.innerText = "Starting.."
                launchButton.removeEventListener("click", eventListen)
                break
            case "stopping":
                btn.classList.remove("bg-rose-500", "bg-palette-500", "text-black", "text-palette-100")
                btn.classList.add("bg-amber-300", "text-black")
                btn.innerText = "Stopping.."
                launchButton.removeEventListener("click", eventListen)
                break
            default:
                return
        }
    }

    const eventListen = async (e: any) => {
        const launchButton = e.target as HTMLButtonElement
        const buttonStatus = launchButton.innerText == "Stop"
        const controller = new AbortController()
        const { signal } = controller

        const data = new FormData()
        data.append("level", `${props.challObject.level}`)
        
        switch(buttonStatus) {
            case true:
                changeBtn(launchButton, "stopping")
                try {			
                    const request = await fetchTimeout("/api/stop", 100000, signal, { 
                        method: "POST",
                        headers: {
                            "Authorization": `Bearer ${Cookies.get('token')}`
                        },
                        body: data
                    })
                    const status = await request.status
                    if (status == 401) {
                        router.push("/logout")
                    }
    
                    const reqJSON = await request.json()
                    show(reqJSON.status, reqJSON.message)
                    if (reqJSON.status == "failure") {
                        changeBtn(launchButton, "running")
                    } else {
                        changeBtn(launchButton, "stopped")
                        setActive(false)
                    }
                } catch (error: any) {
                    if (error.name === "AbortError") {
                        changeBtn(launchButton, "running")
                        show("failure", "Request timed out! please reload")
                    } else {
                        changeBtn(launchButton, "running")
                        show("failure", "Server not responding, contact admin")
                    }
                }
                break
                
            case false:
                changeBtn(launchButton, "starting")

                try {			
                    const requestLanuch = await fetchTimeout("/api/launch", 100000, signal, { 
                        method: "POST",
                        headers: {
                            "Authorization": `Bearer ${Cookies.get('token')}`
                        },
                        body: data
                    })

                    const statusLaunch = await requestLanuch.status
                    if (statusLaunch == 401) {
                        router.push("/logout")
                    }
                    const reqLaunchJSON = await requestLanuch.json()
                    if (reqLaunchJSON.status == "failure") {
                        show(reqLaunchJSON.status, reqLaunchJSON.message)
                        changeBtn(launchButton, "stopped")
                    } else {
                        show(reqLaunchJSON.status, "Instance launched successfully")
                        let returnedData = JSON.parse(atob(reqLaunchJSON.message))
                        port.current = returnedData.port
                        password.current = returnedData.password
                        hostname.current = returnedData.hostname
                        changeBtn(launchButton, "running")
                        setActive(true)
                    }

                } catch (error: any) {
                    if (error.name === "AbortError") {
                        changeBtn(launchButton, "stopped")
                        show("failure", "Request timed out! please reload")
                    } else {
                        changeBtn(launchButton, "stopped")
                        show("failure", "Server not responding, contact admin")
                    }
                }
                break
        }
    }

    const copyAccessString = () => {
        navigator.clipboard.writeText(`ssh level${props.challObject.level}@${hostname.current} -p ${port.current}`)
        show("success", "copied!")
    }

    const copyPasswdString = () => {
        navigator.clipboard.writeText(`${password.current}`)
        show("success", "copied!")
    }

    useEffect(() => {
        let launchButton = document.getElementById(`launch-${props.challObject.level}`) as HTMLButtonElement
        let copyAccessDiv = document.getElementById(`access-${props.challObject.level}`) as HTMLDivElement
        let passwdDiv = document.getElementById(`passwd-${props.challObject.level}`) as HTMLDivElement

        launchButton.addEventListener("click", eventListen)
        copyAccessDiv.addEventListener("click", copyAccessString)
        passwdDiv.addEventListener("click", copyPasswdString)

        return () => {
            launchButton.removeEventListener("click", eventListen)
            copyAccessDiv.removeEventListener("click", copyAccessString)
        }
    }, [])

    return (
        <>
            <div className={`w-11/12 sm:w-9/12 flex flex-col p-3 border border-palette-600 self-center transition duration-300 ease-in-out relative rounded-md`} onClick={ props.onClick } data-level={ props.challObject.level } id={`level-${props.challObject.level}`}>
                <div id={`ping-${props.challObject.level}`} className="absolute -top-1 -right-1" data-level={ props.challObject.level }>
                    <span className={`${isActive ? "": "hidden"} relative flex h-3 w-3`} data-level={ props.challObject.level }>
                        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-palette-500 opacity-75" data-level={ props.challObject.level }></span>
                        <span className="relative inline-flex rounded-full h-3 w-3 bg-palette-500" data-level={ props.challObject.level }></span>
                    </span>
                </div>
                <div className="flex w-full justify-between" data-level={ props.challObject.level }>
                    <div className="flex items-center font-Roboto text-xl font-semibold p-2" data-level={ props.challObject.level }>
                        { props.challObject.name }
                    </div>
                    <div className="flex items-center font-thin p-2 from-neutral-500" data-level={ props.challObject.level }>
                        { props.challObject.solves } solves
                    </div>
                </div> 
                <div id={`submit-${props.challObject.level}`} className={ `${props.isVisible ? "": "hidden"} flex flex-col p-3 font-mono transition duration-300 ease-in-out bg-gray-800 rounded-md` } data-level={ props.challObject.level }>
                    <div data-level={ props.challObject.level } className="flex font-light items-center p-2">
                        { props.challObject.prompt }
                    </div>
                    <div data-level={ props.challObject.level } className="flex justify-between flex-wrap">
                        <div data-level={ props.challObject.level } className="flex gap-2 w-full justify-start flex-wrap">
                            <button id={`launch-${props.challObject.level}`} className={`p-2 w-32 rounded-md ${ isActive ? "text-palette-100": "text-black" } ${ isActive ? "bg-rose-500": "bg-palette-500" }`} data-level={ props.challObject.level }>{ isActive ? "Stop": "Start" }</button>
                            <div data-level={ props.challObject.level } className={`flex rounded-md items-center h-10 justify-center gap-2 ${isActive ? "": "hidden"}`}>
                                <div className="bg-slate-950 p-2 rounded-md" data-level={ props.challObject.level }> {`${password.current.substring(0, 7)}************${password.current.substring(27)}`} </div>
                                <div id={`passwd-${props.challObject.level}`} className="hover:cursor-pointer hover:bg-slate-950 p-3 rounded-md" data-level={ props.challObject.level }>
                                    <svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" fill="#E4EEE7" data-level={ props.challObject.level } >
                                        <path d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 0 1 0 1.5h-1.5a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-1.5a.75.75 0 0 1 1.5 0v1.5A1.75 1.75 0 0 1 9.25 16h-7.5A1.75 1.75 0 0 1 0 14.25Z" data-level={ props.challObject.level }></path><path d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0 1 14.25 11h-7.5A1.75 1.75 0 0 1 5 9.25Zm1.75-.25a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-7.5a.25.25 0 0 0-.25-.25Z" data-level={ props.challObject.level }></path>
                                    </svg>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div className="flex gap-2 py-2 flex-wrap" data-level={ props.challObject.level }>
                        <input id={`flag-${props.challObject.level}`} placeholder="flag" name="flag" type="text" className="border p-2 grow outline-palette-500 rounded-md text-black" data-level={ props.challObject.level } required></input>
                        <button onClick={ handleSubmit } className="p-2 w-24 text-black bg-palette-500 rounded-md hover:bg-palette-400" data-level={ props.challObject.level }>Submit</button>
                    </div>
                    <div className="flex gap-2 py-2 flex-wrap" data-level={ props.challObject.level }>
                        {
                            props.challObject.tags.map((tag, index) => {
                                return (
                                    <div key={ index } className="p-1 px-3 bg-slate-950 rounded-md" data-level={ props.challObject.level }>
                                        { tag }
                                    </div>
                                )
                            })
                        }
                    </div>
                    <div className={`${isActive ? "": "hidden"} flex gap-2 items-center`} data-level={ props.challObject.level }>
                        <div className="bg-slate-950 p-1 px-3 rounded-md text-palette-500" data-level={ props.challObject.level }>
                            { `$ ssh level${props.challObject.level}@${hostname.current} -p ${port.current}` }
                        </div>
                        <div id={ `access-${props.challObject.level}` } className="hover:cursor-pointer hover:bg-slate-950 p-2 rounded-md" data-level={ props.challObject.level }>
                            <svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" fill="#E4EEE7" data-level={ props.challObject.level }>
                                <path d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 0 1 0 1.5h-1.5a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-1.5a.75.75 0 0 1 1.5 0v1.5A1.75 1.75 0 0 1 9.25 16h-7.5A1.75 1.75 0 0 1 0 14.25Z" data-level={ props.challObject.level }></path><path d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0 1 14.25 11h-7.5A1.75 1.75 0 0 1 5 9.25Zm1.75-.25a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-7.5a.25.25 0 0 0-.25-.25Z" data-level={ props.challObject.level }></path>
                            </svg>
                        </div>
                    </div>
                </div>
            </div>
        </>
    )
}

export default Challenge;