import type { ScoreGraphEntry } from "@/api";
import { ChartData, ChartPoint } from "@/models/score";

/**
 * Convert API ScoreGraph into chart-friendly data
 */
export function toChartData(entries: ScoreGraphEntry[]): ChartData {
  if (!entries || entries.length === 0) {
    return { labels: [], points: [] };
  }

  // Step 1: Collect teams
  const teams = entries.map((entry) => entry.team_name);
  const scores: Record<string, number> = {};
  teams.forEach((t) => (scores[t] = 0));

  // Step 2: Collect all records
  type Submission = { team: string; timestamp: number; points: number };
  const submissions: Submission[] = [];

  entries.forEach((entry) => {
    entry.records.forEach((rec) => {
      submissions.push({
        team: entry.team_name,
        timestamp: rec.timestamp,
        points: rec.points,
      });
    });
  });

  // Step 3: Sort by timestamp
  submissions.sort((a, b) => a.timestamp - b.timestamp);

  // Step 4: Build chart points
  const points: ChartPoint[] = [];

  // Initialize with a zero baseline (at first timestamp - optional)
  if (submissions.length > 0) {
    points.push({
      timestamp: submissions[0].timestamp, // unix timestamp
      ...scores,
    });
  }

  submissions.forEach((sub) => {
    scores[sub.team] += sub.points;
    points.push({
      timestamp: sub.timestamp,
      ...scores,
    });
  });

  return {
    labels: teams,
    points,
  };
}
