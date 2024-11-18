'use client'

import localFont from "next/font/local";
import "@/styles/globals.css";
import { ToastContainer } from "react-toastify";
import "react-toastify/ReactToastify.css";
import { useEffect } from "react";
import { useAuthStore } from "@/store/authStore";
import NavBar from "@/components/NavBar";

const geistSans = localFont({
	src: "./fonts/GeistVF.woff",
	variable: "--font-geist-sans",
	weight: "100 900",
});
const geistMono = localFont({
	src: "./fonts/GeistMonoVF.woff",
	variable: "--font-geist-mono",
	weight: "100 900",
});

export default function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {

	const { loggedIn, fetchUser } = useAuthStore();

	useEffect(() => {
		if (!loggedIn) {
			fetchUser();
		}
	}, [loggedIn, fetchUser]);

	return (
		<html lang="en">
			<body
				className={`${geistSans.variable} ${geistMono.variable} antialiased`}
			>
				<ToastContainer />
				<NavBar />
				{children}
				<div className={`${geistSans.variable} ${geistMono.variable} fixed bottom-5 end-5 text-slate-500`}>
					powered by isolet
				</div>
			</body>
		</html>
	);
}
