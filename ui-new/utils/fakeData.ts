import { subDays, addHours, format } from "date-fns";

export interface Team {
	id: number;
	name: string;
	score: number;
	rank: number;
}

export interface ScoreHistory {
	timestamp: string;
	[teamName: string]: number | string;
}

export function generateFakeData(
	numTeams: number,
	numDays: number
): { teams: Team[]; scoreHistory: ScoreHistory[] } {
	const teams: Team[] = [];
	const scoreHistory: ScoreHistory[] = [];

	// Generate teams
	for (let i = 1; i <= numTeams; i++) {
		teams.push({
			id: i,
			name: `Team ${i}`,
			score: 0,
			rank: 0,
		});
	}

	// Generate score history
	const endDate = new Date();
	const startDate = subDays(endDate, numDays);
	const totalHours = numDays * 24;
	const submissionCount = 20;
	const hoursBetweenSubmissions = Math.floor(totalHours / submissionCount);

	for (let s = 0; s < submissionCount; s++) {
		const timestamp = format(
			addHours(startDate, s * hoursBetweenSubmissions),
			"yyyy-MM-dd HH:mm"
		);
		const historyEntry: ScoreHistory = { timestamp };

		teams.forEach((team) => {
			const scoreIncrease = Math.floor(Math.random() * 100) + 50; // Random score increase between 50 and 149
			team.score += scoreIncrease;
			historyEntry[team.name] = team.score;
		});

		scoreHistory.push(historyEntry);
	}

	// Sort teams by final score
	teams.sort((a, b) => b.score - a.score);

	for (let i = 0; i < teams.length; i++) {
		teams[i].rank = i + 1;
	}

	return { teams, scoreHistory };
}

export function getPagedTeams(
	teams: Team[],
	page: number,
	pageSize: number
): Team[] {
	const start = (page - 1) * pageSize;
	const end = start + pageSize;
	return teams.slice(start, end);
}

export function searchTeams(teams: Team[], query: string): Team[] {
	const lowercaseQuery = query.toLowerCase();
	return teams.filter((team) =>
		team.name.toLowerCase().includes(lowercaseQuery)
	);
}
