"use client";

import type { ScoreGraphEntry, ScoreRecord, Submission } from "@/api";
import { CategoryProgress } from "@/components/charts/CategoryProgress";
import { ScoreGraph } from "@/components/charts/ScoreGraph";
import { SubmissionStats } from "@/components/charts/SubmissionStats";
import { TeamManagement } from "@/components/profile/TeamManagement";
import { TeamProfile } from "@/components/profile/TeamProfile";
import { UserProfile } from "@/components/profile/UserProfile";
import { ChartSkeleton } from "@/components/skeletons/profile";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import type {
  CategoryProgress as CatProgress,
  ChartData,
  SubmissionStats as SubStats,
} from "@/models/score";
import { useChallengeStore, useProfileStore } from "@/store";
import { processCategoryData } from "@/utils/categoryProgress";
import { toChartData } from "@/utils/scoreTransform";
import { submissionStats } from "@/utils/submissionStats";
import { useEffect, useState } from "react";

export default function ProfilePage() {
  const [activeTab, setActiveTab] = useState("user");
  const [teamCategoryProgress, setTeamCategoryProgress] = useState<CatProgress[]>([]);
  const [userCategoryProgress, setuserCategoryProgress] = useState<CatProgress[]>([]);
  const [teamSubmissionStats, setTeamSubmissionStats] = useState<SubStats>({
    correct: 0,
    incorrect: 0,
  });
  const [userSubmissionStats, setUserSubmissionStats] = useState<SubStats>({
    correct: 0,
    incorrect: 0,
  });
  const [teamGraph, setTeamGraph] = useState<ChartData>({ labels: [], points: [] });
  // const [solves, setSolves] = useState<ScoreGraphEntry[]>([]);
  const { user, team, submissions, score, rank, members, teamLoading, fetchTeam } =
    useProfileStore();
  const { challengeIdMap, categoryIdMap, categoryChallengeMap, loading, fetchChallenges } = useChallengeStore();

  useEffect(() => {
    const fetchData = async () => {
      await fetchTeam();
      await fetchChallenges();
      const userSubmissions: Submission[] = [];
      const teamScoreRecords: ScoreRecord[] = [];

      // filter solves and submissions
      for (let i = 0; i < submissions.length; i++) {
        if (submissions[i].user_id === user!.id) {
          userSubmissions.push(submissions[i]);
        }
        if (submissions[i].is_correct) {
          teamScoreRecords.push({
            timestamp: submissions[i].timestamp,
            points: challengeIdMap[submissions[i].challenge_id]?.points || 0,
          });
        }
      }

      // score graph entry
      const teamScoreEntry: ScoreGraphEntry = {
        team_id: team!.id,
        team_name: team!.name,
        rank: rank || 0,
        score: score,
        records: teamScoreRecords,
      };

      // process submission stats
      setTeamCategoryProgress(
        processCategoryData(submissions, categoryIdMap, categoryChallengeMap)
      );
      setuserCategoryProgress(
        processCategoryData(userSubmissions, categoryIdMap, categoryChallengeMap)
      );
      setTeamSubmissionStats(submissionStats(submissions));
      setUserSubmissionStats(submissionStats(userSubmissions));
      setTeamGraph(toChartData([teamScoreEntry]));
    };

    fetchData();
  }, [fetchTeam]);

  return (
    <div className="container mx-auto p-4 space-y-8">
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-2">
          <TabsTrigger value="user">User</TabsTrigger>
          <TabsTrigger value="team">Team</TabsTrigger>
        </TabsList>
        <TabsContent value="user" className="space-y-4">
          <UserProfile user={user!} team={team!} score={score} rank={rank} />
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {(teamLoading || loading) ? <ChartSkeleton /> : <SubmissionStats {...userSubmissionStats} />}
            {(teamLoading || loading) ? (
              <ChartSkeleton />
            ) : (
              <CategoryProgress categories={userCategoryProgress} title="Progress" />
            )}
          </div>
        </TabsContent>
        <TabsContent value="team" className="space-y-4">
          <TeamProfile team={team!} members={members} score={score} rank={rank} />
          <TeamManagement user={user!} members={members} />
          {(teamLoading || loading) ? <ChartSkeleton /> : <ScoreGraph data={teamGraph} />}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {(teamLoading || loading) ? <ChartSkeleton /> : <SubmissionStats {...teamSubmissionStats} />}
            {(teamLoading || loading) ? (
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
