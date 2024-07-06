'use client'

import { useEffect, useState } from "react"
import Cookies from "js-cookie"
import { useRouter } from "next/navigation"
import LoginStatus from "@/components/User"
import { toast } from "react-toastify"

function Scoreboard() {
	const { loggedin, respHook } = LoginStatus()
	const router = useRouter()
	const [ scores, setScores ] = useState([{"username": "", "score": 0}])

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

	const getScores = async () => {
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
			const request = await fetchTimeout("/api/scoreboard", 7000, signal, { 
				headers: {
					"Authorization": `Bearer ${Cookies.get("token")}`
				}
			})
	
			const scoreJSON = await request.json()
			setScores(scoreJSON)
		} catch (error: any) {
			if (error.name === "AbortError") {
				show("failure", "Request timed out! please reload")

			} else {
				show("failure", "Scoreboard not available at the moment")
			}
		}
	}

	useEffect(() => {
		if (!respHook) {
			router.push("/")
		} else if (respHook && !loggedin) {
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