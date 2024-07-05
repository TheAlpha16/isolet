'use client'

import Cookies from "js-cookie"
import { useRouter } from "next/navigation"
import { useEffect, useState, useRef } from "react"
import { toast } from "react-toastify"
import Timer from "./Timer"
import "react-toastify/dist/ReactToastify.css"
import { challContext } from "./Contexts"

export interface Challenge {
	chall_id: number;
	level: number;
	name: string;
	prompt: string;
	category: string;
	type: "static" | "dynamic" | "on-demand";
	points: number;
	files: Array<string>;
	hints: Array<Hint>;
	solves: number;
	author: string;
	tags: Array<string>;
	port: number;
	subd: string;
}

export interface Hint {
	hid: number;
	chall_id: number;
	hint: string;
	cost: number;
}

interface Props {
	challObject: Challenge
	isActive: boolean
	isVisible: boolean
	onClick: any
	password: string
	port: string
	hostname: string
	deadline: number
}

// function Challenge(props: Props) {
// 	const [isActive, setActive] = useState(props.isActive)
// 	const port = useRef(props.port)
// 	const password = useRef(props.password)
// 	const hostname = useRef(props.hostname)
// 	const [deadline, setDeadLine] = useState(props.deadline)
// 	const [timeout, setTimeLeft] = useState(false)
// 	const router = useRouter()

// 	const fetchTimeout = (url: string, ms: number, signal: AbortSignal, options = {}) => {
// 		const controller = new AbortController()
// 		const promise = fetch(url, { signal: controller.signal, ...options })
// 		if (signal) signal.addEventListener("abort", () => controller.abort(), true)
// 		const timeout = setTimeout(() => controller.abort(), ms)
// 		return promise.finally(() => clearTimeout(timeout))
// 	}

// 	const show = (status: string, message: string) => {
// 		switch (status) {
// 			case "success":
// 				toast.success(message, {
// 					position: toast.POSITION.TOP_RIGHT,
// 				})
// 				break
// 			case "failure":
// 				toast.error(message, {
// 					position: toast.POSITION.TOP_RIGHT,
// 				})
// 				break
// 			default:
// 				toast.warn(message, {
// 					position: toast.POSITION.TOP_RIGHT,
// 				})
// 		}
// 	}

// 	const changeBtn = (btn: HTMLButtonElement, status: string) => {
// 		switch(status) {
// 			case "stopped":
// 				btn.classList.remove("bg-rose-500", "bg-amber-300", "text-black", "text-palette-100")
// 				btn.classList.add("bg-palette-500", "text-black")
// 				btn.innerText = "Start"
// 				btn.addEventListener("click", eventListen, true)
// 				break
// 			case "running":
// 				btn.classList.add("bg-rose-500", "text-palette-100")
// 				btn.classList.remove("bg-palette-500", "bg-amber-300", "text-black", "text-palette-100")
// 				btn.innerText = "Stop"
// 				btn.addEventListener("click", eventListen, true)
// 				break
// 			case "starting":
// 				btn.classList.remove("bg-rose-500", "bg-palette-500", "text-black", "text-palette-100")
// 				btn.classList.add("bg-amber-300", "text-black")
// 				btn.innerText = "Starting.."
// 				btn.removeEventListener("click", eventListen, true)
// 				break
// 			case "stopping":
// 				btn.classList.remove("bg-rose-500", "bg-palette-500", "text-black", "text-palette-100")
// 				btn.classList.add("bg-amber-300", "text-black")
// 				btn.innerText = "Stopping.."
// 				btn.removeEventListener("click", eventListen, true)
// 				break
// 			case "wait":
// 				btn.classList.remove("bg-palette-500", "bg-amber-300")
// 				btn.classList.add("bg-amber-300")
// 				btn.innerText = "Wait.."
// 				btn.removeEventListener("click", eventExtend, true)
// 				break
// 			case "extend":
// 				btn.classList.remove("bg-palette-500", "bg-amber-300")
// 				btn.classList.add("bg-palette-500")
// 				btn.innerText = "Extend"
// 				btn.addEventListener("click", eventExtend, true)
// 				break
// 			default:
// 				return
// 		}
// 	}

