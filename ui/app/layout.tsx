"use client";

import { HintToastContainer } from "@/components/hints/HintToastContainer";
import NavBar from "@/components/NavBar";
import { NotificationContainer } from "@/components/NotificationContainer";
import { ThemeProvider } from "@/components/theme-provider";
import { useChallengeStore, useEventStore, useProfileStore } from "@/store";
import { useRealtimeClient } from "@/hooks/useRealtimeClient";
import "@/styles/globals.css";
import "@/styles/hint-toast.css";
import "@/styles/notification.css";
import localFont from "next/font/local";
import { useEffect } from "react";

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
  const { user, team, fetchMe } = useProfileStore();
  const { fetchChallenges } = useChallengeStore();
  const { loaded, fetchInfo } = useEventStore();

  useRealtimeClient(team?.id, user);

  // always fetch user profile on mount to ensure proper authentication state
  useEffect(() => {
    fetchMe();
  }, [fetchMe, fetchChallenges]);

  // pre-fetch challenges when user and team are available
  // might need challenges in profile and scoreboard pages
  useEffect(() => {
    if (user && team) {
      fetchChallenges();
    }
  }, [user, team, fetchChallenges]);

  useEffect(() => {
    if (!loaded) {
      fetchInfo();
    }
    return () => {};
  }, [loaded, fetchInfo]);

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
          <div
            className={`${geistSans.variable} ${geistMono.variable} fixed bottom-5 end-5 text-slate-500`}
          >
            powered by isolet
          </div>
        </ThemeProvider>
      </body>
    </html>
  );
}
