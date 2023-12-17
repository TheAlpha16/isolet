'use client'
import { useState, useEffect } from "react"
import { toast } from "react-toastify"

function Register() {
    const [view, setView] = useState(false)

    const handleShowHide = () => {
        setView(!view);
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
        const email = (document.getElementById("email") as HTMLInputElement).value
		const username = (document.getElementById("username") as HTMLInputElement).value
		const password = (document.getElementById("password") as HTMLInputElement).value
		const confirm = (document.getElementById("confirm") as HTMLInputElement).value


		if (email === "" || password === "" || confirm === "" || username === ""){
			show("failure", "All fields are required!")
			return
		}

        if (password !== confirm) {
            show("failure", "Passwords don't match")
            return
        }

        const checkPassword = (str: string) => {
            var re = /^(?=.*\d)(?=.*[$#@%&*.!])(?=.*[a-z])(?=.*[A-Z]).{6,}$/
            return re.test(str)
        }

        if (!checkPassword(password)) {
            show("failure", "Not a strong password")
            show("", "match ^[a-zA-Z0-9$#@%&*.!]{6, 32}$")
            return
        }

        changeBtn(launchButton, "waiting")
        let formData = new FormData()
        formData.append("email", email)
        formData.append("username", username)
        formData.append("password", password)
        formData.append("confirm", confirm)

        const resp = await fetch("/auth/register", {
            method: "POST",
            body: formData,
        })

        const jsonResp = await resp.json()
        show(jsonResp.status, jsonResp.message)
        changeBtn(launchButton, "done")
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
					<button id="reg-btn" className="px-5 py-2 relative duration-300 ease-in bg-palette-500 text-palette-100 rounded-md">Register</button>
				</div>
			</div>
        </>
    )
}

export default Register;