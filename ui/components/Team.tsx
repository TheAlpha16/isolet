"use client";

import { toast } from "react-toastify";
import { ShowPassword } from "@/components/Utils";
import Cookies from "js-cookie";

export interface Team {
	teamid: number;
	teamname: string;
	captain: number;
	members: TeamMember[];
}

interface TeamMember {
	userid: number;
	username: string;
	score: number;
}

interface joinProps {
	router: any;
}

export function JoinTeam(props: joinProps) {
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

	const handleJoin = async () => {
		const teamname = (
			document.getElementById("teamname") as HTMLInputElement
		).value;
		const password = (
			document.getElementById("password") as HTMLInputElement
		).value;
		const controller = new AbortController();
		const { signal } = controller;

		if (teamname == "" || password == "") {
			show("failure", "All fields are required");
			return;
		}

		let formData = new FormData();
		formData.append("teamname", teamname);
		formData.append("password", password);

		try {
			const resp = await fetchTimeout("/api/team/join", 7000, signal, {
				method: "POST",
				body: formData,
				headers: {
					Authorization: `Bearer ${Cookies.get("token")}`,
				},
			});
			const jsonResp = await resp.json();
			show(jsonResp.status, jsonResp.message);

			if (jsonResp.status == "success") {
				props.router.push("/logout");
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
		"px-4 py-2 w-72 rounded-md bg-white text-black outline-palette-500 border border-gray-400";
	const buttonClass =
		"px-4 py-2 duration-300 ease-in bg-palette-500 text-black rounded-md hover:bg-palette-400";

	return (
		<div className="flex flex-col gap-1 px-6 pt-6 pb-4 mt-6 font-mono justify-center self-center border-2 border-palette-600 text-palette-100 rounded-md">
			<div className="grid grid-cols-1 gap-y-4 justify-items-center">
				<label>Create or Join Team</label>
				<input
					id="teamname"
					type="text"
					placeholder="Team Name"
					name="teamname"
					className={inputClass}
				/>
				<div className="w-full relative">
					<input
						id="password"
						type="password"
						placeholder="Password"
						name="password"
						className={inputClass}
					/>
					<ShowPassword />
				</div>
				<div className="flex gap-2">
					<button className={buttonClass}>Create</button>
					<button className={buttonClass} onClick={handleJoin}>
						Join
					</button>
				</div>
			</div>
		</div>
	);
}

export function TeamPage(teamid: number) {
	return (
		<div>
			<h1>Team Page</h1>
			<p>Team ID: {teamid}</p>
		</div>
	);
}
