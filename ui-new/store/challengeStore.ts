import showToast, { ToastStatus } from "@/utils/toastHelper";
import { create } from "zustand";
import fetchTimeout from "@/utils/fetchTimeOut";

enum ChallType {
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

export interface ChallengeData {
    [category: string]: ChallengeType[];
}

interface ChallengeStore {
    challenges: ChallengeData;
    loading: boolean;
    fetchChallenges: () => void;
    submitFlag: (chall_id: number, flag: string) => void;
    // unlockHint: (chall_id: number, hint_id: number) => void;
};

export const useChallengeStore = create<ChallengeStore>((set) => ({

    challenges: {},
    loading: false,

    fetchChallenges: async () => {
        set({ loading: true });

        try {
            const res = await fetch("/api/challs");

            if (res.ok) {
                const rawChallenges = await res.json();
                const processedChallenges: ChallengeData = {};

                for (const category in rawChallenges) {                    
                    processedChallenges[category] = rawChallenges[category].map((chall: any) => ({
                        ...chall,
                        type: ChallType[chall.type.split("-").map((word: string) => word.charAt(0).toUpperCase() + word.slice(1)).join("")]
                    }));
                }

                set({ challenges: processedChallenges }); 

            } else if (res.status === 401) {
                showToast(ToastStatus.Warning, "login to continue");
            } else if (res.status === 503) {
                showToast(ToastStatus.Failure, "event has not yet started");
            } else {
                showToast(ToastStatus.Failure, "failed to fetch challenges");
            }
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

    submitFlag: async (chall_id, flag) => {
        flag = flag.trim();

        if (!flag) {
            showToast(ToastStatus.Warning, "flag cannot be empty");
            return;
        }

        let formData = new FormData();
        formData.append("chall_id", chall_id.toString());
        formData.append("flag", flag);

        try {
            const res = await fetchTimeout("/api/submit", 7000, new AbortController().signal, {
                method: "POST",
                body: formData,
            });

            if (res.ok) {
                showToast(ToastStatus.Success, "correct flag!");
                set((state) => {
                    const updatedChallenges = { ...state.challenges };
                    for (const category in updatedChallenges) {
                        const challenge = updatedChallenges[category].find((c) => c.chall_id === chall_id);
                        if (challenge) {
                            challenge.solves++;
                            challenge.done = true;
                            break;
                        }
                    }

                    return { challenges: updatedChallenges };
                })
            } else if (res.status === 401) {
                showToast(ToastStatus.Warning, "login to continue");
            } else {
                const response = await res.json();
                showToast(ToastStatus.Failure, response.message);
            };

        } catch (error: any) {
            if (error.name === "AbortError") {
                showToast(ToastStatus.Failure, "request timed out!");
            } else {
                showToast(ToastStatus.Warning, "seems offline");
            }
        }
    },

    // unlockHint: async (chall_id, hint_id) => {
    //     const res = await fetch("/api/unlock", {
    //         method: "POST",
    //         body: JSON.stringify({ chall_id, hint_id }),
    //     });
    //     if (res.ok) {
    //         set((state) => {
    //             const challenge = state.challenges.find((c) => c.chall_id === chall_id);
    //             if (challenge) {
    //                 const hint = challenge.hints.find((h) => h.hid === hint_id);
    //                 if (hint) {
    //                     hint.unlocked = true;
    //                 }
    //             }
    //         });
    //     }
    // },
}));