// 	const handleSubmit = async () => {
// 		const data = new FormData()
// 		const flag = document.getElementById(`flag-${props.challObject.level}`) as HTMLInputElement

// 		if (flag.value == "") {
// 			show("failure", "empty string is not the flag!")
// 			return
// 		}

// 		let flagValue = flag.value.trim()

// 		data.append("flag", flagValue)
// 		data.append("level", `${props.challObject.level}`)

// 		const controller = new AbortController()
// 		const { signal } = controller

// 		try {			
// 			const request = await fetchTimeout("/api/submit", 5000, signal, { 
// 				method: "POST",
// 				headers: {
// 					"Authorization": `Bearer ${Cookies.get('token')}`
// 				},
// 				body: data
// 			})
// 			const status = await request.status
// 			if (status == 401) {
// 				show("failure", "not logged in :/")
// 				router.push("/logout")
// 				return
// 			}
// 			const submitJSON = await request.json()
// 			show(submitJSON.status, submitJSON.message)
// 			if (submitJSON.status == "success") {
// 				flag.value = ""
// 			}
// 		} catch (error: any) {
// 			if (error.name === "AbortError") {
// 				show("failure", "Request timed out! please reload")
// 			} else {
// 				show("failure", "Server not responding, contact admin")
// 			}
// 		}
// 	}

// 	const eventListen = async (e: any) => {
// 		const launchButton = e.target as HTMLButtonElement
// 		const buttonStatus = launchButton.innerText == "Stop"
// 		const controller = new AbortController()
// 		const { signal } = controller

// 		const data = new FormData()
// 		data.append("level", `${props.challObject.level}`)
		
// 		switch(buttonStatus) {
// 			case true:
// 				changeBtn(launchButton, "stopping")
// 				try {			
// 					const request = await fetchTimeout("/api/stop", 100000, signal, { 
// 						method: "POST",
// 						headers: {
// 							"Authorization": `Bearer ${Cookies.get('token')}`
// 						},
// 						body: data
// 					})
// 					const status = await request.status
// 					if (status == 401) {
// 						router.push("/logout")
// 					}
	
// 					const reqJSON = await request.json()
// 					show(reqJSON.status, reqJSON.message)
// 					if (reqJSON.status == "failure") {
// 						changeBtn(launchButton, "running")
// 					} else {
// 						changeBtn(launchButton, "stopped")
// 						setActive(false)
// 					}
// 				} catch (error: any) {
// 					if (error.name === "AbortError") {
// 						changeBtn(launchButton, "running")
// 						show("failure", "Request timed out! please reload")
// 					} else {
// 						changeBtn(launchButton, "running")
// 						show("failure", "Server not responding, contact admin")
// 					}
// 				}
// 				break
				
// 			case false:
// 				changeBtn(launchButton, "starting")

// 				try {
// 					const requestLanuch = await fetchTimeout("/api/launch", 100000, signal, { 
// 						method: "POST",
// 						headers: {
// 							"Authorization": `Bearer ${Cookies.get('token')}`
// 						},
// 						body: data
// 					})

// 					const statusLaunch = await requestLanuch.status
// 					if (statusLaunch == 401) {
// 						router.push("/logout")
// 					}
// 					const reqLaunchJSON = await requestLanuch.json()
// 					if (reqLaunchJSON.status == "failure") {
// 						show(reqLaunchJSON.status, reqLaunchJSON.message)
// 						changeBtn(launchButton, "stopped")
// 					} else {
// 						show(reqLaunchJSON.status, "Instance launched successfully")
// 						let returnedData = JSON.parse(atob(reqLaunchJSON.message))
// 						port.current = returnedData.port
// 						password.current = returnedData.password
// 						hostname.current = returnedData.hostname
// 						setDeadLine(returnedData.deadline)
// 						changeBtn(launchButton, "running")
// 						setActive(true)
// 					}

