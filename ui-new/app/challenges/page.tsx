"use client";

import React, { useEffect, useState } from "react";
import { ChallengeType, useChallengeStore } from "@/store/challengeStore";
import { StaticChallenge } from "@/components/challenges/Challenge";

function Challenges() {
	const { challenges, fetchChallenges, loading } = useChallengeStore();
	const [current, setCurrent] = useState<ChallengeType | null>(null);

	const closeChallenge = () => {
		setCurrent(null);
	}

	useEffect(() => {
		fetchChallenges();
	}, []);

	return (
		<>
			<div className={`flex flex-col gap-4 p-4 ${current ? "blur-sm": ""}`}>
				{loading ? (<div>Loading...</div>) : (
					(
						Object.keys(challenges).map((category: string) => {
							return (
								<div key={category} className="flex flex-col gap-2 p-2">
									<div>{category}</div>
									<div className="flex gap-2 flex-wrap">
										{challenges[category].map((challenge) => {
											return (
												<StaticChallenge
													key={challenge.chall_id}
													challenge={challenge}
													onClick={() => setCurrent(challenge)}
												/>
											)
										})}
									</div>
								</div>
							)
					}))
				)}
			</div>
			{current && (
				<div className="fixed top-0 left-0 w-full h-full bg-black bg-opacity-50 flex justify-center items-center" onClick={closeChallenge}>
					<StaticChallenge 
						challenge={current} 
						isFocussed={true} 
						closeChallenge={closeChallenge}
					/>
				</div>
			)}
		</>
	)
}

export default Challenges;
