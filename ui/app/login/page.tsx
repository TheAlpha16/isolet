'use client'

import { useState } from "react"
import User from "@/components/User"
import { useRouter } from "next/navigation"
import { toast } from "react-toastify"
import Link from "next/link"

function Login() {
    const [view, setView] = useState(false)
    const user = User()
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

    const handleSubmit = async () => {
		const email = (document.getElementById("email") as HTMLInputElement).value;
		const password = (document.getElementById("password") as HTMLInputElement).value;

		if (email === "" || password === ""){
			show("failure", "All fields are required!")
			return;
		}

        let formData = new FormData()
        formData.append("email", email)
        formData.append("password", password)

		const resp = await fetch("/auth/login", {
				method: "POST",
				body: formData,
		})
        const jsonResp = await resp.json();
        if (jsonResp.status == "failure") {
            show(jsonResp.status, jsonResp.message)
        } else {
            user.setLoggedin(true);
            router.push("/");
        }
    }

    const handleShowHide = () => {
        setView(!view);
    }

	const inputClass = "px-4 py-2 w-72 bg-transparent border border-gray-400 rounded-md outline-palette-500 text-black bg-white"
	return (
		<>
			<div className="flex flex-col gap-1 px-6 pt-6 pb-4 mt-6 font-mono justify-center self-center border-2 border-palette-600 text-palette-100 rounded-md">
				<div className="grid grid-cols-1 gap-y-4 justify-items-center">
					<label>Creds please!</label>
					<input
						id="email"
						name="email"
						placeholder="Email"
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
					<div className="flex gap-2 w-full">
						<input type="checkbox" onClick={ handleShowHide } className="cursor-pointor"></input>
						<div>Show Password</div>
					</div>
					<button className="px-5 py-2 relative duration-300 ease-in bg-palette-500 text-black rounded-md" onClick={ handleSubmit }>Login</button>
				</div>
				<Link href={ "/register" }>
					<div className="text-blue-400 text-sm text-center underline">
						Sign up
					</div>
				</Link>
			</div>
		</>
	)
}

export default Login;