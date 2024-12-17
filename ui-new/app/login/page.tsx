"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2 } from "lucide-react"
import useLogin from "@/hooks/useLogin"

export default function Login() {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const { loading, loginAPI } = useLogin();
	const router = useRouter();

	async function onSubmit(event: React.SyntheticEvent) {
		event.preventDefault();
		let result = await loginAPI(email, password);
		if (result) {
			router.push("/");
		}
	}

	return (
		<div className="container flex flex-col items-center justify-center h-screen">
		<Card className="w-[350px]">
			<CardHeader className="space-y-1">
			<CardTitle className="text-2xl">Sign in</CardTitle>
			<CardDescription>
				Enter your email/username and password 
			</CardDescription>
			</CardHeader>
			<CardContent className="grid gap-4">
			<div className="grid gap-2">
				<Label htmlFor="email">Email</Label>
				<Input
					id="email" 
					type="email"
					placeholder="titan@titancrew"
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
					placeholder="password"
					onChange={(event) => {
						setPassword(event.target.value);
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
				Sign In
			</Button>
			</CardFooter>
		</Card>
		</div>
	)
}
