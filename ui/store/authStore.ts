import fetchTimeout from "@/utils/fetchTimeOut";
import showToast, { ToastStatus } from "@/utils/toastHelper";
import { UserType } from "@/utils/types";
import { create } from "zustand";

interface AuthState {
  fetching: boolean;
  user: UserType;
  setUser: (user: AuthState["user"]) => void;
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
}));
