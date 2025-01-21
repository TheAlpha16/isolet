"use client";

import React, { useState, useEffect } from "react";
import { ChallengeType, useChallengeStore } from "@/store/challengeStore";
import { ChallengeCard } from "@/components/challenges/ChallengeCard";
import { ChallengeModal } from "@/components/challenges/ChallengeModal";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

function Challenges() {
	const [currentChallenge, setCurrentChallenge] = useState<ChallengeType | null>(null);
	const { challenges, fetchChallenges, loading } = useChallengeStore();
	const categories = Object.keys(challenges);

	useEffect(() => {
		fetchChallenges();
	}, []);

	if (loading) {
		return (
		<div className="container mx-auto p-4">
			<h1 className="text-3xl font-bold mb-6">Challenges</h1>
			<p className="text-center text-lg">Loading challenges...</p>
		</div>
    );
}

	return (
		<div className="container p-4 items-center justify-start h-full flex flex-col">
			<h1 className="text-3xl font-bold mb-6">Challenges</h1>
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
