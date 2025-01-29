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
import { RefreshCw } from 'lucide-react';
import showToast, { ToastStatus } from '@/utils/toastHelper';
import { CopyButton } from '@/components/utils/copy-button';

interface TeamInviteProps {
	isOpen: boolean;
	onClose: () => void;
	onGenerate: () => void;
	token: string;
}

export function TeamInvite({ isOpen, onClose, onGenerate, token }: TeamInviteProps) {
	const [copiedLink, setCopiedLink] = useState<string | null>(null);

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
					>
						<RefreshCw className="h-4 w-4" />
					</Button>
					<Input
						value={token}
						className="truncate font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0"
						readOnly
					/>
					<CopyButton copiedLink={copiedLink} content={token} copyToClipboard={copyToClipboard} />
				</div>
			</DialogContent>
		</Dialog>
	);
}
