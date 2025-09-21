import { ScoreboardEntry, ScoreGraphEntry } from "@/api";
import { ScoreService } from "@/services";
import { create } from "zustand";

const CACHE_TIMEOUT = 1000 * 60; // 1 minute
const PAGE_SIZE = 50;

interface CachedPage {
  entries: ScoreboardEntry[];
  timestamp: number;
}

interface ScoreStore {
  scoresLoading: boolean;
  graphLoading: boolean;
  currentPage: number;
  totalPages: number;
  scores: ScoreboardEntry[];
  graphScores: ScoreGraphEntry[];
  cachedPages: Record<number, CachedPage>;
  fetchPage: (page: number) => Promise<void>;
  preFetchPage: (page: number) => Promise<void>;
  fetchGraph: () => Promise<void>;
}

export const useScoreStore = create<ScoreStore>((set, get) => ({
  scoresLoading: false,
  graphLoading: false,
  currentPage: 1,
  totalPages: 1,
  scores: [],
  graphScores: [],
  cachedPages: {},

  fetchPage: async (page: number) => {
    const { cachedPages } = get();
    const cachedPage = cachedPages[page];

    if (cachedPage && Date.now() - cachedPage.timestamp < CACHE_TIMEOUT) {
      set({ scores: cachedPage.entries, currentPage: page });
      return;
    }

    set({ scoresLoading: true });
    try {
      const res = await ScoreService.getScoreboard(page, PAGE_SIZE);
      if (res) {
        set((state) => ({
          scores: res.entries,
          currentPage: page,
          totalPages: res.total_pages,
          cachedPages: {
            ...state.cachedPages,
            [page]: { entries: res.entries, timestamp: Date.now() },
          },
        }));
      }
    } catch {
    } finally {
      set({ scoresLoading: false });
    }
  },

  preFetchPage: async (page: number) => {
    const { totalPages, cachedPages } = get();
    if (page < 1 || page > totalPages) {
      return;
    }

    const cachedPage = cachedPages[page];
    if (cachedPage && Date.now() - cachedPage.timestamp < CACHE_TIMEOUT) {
      return;
    }

    set({ scoresLoading: true });
    try {
      const res = await ScoreService.getScoreboard(page, PAGE_SIZE);
      if (res) {
        set((state) => ({
          totalPages: res.total_pages,
          cachedPages: {
            ...state.cachedPages,
            [page]: { entries: res.entries, timestamp: Date.now() },
          },
        }));
      }
    } catch {
    } finally {
      set({ scoresLoading: false });
    }
  },

  fetchGraph: async () => {
    set({ graphLoading: true });
    try {
      const res = await ScoreService.getScoreGraph();
      if (res) {
        set({ graphScores: res.entries });
      }
    } catch {
    } finally {
      set({ graphLoading: false });
    }
  },
}));
