"use client";

import React, { useState, useEffect } from "react";
import { useChallengeStore } from "@/store/challengeStore";
import { ChallengeCard } from "@/components/challenges/ChallengeCard";
import { ChallengeModal } from "@/components/challenges/ChallengeModal";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type { ChallengeType } from "@/utils/types";
import ChallengeSkeleton from "@/components/skeletons/challenge"

function Challenges() {
	const [currentChallenge, setCurrentChallenge] = useState<ChallengeType | null>(null);
	const { challenges, fetchChallenges, loading } = useChallengeStore();
	const categories = Object.keys(challenges);

	useEffect(() => {
		fetchChallenges();
	}, [fetchChallenges]);

	if (loading) {
		return <ChallengeSkeleton />;
	}

	return (
		<div className="container p-4 justify-start h-full flex flex-col">
			<Tabs defaultValue={categories[0]} className="flex flex-col w-full items-center sm:items-start">
				<TabsList className="mb-4 flex flex-wrap max-w-fit">
					{categories.map((category) => (
						<TabsTrigger key={category} value={category}>
							{category}
						</TabsTrigger>
					))}
				</TabsList>
				{categories.map((category) => (
					<TabsContent key={category} value={category}>
						<div className="flex flex-wrap gap-4">
							{challenges[category].map((challenge) => (
								<ChallengeCard
									key={challenge.chall_id}
									challenge={challenge}
									onClick={() => setCurrentChallenge(challenge)}
								/>
							))}
						</div>
					</TabsContent>
				))}
			</Tabs>
			{currentChallenge && (
				<ChallengeModal
					challenge={currentChallenge}
					onClose={() => setCurrentChallenge(null)}
				/>
			)}
		</div>
	);
}

export default Challenges;
