import { Submission } from "@/api";
import type { SubmissionStats } from "@/models/score";

export function submissionStats(submissions: Submission[]): SubmissionStats {
  const stats: SubmissionStats = { correct: 0, incorrect: 0 };

  submissions.forEach((submission) => {
    if (submission.is_correct) {
      stats.correct += 1;
    } else {
      stats.incorrect += 1;
    }
  });

  return stats;
}
