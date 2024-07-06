"use client";

import { useState, useEffect } from "react";
import LoginStatus from "@/components/User";
import { useRouter } from "next/navigation";
import { toast } from "react-toastify";
import Link from "next/link";
import { ShowPassword } from "@/components/Utils";

function Login() {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const user = LoginStatus();
	const router = useRouter();

	useEffect(() => {
		if (user.respHook) {
			if (user.loggedin) {
				router.push("/");
			}
		}
	}, [user.respHook]);

	const show = (status: string, message: string) => {
		switch (status) {
			case "success":
				toast.success(message, {
					position: toast.POSITION.TOP_RIGHT,
				});
				break;
			case "failure":
				toast.error(message, {
					position: toast.POSITION.TOP_RIGHT,
				});
				break;
			default:
				toast.warn(message, {
					position: toast.POSITION.TOP_RIGHT,
				});
		}
	};

	const fetchTimeout = (
		url: string,
		ms: number,
		signal: AbortSignal,
		options = {}
	) => {
		const controller = new AbortController();
		const promise = fetch(url, { signal: controller.signal, ...options });
		if (signal) signal.addEventListener("abort", () => controller.abort());
		const timeout = setTimeout(() => controller.abort(), ms);
		return promise.finally(() => clearTimeout(timeout));
	};

	const handleSubmit = async () => {
		const email = (document.getElementById("email") as HTMLInputElement)
			.value;
		const password = (
			document.getElementById("password") as HTMLInputElement
		).value;
		const controller = new AbortController();
		const { signal } = controller;

		if (email === "" || password === "") {
			show("failure", "All fields are required!");
			return;
		}

		let formData = new FormData();
		formData.append("email", email);
		formData.append("password", password);

		try {
			const resp = await fetchTimeout("/auth/login", 5000, signal, {
				method: "POST",
				body: formData,
			});
			const jsonResp = await resp.json();
			if (jsonResp.status == "failure") {
				show(jsonResp.status, jsonResp.message);
			} else {
				user.setLoggedin(true);
				router.push("/");
			}
		} catch (error: any) {
			if (error.name === "AbortError") {
				show("failure", "Request timed out! please reload");
			} else {
				show("failure", "Server not responding, contact admin");
			}
		}
	};

	const inputClass =
		"px-4 py-2 w-72 bg-transparent border border-gray-400 rounded-md outline-palette-500 text-black bg-white";
	return (
		<>
			<div className="flex flex-col gap-1 px-6 pt-6 pb-4 mt-6 font-mono justify-center self-center border-2 border-palette-600 text-palette-100 rounded-md">
				<div>
					<div className="grid grid-cols-1 gap-y-4 justify-items-center">
						<label>Creds please!</label>
						<input
							id="email"
							name="email"
							placeholder="Email"
							type="text"
							value={email}
							onChange={(e) => setEmail(e.target.value)}
							className={inputClass}
						></input>
						<div className="relative w-full">
							<input
								id="password"
								name="password"
								placeholder="Password"
								type="password"
								value={password}
								onChange={(e) => setPassword(e.target.value)}
								onKeyDown={(e) => {
									if (e.key == "Enter") {
										handleSubmit();
									}
								}}
								className={`${inputClass} w-full`}
							></input>
							<ShowPassword />
						</div>
						<button
							type="submit"
							className="px-5 py-2 relative duration-300 ease-in bg-palette-500 text-black rounded-md hover:bg-palette-400"
							onClick={handleSubmit}
						>
							Login
						</button>
					</div>
				</div>
				<Link href={"/register"}>
					<div className="text-blue-400 text-sm text-center underline">
						Sign up
					</div>
				</Link>
			</div>
		</>
	);
}

export default Login;
