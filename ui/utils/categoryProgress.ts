import { Category, Submission } from "@/api";
import { CategoryProgress } from "@/models/score";

export function processCategoryData(
  submissions: Submission[],
  categoryIdMap: Record<number, Category>,
  categoryChallengeMap: Record<number, number[]>
): CategoryProgress[] {
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
