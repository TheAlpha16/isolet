import { create } from "zustand";
import fetchTimeout from "@/utils/fetchTimeOut";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { UserType } from "@/utils/types";

interface AuthState {
	fetching: boolean;
	user: UserType;
	setUser: (user: AuthState["user"]) => void;
	logout: () => void;
	fetchUser: () => void;
}

const LOCAL_STORAGE_KEY = "userData";
const EXPIRY_KEY = "userExpiry";
const expiryHours = 24;

export const useAuthStore = create<AuthState>((set) => ({
	fetching: true,
	user: {
		userid: -1,
		username: "",
		email: "",
		teamid: -1,
		teamname: "",
		rank: 0,
		score: 0,
	},

	setUser: (user) => {
		const expiry = Date.now() + 1000 * 60 * 60 * expiryHours;
		localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(user));
		localStorage.setItem(EXPIRY_KEY, expiry.toString());
		set({ user });
	},

	logout: async () => {
		await fetchTimeout("/api/logout", 60000, new AbortController().signal);
		localStorage.removeItem(LOCAL_STORAGE_KEY);
		localStorage.removeItem(EXPIRY_KEY);
		set({
			user: {
				userid: -1,
				username: "",
				email: "",
				teamid: -1,
				teamname: "",
				rank: 0,
				score: 0,
			},
		});
	},

	fetchUser: async () => {
		let storedData = localStorage.getItem(LOCAL_STORAGE_KEY) || "";
		let expiry = localStorage.getItem(EXPIRY_KEY) || "";

		if (expiry && Number(expiry) < Date.now()) {
			expiry = "";
		}

		if (storedData && expiry && JSON.parse(storedData).teamid !== -1) {
			const userData = JSON.parse(storedData);
			set({ user: userData, fetching: false });
			return;
		}

		try {
			const res = await fetchTimeout(
				"/api/identify",
				60000,
				new AbortController().signal
			);
			if (res.ok) {
				const user = await res.json();
				const newExpiry = Date.now() + 1000 * 60 * 60 * expiryHours;
				localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(user));
				localStorage.setItem(EXPIRY_KEY, newExpiry.toString());
				set({ user, fetching: false });
			} else {
				set({ fetching: false });
			}
		} catch (error: any) {
			set({ fetching: false });

			if (error.name === "AbortError") {
				showToast(
					ToastStatus.Failure,
					"verification timed out, reload!"
				);
			} else {
				showToast(ToastStatus.Warning, "seems offline");
			}
		}
	},
}));
