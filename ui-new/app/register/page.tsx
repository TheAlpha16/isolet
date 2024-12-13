"use client";

import React, { useState } from "react";
import useRegister from "@/hooks/useRegister";
import FormButton from "@/components/extras/buttons";
import Link from "next/link";

function Register() {
    const [username, setUsername] = useState("");
	const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [confirm, setConfirm] = useState("");
    const { loading, registerAPI } = useRegister();

    const handleSubmit = async () => {
        await registerAPI(username, email, password, confirm);
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
					id="username"
					type="text"
					name="username"
					placeholder="username"
					value={username}
					onChange={(event) => {
						setUsername(event.target.value);
					}}
					className={inputClass}
					required
				></input>
				<input
					id="email"
					type="text"
					name="email"
					placeholder="email"
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
				<input
					id="confirm"
					type="password"
					name="confirm"
					placeholder="confirm password"
					value={confirm}
					onChange={(event) => {
						setConfirm(event.target.value);
					}}
					className={inputClass}
					required
				></input>
				<div className="flex gap-2">
					<FormButton type="submit" disabled={loading}>
						{loading ? "Waiting..." : "Register"}
					</FormButton>
					<Link href="/login">
						<FormButton type="button" variant="secondary">
							Login
						</FormButton>
					</Link>
				</div>
			</form>
		</div>
	);
}

export default Register;
