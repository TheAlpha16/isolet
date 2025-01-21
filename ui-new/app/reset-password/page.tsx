"use client"

import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "@/components/ui/card"
import { Loader2 } from "lucide-react"
import { Eye, EyeClosed } from "lucide-react"
import useLogin from "@/hooks/useLogin"
import { useRouter, useSearchParams } from "next/navigation"
import showToast, { ToastStatus } from "@/utils/toastHelper"

export default function Register() {
	const [password, setPassword] = useState("");
	const [confirm, setConfirm] = useState("");
	const { loading, resetPasswordAPI } = useLogin();
	const [showPasswd, setShowPasswd] = useState(false);
	const [showConfirm, setShowConfirm] = useState(false);
	const [isClient, setIsClient] = useState(false);
	const router = useRouter();

	const searchParams = useSearchParams();
	const token = isClient ? searchParams.get("token") : null;

	useEffect(() => {
		setIsClient(true);
	}, []);

	useEffect(() => {
		if (isClient && !token) {
			showToast(ToastStatus.Failure, "missing token");
			router.push("/");
		}
	}, [token, isClient]);

	async function onSubmit(event: React.SyntheticEvent) {
		if (!token) {
			showToast(ToastStatus.Failure, "missing token");
			return;
		}

		event.preventDefault();
		await resetPasswordAPI(password, confirm, token);
	}

	if (!isClient) {
		return null;
	}

	return (
		<div className="container flex flex-col items-center justify-center h-full">
			<Card className="w-[350px]">
				<CardHeader className="space-y-1">
					<CardTitle className="text-2xl">Reset Password</CardTitle>
					<CardDescription>
						Enter your new password
					</CardDescription>
				</CardHeader>
				<CardContent className="grid gap-4">
					<div className="grid gap-2">
						<Label htmlFor="password">Password</Label>
						<div className="relative">
							<Input
								id="password"
								type={showPasswd ? "text" : "password"}
								name="password"
								placeholder="password"
								autoComplete="new-password"
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
					<div className="grid gap-2">
						<Label htmlFor="confirm-password">Confirm</Label>
						<div className="relative">
							<Input
								id="confirm-password"
								type={showConfirm ? "text" : "password"}
								name="confirm-password"
								placeholder="confirm password"
								autoComplete="on"
								onChange={(event) => {
									setConfirm(event.target.value);
								}}
								className="pr-10"
								required
							/>
							<Button variant={"ghost"} size="icon" className="absolute inset-y-0 right-0" onClick={() => setShowConfirm(!showConfirm)}>
								{showConfirm ? <Eye className="h-5 w-5" /> : <EyeClosed className="h-5 w-5" />}
							</Button>
						</div>
					</div>
				</CardContent>
				<CardFooter>
					<Button className="w-full" onClick={onSubmit}>
						{loading && (
							<Loader2 className="mr-2 h-4 w-4 animate-spin" />
						)}
						Reset
					</Button>
				</CardFooter>
			</Card>
		</div>
	)
}
