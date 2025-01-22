import type { User, Team, Submission, CategoryProgress } from "./types";
import { subDays, format } from "date-fns";

export const generateMockUser = (id: string): User => ({
	id,
	username: `user${id}`,
	email: `user${id}@example.com`,
	avatarUrl: `/placeholder.svg?height=100&width=100`,
	totalPoints: Math.floor(Math.random() * 1000),
});

export const generateMockTeam = (id: string, captainId: string): Team => {
	const members = [
		{ ...generateMockUser(captainId), role: "captain" as const },
		{
			...generateMockUser(`${Number.parseInt(id) + 1}`),
			role: "member" as const,
		},
		{
			...generateMockUser(`${Number.parseInt(id) + 2}`),
			role: "member" as const,
		},
	];
	return {
		id,
		name: `Team ${id}`,
		members,
		rank: Math.floor(Math.random() * 100) + 1,
		totalPoints: members.reduce(
			(sum, member) => sum + member.totalPoints,
			0
		),
	};
};

export const generateMockSubmissions = (count: number): Submission[] => {
	const categories = ["Web", "Crypto", "Pwn", "Reverse", "Forensics"];
	return Array.from({ length: count }, (_, i) => ({
		id: i + 1,
		challengeName: `Challenge ${i + 1}`,
		category: categories[Math.floor(Math.random() * categories.length)],
		points: Math.floor(Math.random() * 500) + 100,
		timestamp: format(
			subDays(new Date(), Math.floor(Math.random() * 7)),
			"yyyy-MM-dd'T'HH:mm:ss'Z'"
		),
		isCorrect: Math.random() > 0.3,
	}));
};

export const generateMockCategoryProgress = (): CategoryProgress[] => {
	const categories = ["Web", "Crypto", "Pwn", "Reverse", "Forensics"];
	return categories.map((category) => ({
		category,
		solved: Math.floor(Math.random() * 10),
		total: Math.floor(Math.random() * 10) + 10,
	}));
};
