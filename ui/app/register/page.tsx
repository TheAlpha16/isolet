'use client'
import { useState, useEffect } from "react"
import { toast } from "react-toastify"

function Register() {
	const [view, setView] = useState(false)

	const handleShowHide = () => {
		setView(!view)
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


	const changeBtn = (btn: HTMLButtonElement, status: string) => {
		switch(status) {
			case "done":
				btn.classList.remove("bg-amber-300", "bg-palette-500", "text-black", "text-palette-100")
				btn.classList.add("bg-palette-500", "text-palette-100")
				btn.innerText = "Register"
				break
			case "waiting":
				btn.classList.remove("bg-palette-500", "bg-amber-300", "text-black", "text-palette-100")
				btn.classList.add("bg-amber-300", "text-black")
				btn.innerText = "Loading.."
				break
			default:
				return
		}
	}

	const eventListen = async (e: any) => {
		const launchButton = e.target as HTMLButtonElement
		const prevemail = (document.getElementById("email") as HTMLInputElement).value
		const prevusername = (document.getElementById("username") as HTMLInputElement).value
		const prevpassword = (document.getElementById("password") as HTMLInputElement).value
		const prevconfirm = (document.getElementById("confirm") as HTMLInputElement).value


		if (prevemail === "" || prevpassword === "" || prevconfirm === "" || prevusername === ""){
			show("failure", "All fields are required!")
			return
		}

		const email = prevemail.trim()
		const username = prevusername.trim()
		const password = prevpassword.trim()
		const confirm = prevconfirm.trim()

		if (password !== confirm) {
			show("failure", "Passwords don't match")
			return
		}

		if (!(password.length >= 8)) {
			show("failure", "password length should be 8 or more")
		}
		// const checkPassword = (str: string) => {
		//     var re = /^(?=.*\d)(?=.*[$#@%&*.!])(?=.*[a-z])(?=.*[A-Z]).{6,}$/
		//     return re.test(str)
		// }

		// if (!checkPassword(password)) {
		//     show("failure", "Not a strong password")
		//     show("", "match ^[a-zA-Z0-9$#@%&*.!]{6, 32}$")
		//     return
		// }

		changeBtn(launchButton, "waiting")
		let formData = new FormData()
		formData.append("email", email)
		formData.append("username", username)
		formData.append("password", password)
		formData.append("confirm", confirm)

		const fetchTimeout = (url: string, ms: number, signal: AbortSignal, options = {}) => {
			const controller = new AbortController()
			const promise = fetch(url, { signal: controller.signal, ...options })
			if (signal) signal.addEventListener("abort", () => controller.abort())
			const timeout = setTimeout(() => controller.abort(), ms)
			return promise.finally(() => clearTimeout(timeout))
		}

		const controller = new AbortController()
		const { signal } = controller

		try {			
			const resp = await fetchTimeout("/auth/register", 5000, signal, { 
					method: "POST",
					body: formData,
			})
			const jsonResp = await resp.json()
			show(jsonResp.status, jsonResp.message)
			changeBtn(launchButton, "done")
		} catch (error: any) {
			if (error.name === "AbortError") {
				show("failure", "Request timed out! please reload")
				changeBtn(launchButton, "done")

			} else {
				show("failure", "Server not responding, contact admin")
				changeBtn(launchButton, "done")
			}
		}
	}

	useEffect(() => {
		let launchButton = document.getElementById("reg-btn") as HTMLButtonElement

		launchButton.addEventListener("click", eventListen)
		return () => launchButton.removeEventListener("click", eventListen);
	}, [])

	const inputClass = "px-4 py-2 w-72 bg-transparent border border-gray-400 rounded-md outline-palette-500 text-black bg-white"
	return (
		<>
			<div className="flex flex-col gap-1 px-6 pt-6 pb-4 mt-6 font-mono justify-center self-center border-2 border-palette-600 text-palette-100 rounded-md">
				<div className="grid grid-cols-1 gap-y-2 justify-items-center">
					<label>Enter your details</label>
					<input
						id="email"
						name="email"
						placeholder="Email"
						type="text"
						className={ inputClass } 
					></input>
					<input
						id="username"
						name="username"
						placeholder="Username"
						type="text"
						className={ inputClass } 
					></input>
					<input
						id="password"
						name="password"
						placeholder="Password"
						type={ view ? "text": "password" } 
						className={ inputClass }
					></input>
					<input
						id="confirm"
						name="confirm"
						placeholder="Confirm Password"
						type={ view ? "text": "password" } 
						className={ inputClass }
					></input>
					<div className="flex gap-2">
						<input type="checkbox" onClick={ handleShowHide } className="cursor-pointor"></input>
						<div>Show Password</div>
					</div>
					<button id="reg-btn" className="px-5 py-2 relative duration-300 ease-in bg-palette-500 text-black rounded-md hover:bg-palette-400">Register</button>
				</div>
			</div>
		</>
	)
}

export default Register