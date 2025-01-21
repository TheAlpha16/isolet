import React, { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Loader2, Play, StopCircle, Copy, RefreshCw, Check, KeyRound, Terminal } from "lucide-react";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { Input } from "@/components/ui/input";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { useInstanceStore } from "@/store/instanceStore";

interface InstanceCardProps {
	chall_id: number;
}

export function InstanceCard({ chall_id }: InstanceCardProps) {
	const [timeLeft, setTimeLeft] = useState(0);
	const {
		instances,
		loading,
		startInstance,
		stopInstance,
		extendInstance,

		// DEBUG
		updateInstance,
		setLoading,

	} = useInstanceStore();
	const instance = instances[chall_id];
	const [copiedLink, setCopiedLink] = useState<string | null>(null);

	useEffect(() => {
		if (instance && instance.deadline) {
			setTimeLeft(Math.max(0, Math.floor((instance.deadline - Date.now()) / 1000)));
		}

		let interval: NodeJS.Timeout | undefined;
		if (instance?.active) {
			interval = setInterval(() => {
				setTimeLeft((prev) => Math.max(0, prev - 1));
			}, 1000);
		}

		return () => {
			if (interval) clearInterval(interval);
		};
	}, [instance]);

	const handleStart = async () => {
		setLoading(true);
		try {

			// DEBUG
			await new Promise((resolve) => setTimeout(resolve, 2000));
			updateInstance(chall_id, {
				active: true,
				deadline: Date.now() + 60 * 60 * 1000,
				connString: "ssh hacker@ctf.infosec.org.in",
				password: "ab390b1e003f027ca48d926fa16"
			});

			// await startInstance(chall_id);
		} finally {
			setLoading(false);
		}
	};

	const handleStop = async () => {
		setLoading(true);
		try {

			// DEBUG
			await new Promise((resolve) => setTimeout(resolve, 2000));
			updateInstance(chall_id, { active: false, deadline: 0 });

			// await stopInstance(chall_id);
		}
		finally {
			setLoading(false);
		}
	};

	const handleExtend = async () => {
		setLoading(true);
		try {

			// DEBUG
			await new Promise((resolve) => setTimeout(resolve, 2000));
			updateInstance(chall_id, { deadline: instance.deadline + 60 * 60 * 1000 });

			// await extendInstance(chall_id);

		} finally {
			setLoading(false);
		}
	};

	const formatTime = (seconds: number) => {
		const minutes = Math.floor(seconds / 60);
		const remainingSeconds = seconds % 60;
		return `${minutes}:${remainingSeconds.toString().padStart(2, "0")}`;
	};

	const copyToClipboard = (text: string) => {
		try {
			navigator.clipboard.writeText(text);
			setCopiedLink(text);
			setTimeout(() => setCopiedLink(null), 2000);
		} catch (error) {
			showToast(ToastStatus.Failure, "Failed to copy to clipboard");
		}
	};

	const getStatusColor = () => {
		if (!instance) return "text-gray-500";
		if (instance.active) return "text-green-500";
		return "text-red-500";
	};

	return (
		<div className="flex flex-col bg-card p-3 rounded-lg shadow-sm border space-y-2">
			<div className="flex items-center space-x-4 justify-between">
				<div className="flex space-x-2 items-center">
					<Button
						onClick={instance?.active ? handleStop : handleStart}
						disabled={loading}
						variant="outline"
						size="sm"
					>
						{loading ? (
							<Loader2 className="animate-spin text-yellow-500" />
						) : instance?.active ? (
							<StopCircle className="text-red-500" />
						) : (
							<Play className="text-green-500" />
						)}
					</Button>

					{instance?.active && (
						<Button
							size="sm"
							variant="outline"
							onClick={handleExtend}
							disabled={loading}
						>
							<RefreshCw className="h-4 w-4 mr-2" />
							Extend
						</Button>
					)}

					{instance?.active && (
						<span className="text-sm font-mono border h-9 p-2 rounded-lg">
							{formatTime(timeLeft)}
						</span>
					)}
				</div>

				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger asChild>
							<Badge variant="outline" className={getStatusColor()}>
								{instance?.active ? "Running" : "Stopped"}
							</Badge>
						</TooltipTrigger>
						<TooltipContent>
							<p>
								{instance?.active
									? "Access the instance using the details below"
									: "Start the instance to access it"}
							</p>
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>
			</div>

			{instance?.active && (
				<div className="flex items-center space-x-2">
					<div className="relative flex-grow">
						<Terminal className="h-5 w-5 absolute left-2 top-1/2 transform -translate-y-1/2 text-gray-500" />
						<Input
							value={instance.connString}
							readOnly
							className="pl-8 truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
						/>
					</div>
					<TooltipProvider>
						<Tooltip>
							<TooltipTrigger asChild>
								<Button
									variant="outline"
									size="icon"
									onClick={() => copyToClipboard(instance.connString)}
								>
									{copiedLink === instance.connString ? (
										<Check className="h-4 w-4 text-green-500" />
									) : (
										<Copy className="h-4 w-4" />
									)}
								</Button>
							</TooltipTrigger>
							<TooltipContent>
								<p>
									{copiedLink === instance.connString ? "Copied!" : "Copy to clipboard"}
								</p>
							</TooltipContent>
						</Tooltip>
					</TooltipProvider>
				</div>)}
			{instance?.active && instance?.password && (
				<div className="flex items-center space-x-2">
					<div className="relative flex-grow">
						<KeyRound className="h-5 w-5 absolute left-2 top-1/2 transform -translate-y-1/2 text-gray-500" />
						<Input
							value={instance.password}
							readOnly
							className="pl-8 truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
						/>
					</div>
					<TooltipProvider>
						<Tooltip>
							<TooltipTrigger asChild>
								<Button
									variant="outline"
									size="icon"
									onClick={() => copyToClipboard(instance.password)}
								>
									{copiedLink === instance.connString ? (
										<Check className="h-4 w-4 text-green-500" />
									) : (
										<Copy className="h-4 w-4" />
									)}
								</Button>
							</TooltipTrigger>
							<TooltipContent>
								<p>
									{copiedLink === instance.password ? "Copied!" : "Copy to clipboard"}
								</p>
							</TooltipContent>
						</Tooltip>
					</TooltipProvider>
				</div>
			)}
		</div>
	);
}
