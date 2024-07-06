"use client";

import { useState } from "react";

export function ShowPassword() {
	const [isPasswordVisible, setIsPasswordVisible] = useState(false);

	const switchPassword = (event: any) => {
		setIsPasswordVisible((prevState) => !prevState);
		const parentDiv = (event.target as HTMLElement).closest(".relative");
		const inputField = parentDiv?.querySelector("input");
		if (inputField) {
			inputField.type =
				inputField.type === "password" ? "text" : "password";
		}
	};

	const iconSrc = isPasswordVisible
		? "/static/assets/eye.svg"
		: "/static/assets/eye-closed.svg";

	return (
		<button
			type="button"
			onClick={switchPassword}
			className="absolute inset-y-0 right-0 flex items-center pr-3"
		>
			<img
				src={iconSrc}
				alt="Toggle password visibility"
				className="h-5 w-5"
			/>
		</button>
	);
}

