"use client";

import React, { useState } from "react";
import useJoinTeam from "@/hooks/useJoinTeam";
import { Button } from "@/components/extras/buttons";

function TeamInit() {
	const [teamname, setTeamName] = useState("");
	const [password, setPassword] = useState("");
	const { loading, teamJoin } = useJoinTeam();

	const handleSubmit = async (action: string) => {
		await teamJoin(teamname, password, action);
	};

	let inputClass =
		"px-4 py-2 w-72 border border-gray-600 rounded-md bg-background text-foreground";

	return (
		<div className="flex flex-col gap-4">
			<div className="flex justify-center">
				Please join a team or create a new team to continue.
			</div>
			<div>
				<form
					onSubmit={(event) => {
						event.preventDefault();
						const action = (event.nativeEvent as SubmitEvent).submitter?.id || "";
						handleSubmit(action);
					}}
					className="flex flex-col gap-2 justify-center items-center"
				>
					<input
						id="teamname"
						type="text"
						name="teamname"
						placeholder="teamname"
						value={teamname}
						onChange={(event) => {
							setTeamName(event.target.value);
						}}
						className={inputClass}
						required
					></input>
					<input
						id="password"
						type="password"
						name="password"
						placeholder="password"
						value={password}
						onChange={(event) => {
							setPassword(event.target.value);
						}}
						className={inputClass}
						required
					></input>
					<div className="flex gap-2">
						<Button type="submit" id="join">
							{loading ? "Joining..." : "Join"}
						</Button>
						<Button type="submit" variant="secondary" id="create">
							{loading ? "Creating..." : "Create"}
						</Button>
					</div>
				</form>
			</div>
		</div>
	);
}

export default TeamInit;
