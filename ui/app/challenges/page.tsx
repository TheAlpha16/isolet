"use client";

import { Challenge } from "@/api";
import { ChallengeCard } from "@/components/challenges/ChallengeCard";
import { ChallengeModal } from "@/components/challenges/ChallengeModal";
import { ChallengeSkeleton } from "@/components/skeletons/challenge";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useChallengeStore, useInstanceStore } from "@/store";
import { useEffect, useState } from "react";

function Challenges() {
  const [currentChallenge, setCurrentChallenge] = useState<Challenge | null>(null);
  const { challengeIdMap, categoryIdMap, categoryChallengeMap, fetchChallenges, loading } =
    useChallengeStore();
  const { fetchInstances } = useInstanceStore();
  const categoryIds = Object.keys(categoryIdMap).map(Number);

  useEffect(() => {
    fetchChallenges();
    fetchInstances();
  }, [fetchChallenges, fetchInstances]);

  if (loading) {
    return <ChallengeSkeleton />;
  }

  return (
    <div className="container p-4 justify-start h-full flex flex-col">
      {categoryIds.length !== 0 ? (
        <Tabs
          defaultValue={categoryIdMap[categoryIds[0]].name}
          className="flex flex-col w-full items-center sm:items-start"
        >
          <TabsList className="mb-4 flex flex-wrap max-w-fit">
            {categoryIds.map((categoryId) => (
              <TabsTrigger key={categoryId} value={categoryIdMap[categoryId].name}>
                {categoryIdMap[categoryId].name}
              </TabsTrigger>
            ))}
          </TabsList>
          {categoryIds.map((categoryId) => (
            <TabsContent key={categoryId} value={categoryIdMap[categoryId].name}>
              <div className="flex flex-wrap gap-4">
                {(categoryChallengeMap[categoryId] || []).map((challengeId) => {
                  const challenge = challengeIdMap[challengeId];
                  return (
                    <ChallengeCard
                      key={challenge.id}
                      challenge={challenge}
                      onClick={() => setCurrentChallenge(challenge)}
                    />
                  );
                })}
              </div>
            </TabsContent>
          ))}
        </Tabs>
      ) : (
        <div className="flex justify-center items-center h-full">
          <p className="text-2xl text-gray-500">No challenges available</p>
        </div>
      )}
      {currentChallenge && (
        <ChallengeModal challenge={currentChallenge} onClose={() => setCurrentChallenge(null)} />
      )}
    </div>
  );
}

export default Challenges;