// 				} catch (error: any) {
// 					if (error.name === "AbortError") {
// 						changeBtn(launchButton, "stopped")
// 						show("failure", "Request timed out! please reload")
// 					} else {
// 						changeBtn(launchButton, "stopped")
// 						show("failure", "Server not responding, contact admin")
// 					}
// 				}
// 				break
// 		}
// 	}

// 	const eventExtend = async (e: any) => {
// 		const extendButton = e.target as HTMLButtonElement
// 		const controller = new AbortController()
// 		const { signal } = controller

// 		const data = new FormData()
// 		data.append("level", `${props.challObject.level}`)
		
// 		changeBtn(extendButton, "wait")

// 		try {			
// 			const requestExtend = await fetchTimeout("/api/extend", 100000, signal, { 
// 				method: "POST",
// 				headers: {
// 					"Authorization": `Bearer ${Cookies.get('token')}`
// 				},
// 				body: data
// 			})

// 			const statusLaunch = await requestExtend.status
// 			if (statusLaunch == 401) {
// 				router.push("/logout")
// 			}
// 			const reqLaunchJSON = await requestExtend.json()
// 			if (reqLaunchJSON.status == "failure") {
// 				show(reqLaunchJSON.status, reqLaunchJSON.message)
// 				changeBtn(extendButton, "extend")
// 			} else {
// 				show(reqLaunchJSON.status, "Instance extention successful")
// 				let returnedData = JSON.parse(atob(reqLaunchJSON.message))
// 				setDeadLine(returnedData.deadline)
// 				changeBtn(extendButton, "extend")
// 			}

// 		} catch (error: any) {
// 			if (error.name === "AbortError") {
// 				changeBtn(extendButton, "extend")
// 				show("failure", "Request timed out! please reload")
// 			} else {
// 				changeBtn(extendButton, "extend")
// 				show("failure", "Server not responding, contact admin")
// 			}
// 		}
// 	}

// 	const instanceEnded = async () => {
// 		const controller = new AbortController()
// 		const { signal } = controller
// 		// const launchButton = document.getElementById(`launch-${props.challObject.level}`) as HTMLButtonElement
// 		// changeBtn(launchButton, "stopping")

// 		for (let count = 0; count < 10; count++) {
// 			try {			
// 				const instancesRequest = await fetchTimeout("/api/status", 7000, signal, { 
// 					headers: {
// 						"Authorization": `Bearer ${Cookies.get('token')}`
// 					}
// 				})
// 				const instancesStatus = await instancesRequest.status
// 				if (instancesStatus != 200) {
// 					show("info", "User not logged in!")
// 					router.push("/logout")
// 				}
				
// 				let foundLevel = false
// 				const sleep = (ms: number) => new Promise((r) => setTimeout(r, ms))
// 				const instancesJSON = await instancesRequest.json()

// 				for (let index = 0; index < instancesJSON.length; index++) {
// 					const level = instancesJSON[index]["level"]
// 					if (level == props.challObject.level) {
// 						foundLevel = true
// 						break
// 					}
// 				}

// 				if (foundLevel) {
// 					await sleep(10000)
// 				} else {
// 					break
// 				}

// 			} catch (error: any) {
// 				if (error.name === "AbortError") {
// 					show("failure", "Request timed out! please reload")
// 				} else {
// 					show("failure", "Server not responding, contact admin")
// 				}
// 			}
// 		}
// 		setTimeLeft(false)
// 		setActive(false)
// 	}

// 	const copyAccessString = () => {
// 		navigator.clipboard.writeText(`ssh level${props.challObject.level}@${hostname.current} -p ${port.current}`)
// 		show("success", "copied!")
// 	}

