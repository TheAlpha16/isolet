import React from "react";
import type { ChallengeType } from "@/utils/types";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Trophy, Flag, Check } from 'lucide-react';

interface ChallengeCardProps {
	challenge: ChallengeType;
	onClick: () => void;
}

export function ChallengeCard({ challenge, onClick }: ChallengeCardProps) {
	return (
		<Card className={`hover:shadow-lg dark:hover:shadow-zinc-900 transition-shadow w-[300px]`}>
			<CardHeader>
				<CardTitle className="flex justify-between items-center">
					<span className="truncate">{challenge.name}</span>
					<Badge variant={challenge.done ? "secondary" : "default"}>{challenge.points} pts</Badge>
				</CardTitle>
			</CardHeader>
			<CardContent>
				<p className="text-sm text-muted-foreground mb-2 truncate">{challenge.prompt}</p>
				<div className="flex flex-wrap gap-1">
					{challenge.tags.map((tag) => (
						<Badge key={tag} variant="outline">
						{tag}
						</Badge>
					))}
				</div>
			</CardContent>
			<CardFooter className="flex justify-between items-center">
				<div className="flex items-center text-sm text-muted-foreground">
					<Trophy className="w-4 h-4 mr-1" />
					<span>{challenge.solves} solves</span>
				</div>
				<Button onClick={onClick} size="sm" variant={challenge.done ? "secondary" : "default"}>
					{challenge.done ? (
						<>
							<Check className="w-4 h-4 mr-2 text-green-600 dark:text-green-500" />
							Solved
						</>
					) : (
						<>
							<Flag className="w-4 h-4 mr-2" />
							Solve
						</>
					)}
				</Button>
			</CardFooter>
		</Card>
	);
}

