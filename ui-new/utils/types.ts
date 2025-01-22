export interface User {
	id: string;
	username: string;
	email: string;
	avatarUrl: string;
	totalPoints: number;
}

export interface TeamMember extends User {
	role: "captain" | "member";
}

export interface Team {
	id: string;
	name: string;
	members: TeamMember[];
	rank: number;
	totalPoints: number;
}

export interface Submission {
	id: number;
	challengeName: string;
	category: string;
	points: number;
	timestamp: string;
	isCorrect: boolean;
}

export interface CategoryProgress {
	category: string;
	solved: number;
	total: number;
}
