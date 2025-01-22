import { TopScore } from "@/store/scoreboardStore";

interface Submission {
	rank: number;
	teamname: string;
	timestamp: string;
	points: number;
}

interface TeamPlot {
	timestamp: string;
	[key: string]: number | string;
}

function prepareSubmissions(data: TopScore[]) {
	return data.flatMap((team) =>
		team.submissions.map((sub) => ({
			rank: team.rank,
			teamname: team.teamname,
			timestamp: sub.timestamp,
			points: sub.points,
		}))
	);
}

function buildGraphData(
	preparedData: Submission[],
	startTime: string
): TeamPlot[] {
	const scoresTillNow: { [teamname: string]: number } = {};
	preparedData.forEach(({ teamname }) => (scoresTillNow[teamname] = 0));

	const finalData = [{ timestamp: startTime, ...scoresTillNow }];
	preparedData.forEach((submission) => {
		scoresTillNow[submission.teamname] += submission.points;
		finalData.push({ timestamp: submission.timestamp, ...scoresTillNow });
	});

	return finalData;
}

export function processTopScores(data: TopScore[], startTime: string) {
	const preparedData = prepareSubmissions(data);

	preparedData.sort(
		(a, b) =>
			new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
	);

	return buildGraphData(preparedData, startTime);
}
