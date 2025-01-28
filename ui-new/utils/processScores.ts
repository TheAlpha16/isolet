import type { ScoreGraphEntryType, ScoreGraphInputType } from "@/utils/types";

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

function buildGraphData(
	preparedData: Submission[],
	startTime: string
): ScoreGraphEntryType[] {
	const scoresTillNow: { [label: string]: number } = {};
	preparedData.forEach(({ label }) => (scoresTillNow[label] = 0));

	const finalData = [{ timestamp: startTime, ...scoresTillNow }];
	preparedData.forEach((submission) => {
		scoresTillNow[submission.label] += submission.points;
		finalData.push({ timestamp: submission.timestamp, ...scoresTillNow });
	});

	return finalData;
}

export function processScores(data: ScoreGraphInputType[], startTime: string): ScoreGraphEntryType[] {
	const preparedData = prepareSubmissions(data);

	preparedData.sort(
		(a, b) =>
			new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
	);

	return buildGraphData(preparedData, startTime);
}
