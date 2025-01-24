"use client"

import React from 'react';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Copy, RefreshCw } from 'lucide-react';

interface TeamInviteProps {
	isOpen: boolean;
	onClose: () => void;
	onGenerate: () => void;
	onCopy: () => void;
	token: string;
}

export function TeamInvite({ isOpen, onClose, onGenerate, onCopy, token }: TeamInviteProps) {
	return (
		<Dialog open={isOpen} onOpenChange={onClose}>
			<DialogContent className="sm:max-w-[425px]">
				<DialogHeader>
					<DialogTitle>Invite</DialogTitle>
				</DialogHeader>
				<DialogDescription>
					Share this token with others to invite.
				</DialogDescription>
				<div className="flex flex-row  space-x-2 w-full max-h-fit">
					<Button onClick={onGenerate} 
					variant="outline">
						<RefreshCw size={20} />
					</Button>
					<div className="relative flex-grow">
						<Input
							value={token}
							className="pr-10"
							readOnly
						/>
						<Button variant={"ghost"} size="icon" className="absolute inset-y-0 right-0" onClick={onCopy}>
							<Copy size={20} />
						</Button>
					</div>
				</div>
			</DialogContent>
		</Dialog>
	);
}