// 	const copyPasswdString = () => {
// 		navigator.clipboard.writeText(`${password.current}`)
// 		show("success", "copied!")
// 	}

// 	useEffect(() => {
// 		let launchButton = document.getElementById(`launch-${props.challObject.level}`) as HTMLButtonElement
// 		let copyAccessDiv = document.getElementById(`access-${props.challObject.level}`) as HTMLDivElement
// 		let passwdDiv = document.getElementById(`passwd-${props.challObject.level}`) as HTMLDivElement
// 		let extendButton = document.getElementById(`extend-${props.challObject.level}`) as HTMLButtonElement

// 		launchButton.addEventListener("click", eventListen, true)
// 		copyAccessDiv.addEventListener("click", copyAccessString, true)
// 		passwdDiv.addEventListener("click", copyPasswdString, true)
// 		extendButton.addEventListener("click", eventExtend, true)

// 		return () => {
// 			launchButton.removeEventListener("click", eventListen, true)
// 			copyAccessDiv.removeEventListener("click", copyAccessString, true)
// 			passwdDiv.removeEventListener("click", copyPasswdString, true)
// 			extendButton.removeEventListener("click", eventExtend, true)
// 		}
// 	}, [])

// 	useEffect(() => {
// 		setActive(props.isActive)
// 		setDeadLine(props.deadline)
// 		port.current = props.port
// 		password.current = props.password
// 		hostname.current = props.hostname
// 	}, [props])

// 	useEffect(() => {
// 		if (timeout) {
// 			instanceEnded()
// 		}
// 	}, [timeout])

