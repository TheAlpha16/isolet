export interface UserType {
	userid: number;
	email: string;
	username: string;
	rank: number;
	teamid: number;
	teamname: string;
	score: number;
};

export interface TeamType {
	teamid: number;
	teamname: string;
	members: UserType[];
	captain: number;
	rank: number;
	score: number;
	submissions: SubmissionType[];
};

export enum ChallType {
	Static,
	Dynamic,
	OnDemand
};

export interface HintType {
	hid: number;
	chall_id: number;
	hint: string;
	cost: number;
	unlocked: boolean;
};

export interface ChallengeType {
	chall_id: number;
	name: string;
	prompt: string;
	type: ChallType;
	points: number;
	files: string[];
	hints: HintType[];
	solves: number;
	author: string;
	tags: string[];
	links: string[];
	done: boolean;
};

export interface InstanceType {
	chall_id: number;
	password: string;
	port: number;
	hostname: string;
	deadline: number;
	deployment: string;
	active: boolean;
};

export interface SubmissionType {
	sid: number;
	chall_name: string;
	chall_id: number;
	userid: number;
	teamid: number;
	correct: boolean;
	timestamp: string;
	points: number;
}

export interface CategoryProgress {
	category: string;
	solved: number;
	total: number;
}

export interface ScoreGraphPointType {
	timestamp: string;
	points: number;
}

export interface ScoreGraphInputType {
	label: string;
	scores: ScoreGraphPointType[];
}

export interface ScoreGraphEntryType {
	timestamp: string;
	[key: string]: number | string;
}
