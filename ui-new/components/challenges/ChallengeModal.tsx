import React, { useState } from "react";
import { ChallengeType, ChallType, HintType, useChallengeStore } from "@/store/challengeStore";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import Hint from "@/components/hints/Hint";
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "@/components/ui/tooltip";
import { Download, ExternalLink, Copy, Check, Users, Play, StopCircle } from 'lucide-react';
import { InstanceType } from "@/store/instanceStore";
import { InstanceCard } from "./InstanceCard";

interface ChallengeModalProps {
	challenge: ChallengeType;
	onClose: () => void;
}

export function ChallengeModal({ challenge, onClose }: ChallengeModalProps) {
	const [flag, setFlag] = useState('');
	const [copiedLink, setCopiedLink] = useState<string | null>(null);
	const { submitFlag } = useChallengeStore();

	const flagSubmit = async () => {
		await submitFlag(challenge.chall_id, flag);
	}

	const handleSubmit = (e: React.FormEvent) => {
		e.preventDefault();
		flagSubmit();
	};

	const copyToClipboard = (text: string) => {
		navigator.clipboard.writeText(text);
		setCopiedLink(text);
		setTimeout(() => setCopiedLink(null), 2000);
	};

	return (
		<Dialog open={true} onOpenChange={onClose}>
		<DialogContent className="sm:max-w-[600px]">
			<DialogHeader>
				<DialogTitle className="text-2xl font-bold flex items-center justify-between">
					<div className="flex items-center gap-2">
						<span>{challenge.name}</span>
						{challenge.done && <Check className="w-5 h-5 text-green-500" strokeWidth={2.5}/>}
					</div>
					<Badge variant="secondary" className="ml-2">{challenge.points} pts</Badge>
				</DialogTitle>
				<div className="flex items-center justify-between text-sm text-muted-foreground">
					<span>by {challenge.author}</span>
					<div className="flex items-center">
						<Users className="w-4 h-4 mr-1" />
						<span>{challenge.solves} solves</span>
					</div>
				</div>
			</DialogHeader>
			
			<div className="space-y-4">
				<p className="text-lg">{challenge.prompt}</p>
				
				<div className="flex flex-wrap gap-2">
					{challenge.tags.map((tag) => (
						<Badge key={tag} variant="secondary">{tag}</Badge>
					))}
				</div>
				
				{challenge.files.length > 0 && (
					<div>
						<h3 className="text-lg font-semibold mb-2">Files</h3>
						<div className="flex flex-wrap gap-2">
							{challenge.files.map((file) => (
								<Button key={file} variant="outline" size="sm">
									<Download className="mr-2 h-4 w-4" />
									{file}
								</Button>
							))}
						</div>
					</div>
				)}
				
				{challenge.links.length > 0 && (
					<div><h3 className="text-lg font-semibold mb-2">Links</h3>
						<div className="space-y-2">
							{challenge.links.map((link) => (
							<div key={link} className="flex items-center space-x-2">
									<Input value={link} readOnly className="flex-grow font-mono focus:outline-none focus-visible:ring-0 focus-visible:ring-offset-0" />
								<TooltipProvider>
									<Tooltip>
										<TooltipTrigger asChild>
											<Button variant="outline" size="icon" onClick={() => copyToClipboard(link)}>
												{copiedLink === link ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
											</Button>
										</TooltipTrigger>
										<TooltipContent>
											<p>{copiedLink === link ? 'Copied!' : 'Copy to clipboard'}</p>
										</TooltipContent>
									</Tooltip>
								</TooltipProvider>
								<Button variant="outline" size="icon" asChild>
									<a href={link} target="_blank" rel="noopener noreferrer">
										<ExternalLink className="h-4 w-4" />
									</a>
								</Button>
							</div>
							))}
						</div>
					</div>
				)}
				
                {challenge.hints.length > 0 && (
                    <div>
                        <h3 className="text-lg font-semibold mb-2">Hints</h3>
                        <div className="flex gap-2">
							{challenge.hints.map((hint: HintType) => (
								<Hint key={hint.hid} {...hint} chall_id={challenge.chall_id} />
							))}
                        </div>
                    </div>
                )}
				
				{challenge.type === ChallType.OnDemand && (
					<InstanceCard chall_id={challenge.chall_id} />
				)}
			
				<form onSubmit={handleSubmit} className="flex gap-2">
					<Input
					type="text"
					placeholder="Enter flag"
					value={flag}
					onChange={(e) => setFlag(e.target.value)}
					/>
					<Button type="submit" >Submit</Button>
				</form>
			</div>
		</DialogContent>
		</Dialog>
	);
}

