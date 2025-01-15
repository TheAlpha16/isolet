'use client'

import showToast, { ToastStatus } from "@/utils/toastHelper";
import Image from "next/image";
import { useAuthStore } from "@/store/authStore";

export default function Home() {

	const { loggedIn, user, logout } = useAuthStore();

	return (
		<div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
		<main className="flex flex-col gap-3 row-start-2 items-center sm:items-start">
			<div className="font-mono font-bold text-2xl">isolet v2</div>
			<p>version 2 is currently under development. </p>
			<p>If you have any thoughts, I would be happy to take suggestions</p>
			<p>Feel free to open an issue on the GitHub repository.</p>
			{!loggedIn ? (
				<div className="mt-2"> You can login with sample creds
					<code className="flex gap-1 items-center">
					username:
						<pre className="flex gap-1 bg-secondary py-1 px-1 rounded-md max-w-fit justify-center">
							{`glimpse`}
						</pre>
					</code>
					<code className="flex gap-1 items-center">
					password:
						<pre className="bg-secondary py-1 px-1 rounded-md max-w-fit">
							{`Strongpasswd123.`}
						</pre>
					</code>	
				</div>
			) : (
				<div>
					Welcome back, {user?.username}!
				</div>
			)}
		</main>
		<footer className="row-start-3 flex gap-6 flex-wrap items-center justify-center">
			<a
			className="flex items-center gap-2 hover:underline hover:underline-offset-4"
			href="https://github.com/TheAlpha16/isolet"
			target="_blank"
			rel="noopener noreferrer"
			>
			<Image
				aria-hidden
				src="/window.svg"
				alt="Window icon"
				width={16}
				height={16}
			/>
			Source code
			</a>
		</footer>
		</div>
	);
}
