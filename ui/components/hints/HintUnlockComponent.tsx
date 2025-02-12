import React from 'react';
import {
	Dialog,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";

interface HintUnlockConfirmationProps {
	isOpen: boolean;
	onClose: () => void;
	onConfirm: () => void;
	cost: number;
}

export function HintUnlockConfirmation({ isOpen, onClose, onConfirm, cost }: HintUnlockConfirmationProps) {
	return (
		<Dialog open={isOpen} onOpenChange={onClose}>
			<DialogContent className="sm:max-w-[425px]">
				<DialogHeader>
					<DialogTitle>Unlock Hint</DialogTitle>
					<DialogDescription>
						Are you sure you want to unlock this hint for {cost} points?
					</DialogDescription>
				</DialogHeader>
				<DialogFooter className="flex flex-col space-y-2 sm:space-y-0 sm:flex-row sm:justify-end">
					<Button variant="outline" onClick={onClose}>Cancel</Button>
					<Button onClick={onConfirm}>Unlock</Button>
				</DialogFooter>
			</DialogContent>
		</Dialog>
	);
}

