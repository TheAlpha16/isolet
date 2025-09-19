// Auth
export interface LoginInput {
  identifier: string; // email or username
  password: string;
}

export interface RegisterInput {
  username: string;
  email: string;
  password: string;
}

export interface ForgotPasswordInput {
  email: string;
}

export interface ResetPasswordInput {
  token: string;
  password: string;
}

// Responses
export interface Session {
  user_id: number;
  team_id: number | null;
  expires_at: number;
}

export enum ResponseStatus {
  Success = "success",
  Error = "error",
}

export interface Response<T = unknown> {
  status: ResponseStatus;
  message: string;
  data: T | null;
}

// Team
export interface CreateTeamInput {
  team_name: string;
  password: string;
}

export interface JoinTeamInput {
  team_name: string;
  password: string;
}

export interface GenerateInviteOutput {
  invite_link: string;
}

export interface Team {
  id: number;
  name: string;
  captain_id: number;
}

export interface ProfileTeam extends Team {
  score: number;
  rank: number | null;
  submissions: Submission[];
}

// User
export enum UserRole {
  Admin = "admin",
  Author = "author",
  Captain = "captain",
  Player = "player",
}

export interface User {
  id: number;
  username: string;
  email: string;
  role: UserRole;
}

export interface ProfileMe {
  user: User;
  team: Team | null;
}

// Event
export interface EventInfo {
  name: string;
  event_start: number;
  event_end: number;
  post_event: boolean;
  team_length: number;
}

// Challenge
export interface Category {
  id: number;
  name: string;
}

export interface Hint {
  id: number;
  text: string;
  cost: number;
  unlocked: boolean;
}

export enum ChallengeType {
  Static = "static",
  Dynamic = "dynamic",
  OnDemand = "on-demand",
}

export interface Challenge {
  id: number;
  name: string;
  prompt: string;
  category: Category;
  type: ChallengeType;
  points: number;
  files: string[];
  hints: Hint[];
  author: string;
  tags: string[];
  links: string[];
  max_attempts: number;
  total_solves: number;
  solved: boolean;
  attempt_count: number;
}

export interface SubmitFlagInput {
  challenge_id: number;
  flag: string;
}

export interface SubmitFlagOutput {
  is_correct: boolean;
}

export interface UnlockHintInput {
  hint_id: number;
}

// Scoreboard
export interface Submission {
  challenge_id: number;
  user_id: number;
  team_id: number;
  is_correct: boolean;
  timestamp: number;
}

export interface ScoreboardEntry {
  team_id: number;
  team_name: string;
  rank: number;
  score: number;
}

export interface Scoreboard {
  page: number;
  page_size: number;
  total_pages: number;
  entries: ScoreboardEntry[];
}

export interface ScoreRecord {
  points: number;
  timestamp: number;
}

export interface ScoreGraphEntry extends ScoreboardEntry {
  records: ScoreRecord[];
}

export interface ScoreGraph {
  count: number;
  entries: ScoreGraphEntry[];
}
