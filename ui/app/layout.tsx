"use client";

import NavBar from "@/components/Navigation";
import "../public/static/css/globals.css";
import { useState, useEffect } from "react"
import Cookies from "js-cookie"
import { ToastContainer } from "react-toastify"
import React from "react"
import { Context } from "@/components/User"


export default function RootLayout({
  	children,
}: {
  	children: React.ReactNode;
}) {
	const [loggedin, setLoggedin] = useState(false);
	const [respHook, setRespHook] = useState(false);

	useEffect(() => {
		const verify = async () => {
			const data = await fetch("/api/status", {
				headers: {
					Authorization: `Bearer ${Cookies.get("token")}`,
				},
			})
			const statusCode = await data.status;
			if (statusCode == 200) {
				setLoggedin(true)
			}
			setRespHook(true)
		}
		verify()
		return
	}, [])

	return (
		<>
			<html lang="en">
				<head>
					<title>
						ISOLET
					</title>
				</head>
				<body className="flex flex-col bg-palette-200 text-palette-100 h-screen">
					<NavBar loggedin={loggedin} />
					<ToastContainer></ToastContainer>
					<Context.Provider value={{ loggedin, setLoggedin, respHook }}>
						{ children }
					</Context.Provider>
					<div className="z-40 fixed bottom-5 end-5 text-slate-500">
						made by titans@titancrew 👀
					</div>
				</body>
			</html>
		</>
	)
}
