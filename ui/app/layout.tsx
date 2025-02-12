'use client'

import localFont from "next/font/local";
import "@/styles/globals.css";
import "@/styles/hint-toast.css";
import "@/styles/notification.css";
import { useEffect } from "react";
import { useAuthStore } from "@/store/authStore";
import { useChallengeStore } from "@/store/challengeStore";
import { useMetadataStore } from "@/store/metadataStore";
import NavBar from "@/components/NavBar";
import { ThemeProvider } from "@/components/theme-provider"
import { HintToastContainer } from "@/components/hints/HintToastContainer";
import { NotificationContainer } from "@/components/NotificationContainer";

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

	const { user, fetchUser } = useAuthStore();
	const { fetchChallenges } = useChallengeStore();
	const { metadataLoaded, fetchMetadata } = useMetadataStore();

	useEffect(() => {
		if (user.userid !== -1 && user.teamid !== -1) {
			fetchChallenges();
		}

		if (user.userid === -1) {
			fetchUser();
		}
	}, [user, fetchUser, fetchChallenges]);

	useEffect(() => {
		if (!metadataLoaded) {
			fetchMetadata();
		}

		return () => { };
	}, [metadataLoaded, fetchMetadata]);

	return (
		<html lang="en" suppressHydrationWarning>
			<body
				className={`${geistSans.variable} ${geistMono.variable} antialiased flex flex-col h-screen`}
			>
				<ThemeProvider attribute="class" defaultTheme="dark" enableSystem>
					<NotificationContainer />
					<HintToastContainer />
					<NavBar />
					{children}
					<div className={`${geistSans.variable} ${geistMono.variable} fixed bottom-5 end-5 text-slate-500`}>
						powered by isolet
					</div>
				</ThemeProvider>
			</body>
		</html>
	);
}
