"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2 } from "lucide-react"
import useRegister from "@/hooks/useRegister"

export default function Register() {
	const [username, setUsername] = useState("");
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const [confirm, setConfirm] = useState("");
	const { loading, registerAPI } = useRegister();

	async function onSubmit(event: React.SyntheticEvent) {
		event.preventDefault();
		await registerAPI(username, email, password, confirm);
	}

	return (
		<div className="container flex flex-col items-center justify-center h-full">
		<Card className="w-[350px]">
			<CardHeader className="space-y-1">
			<CardTitle className="text-2xl">Register</CardTitle>
			<CardDescription>
				Sign up for a new account 
			</CardDescription>
			</CardHeader>
			<CardContent className="grid gap-4">
			<div className="grid gap-2">
				<Label htmlFor="username">Username</Label>
				<Input
					id="username" 
					type="username"
					placeholder="username"
					autoComplete="username"
					onChange={(event) => {
						setUsername(event.target.value);
					}}
					required
				/>
			</div>
			<div className="grid gap-2">
				<Label htmlFor="email">Email</Label>
				<Input
					id="email" 
					type="email"
					placeholder="titan@titancrew"
					autoComplete="email"
					onChange={(event) => {
						setEmail(event.target.value);
					}}
					required
				/>
			</div>
			<div className="grid gap-2">
				<Label htmlFor="password">Password</Label>
				<Input 
					id="password" 
					type="password" 
					name="password"
					placeholder="password"
					autoComplete="new-password"
					onChange={(event) => {
						setPassword(event.target.value);
					}}
					required
				/>
			</div>
			<div className="grid gap-2">
				<Label htmlFor="confirm-password">Confirm</Label>
				<Input 
					id="confirm-password" 
					type="password" 
					name="confirm-password"
					placeholder="confirm password"
					autoComplete="on"
					onChange={(event) => {
						setConfirm(event.target.value);
					}}
					required
				/>
			</div>
			</CardContent>
			<CardFooter>
			<Button className="w-full" onClick={onSubmit}>
				{loading && (
					<Loader2 className="mr-2 h-4 w-4 animate-spin" />
				)}
				Register
			</Button>
			</CardFooter>
		</Card>
		</div>
	)
}
