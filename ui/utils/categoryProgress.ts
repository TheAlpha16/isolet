import { Submission } from "@/api";
import { CategoryProgress } from "@/models/score";
import { useChallengeStore } from "@/store/challenge";

export function processCategoryData(submissions: Submission[]): CategoryProgress[] {
  const { categoryIdMap, categoryChallengeMap } = useChallengeStore();
  const categoryIds = Object.keys(categoryIdMap).map(Number);
  const categoryProgress: CategoryProgress[] = [];

  categoryIds.forEach((categoryId) => {
    const total = categoryChallengeMap[categoryId].length;
    const solved = submissions.filter((sub) => {
      const challengeIds = categoryChallengeMap[categoryId] || [];
      return sub.is_correct && challengeIds.includes(sub.challenge_id);
    }).length;

    categoryProgress.push({ category: categoryIdMap[categoryId].name, solved, total });
  });

  return categoryProgress;
}
