import showToast, { ToastStatus } from "@/utils/toastHelper";
import { create } from "zustand";
import fetchTimeout from "@/utils/fetchTimeOut";

export interface TeamType {
	teamid: number;
	teamname: string;
	score: number;
	rank: number;
}

interface CachedPage {
	data: TeamType[];
	timestamp: number;
}

interface ScoreboardStore {
	scores: TeamType[];
	pages: Record<number, CachedPage>;
	currentPage: number;
	totalPages: number;
	loading: boolean;
	fetchPage: (page: number) => void;
	prefetchPage: (page: number) => void;
}

export const useScoreboardStore = create<ScoreboardStore>((set) => ({
	scores: [],
	pages: {},
	currentPage: 1,
	totalPages: 1,
	loading: false,

	fetchPage: async (page: number) => {
		set({ loading: true });

		const { pages } = useScoreboardStore.getState();
		const cacheTimeOut = 1000 * 60;
		const cachedPage = pages[page];

		if (cachedPage && Date.now() - cachedPage.timestamp < cacheTimeOut) {
			set({ scores: cachedPage.data, currentPage: page, loading: false });
		}

		try {
			const res = await fetchTimeout(
				`/api/scoreboard?page=${page}`,
				10000,
				new AbortController().signal,
				{
					method: "GET",
				}
			);

			if (res.ok) {
				const data = await res.json();

				await set((state) => {
					const newPages = {
						...state.pages,
						[page]: {
							data: data.scores,
							timestamp: Date.now(),
						},
					};

					return {
						scores: data.scores,
						pages: newPages,
						currentPage: page,
						totalPages: data.page_count,
						loading: false,
					};
				});
			} else if (res.status === 401) {
				showToast(ToastStatus.Failure, "login to continue");
			} else {
				const response = await res.json();
				showToast(ToastStatus.Failure, response.message);
			}

			useScoreboardStore.getState().prefetchPage(page + 1);
			useScoreboardStore.getState().prefetchPage(page - 1);

		} catch (error: any) {
			if (error.name === "AbortError") {
				showToast(ToastStatus.Failure, "request timed out!");
			} else {
				showToast(ToastStatus.Warning, "seems offline");
			}
		} finally {
			set({ loading: false });
		}
	},

	prefetchPage: async (page: number) => {
		const { pages, totalPages } = useScoreboardStore.getState();

        if (page < 1 || page > totalPages || pages[page]) return;
        
        try {
            const res = await fetchTimeout(`/api/scoreboard?page=${page}`, 10000, new AbortController().signal, {
                method: "GET",
            });

            if (res.ok) {
                const data = await res.json();

                set((state) => {
                    const newPages = {
                        ...state.pages,
                        [page]: {
                            data: data.scores,
                            timestamp: Date.now(),
                        },
                    };

                    return {
                        pages: newPages,
                    };
                });
            }
        } finally {
            // do nothing
        }
	},
}));
