"use client";

import React, { useState } from "react";
import useJoinTeam from "@/hooks/useJoinTeam";
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2 } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Eye, EyeClosed } from "lucide-react"

export default function TeamInit() {
	const [teamname, setTeamName] = useState("");
	const [password, setPassword] = useState("");
	const { loading, teamJoin } = useJoinTeam();
	const [showPasswd, setShowPasswd] = useState(false);

	async function onSubmit(action: string) {
		await teamJoin(teamname, password, action);
	}

	return (
		<div className="container flex flex-col items-center justify-center h-full">
			<Card className="w-[350px]">
				<CardHeader className="space-y-1">
					<CardTitle className="text-2xl">Team</CardTitle>
					<CardDescription>
						Please create/join a team to continue
					</CardDescription>
				</CardHeader>
				<CardContent className="grid gap-4">
					<div className="grid gap-2">
						<Label htmlFor="teamname">Team Name</Label>
						<Input
							id="teamname"
							type="text"
							placeholder="teamname"
							name="teamname"
							autoComplete="teamname"
							onChange={(event) => {
								setTeamName(event.target.value);
							}}
							required
						/>
					</div>
					<div className="grid gap-2">
						<Label htmlFor="password">Password</Label>
						<div className="relative">
							<Input
								id="password"
								type={showPasswd ? "text" : "password"}
								placeholder="password"
								name="password"
								autoComplete="current-password"
								onChange={(event) => {
									setPassword(event.target.value);
								}}
								className="pr-10"
								required
							/>
							<Button variant={"ghost"} size="icon" className="absolute inset-y-0 right-0" onClick={() => setShowPasswd(!showPasswd)}>
								{showPasswd ? <Eye className="h-5 w-5" /> : <EyeClosed className="h-5 w-5" />}
							</Button>
						</div>
					</div>
				</CardContent>
				<CardFooter>
					<div className="flex gap-2">
						<Button className="w-full" onClick={(event) => {
							event.preventDefault();
							onSubmit("join")
						}}>
							{loading && (
								<Loader2 className="mr-2 h-4 w-4 animate-spin" />
							)}
							Join
						</Button>
						<Button className="w-full" variant={"secondary"} onClick={(event) => {
							event.preventDefault();
							onSubmit("create")
						}}>
							{loading && (
								<Loader2 className="mr-2 h-4 w-4 animate-spin" />
							)}
							Create
						</Button>
					</div>
				</CardFooter>
			</Card>
		</div>
	)
}
