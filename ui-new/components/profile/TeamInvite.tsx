"use client"

import React, { useState } from 'react';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { RefreshCw, Loader2 } from 'lucide-react';
import showToast, { ToastStatus } from '@/utils/toastHelper';
import { CopyButton } from '@/components/utils/copy-button';
import useInvite from '@/hooks/useInviteToken';

interface TeamInviteProps {
	isOpen: boolean;
	onClose: () => void;
}

export function TeamInvite({ isOpen, onClose }: TeamInviteProps) {
	const [copiedLink, setCopiedLink] = useState<string | null>(null);
	const { inviteToken, loading, generateInviteTokenAPI } = useInvite();

	const onGenerate = async () => {
		await generateInviteTokenAPI();
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

	return (
		<Dialog open={isOpen} onOpenChange={onClose}>
			<DialogContent className="sm:max-w-[425px]">
				<DialogHeader>
					<DialogTitle>Invite</DialogTitle>
				</DialogHeader>
				<DialogDescription>
					Share this token with others to invite.
				</DialogDescription>
				<div className="flex items-center space-x-2">
					<Button
						onClick={onGenerate}
						variant="outline"
						size={"icon"}
						disabled={loading}
					>
						{loading ? <Loader2 className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
					</Button>
					<Input
						value={inviteToken}
						className="truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
						readOnly
					/>
					<CopyButton copiedLink={copiedLink} content={inviteToken} copyToClipboard={copyToClipboard} />
				</div>
			</DialogContent>
		</Dialog>
	);
}
