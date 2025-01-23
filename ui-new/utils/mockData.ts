import type {
	UserType,
	TeamType,
	SubmissionType,
	CategoryProgress,
} from "./types";
import { subDays, format } from "date-fns";

// Generate a mock user
export const generateMockUser = (userid: number, teamid: number, rank: number): UserType => ({
	userid: userid,
	username: `user${userid}`,
	email: `user${userid}@example.com`,
	rank: rank,
	teamid: teamid,
	teamname: `Team ${teamid}`,
	score: Math.floor(Math.random() * 1000),
});

// Generate a mock team
export const generateMockTeam = (
	teamid: number,
	captainId: number
): TeamType => {
	const members: UserType[] = [
		generateMockUser(captainId, teamid, 2), // Captain
		// generate random number of team members
		...Array.from({ length: Math.floor(Math.random() * 3) }, (_, i) =>
			generateMockUser(i + 2, teamid, 3)
		),
	];
	return {
		teamid: teamid,
		teamname: `Team ${teamid}`,
		members,
		captain: captainId,
		rank: Math.floor(Math.random() * 500) + 1,
		score: members.reduce((sum, member) => sum + member.score, 0),
		submissions: generateMockSubmissions(
			Math.floor(Math.random() * 10) + 5
		),
	};
};

// Generate mock submissions
export const generateMockSubmissions = (count: number): SubmissionType[] => {
	return Array.from({ length: count }, (_, i) => ({
		sid: i + 1,
		chall_name: `Challenge ${i + 1}`,
		chall_id: Math.floor(Math.random() * 100) + 1,
		userid: Math.floor(Math.random() * 100) + 1,
		teamid: Math.floor(Math.random() * 10) + 1,
		correct: Math.random() > 0.3,
		timestamp: format(
			subDays(new Date(), Math.floor(Math.random() * 7)),
			"yyyy-MM-dd'T'HH:mm:ss'Z'"
		),
		points: Math.floor(Math.random() * 500) + 100,
	}));
};

// Generate mock category progress
export const generateMockCategoryProgress = (): CategoryProgress[] => {
	const categories = ["Web", "Crypto", "Pwn", "Reverse", "Forensics"];
	return categories.map((category) => ({
		category,
		solved: Math.floor(Math.random() * 10),
		total: Math.floor(Math.random() * 10) + 10,
	}));
};