// 	return (
// 		<>
// 			<div className={`w-11/12 sm:w-9/12 flex flex-col p-3 border border-palette-600 self-center transition duration-300 ease-in-out relative rounded-md`} onClick={ props.onClick } data-level={ props.challObject.level } id={`level-${props.challObject.level}`}>
// 				<div id={`ping-${props.challObject.level}`} className="absolute -top-1 -right-1" data-level={ props.challObject.level }>
// 					<span className={`${isActive ? "": "hidden"} relative flex h-3 w-3`} data-level={ props.challObject.level }>
// 						<span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-palette-500 opacity-75" data-level={ props.challObject.level }></span>
// 						<span className="relative inline-flex rounded-full h-3 w-3 bg-palette-500" data-level={ props.challObject.level }></span>
// 					</span>
// 				</div>
// 				<div className="flex w-full justify-between" data-level={ props.challObject.level }>
// 					<div className="flex items-center font-Roboto text-xl font-semibold p-2" data-level={ props.challObject.level }>
// 						{ props.challObject.name }
// 					</div>
// 					<div className="flex items-center font-thin p-2 from-neutral-500" data-level={ props.challObject.level }>
// 						{ props.challObject.solves } solves
// 					</div>
// 				</div> 
// 				<div id={`submit-${props.challObject.level}`} className={ `${props.isVisible ? "": "hidden"} flex flex-col p-3 font-mono transition duration-300 ease-in-out bg-gray-800 rounded-md` }>
// 					<div data-level={ props.challObject.level } className="flex font-light items-center p-2 max-w-max">
// 						{ props.challObject.prompt }
// 					</div>
// 					<div className="flex justify-between flex-wrap">
// 						<div className="flex gap-2 w-full justify-start flex-wrap">
// 							<button id={`launch-${props.challObject.level}`} className={`p-2 w-32 rounded-md ${ isActive ? "text-palette-100": "text-black" } ${ isActive ? "bg-rose-500": "bg-palette-500" }`} data-level={ props.challObject.level }>{ isActive ? "Stop": "Start" }</button>
// 							<button id={`extend-${props.challObject.level}`} className={`p-2 w-32 rounded-md bg-palette-500 text-black ${ isActive ? "": "hidden" }`} data-level={ props.challObject.level }>{ "Extend" }</button>
// 							{ isActive ? <challContext.Provider value={{deadline}}><Timer setTimeLeft={setTimeLeft} level={ props.challObject.level } classes="flex items-center"/></challContext.Provider>: <div className="hidden"></div>}
// 						</div>
// 					</div>
// 					<div className="flex gap-2 py-2 flex-wrap" data-level={ props.challObject.level }>
// 						<input id={`flag-${props.challObject.level}`} placeholder="flag" name="flag" type="text" className="border p-2 grow outline-palette-500 rounded-md text-black" data-level={ props.challObject.level } required></input>
// 						<button onClick={ handleSubmit } className="p-2 w-24 text-black bg-palette-500 rounded-md hover:bg-palette-400" data-level={ props.challObject.level }>Submit</button>
// 					</div>
// 					<div className="flex gap-2 pb-2 flex-wrap">
// 						{
// 							props.challObject.tags.map((tag, index) => {
// 								return (
// 									<div key={ index } className="p-1 px-3 bg-slate-950 rounded-md" data-level={ props.challObject.level }>
// 										{ tag }
// 									</div>
// 								)
// 							})
// 						}
// 					</div>
// 					<div className={`${isActive ? "": "hidden"} flex gap-2 pb-2 items-center`}>
// 						<div className="bg-slate-950 p-1 px-3 rounded-md text-palette-500" data-level={ props.challObject.level }>
// 							{ `$ ssh level${props.challObject.level}@${hostname.current} -p ${port.current}` }
// 						</div>
// 						<div id={ `access-${props.challObject.level}` } className="hover:cursor-pointer hover:bg-slate-950 p-2 rounded-md" data-level={ props.challObject.level }>
// 							<svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" fill="#E4EEE7" data-level={ props.challObject.level }>
// 								<path d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 0 1 0 1.5h-1.5a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-1.5a.75.75 0 0 1 1.5 0v1.5A1.75 1.75 0 0 1 9.25 16h-7.5A1.75 1.75 0 0 1 0 14.25Z" data-level={ props.challObject.level }></path><path d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0 1 14.25 11h-7.5A1.75 1.75 0 0 1 5 9.25Zm1.75-.25a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-7.5a.25.25 0 0 0-.25-.25Z" data-level={ props.challObject.level }></path>
// 							</svg>
// 						</div>
// 					</div>
// 					<div className={`flex gap-2 items-center ${isActive ? "": "hidden"}`}>
// 						<div className="bg-slate-950 p-1 px-3 rounded-md" data-level={ props.challObject.level }>
// 							{`${password.current.substring(0, 7)}************${password.current.substring(27)}`}
// 						</div>
// 						<div id={`passwd-${props.challObject.level}`} className="hover:cursor-pointer hover:bg-slate-950 p-2 rounded-md" data-level={ props.challObject.level }>
// 							<svg aria-hidden="true" height="16" viewBox="0 0 16 16" version="1.1" width="16" data-view-component="true" fill="#E4EEE7" data-level={ props.challObject.level } >
// 								<path d="M0 6.75C0 5.784.784 5 1.75 5h1.5a.75.75 0 0 1 0 1.5h-1.5a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-1.5a.75.75 0 0 1 1.5 0v1.5A1.75 1.75 0 0 1 9.25 16h-7.5A1.75 1.75 0 0 1 0 14.25Z" data-level={ props.challObject.level }></path><path d="M5 1.75C5 .784 5.784 0 6.75 0h7.5C15.216 0 16 .784 16 1.75v7.5A1.75 1.75 0 0 1 14.25 11h-7.5A1.75 1.75 0 0 1 5 9.25Zm1.75-.25a.25.25 0 0 0-.25.25v7.5c0 .138.112.25.25.25h7.5a.25.25 0 0 0 .25-.25v-7.5a.25.25 0 0 0-.25-.25Z" data-level={ props.challObject.level }></path>
// 							</svg>
// 						</div>
// 					</div>
// 				</div>
// 			</div>
// 		</>
// 	)
// }

// export default Challenge