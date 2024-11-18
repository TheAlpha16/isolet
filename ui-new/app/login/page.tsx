"use client";

import React, { useEffect, useState } from "react";
import useLogin from "@/hooks/useLogin";
import { useRouter } from "next/navigation";
import FormButton from "@/components/extras/buttons";

function Login() {
	const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const { loading, loginAPI } = useLogin();
    const router = useRouter();

    const handleSubmit = async () => {
        let result = await loginAPI(email, password);
        if (result) {
            router.push('/');
        }
    };

    let inputClass = "px-4 py-2 w-72 border border-gray-600 rounded-md bg-background text-foreground";

	return (
		<div>
			<form
				onSubmit={(event) => {
					event.preventDefault();
					handleSubmit();
				}}
				className="flex flex-col gap-2 justify-center items-center"
			>
				<input
					id="email"
					type="text"
					name="email"
					placeholder="email/username"
					value={email}
					onChange={(event) => {
						setEmail(event.target.value);
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
					<FormButton type="submit">
						{loading ? "Logging in..." : "Login"}
					</FormButton>
					<FormButton type="button" variant="secondary">
						<a href="/register">Register</a>
					</FormButton>
				</div>
			</form>
		</div>
	);
}

export default Login;
