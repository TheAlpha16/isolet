import { useChallengeStore } from "@/store/challenge";
import type {
  CategoryProgress,
  ScoreGraphEntryType,
  ScoreGraphInputType,
  SubmissionType,
} from "@/utils/types";

interface Submission {
  label: string;
  timestamp: string;
  points: number;
}

function prepareSubmissions(data: ScoreGraphInputType[]): Submission[] {
  return data.flatMap((plot) =>
    plot.scores.map((sub) => ({
      label: plot.label,
      timestamp: sub.timestamp,
      points: sub.points,
    }))
  );
}

function buildGraphData(preparedData: Submission[], startTime: string): ScoreGraphEntryType[] {
  const scoresTillNow: { [label: string]: number } = {};
  preparedData.forEach(({ label }) => (scoresTillNow[label] = 0));

  const finalData = [{ timestamp: startTime, ...scoresTillNow }];
  preparedData.forEach((submission) => {
    scoresTillNow[submission.label] += submission.points;
    finalData.push({ timestamp: submission.timestamp, ...scoresTillNow });
  });

  return finalData;
}

export function processScores(
  data: ScoreGraphInputType[],
  startTime: string
): ScoreGraphEntryType[] {
  const preparedData = prepareSubmissions(data);

  preparedData.sort((a, b) => new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime());

  return buildGraphData(preparedData, startTime);
}

export function processCategoryData(submissions: SubmissionType[]): CategoryProgress[] {
  const { categoryIdMap, categoryChallengeMap } = useChallengeStore();
  const categoryIds = Object.keys(categoryIdMap).map(Number);
  const categoryProgress: CategoryProgress[] = [];

  categoryIds.forEach((categoryId) => {
    const total = categoryChallengeMap[categoryId].length;
    const solved = submissions.filter((sub) => {
      const challengeIds = categoryChallengeMap[categoryId] || [];
      return sub.correct && challengeIds.includes(sub.chall_id);
    }).length;

    categoryProgress.push({ category: categoryIdMap[categoryId].name, solved, total });
  });

  return categoryProgress;
}
