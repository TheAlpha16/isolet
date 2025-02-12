import { create } from "zustand";
import type { TeamType, ScoreGraphEntryType, CategoryProgress, SubmissionType } from "@/utils/types";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { useMetadataStore } from "@/store/metadataStore";
import { processScores, processCategoryData } from "@/utils/processScores";
import { useAuthStore } from "@/store/authStore";

interface CorrectVIncorrect {
	correct: number;
	incorrect: number;
}
interface ProfileStore {
	teamLoading: boolean;
	team: TeamType;
	userSubmissionsProgress: CorrectVIncorrect;
	teamSubmissionsProgress: CorrectVIncorrect;
	teamGraph: ScoreGraphEntryType[];
	userCategoryProgress: CategoryProgress[];
	teamCategoryProgress: CategoryProgress[];
	fetchSelfTeam(): void;
}

export const useProfileStore = create<ProfileStore>((set) => ({
	teamLoading: true,
	team: {
		teamid: -1,
		teamname: "",
		members: [],
		captain: -1,
		rank: -1,
		score: -1,
		submissions: [],
	},
	userSubmissionsProgress: {
		correct: 0,
		incorrect: 0,
	},
	teamSubmissionsProgress: {
		correct: 0,
		incorrect: 0,
	},
	teamGraph: [],
	userCategoryProgress: [],
	teamCategoryProgress: [],

	fetchSelfTeam: async () => {
		set({ teamLoading: true });

		try {
			const res = await fetch("/api/profile/team/self");

			if (res.ok) {
				const response = await res.json();

				const team = response.team;
				const submissions = response.submissions;
				const rank = response.rank;
				const score = response.score;

				const teamData = {
					label: team.teamname,
					scores: submissions.filter((sub: SubmissionType) => sub.correct === true).map((sub: SubmissionType) => ({
						timestamp: sub.timestamp,
						points: sub.points,
					})),
				};

				const teamGraph = processScores([teamData], useMetadataStore.getState().eventStart);
				const userSubmissions = submissions.filter((sub: SubmissionType) => sub.userid === useAuthStore.getState().user.userid);

				var userCorrect = userSubmissions.filter((sub: SubmissionType) => sub.correct === true).length;
				var userScore = userSubmissions.filter((sub: SubmissionType) => sub.correct === true && sub.userid === useAuthStore.getState().user.userid).reduce((acc: number, sub: SubmissionType) => acc + sub.points, 0);

				var teamCorrect = submissions.filter((sub: SubmissionType) => sub.correct === true).length;

				useAuthStore.getState().user.score = userScore;
				team.score = score;
				team.rank = rank;

				const userCategoryProgress = processCategoryData(userSubmissions);
				const teamCategoryProgress = processCategoryData(submissions);

				set({ team: { ...team, submissions: submissions }, teamGraph: teamGraph, userSubmissionsProgress: { correct: userCorrect, incorrect: userSubmissions.length - userCorrect }, teamSubmissionsProgress: { correct: teamCorrect, incorrect: submissions.length - teamCorrect }, userCategoryProgress: userCategoryProgress, teamCategoryProgress: teamCategoryProgress });
			} else {
				const response = await res.json();
				showToast(ToastStatus.Failure, response.message);
			}

		} catch (error: any) {
			if (error.name === "AbortError") {
				showToast(ToastStatus.Failure, "request timed out!");
			} else {
				showToast(ToastStatus.Warning, "seems offline");
			}
		} finally {
			set({ teamLoading: false });
		}
	},
}));
