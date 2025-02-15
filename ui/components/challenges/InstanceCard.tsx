import React, { useEffect, useState, useMemo } from "react";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Loader2, Play, StopCircle, RefreshCw, KeyRound, Terminal } from "lucide-react";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { Input } from "@/components/ui/input";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { useInstanceStore } from "@/store/instanceStore";
import { CopyButton } from "@/components/utils/copy-button";
import { GenerateChallengeEndpoint } from "@/utils/parser";

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
		setLoading,
	} = useInstanceStore();
	const instance = instances[chall_id];
	const [copiedLink, setCopiedLink] = useState<string | null>(null);

	const connectionLink = useMemo(() => {
		if (!instance) return "";
		return GenerateChallengeEndpoint(instance.deployment, instance.hostname, instance.port);
	}, [instance]);

	useEffect(() => {
		if (!instance || !instance.deadline) return;

		setTimeLeft(Math.max(0, Math.floor((instance.deadline - Date.now()) / 1000)));
		let timeout: NodeJS.Timeout | undefined;

		const tick = () => {
			setTimeLeft(Math.max(0, Math.floor((instance.deadline - Date.now()) / 1000)));

			if (instance.active && timeLeft > 0) {
				timeout = setTimeout(tick, 1000);
			}
		};

		tick();

		return () => {
			if (timeout) clearTimeout(timeout);
		};
	}, [instance]);

	const handleStart = async () => {
		setLoading(true);
		try {
			await startInstance(chall_id);
		} finally {
			setLoading(false);
		}
	};

	const handleStop = async () => {
		setLoading(true);
		try {
			await stopInstance(chall_id);
		}
		finally {
			setLoading(false);
		}
	};

	const handleExtend = async () => {
		setLoading(true);
		try {
			await extendInstance(chall_id);
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
			setTimeout(() => setCopiedLink(null), 4000);
		} catch {
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
						onClick={instance.active ? handleStop : handleStart}
						disabled={loading}
						variant="outline"
						size="sm"
					>
						{loading ? (
							<Loader2 className="animate-spin text-yellow-500" />
						) : instance.active ? (
							<StopCircle className="text-red-500" />
						) : (
							<Play className="text-green-500" />
						)}
					</Button>

					{instance.active && (
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

					{instance.active && (
						<span className="text-sm font-mono border h-9 p-2 rounded-lg">
							{formatTime(timeLeft)}
						</span>
					)}
				</div>

				<TooltipProvider>
					<Tooltip>
						<TooltipTrigger asChild>
							<Badge variant="outline" className={getStatusColor()}>
								{instance.active ? "Running" : "Stopped"}
							</Badge>
						</TooltipTrigger>
						<TooltipContent>
							<p>
								{instance.active
									? "Access the instance using the details below"
									: "Start the instance to access it"}
							</p>
						</TooltipContent>
					</Tooltip>
				</TooltipProvider>
			</div>

			{instance.active && (
				<div className="flex items-center space-x-2">
					<div className="relative flex-grow">
						<Terminal className="h-5 w-5 absolute left-2 top-1/2 transform -translate-y-1/2 text-gray-500" />
						<Input
							value={connectionLink}
							readOnly
							className="pl-8 truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
						/>
					</div>
					<CopyButton copiedLink={copiedLink} content={connectionLink} copyToClipboard={copyToClipboard} />
				</div>)}
			{instance.active && instance.password && (
				<div className="flex items-center space-x-2">
					<div className="relative flex-grow">
						<KeyRound className="h-5 w-5 absolute left-2 top-1/2 transform -translate-y-1/2 text-gray-500" />
						<Input
							value={instance.password}
							readOnly
							className="pl-8 truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
						/>
					</div>
					<CopyButton copiedLink={copiedLink} content={instance.password} copyToClipboard={copyToClipboard} />
				</div>
			)}
		</div>
	);
}
