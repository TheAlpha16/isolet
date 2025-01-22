import { TopScore } from "@/store/scoreboardStore";

export function processTopScores(data: TopScore[], startTime: string) {
	let scoresTillNow = Object.fromEntries(
		data.map((team) => [team.teamname, 0])
	);
	let preparedData: {
		rank: number;
		teamname: string;
		timestamp: string;
		points: number;
	}[] = [];
	let finalData: {
		timestamp: string;
		[key: string]: number | string;
	}[] = [];

	for (let i = 0; i < data.length; i++) {
		let team = data[i];

		for (let submission of team.submissions) {
			preparedData.push({
				rank: team.rank,
				teamname: team.teamname,
				timestamp: submission.timestamp,
				points: submission.points,
			});
		}
	}

	preparedData.sort(
		(a, b) =>
			new Date(a.timestamp).getTime() - new Date(b.timestamp).getTime()
	);

	finalData.push({
		timestamp: startTime,
		...scoresTillNow,
	});

	for (let i = 0; i < preparedData.length; i++) {
		let submission = preparedData[i];
		scoresTillNow[submission.teamname] += submission.points;
		finalData.push({
			timestamp: submission.timestamp,
			...scoresTillNow,
		});
	}

	return finalData;
}
