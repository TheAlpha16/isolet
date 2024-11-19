import { create } from "zustand";
import fetchTimeout from "@/utils/fetchTimeOut";
import showToast, { ToastStatus } from "@/utils/toastHelper";

interface AuthState {
    loggedIn: boolean;
    fetching: boolean;
    user: {
        userid: number;
        email: string;
        username: string;
        rank: number;
        teamid: number;
        teamname: string;
    } | null;
    setUser: (user: AuthState["user"]) => void;
    logout: () => void;
    fetchUser: () => void;
};

const LOCAL_STORAGE_KEY = "userData";
const EXPIRY_KEY = "userExpiry";
const expiryHours = 24;

export const useAuthStore = create<AuthState>((set) => ({
    loggedIn: false,
    fetching: true,
    user: null, 

    setUser: (user) => {
        const expiry = Date.now() + 1000 * 60 * 60 * expiryHours;
        localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(user));
        localStorage.setItem(EXPIRY_KEY, expiry.toString());
        set({ user });
        set({ loggedIn: true });
    },

    logout: async () => {
        await fetchTimeout("/api/logout", 5000, new AbortController().signal);
        localStorage.removeItem(LOCAL_STORAGE_KEY);
        localStorage.removeItem(EXPIRY_KEY);
        set({ loggedIn: false, user: null });
    },

    fetchUser: async () => {
        let storedData = localStorage.getItem(LOCAL_STORAGE_KEY) || "";
        let expiry = localStorage.getItem(EXPIRY_KEY) || "";

        if (expiry && Number(expiry) < Date.now()) {
            expiry = "";
        }

        if (storedData && expiry) {
            const userData = JSON.parse(storedData);
            set({ loggedIn: true, user: userData, fetching: false });
            return;
        }

        try {
            const res = await fetchTimeout("/api/identify", 5000, new AbortController().signal);
            if (res.ok) {
                const user = await res.json();
                const newExpiry = Date.now() + 1000 * 60 * 60 * expiryHours;
                localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(user));
                localStorage.setItem(EXPIRY_KEY, newExpiry.toString());
                set({ user, loggedIn: true, fetching: false });
            } else {
                set({ loggedIn: false, user: null, fetching: false });
            }
        } catch (error: any) {
            set({ loggedIn: false, user: null, fetching: false });

            if (error.name === "AbortError") {
                showToast(ToastStatus.Failure, "verification timed out, reload!");
            } else {
                showToast(ToastStatus.Warning, "seems offline");
            }
        }
    },
}));