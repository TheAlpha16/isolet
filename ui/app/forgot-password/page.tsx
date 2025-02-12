"use client"

import { useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2 } from "lucide-react"
import useLogin from "@/hooks/useLogin"

export default function ForgotPassword() {
	const [email, setEmail] = useState("");
	const { loading, forgotPasswordAPI } = useLogin();

	return (
		<div className="container flex flex-col items-center justify-center h-full">
		<Card className="w-[350px]">
			<CardHeader className="space-y-1">
			<CardTitle className="text-2xl">Forgot Password</CardTitle>
			<CardDescription>
				Enter your email address to get the link to reset your password
			</CardDescription>
			</CardHeader>
			<CardContent className="grid gap-4">
			<div className="grid gap-2">
				<Label htmlFor="email">Email</Label>
				<Input
					id="email" 
					type="email"
					placeholder="titan@titancrew"
					name="email"
					autoComplete="email"
					onChange={(event) => {
						setEmail(event.target.value);
					}}
					required
				/>
			</div>
			</CardContent>
			<CardFooter>
			<Button className="w-full" onClick={() => forgotPasswordAPI(email)} disabled={loading}>
				{loading && (
					<Loader2 className="mr-2 h-4 w-4 animate-spin" />
				)}
				Submit
			</Button>
			</CardFooter>
		</Card>
		</div>
	)
}
