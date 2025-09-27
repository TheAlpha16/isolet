import { Category, Challenge } from "@/api";
import { showHint } from "@/components/hints/HintToastContainer";
import { ChallengeService } from "@/services";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { create } from "zustand";

interface ChallengeStore {
  loading: boolean;
  challengeIdMap: Record<number, Challenge>;
  categoryIdMap: Record<number, Category>;
  categoryChallengeMap: Record<number, number[]>;
  fetchChallenges: () => Promise<void>;
  submitFlag: (challenge_id: number, flag: string) => Promise<void>;
  unlockHint: (challenge_id: number, hint_id: number) => Promise<void>;
}

export const useChallengeStore = create<ChallengeStore>((set, get) => ({
  loading: false,
  challengeIdMap: {},
  categoryIdMap: {},
  categoryChallengeMap: {},

  fetchChallenges: async () => {
    set({ loading: true });
    try {
      const challenges = await ChallengeService.getChallenges();
      if (!challenges) {
        return;
      }
      const challengeIdMap: Record<number, Challenge> = {};
      const categoryIdMap: Record<number, Category> = {};
      const categoryChallengeMap: Record<number, number[]> = {};

      // organize data
      for (const challenge of challenges) {
        const category = challenge.category;

        // store challenge
        challengeIdMap[challenge.id] = challenge;

        // store category
        categoryIdMap[category.id] = category;

        // link challenge -> category
        if (!categoryChallengeMap[category.id]) {
          categoryChallengeMap[category.id] = [];
        }
        categoryChallengeMap[category.id].push(challenge.id);
      }

      // update store
      set({ challengeIdMap, categoryIdMap, categoryChallengeMap });
    } catch {
    } finally {
      set({ loading: false });
    }
  },

  submitFlag: async (challenge_id: number, flag: string) => {
    flag = flag.trim();
    if (!flag) {
      showToast(ToastStatus.Warning, "flag cannot be empty");
      return;
    }
    const challenge = get().challengeIdMap[challenge_id];
    if (!challenge) {
      showToast(ToastStatus.Failure, "challenge not found!");
      return;
    }
    if (challenge.solved) {
      showToast(ToastStatus.Failure, "challenge already solved!");
      return;
    }

    try {
      const res = await ChallengeService.submitFlag({ challenge_id, flag });
      if (res) {
        if (res.is_correct) {
          showToast(ToastStatus.Success, "correct flag!");
        } else {
          showToast(ToastStatus.Failure, "incorrect flag!");
        }
      }

      // update attempts and solved status
      challenge.attempt_count++;
      if (res?.is_correct) {
        challenge.solved = true;
        challenge.total_solves++;
      }

      set((state) => ({
        challengeIdMap: {
          ...state.challengeIdMap,
          [challenge_id]: challenge,
        },
      }));
    } catch {}
  },

  unlockHint: async (challenge_id: number, hint_id: number) => {
    const challenge = get().challengeIdMap[challenge_id];
    const localHint = challenge.hints.find((h) => h.id === hint_id);
    if (!localHint) {
      showToast(ToastStatus.Failure, "hint not found!");
      return;
    }
    if (localHint.unlocked) {
      showToast(ToastStatus.Warning, "hint already unlocked!");
      return;
    }

    try {
      const hint = await ChallengeService.unlockHint({ hint_id });
      if (hint) {
        const updatedHints = challenge.hints.map((h) => (h.id === hint.id ? hint : h));
        challenge.hints = updatedHints;
        set((state) => ({
          challengeIdMap: {
            ...state.challengeIdMap,
            [challenge_id]: challenge,
          },
        }));
        showHint(hint.text);
      }
    } catch {}
  },
}));
