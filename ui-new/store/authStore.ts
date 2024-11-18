import { create } from "zustand";
import fetchTimeout from "@/utils/fetchTimeOut";
import showToast, { ToastStatus } from "@/utils/toastHelper";

interface AuthState {
    loggedIn: boolean;
    user: {
        userid: number;
        email: string;
        username: string;
        rank: number;
        teamid: number;
    } | null;
    setLoggedIn: (status: boolean) => void;
    setUser: (user: AuthState["user"]) => void;
    logout: () => void;
    fetchUser: () => void;
};

export const useAuthStore = create<AuthState>((set) => ({
    loggedIn: false,
    user: null, 
    setLoggedIn: (status) => set({ loggedIn: status }),
    setUser: (user) => set({ user }),

    logout: async () => {
        await fetchTimeout("/api/logout", 5000, new AbortController().signal);
        set({ loggedIn: false, user: null });
    },

    fetchUser: async () => {
        try {
            const res = await fetchTimeout("/api/identify", 5000, new AbortController().signal);
            if (res.ok) {
                const user = await res.json();
                set({ user, loggedIn: true });
                showToast(ToastStatus.Success, `Welcome back ${user.username}!`);
            } else {
                set({ loggedIn: false, user: null });
            }
        } catch (error: any) {
            set({ loggedIn: false, user: null });

            if (error.name === "AbortError") {
                showToast(ToastStatus.Failure, "verification timed out, reload!");
            } else {
                showToast(ToastStatus.Warning, "seems offline");
            }
        }
    },
}));