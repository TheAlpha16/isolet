"use client";

import type { ScoreGraphEntry } from "@/api";
import { CategoryProgress } from "@/components/charts/CategoryProgress";
import { ScoreGraph } from "@/components/charts/ScoreGraph";
import { SubmissionStats } from "@/components/charts/SubmissionStats";
import { TeamManagement } from "@/components/profile/TeamManagement";
import { TeamProfile } from "@/components/profile/TeamProfile";
import { UserProfile } from "@/components/profile/UserProfile";
import { ChartSkeleton } from "@/components/skeletons/profile";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useChallengeStore, useProfileStore } from "@/store";
import { processCategoryData } from "@/utils/categoryProgress";
import { toChartData } from "@/utils/scoreTransform";
import { submissionStats } from "@/utils/submissionStats";
import { useEffect, useMemo, useState } from "react";

export default function ProfilePage() {
  const [activeTab, setActiveTab] = useState("user");

  const { user, team, submissions, score, rank, members, teamLoading, fetchTeam } =
    useProfileStore();

  const {
    challengeIdMap,
    categoryIdMap,
    categoryChallengeMap,
    loading: challengeLoading,
    fetchChallenges,
  } = useChallengeStore();

  // Initial fetch on mount
  useEffect(() => {
    (async () => {
      await fetchTeam();
      await fetchChallenges();
    })();
  }, [fetchTeam, fetchChallenges]);

  // Filter submissions
  const userSubmissions = useMemo(
    () => submissions.filter((s) => s.user_id === user?.id),
    [submissions, user]
  );

  // Build team score entry
  const teamScoreEntry = useMemo<ScoreGraphEntry>(
    () => ({
      team_id: team?.id ?? 0,
      team_name: team?.name ?? "",
      rank: rank ?? 0,
      score,
      records: submissions
        .filter((s) => s.is_correct)
        .map((s) => ({
          timestamp: s.timestamp,
          points: challengeIdMap[s.challenge_id]?.points || 0,
        })),
    }),
    [submissions, challengeIdMap, team, rank, score]
  );

  // Derived data
  const teamGraph = useMemo(() => toChartData([teamScoreEntry]), [teamScoreEntry]);

  const teamCategoryProgress = useMemo(
    () => processCategoryData(submissions, categoryIdMap, categoryChallengeMap),
    [submissions, categoryIdMap, categoryChallengeMap]
  );

  const userCategoryProgress = useMemo(
    () => processCategoryData(userSubmissions, categoryIdMap, categoryChallengeMap),
    [userSubmissions, categoryIdMap, categoryChallengeMap]
  );

  const teamSubmissionStats = useMemo(() => submissionStats(submissions), [submissions]);

  const userSubmissionStats = useMemo(() => submissionStats(userSubmissions), [userSubmissions]);

  const loading = teamLoading || challengeLoading;

  return (
    <div className="container mx-auto p-4 space-y-8">
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="user">User</TabsTrigger>
          <TabsTrigger value="team">Team</TabsTrigger>
        </TabsList>

        {/* USER TAB */}
        <TabsContent value="user" className="space-y-4">
          <UserProfile user={user!} team={team!} score={score} rank={rank} />
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {loading ? <ChartSkeleton /> : <SubmissionStats {...userSubmissionStats} />}
            {loading ? (
              <ChartSkeleton />
            ) : (
              <CategoryProgress categories={userCategoryProgress} title="Progress" />
            )}
          </div>
        </TabsContent>

        {/* TEAM TAB */}
        <TabsContent value="team" className="space-y-4">
          <TeamProfile team={team!} members={members} score={score} rank={rank} />
          <TeamManagement user={user!} members={members} />
          {loading ? <ChartSkeleton /> : <ScoreGraph data={teamGraph} />}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {loading ? <ChartSkeleton /> : <SubmissionStats {...teamSubmissionStats} />}
            {loading ? (
              <ChartSkeleton />
            ) : (
              <CategoryProgress categories={teamCategoryProgress} title="Progress" />
            )}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}
