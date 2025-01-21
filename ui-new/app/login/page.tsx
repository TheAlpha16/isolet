"use client"

import { useState } from "react"
import { useRouter } from "next/navigation"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2 } from "lucide-react"
import useLogin from "@/hooks/useLogin"
import { Eye, EyeClosed } from "lucide-react"
import Link from "next/link"

export default function Login() {
	const [email, setEmail] = useState("");
	const [password, setPassword] = useState("");
	const { loading, loginAPI } = useLogin();
	const router = useRouter();
	const [showPasswd, setShowPasswd] = useState(false);

	async function onSubmit(event: React.SyntheticEvent) {
		event.preventDefault();
		let result = await loginAPI(email, password);
		if (result) {
			router.push("/");
		}
	}

	return (
		<div className="container flex flex-col items-center justify-center h-full">
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
							name="email"
							autoComplete="email"
							onChange={(event) => {
								setEmail(event.target.value);
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
					<div className="text-sm text-right">
						<Link href="/forgot-password" className="text-primary hover:underline">
							Forgot password?
						</Link>
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
