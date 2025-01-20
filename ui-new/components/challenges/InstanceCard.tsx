import React, { useState, useEffect } from "react"
import { Button } from "@/components/ui/button"
import { Badge } from "@/components/ui/badge"
import { Loader2, Play, StopCircle, Copy, RefreshCw, Check } from "lucide-react"
import showToast, { ToastStatus } from "@/utils/toastHelper"
import { InstanceType } from "@/store/instanceStore"
import { Input } from "@/components/ui/input"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip"

interface InstanceCardProps {
	instance: InstanceType
}

export function InstanceCard({ instance }: InstanceCardProps) {
	const [containerStatus, setContainerStatus] = useState<"stopped" | "starting" | "running" | "stopping">("stopped")
	const [isLoading, setIsLoading] = useState(false)
	const [timeLeft, setTimeLeft] = useState(3600)
	const [copiedLink, setCopiedLink] = useState<string | null>(null);

	useEffect(() => {
		let interval: NodeJS.Timeout
		if (containerStatus === "running") {
			interval = setInterval(() => {
				setTimeLeft((prev) => Math.max(0, prev - 1))
			}, 1000)
		}
		return () => clearInterval(interval)
	}, [containerStatus])

	const handleContainerAction = async () => {
		setIsLoading(true)
		if (containerStatus === "stopped") {
			setContainerStatus("starting")
			await new Promise((resolve) => setTimeout(resolve, 2000))
			setContainerStatus("running")
			setTimeLeft(3600)
			showToast(ToastStatus.Success, "Container started successfully")
		} else if (containerStatus === "running") {
			setContainerStatus("stopping")
			await new Promise((resolve) => setTimeout(resolve, 2000))
			setContainerStatus("stopped")
			showToast(ToastStatus.Success, "Container stopped successfully")
		}
		setIsLoading(false)
	}

	const handleExtend = async () => {
		setIsLoading(true)
		await new Promise((resolve) => setTimeout(resolve, 1000))
		setTimeLeft((prev) => prev + 1800) // Extend by 30 minutes
		setIsLoading(false)
		showToast(ToastStatus.Success, "Time extended by 30 minutes")
	}

	const getStatusColor = () => {
		switch (containerStatus) {
			case "running":
				return "text-green-500"
			case "starting":
			case "stopping":
				return "text-yellow-500"
			case "stopped":
				return "text-red-500"
		}
	}

	const formatTime = (seconds: number) => {
		const minutes = Math.floor(seconds / 60)
		const remainingSeconds = seconds % 60
		return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`
	}


	const copyToClipboard = (text: string) => {
		// navigator.clipboard.writeText(text);
		setCopiedLink(text);
		setTimeout(() => setCopiedLink(null), 2000);
	};

	return (
		<div className="flex flex-col bg-card p-4 rounded-lg shadow-sm border space-y-2">
			<div className="flex items-center space-x-4 justify-between">
				<div className="flex space-x-2 items-center">
					<Button
						onClick={handleContainerAction}
						disabled={isLoading || containerStatus === "starting" || containerStatus === "stopping"}
						variant={"outline"}
						size={"sm"}
					>
						{isLoading ? (
							<Loader2 className="animate-spin text-yellow-500" />
						) : containerStatus === "running" ? (
							<StopCircle className="text-red-500" />
						) : (
							<Play className="text-green-500" />
						)}
					</Button>

					<Button
						size={"sm"}
						variant={"outline"}
						onClick={handleExtend}
						disabled={containerStatus !== "running" || isLoading}
					>
						<RefreshCw className="h-4 w-4 mr-2" />
						Extend
					</Button>
					{containerStatus === "running" && <span className="text-sm font-mono border h-9 p-2 rounded-lg">{formatTime(timeLeft)}</span>}
				</div>
				<Badge variant="outline" className={`${getStatusColor()}`}>
					{containerStatus.charAt(0).toUpperCase() + containerStatus.slice(1)}
				</Badge>
			</div>
			{containerStatus === "running" && (
				<div className="flex items-center space-x-2">
					<Input value={"ssh hacker@ctf.infosec.org.in"} readOnly className="flex-grow truncate font-mono" />
					<TooltipProvider>
						<Tooltip>
							<TooltipTrigger asChild>
								<Button variant="outline" size="icon" onClick={() => copyToClipboard("ssh hacker@ctf.infosec.org.in")}>
									{copiedLink === "ssh hacker@ctf.infosec.org.in" ? <Check className="h-4 w-4 text-green-500" /> : <Copy className="h-4 w-4" />}
								</Button>
							</TooltipTrigger>
							<TooltipContent>
								<p>{copiedLink === "ssh hacker@ctf.infosec.org.in" ? 'Copied!' : 'Copy to clipboard'}</p>
							</TooltipContent>
						</Tooltip>
					</TooltipProvider>
				</div>)}
		</div>
	)
}

